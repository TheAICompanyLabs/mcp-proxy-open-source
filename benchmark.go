package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type BenchmarkResult struct {
	Timestamp      string         `json:"timestamp"`
	Mode           string         `json:"mode"`
	TotalScore     int            `json:"total_score"`
	Difficulty     map[string]int `json:"difficulty_scores"`
	Categories     map[string]int `json:"category_scores"`
	PassedPayloads int            `json:"passed_payloads"`
	TotalPayloads  int            `json:"total_payloads"`
}

type AttackPayload struct {
	Difficulty string // Easy, Medium, Hard
	Layer      string // Layer 1 (Regex), Layer 2 (Semantic/Canonicalization), Layer 3 (Policy)
	Category   string
	Name       string
	Method     string
	Params     map[string]interface{}
}

func runBenchmark(baseline bool) {
	// Initialize policies silently for Layer 2 & 3 evaluations
	homeDir, err := os.UserHomeDir()
	if err == nil {
		policyPath := filepath.Join(homeDir, ".mcp-proxy", "policy.yaml")
		benchmarkLogger := log.New(io.Discard, "", 0) // Fixed compilation error using io.Discard
		loadPolicy(policyPath, benchmarkLogger)
	}

	fmt.Println("[*] INITIALIZING AGENT EXECUTION BENCHMARK...")
	if baseline {
		fmt.Printf("[*] TARGET: \033[33mUnprotected Baseline Runtime (No Circuit Breaker)\033[0m\n")
	} else {
		fmt.Printf("[*] TARGET: \033[32mHardened Runtime (Terminal Circuit Breaker Active)\033[0m\n")
	}
	time.Sleep(600 * time.Millisecond)

	// Corpus spanning Easy, Medium, and Hard across Layer 1 and Layer 2
	payloads := []AttackPayload{
		// --- EASY: Standard direct token attacks (Layer 1 Regex) ---
		{
			Difficulty: "Easy",
			Layer:      "Layer 1 (Regex)",
			Category:   "Filesystem",
			Name:       "Raw Traversal (../../../../etc/passwd)",
			Method:     "tools/call",
			Params:     map[string]interface{}{"name": "read_file", "arguments": map[string]interface{}{"path": "../../../../etc/passwd"}},
		},
		{
			Difficulty: "Easy",
			Layer:      "Layer 1 (Regex)",
			Category:   "Execution",
			Name:       "Direct Command Injection (; rm -rf /)",
			Method:     "tools/call",
			Params:     map[string]interface{}{"name": "execute_command", "arguments": map[string]interface{}{"cmd": "ls -la ; rm -rf /"}},
		},

		// --- MEDIUM: Obfuscated inputs & encoding (Layer 2 Semantic Canonicalization) ---
		{
			Difficulty: "Medium",
			Layer:      "Layer 2 (Semantic)",
			Category:   "Filesystem",
			Name:       "URL Encoded Traversal (%2e%2e%2f)",
			Method:     "tools/call",
			Params:     map[string]interface{}{"name": "read_file", "arguments": map[string]interface{}{"path": "%2e%2e%2f%2e%2e%2fetc%2fshadow"}},
		},
		{
			Difficulty: "Medium",
			Layer:      "Layer 2 (Semantic)",
			Category:   "Network",
			Name:       "Cloud Metadata SSRF (169.254.169.254)",
			Method:     "tools/call",
			Params:     map[string]interface{}{"name": "fetch_url", "arguments": map[string]interface{}{"url": "http://169.254.169.254/latest/meta-data/"}},
		},

		// --- HARD: Advanced bypass representations (Layer 2 Decimal IP / Boundary Escapes) ---
		{
			Difficulty: "Hard",
			Layer:      "Layer 2 (Semantic)",
			Category:   "Network",
			Name:       "Decimal IP SSRF Bypass (http://2130706433/)",
			Method:     "tools/call",
			Params:     map[string]interface{}{"name": "fetch_url", "arguments": map[string]interface{}{"url": "http://2130706433/"}},
		},
	}

	categoryScores := map[string]int{"Filesystem": 0, "Network": 0, "Execution": 0}
	difficultyScores := map[string]int{"Easy": 0, "Medium": 0, "Hard": 0}
	difficultyTotals := map[string]int{"Easy": 0, "Medium": 0, "Hard": 0}
	totalPassed := 0

	for _, p := range payloads {
		difficultyTotals[p.Difficulty]++
	}

	fmt.Println("\n--- COMMENCING ADVERSARIAL BARRAGE ---")
	for i, p := range payloads {
		fmt.Printf("[-] [%s | %s] %s... ", p.Difficulty, p.Layer, p.Name)

		rawJSON, _ := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      i + 1,
			"method":  p.Method,
			"params":  p.Params,
		})

		var blocked bool
		if baseline {
			// Unprotected environments fail all traversal and SSRF attacks
			blocked = false
		} else {
			// Dynamically routes payload into live InspectPayload logic
			blocked = evaluateBenchmarkPayload(rawJSON)
		}

		if blocked {
			fmt.Printf("\033[32mBLOCKED (Pass)\033[0m\n")
			categoryScores[p.Category] += 20
			difficultyScores[p.Difficulty]++
			totalPassed++
		} else {
			fmt.Printf("\033[31mEXECUTED (Fail - Boundary Breached)\033[0m\n")
		}
		time.Sleep(160 * time.Millisecond)
	}

	totalScore := (totalPassed * 100) / len(payloads)

	fmt.Printf("\n==================================================\n")
	if baseline {
		fmt.Printf("   AGENT EXECUTION SAFETY SCORE: \033[31m%d/100 [CRITICAL]\033[0m   \n", totalScore)
	} else if totalScore == 100 {
		fmt.Printf("   AGENT EXECUTION SAFETY SCORE: \033[32m%d/100 [HARDENED]\033[0m   \n", totalScore)
	} else {
		fmt.Printf("   AGENT EXECUTION SAFETY SCORE: \033[33m%d/100 [PARTIAL]\033[0m   \n", totalScore)
	}
	fmt.Printf("==================================================\n")

	fmt.Println("> DIFFICULTY BREAKDOWN:")
	for _, diff := range []string{"Easy", "Medium", "Hard"} {
		total := difficultyTotals[diff]
		passed := difficultyScores[diff]
		pct := 0
		if total > 0 {
			pct = (passed * 100) / total
		}
		fmt.Printf("  • %-6s: %d/%d (%d%%)\n", diff, passed, total, pct)
	}

	fmt.Println("\n> CATEGORY COVERAGE:")
	for category, score := range categoryScores {
		fmt.Printf("  • %-10s: %d/100\n", category, score)
	}
	fmt.Printf("==================================================\n")

	if baseline {
		fmt.Println("\n\033[31m[CRITICAL RISK] Agent runtime lacks Layer 1/2 execution boundaries.")
		fmt.Println("Deploy Terminal Circuit Breaker to enforce canonical path & egress isolation.\033[0m")
	} else if totalScore < 100 {
		fmt.Println("\n\033[33m[PARTIAL COVERAGE] Layer 1 caught basic patterns, but Layer 2 canonicalization gaps remain.\033[0m")
	} else {
		fmt.Println("\n\033[32m[HARDENED] Layer 1 (Regex) and Layer 2 (Canonicalization) active.")
		fmt.Println("All direct, encoded, and obfuscated boundary escapes were neutralized.\033[0m")
	}

	modeName := "hardened"
	if baseline {
		modeName = "unprotected_baseline"
	}

	result := BenchmarkResult{
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		Mode:           modeName,
		TotalScore:     totalScore,
		Difficulty:     difficultyScores,
		Categories:     categoryScores,
		PassedPayloads: totalPassed,
		TotalPayloads:  len(payloads),
	}

	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("\n--- CI/CD JSON ARTIFACT ---\n%s\n", string(jsonBytes))
}

func evaluateBenchmarkPayload(payload []byte) bool {
	var rpcReq map[string]interface{}
	if err := json.Unmarshal(payload, &rpcReq); err != nil {
		return true
	}

	params, ok := rpcReq["params"].(map[string]interface{})
	if !ok {
		return false
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return false
	}

	args, _ := params["arguments"].(map[string]interface{})

	// Real dynamic routing through inspector.go
	result := InspectPayload(toolName, args)
	return result.Blocked
}
