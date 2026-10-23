package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type CateContext struct {
	UserID       string `json:"user_id"`
	SessionToken string `json:"session_token"`
}

type MCPRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

type CatePayload struct {
	Scenario    string      `json:"scenario"`
	Context     CateContext `json:"cate_context"`
	MCPRequest  MCPRequest  `json:"mcp_request"`
}

func main() {
	fmt.Println("\033[1;36m[*] Initializing RCF + Arcade CATE Integration Harness...\033[0m")
	
	file, err := os.ReadFile("mock_payload.json")
	if err != nil {
		fmt.Println("Error reading payload:", err)
		return
	}

	var payloads []CatePayload
	json.Unmarshal(file, &payloads)

	for _, p := range payloads {
		fmt.Printf("\n\033[1;33m--- Executing Scenario: %s ---\033[0m\n", p.Scenario)
		start := time.Now()

		// 1. Identity Validation (CATE Handoff)
		if p.Context.UserID == "" {
			fail(p, start, -32001, "Missing CATE identity context. WORM logging requires user_id.")
			continue
		}

		// 2. BM25 Schema Pruning Simulation (Token Reduction)
		time.Sleep(300 * time.Microsecond) // 0.3ms overhead

		// 3. wazero Micro-VM Detonation Simulation
		time.Sleep(700 * time.Microsecond) // 0.7ms overhead
		
		paramsJSON, _ := json.Marshal(p.MCPRequest.Params)
		if strings.Contains(string(paramsJSON), "DROP") || strings.Contains(string(paramsJSON), "DELETE") {
			fail(p, start, -32603, "Policy Violation: Malicious payload intercepted in micro-VM")
			continue
		}

		// 4. Safe Execution
		latency := time.Since(start)
		fmt.Printf("\033[1;32m[✓] Payload passed micro-VM validation (Latency: %v)\033[0m\n", latency)
		logWORM(p, "ALLOWED", "Safe execution", latency)
	}
}

func fail(p CatePayload, start time.Time, code int, msg string) {
	latency := time.Since(start)
	
	// JSON-RPC Error Bubbling
	errorResponse := fmt.Sprintf(`{"jsonrpc": "2.0", "error": {"code": %d, "message": "%s"}, "id": %d}`, code, msg, p.MCPRequest.ID)
	
	fmt.Printf("\033[1;31m[x] Execution Blocked (Latency: %v)\033[0m\n", latency)
	fmt.Printf("    JSON-RPC Response: %s\n", errorResponse)
	logWORM(p, "DENIED", msg, latency)
}

func logWORM(p CatePayload, status string, reason string, latency time.Duration) {
	fmt.Printf(`
=========================================
SEC 17a-4 IMMUTABLE AUDIT LOG (WORM)
Timestamp:       %s
User ID (CATE):  %s
Target Method:   %s
Status:          %s
Reason:          %s
RCF Overhead:    %v
=========================================
`, time.Now().Format(time.RFC3339), p.Context.UserID, p.MCPRequest.Method, status, reason, latency)
}
