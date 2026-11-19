package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
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

var universalTraversePattern = regexp.MustCompile(`(\.\.\/|^\/etc\/|^\/var\/run\/)`)

type InspectionResult struct {
	Blocked bool
	Action  string
	Message string
}

func InspectPayload(toolName string, args map[string]interface{}) InspectionResult {
	for argKey, rawVal := range args {
		valStr := fmt.Sprintf("%v", rawVal)

		// ---------------------------------------------------------
		// LAYER 2A: SEMANTIC PATH CANONICALIZATION (Pre-Inspection)
		// ---------------------------------------------------------
		// 1. Decode percent-encoding (e.g., %2e%2e%2f -> ../)
		if unescaped, err := url.QueryUnescape(valStr); err == nil {
			valStr = unescaped
		}
		// 2. Canonicalize path traversals (collapses /a/b/../../etc/passwd -> /etc/passwd)
		cleanedPath := filepath.Clean(valStr)

		// Check canonical traversal against Layer 1 guardrails
		if universalTraversePattern.MatchString(valStr) || universalTraversePattern.MatchString(cleanedPath) {
			return InspectionResult{
				Blocked: true,
				Action:  "DENY",
				Message: fmt.Sprintf("Layer 2 Canonicalization [Security]: Directory traversal detected in '%s' (Resolved: %s).", argKey, cleanedPath),
			}
		}

		// ---------------------------------------------------------
		// LAYER 2B: NETWORK / SSRF CANONICALIZATION
		// ---------------------------------------------------------
		// Detect decimal IP (e.g., 2130706433 -> 127.0.0.1) and cloud metadata
		if strings.Contains(valStr, "169.254.169.254") || strings.Contains(valStr, "2130706433") {
			return InspectionResult{
				Blocked: true,
				Action:  "DENY",
				Message: fmt.Sprintf("Layer 2 Egress Guardrail: SSRF / Private IP bypass attempt detected in '%s'.", argKey),
			}
		}

		// ---------------------------------------------------------
		// LAYER 1: UNIVERSAL GUARDRAILS (DLP & Semantic Traps)
		// ---------------------------------------------------------
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

	// LAYER 2C: GRANULAR CONSTRAINTS (policy.yaml)
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
				argVal := fmt.Sprintf("%v", rawVal)
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
