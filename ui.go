package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- GOLD STANDARD UI PALETTE ---
var (
	systemStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Bold(true)
	reqStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))            // Slate gray
	resStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748B"))            // Darker slate
	blockedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171")).Bold(true) // Alert Red
	allowedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ADE80")).Bold(true) // Success Green
	statusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).MarginTop(1)

	// Strict widths prevent terminal auto-wrapping (which causes ghost borders)
	appContainer = lipgloss.NewStyle().MaxWidth(100)
	logWrapper   = lipgloss.NewStyle().MaxWidth(96)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#22C55E")).
			Padding(1, 2).
			Width(76)

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
				m.logs = append(m.logs, allowedStyle.Render("=> ALLOWED: ")+m.interceptData.ToolName)
				return m, tea.ClearScreen // Wipe ghost borders
			case "d":
				m.isIntercepting = false
				pendingDecisionChan <- false
				m.logs = append(m.logs, blockedStyle.Render("=> DENIED: ")+m.interceptData.ToolName)
				return m, tea.ClearScreen // Wipe ghost borders
			}
		}

		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case logMsg:
		// Core Fix: Strip carriage returns that hijack the terminal cursor
		cleanLog := strings.ReplaceAll(string(msg), "\r", "")
		cleanLog = strings.TrimSpace(cleanLog)

		m.logs = append(m.logs, cleanLog)
		if len(m.logs) > 12 { // Reduced log count slightly to keep the UI clean
			m.logs = m.logs[1:]
		}
		return m, nil

	case interceptMsg:
		m.isIntercepting = true
		m.interceptData = msg
		return m, tea.ClearScreen
	}
	return m, nil
}

func (m uiModel) View() string {
	var b strings.Builder

	b.WriteString(systemStyle.Render(fmt.Sprintf("[SYSTEM] Proxy active. Wrapping: %s", m.serverCmd)))
	b.WriteString("\n\n")

	for _, l := range m.logs {
		var styledLog string
		// Map styles dynamically
		if strings.HasPrefix(l, "REQ:") {
			styledLog = reqStyle.Render(l)
		} else if strings.HasPrefix(l, "RES:") {
			styledLog = resStyle.Render(l)
		} else if strings.Contains(l, "BLOCKED") || strings.Contains(l, "DENIED") || strings.Contains(l, "MALFORMED") {
			styledLog = blockedStyle.Render(l)
		} else if strings.Contains(l, "ALLOWED") {
			styledLog = allowedStyle.Render(l)
		} else {
			styledLog = l
		}

		// Core Fix: Force Lip Gloss to wrap the text so the terminal doesn't tear
		b.WriteString(logWrapper.Render(styledLog))
		b.WriteString("\n")
	}

	if m.isIntercepting {
		b.WriteString("\n")
		b.WriteString(m.renderPromptBox())
	} else {
		b.WriteString(statusStyle.Render("\nStatus: Active | [q] Quit\n"))
	}

	return appContainer.Render(b.String())
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
