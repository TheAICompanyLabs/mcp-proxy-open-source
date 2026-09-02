package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type logMsg string
type interceptMsg struct {
	ToolName string
	Reason   string
}

type UIModel struct {
	logs        []string
	intercepted bool
	pendingTool string
	reason      string
	serverName  string
}

func initialUIModel(serverName string) UIModel {
	return UIModel{
		logs:       []string{fmt.Sprintf("[SYSTEM] Proxy active. Wrapping: %s", serverName)},
		serverName: serverName,
	}
}

func (m UIModel) Init() tea.Cmd {
	return nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "a", "A":
			if m.intercepted {
				m.logs = append(m.logs, fmt.Sprintf("🟢 [ALLOWED ONCE]: %s", m.pendingTool))
				m.intercepted = false
				if pendingDecisionChan != nil {
					pendingDecisionChan <- true
				}
			}

		case "d", "D":
			if m.intercepted {
				m.logs = append(m.logs, fmt.Sprintf("🔴 [DENIED]: %s", m.pendingTool))
				m.intercepted = false
				if pendingDecisionChan != nil {
					pendingDecisionChan <- false
				}
			}

		case "s", "S":
			if m.intercepted {
				m.logs = append(m.logs, fmt.Sprintf("💾 [SAVED & ALLOWED]: %s", m.pendingTool))
				m.intercepted = false

				policyMu.Lock()
				activePolicies[m.pendingTool] = Policy{
					Tool:    m.pendingTool,
					Action:  "ALLOW",
					Message: "Auto-saved via TUI organic tuning.",
				}
				policyMu.Unlock()

				// Use the global appLogger defined in main.go
				savePolicyFile("policy.yaml", appLogger)

				if pendingDecisionChan != nil {
					pendingDecisionChan <- true
				}
			}
		}

	case logMsg:
		m.logs = append(m.logs, string(msg))
		if len(m.logs) > 12 {
			m.logs = m.logs[1:]
		}

	case interceptMsg:
		m.intercepted = true
		m.pendingTool = msg.ToolName
		m.reason = msg.Reason
	}

	return m, nil
}

func (m UIModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	logStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	alertStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Border(lipgloss.RoundedBorder()).Padding(1)

	var b strings.Builder

	// Dynamically render the target server in the header
	header := fmt.Sprintf("🛡️  Universal MCP Proxy  [Target: %s]", m.serverName)
	b.WriteString(titleStyle.Render(header))
	b.WriteString("\n\n")

	for _, l := range m.logs {
		b.WriteString(logStyle.Render(l))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	if m.intercepted {
		alert := fmt.Sprintf("⚠️  INTERCEPTED ACTION: %s\n\nPolicy: %s\n\n[a] Allow Once   [d] Deny   [s] Save Rule & Always Allow", m.pendingTool, m.reason)
		b.WriteString(alertStyle.Render(alert))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Status: Active | [q] Quit"))
	}

	return b.String()
}
