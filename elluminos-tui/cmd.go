package main

import (
	"os"
	"os/user"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pt-tui/internal/parser"
)

var (
	stylePromptUser = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	stylePromptDir  = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	stylePromptSym  = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
)

func (m Model) cmdPromptPrefix() string {
	u, _ := user.Current()
	hostname, _ := os.Hostname()
	if i := strings.Index(hostname, "."); i > 0 {
		hostname = hostname[:i]
	}
	username := "user"
	if u != nil {
		username = u.Username
	}
	dir := m.workDir
	home, _ := os.UserHomeDir()
	switch {
	case dir == home:
		dir = "~"
	case strings.HasPrefix(dir, home+"/"):
		dir = "~" + dir[len(home):]
	}
	return stylePromptUser.Render(username+"@"+hostname) + " " +
		stylePromptDir.Render(dir) + " " +
		stylePromptSym.Render("%") + " "
}

func (m *Model) cmdUpdateViewport() {
	runes := []rune(m.cmdCurrentLine)
	col := m.cmdCursorPos
	if col > len(runes) {
		col = len(runes)
	}
	var lineDisplay string
	if col < len(runes) {
		lineDisplay = string(runes[:col]) +
			styleEditorCursor.Render(string(runes[col:col+1])) +
			string(runes[col+1:])
	} else {
		lineDisplay = string(runes) + styleEditorCursor.Render(" ")
	}
	content := m.cmdScrollback + m.cmdPromptPrefix() + lineDisplay
	if m.cmdPromptAtTop {
		h := m.cmdVp.Height
		if h < 1 {
			h = 8
		}
		content += strings.Repeat("\n", h-1)
	}
	m.cmdVp.SetContent(content)
	m.cmdVp.GotoBottom()
	// After a clear, don't let GotoBottom reveal pre-clear history.
	if m.cmdVp.YOffset < m.cmdClearLine {
		m.cmdVp.YOffset = m.cmdClearLine
	}
}

func (m *Model) cmdInsertRune(r rune) {
	runes := []rune(m.cmdCurrentLine)
	col := m.cmdCursorPos
	if col > len(runes) {
		col = len(runes)
	}
	newRunes := make([]rune, len(runes)+1)
	copy(newRunes, runes[:col])
	newRunes[col] = r
	copy(newRunes[col+1:], runes[col:])
	m.cmdCurrentLine = string(newRunes)
	m.cmdCursorPos++
}

func (m *Model) cmdDeleteBackward() {
	runes := []rune(m.cmdCurrentLine)
	if m.cmdCursorPos > 0 {
		col := m.cmdCursorPos
		m.cmdCurrentLine = string(append(runes[:col-1], runes[col:]...))
		m.cmdCursorPos--
	}
}

func (m *Model) cmdDeleteForward() {
	runes := []rune(m.cmdCurrentLine)
	if m.cmdCursorPos < len(runes) {
		col := m.cmdCursorPos
		m.cmdCurrentLine = string(append(runes[:col], runes[col+1:]...))
	}
}

func (m *Model) cmdClear() {
	m.cmdClearLine = strings.Count(m.cmdScrollback, "\n")
	m.cmdPromptAtTop = true
	m.cmdUpdateViewport()
}

func handleCmdKey(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "escape", "esc":
		m.focus = paneContent
		return m, nil

	case "ctrl+l":
		m.cmdClear()
		return m, nil

	case "ctrl+shift+c":
		_ = clipboard.WriteAll(m.cmdCurrentLine)
		return m, nil

	case "ctrl+shift+v":
		if text, err := clipboard.ReadAll(); err == nil {
			for _, r := range []rune(text) {
				if r != '\n' && r != '\r' {
					m.cmdInsertRune(r)
				}
			}
			m.cmdUpdateViewport()
		}
		return m, nil

	case "enter":
		line := m.cmdCurrentLine
		trimmed := strings.TrimSpace(line)
		m.cmdPromptAtTop = false
		m.cmdScrollback += m.cmdPromptPrefix() + line + "\n"
		if trimmed != "" {
			if len(m.cmdHistory) == 0 || m.cmdHistory[len(m.cmdHistory)-1] != trimmed {
				m.cmdHistory = append(m.cmdHistory, trimmed)
			}
			m.cmdHistIdx = len(m.cmdHistory)
		}
		m.cmdCurrentLine = ""
		m.cmdCursorPos = 0

		if trimmed == "" {
			m.cmdUpdateViewport()
			return m, nil
		}
		if trimmed == "clear" || trimmed == "reset" {
			m.cmdClear()
			return m, nil
		}
		if isCdCmd(trimmed) {
			arg := ""
			if len(trimmed) > 2 {
				arg = strings.TrimSpace(trimmed[2:])
			}
			if newDir, ok := resolveCd(arg, m.workDir); ok {
				m.workDir = newDir
				if m.mode == modeEditor {
					m.editor.loadDir(m.workDir)
				}
			} else {
				m.cmdScrollback += "cd: " + arg + ": No such file or directory\n"
			}
			m.cmdUpdateViewport()
			return m, nil
		}
		m.cmdUpdateViewport()
		return m, runCmd(trimmed, m.workDir)

	case "backspace":
		m.cmdDeleteBackward()
		m.cmdUpdateViewport()

	case "delete":
		m.cmdDeleteForward()
		m.cmdUpdateViewport()

	case "left":
		if m.cmdCursorPos > 0 {
			m.cmdCursorPos--
		}
		m.cmdUpdateViewport()

	case "right":
		if m.cmdCursorPos < len([]rune(m.cmdCurrentLine)) {
			m.cmdCursorPos++
		}
		m.cmdUpdateViewport()

	case "home", "ctrl+a":
		m.cmdCursorPos = 0
		m.cmdUpdateViewport()

	case "end":
		m.cmdCursorPos = len([]rune(m.cmdCurrentLine))
		m.cmdUpdateViewport()

	case "alt+left":
		m.cmdCursorPos = parser.PrevWord(m.cmdCurrentLine, m.cmdCursorPos)
		m.cmdUpdateViewport()

	case "alt+right":
		m.cmdCursorPos = parser.NextWord(m.cmdCurrentLine, m.cmdCursorPos)
		m.cmdUpdateViewport()

	case "ctrl+w":
		newPos := parser.PrevWord(m.cmdCurrentLine, m.cmdCursorPos)
		runes := []rune(m.cmdCurrentLine)
		m.cmdCurrentLine = string(append(runes[:newPos], runes[m.cmdCursorPos:]...))
		m.cmdCursorPos = newPos
		m.cmdUpdateViewport()

	case "up":
		if m.cmdHistIdx > 0 {
			m.cmdHistIdx--
			m.cmdCurrentLine = m.cmdHistory[m.cmdHistIdx]
			m.cmdCursorPos = len([]rune(m.cmdCurrentLine))
		}
		m.cmdUpdateViewport()

	case "down":
		if m.cmdHistIdx < len(m.cmdHistory) {
			m.cmdHistIdx++
			if m.cmdHistIdx == len(m.cmdHistory) {
				m.cmdCurrentLine = ""
			} else {
				m.cmdCurrentLine = m.cmdHistory[m.cmdHistIdx]
			}
			m.cmdCursorPos = len([]rune(m.cmdCurrentLine))
		}
		m.cmdUpdateViewport()

	case "shift+down":
		m.cmdHeight = max(3, m.cmdHeight-1)
		m.vp.Height = m.bodyH() - 2
		m.cmdVp.Height = max(1, m.cmdH()-2)
		m.cmdUpdateViewport()

	case "shift+up":
		m.cmdHeight = min(m.height*2/3, m.cmdHeight+1)
		m.vp.Height = m.bodyH() - 2
		m.cmdVp.Height = max(1, m.cmdH()-2)
		m.cmdUpdateViewport()

	default:
		if len(msg.Runes) == 1 {
			m.cmdInsertRune(msg.Runes[0])
			m.cmdUpdateViewport()
		}
	}
	return m, nil
}
