package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	systemStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Bold(true)
	reqStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E8F0"))
	blockedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171")).Bold(true)
	statusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).MarginTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#22C55E")).
			Padding(1, 2).
			Width(72).
			MaxWidth(72)

	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FB923C")).Bold(true)
	keyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FACC15")).Bold(true)
)

type uiModel struct {
	serverCmd      string
	logs           []string
	isIntercepting bool
	interceptData  interceptMsg
}

func initialUIModel(cmd string) uiModel {
	return uiModel{
		serverCmd: cmd,
		logs:      []string{},
	}
}

func (m uiModel) Init() tea.Cmd {
	return nil
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.isIntercepting {
			switch msg.String() {
			case "a", "s":
				m.isIntercepting = false
				pendingDecisionChan <- true
				m.logs = append(m.logs, blockedStyle.Render("=> ALLOWED: ")+m.interceptData.ToolName)
				return m, nil
			case "d":
				m.isIntercepting = false
				pendingDecisionChan <- false
				m.logs = append(m.logs, blockedStyle.Render("=> DENIED: ")+m.interceptData.ToolName)
				return m, nil
			}
		}

		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case logMsg:
		m.logs = append(m.logs, string(msg))
		if len(m.logs) > 15 {
			m.logs = m.logs[1:]
		}
		return m, nil

	case interceptMsg:
		m.isIntercepting = true
		m.interceptData = msg
		return m, nil
	}
	return m, nil
}

func (m uiModel) View() string {
	var b strings.Builder

	// Header
	b.WriteString(systemStyle.Render(fmt.Sprintf("[SYSTEM] Proxy active. Wrapping: %s", m.serverCmd)))
	b.WriteString("\n\n")

	// Logs
	for _, l := range m.logs {
		if strings.HasPrefix(l, "REQ:") {
			b.WriteString(reqStyle.Render(l))
		} else if strings.Contains(l, "BLOCKED") || strings.Contains(l, "DENIED") || strings.Contains(l, "MALFORMED") {
			b.WriteString(blockedStyle.Render(l))
		} else {
			b.WriteString(l)
		}
		b.WriteString("\n") // Append newline separately
	}

	// Dynamic Footer
	if m.isIntercepting {
		b.WriteString("\n")
		b.WriteString(m.renderPromptBox())
	} else {
		b.WriteString(statusStyle.Render("\nStatus: Active | [q] Quit\n"))
	}

	return b.String()
}

func (m uiModel) renderPromptBox() string {
	content := fmt.Sprintf("%s INTERCEPTED ACTION: %s\n\nPolicy: %s\n\n%s Allow Once   %s Deny   %s Save Rule & Always Allow",
		titleStyle.Render("[!]"),
		m.interceptData.ToolName,
		m.interceptData.Reason,
		keyStyle.Render("[a]"),
		keyStyle.Render("[d]"),
		keyStyle.Render("[s]"),
	)
	return boxStyle.Render(content)
}
