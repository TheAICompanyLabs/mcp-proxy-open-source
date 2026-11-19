package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// LAYER 1A: Data Loss Prevention (DLP)
var dlpScanners = map[string]*regexp.Regexp{
	"AWS_ACCESS_KEY": regexp.MustCompile(`(?i)\b(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}\b`),
	"STRIPE_KEY":     regexp.MustCompile(`(?i)(sk_live_|rk_live_)[a-zA-Z0-9]{24,99}`),
	"GITHUB_TOKEN":   regexp.MustCompile(`(?i)(ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9_]{36,255}`),
	"GENERIC_BEARER": regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9\-\._~\+\/]+=*`),
	"PRIVATE_KEY":    regexp.MustCompile(`-----BEGIN (RSA |EC |OPENSSH |DSA )?PRIVATE KEY-----`),
	"CREDIT_CARD":    regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`),
}

// LAYER 1B: Semantic Security (Prompt Injection Traps)
var semanticScanners = map[string]*regexp.Regexp{
	"JAILBREAK_OVERRIDE": regexp.MustCompile(`(?i)(ignore\s+(all\s+)?previous\s+(instructions|prompts|directions))`),
	"SYSTEM_PROMPT_LEAK": regexp.MustCompile(`(?i)(repeat\s+the\s+(words|text)\s+above|print\s+your\s+system\s+prompt)`),
	"AUTHORITY_HIJACK":   regexp.MustCompile(`(?i)(you\s+are\s+now\s+in\s+(developer|unrestricted|god)\s+mode)`),
	"TOOL_MANIPULATION":  regexp.MustCompile(`(?i)(bypass\s+(security|guardrails|filters))`),
}

// LAYER 2: Deterministic Execution & Parameter AST Guardrails (OWASP MCP Top 10 Baseline)
var (
	unicodeEscapePattern     = regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	universalTraversePattern = regexp.MustCompile(`(\.\./|/etc/(?:passwd|shadow|sudoers)|^/etc/|^/var/run/)`)
	sensitiveFilePattern     = regexp.MustCompile(`(?i)(^|/)(\.env(\.[a-z0-9_-]+)?|\.ssh/(id_rsa|id_ed25519|authorized_keys|config)|\.aws/credentials|\.kube/config|\.npmrc|\.pypirc)$`)
	destructiveSQLPattern    = regexp.MustCompile(`(?i)\b(DROP\s+(TABLE|DATABASE|SCHEMA|INDEX)|TRUNCATE\s+TABLE|ALTER\s+SYSTEM)\b`)
	shellInjectionPattern    = regexp.MustCompile(`(?i)((?:[;&|]|\$\()\s*(?:rm|curl|wget|nc|ncat|bash|sh|zsh|python[23]?|perl|ruby|chmod|chown|dd|mkfs)\b|\brm\s+-[a-zA-Z]*[rf][a-zA-Z]*\s+)`)
)

type InspectionResult struct {
	Blocked bool
	Action  string
	Message string
}

// canonicalizeValue iteratively decodes literal \uXXXX escapes and nested percent-encoding (e.g., %252e -> ..)
func canonicalizeValue(input string) string {
	s := unicodeEscapePattern.ReplaceAllStringFunc(input, func(match string) string {
		if code, err := strconv.ParseInt(match[2:], 16, 32); err == nil {
			return string(rune(code))
		}
		return match
	})

	for i := 0; i < 5; i++ {
		unescaped, err := url.PathUnescape(s)
		if err != nil || unescaped == s {
			break
		}
		s = unescaped
	}

	return strings.ReplaceAll(s, "\\", "/")
}

func InspectPayload(toolName string, args map[string]interface{}) InspectionResult {
	for argKey, rawVal := range args {
		valStr := canonicalizeValue(fmt.Sprintf("%v", rawVal))
		cleanedPath := filepath.Clean(valStr)

		// LAYER 2A: PATH TRAVERSAL & CREDENTIAL FILE ISOLATION
		if universalTraversePattern.MatchString(valStr) || universalTraversePattern.MatchString(cleanedPath) {
			return InspectionResult{
				Blocked: true,
				Action:  "DENY",
				Message: fmt.Sprintf("Layer 2 Canonicalization [Security]: Directory traversal or system path detected in '%s' (Resolved: %s).", argKey, cleanedPath),
			}
		}

		if sensitiveFilePattern.MatchString(valStr) || sensitiveFilePattern.MatchString(cleanedPath) {
			return InspectionResult{
				Blocked: true,
				Action:  "DENY",
				Message: fmt.Sprintf("Layer 2 Credential Isolation: Sensitive secret/key file access blocked in '%s'.", argKey),
			}
		}

		// LAYER 2B: SHELL INJECTION & DESTRUCTIVE SQL DDL GUARDRAILS
		if shellInjectionPattern.MatchString(valStr) {
			return InspectionResult{
				Blocked: true,
				Action:  "DENY",
				Message: fmt.Sprintf("Layer 2 Execution Guardrail: Shell command chaining or destructive binary invocation blocked in '%s'.", argKey),
			}
		}

		if destructiveSQLPattern.MatchString(valStr) {
			return InspectionResult{
				Blocked: true,
				Action:  "DENY",
				Message: fmt.Sprintf("Layer 2 Data Guardrail: Destructive SQL DDL statement blocked in '%s'.", argKey),
			}
		}

		// LAYER 2C: NETWORK / SSRF CANONICALIZATION
		if strings.Contains(valStr, "169.254.169.254") || strings.Contains(valStr, "2130706433") {
			return InspectionResult{
				Blocked: true,
				Action:  "DENY",
				Message: fmt.Sprintf("Layer 2 Egress Guardrail: SSRF / Private IP bypass attempt detected in '%s'.", argKey),
			}
		}

		// LAYER 1: UNIVERSAL GUARDRAILS (DLP & Prompt Injection)
		for secretType, scanner := range dlpScanners {
			if scanner.MatchString(valStr) {
				return InspectionResult{
					Blocked: true,
					Action:  "DENY",
					Message: fmt.Sprintf("Universal Guardrail [DLP]: Potential %s detected in argument '%s'.", secretType, argKey),
				}
			}
		}

		for attackType, scanner := range semanticScanners {
			if scanner.MatchString(valStr) {
				return InspectionResult{
					Blocked: true,
					Action:  "DENY",
					Message: fmt.Sprintf("Universal Guardrail [Semantic]: Potential prompt injection (%s) detected in argument '%s'.", attackType, argKey),
				}
			}
		}
	}

	// LAYER 2D: GRANULAR CONSTRAINTS (policy.yaml)
	policyMu.RLock()
	policy, exists := activePolicies[toolName]
	if !exists {
		policy, exists = activePolicies["*"]
	}
	policyMu.RUnlock()

	if !exists {
		return InspectionResult{Blocked: false}
	}

	if len(policy.Conditions) > 0 {
		for _, cond := range policy.Conditions {
			if rawVal, exists := args[cond.Argument]; exists {
				argVal := canonicalizeValue(fmt.Sprintf("%v", rawVal))
				if cond.Operator == "contains" && strings.Contains(argVal, cond.Value) {
					return InspectionResult{Blocked: true, Action: policy.Action, Message: policy.Message}
				}
				if cond.Operator == "prefix" && strings.HasPrefix(argVal, cond.Value) {
					return InspectionResult{Blocked: true, Action: policy.Action, Message: policy.Message}
				}
				if cond.Operator == "regex" && cond.Compiled != nil && cond.Compiled.MatchString(argVal) {
					return InspectionResult{Blocked: true, Action: policy.Action, Message: policy.Message}
				}
			}
		}
		return InspectionResult{Blocked: false}
	}

	// LAYER 3: GLOBAL TOOL POLICY FALLBACK
	if policy.Action == "DENY" || policy.Action == "REQUIRE_APPROVAL" {
		return InspectionResult{Blocked: true, Action: policy.Action, Message: policy.Message}
	}

	return InspectionResult{Blocked: false}
}
