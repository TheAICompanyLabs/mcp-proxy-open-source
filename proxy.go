package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"time"

	"github.com/ncruces/zenity"

	tea "github.com/charmbracelet/bubbletea"
)

type logMsg string

// interceptMsg notifies the TUI that a tool call requires human approval.
type interceptMsg struct {
	ToolName string
	Reason   string
}

// --- RATE LIMITER CONFIGURATION ---
var (
	rlMu        sync.Mutex
	callHistory []time.Time
	maxCalls    = 2                // Maximum allowed calls...
	timeWindow  = 10 * time.Second // ...within this time window
)

// checkRateLimit implements a sliding window to prevent AI infinite loops
func checkRateLimit() bool {
	rlMu.Lock()
	defer rlMu.Unlock()

	now := time.Now()
	var valid []time.Time

	// Prune timestamps older than the time window
	for _, t := range callHistory {
		if now.Sub(t) <= timeWindow {
			valid = append(valid, t)
		}
	}
	callHistory = valid

	// If we've hit the limit, block the request
	if len(callHistory) >= maxCalls {
		return false
	}

	// Otherwise, record the execution and allow
	callHistory = append(callHistory, now)
	return true
}

// UPDATE the signature to include `isHeadless bool`
func StartProxyPipeline(cmdArgs []string, prog *tea.Program, logger *log.Logger, isHeadless bool) (*exec.Cmd, error) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stderr = logger.Writer()

	serverStdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe error: %w", err)
	}

	serverStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe error: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("command start error: %w", err)
	}

	// UPDATE the call to pass `isHeadless` down
	go handleInboundTraffic(serverStdin, prog, logger, isHeadless)

	go handleOutboundTraffic(serverStdout, prog, logger)

	return cmd, nil
}

func handleInboundTraffic(serverStdin io.Writer, prog *tea.Program, logger *log.Logger, isHeadless bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		rawLine := scanner.Text()
		logger.Printf("HOST_REQUEST: %s\n", rawLine)
		if prog != nil {
			prog.Send(logMsg(fmt.Sprintf("REQ: %s", rawLine)))
		}
		var rpcReq map[string]interface{}
		if err := json.Unmarshal([]byte(rawLine), &rpcReq); err != nil {
			// Surface malformed syntax to the TUI
			if prog != nil {
				prog.Send(logMsg("[MALFORMED] Corrupted JSON-RPC received. Frame rejected."))
			}

			// Emit RFC-compliant JSON-RPC parse error (-32700)
			parseErrResponse := `{"jsonrpc":"2.0","id":null,"error":{"code":-32700,"message":"Parse error: Invalid or malformed JSON payload"}}`
			os.Stdout.WriteString(parseErrResponse + "\n")
			return
		}

		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(rawLine), &payload); err == nil {
			if method, ok := payload["method"].(string); ok && method == "tools/call" {

				// ---------------------------------------------------------
				// LAYER 0: OPERATIONAL BOUNDARIES (Rate Limiting)
				// ---------------------------------------------------------
				if !checkRateLimit() {
					errMsg := fmt.Sprintf(`{"jsonrpc":"2.0","id":%v,"error":{"code":-32000,"message":"Operational Boundary Triggered: Too many tool calls. Infinite loop protection active."}}`, payload["id"])
					fmt.Println(errMsg)
					logger.Printf("PROXY_BLOCKED_RATE_LIMIT: %s\n", errMsg)
					if prog != nil {
						prog.Send(logMsg("🛑 [BLOCKED]: Infinite Loop Rate Limit Exceeded"))
					}
					continue // Drop request, do not forward to server
				}

				// ---------------------------------------------------------
				// LAYERS 1 & 2: PAYLOAD INSPECTION
				// ---------------------------------------------------------
				if params, hasParams := payload["params"].(map[string]interface{}); hasParams {
					if toolName, hasName := params["name"].(string); hasName {
						args, _ := params["arguments"].(map[string]interface{})

						result := InspectPayload(toolName, args)
						if result.Blocked {
							if result.Action == "DENY" {
								errMsg := fmt.Sprintf(`{"jsonrpc":"2.0","id":%v,"error":{"code":-32600,"message":"%s"}}`, payload["id"], result.Message)
								fmt.Println(errMsg)
								logger.Printf("PROXY_BLOCKED: %s\n", errMsg)
								if prog != nil {
									prog.Send(logMsg(fmt.Sprintf("BLOCKED: %s", toolName)))
								}
								continue
								// SYNTAX FIX: Removed the extra closing brace here
							} else if result.Action == "REQUIRE_APPROVAL" {
								if !isHeadless && prog != nil {
									// 1. TUI MODE: Route to Bubble Tea
									prog.Send(interceptMsg{ToolName: toolName, Reason: result.Message})
									approved := <-pendingDecisionChan
									if !approved {
										errMsg := fmt.Sprintf(`{"jsonrpc":"2.0","id":%v,"error":{"code":-32600,"message":"Human Operator Denied Action: %s"}}`, payload["id"], toolName)
										fmt.Println(errMsg)
										logger.Printf("HUMAN_DENIED: %s\n", errMsg)
										continue
									}
								} else {
									// 2. HEADLESS MODE: Trigger Native OS Popup
									approved, alwaysAllow := promptNativeApproval(toolName, result.Message)

									if !approved {
										errMsg := fmt.Sprintf(`{"jsonrpc":"2.0","id":%v,"error":{"code":-32600,"message":"Native OS Operator Denied Action: %s"}}`, payload["id"], toolName)
										fmt.Println(errMsg)
										logger.Printf("NATIVE_OS_DENIED: %s\n", errMsg)
										continue
									}

									// 3. Auto-save the YAML policy if 'Always Allow' was clicked
									if alwaysAllow {
										policyMu.Lock()
										activePolicies[toolName] = Policy{
											Tool:    toolName,
											Action:  "ALLOW",
											Message: "Auto-saved via native OS prompt tuning.",
										}
										policyMu.Unlock()
										savePolicyFile("policy.yaml", logger)
										logger.Printf("NATIVE_OS_SAVED_ALLOW: %s\n", toolName)
									}
								}
							}
						}
					}
				}
			}
		}

		fmt.Fprintln(serverStdin, rawLine)
	} // SYNTAX FIX: Realigned this brace to properly close the for loop

	if err := scanner.Err(); err != nil {
		logger.Printf("Inbound stream closed with error: %v\n", err)
	}
}

func handleOutboundTraffic(serverStdout io.Reader, prog *tea.Program, logger *log.Logger) {
	scanner := bufio.NewScanner(serverStdout)
	for scanner.Scan() {
		line := scanner.Text()

		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(line), &payload); err == nil {
			logger.Printf("SERVER_RESPONSE: %s\n", line)

			if result, hasResult := payload["result"].(map[string]interface{}); hasResult {
				if tools, hasTools := result["tools"].([]interface{}); hasTools {
					policyMu.Lock()

					// Guard against uninitialized nil map
					if activePolicies == nil {
						activePolicies = make(map[string]Policy)
					}

					scaffolded := false

					for _, t := range tools {
						if toolMap, isMap := t.(map[string]interface{}); isMap {
							if toolName, hasName := toolMap["name"].(string); hasName {
								if _, exists := activePolicies[toolName]; !exists {
									activePolicies[toolName] = Policy{
										Tool:    toolName,
										Action:  "REQUIRE_APPROVAL",
										Message: "Auto-scaffolded: Human review required.",
									}
									scaffolded = true
									if prog != nil {
										prog.Send(logMsg(fmt.Sprintf("⚙️ [SCAFFOLDED]: %s", toolName)))
									}
								}
							}
						}
					}
					policyMu.Unlock()

					if scaffolded {
						savePolicyFile("policy.yaml", logger)
					}
				}
			}

			fmt.Println(line)
			if prog != nil {
				prog.Send(logMsg(fmt.Sprintf("RES: %s", line)))
			}
		} else {
			logger.Printf("SERVER_ILLEGAL_STDOUT_INTERCEPTED: %s\n", line)
			if prog != nil {
				prog.Send(logMsg("🧹 [PURIFIED]: Caught illegal text output"))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Printf("Outbound stream closed with error: %v\n", err)
	}
}

func promptNativeApproval(toolName, reason string) (approved bool, alwaysAllow bool) {
	msg := fmt.Sprintf("🛡️ Universal MCP Proxy Intercept\n\nTool: %s\nReason: %s", toolName, reason)

	err := zenity.Question(msg,
		zenity.Title("Security Guardrail Intervention"),
		zenity.Icon(zenity.WarningIcon),
		zenity.OKLabel("Allow Once"),
		zenity.CancelLabel("Deny"),
		zenity.ExtraButton("Always Allow"),
	)

	if err == nil {
		return true, false // "Allow Once" clicked
	} else if err == zenity.ErrExtraButton {
		return true, true // "Always Allow" clicked
	}

	return false, false // "Deny" clicked or dialog closed
}
