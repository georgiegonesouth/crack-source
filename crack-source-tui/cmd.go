package main

import (
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"crack-source/internal/parser"
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

var interactiveCmds = map[string]bool{
	"vim": true, "vi": true, "nvim": true,
	"nano": true, "emacs": true, "pico": true,
	"less": true, "more": true, "man": true,
	"top": true, "htop": true, "btop": true,
	"ssh": true, "ftp": true, "telnet": true,
	// sudo and su open /dev/tty directly for password prompts; they need full
	// terminal control via tea.ExecProcess, not the captured-output runCmd path.
	"sudo": true, "su": true,
}

func isInteractiveCmd(line string) bool {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false
	}
	name := fields[0]
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	return interactiveCmds[name]
}

func handleCmdKey(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := m.keys
	switch {
	case key.Matches(msg, k.Escape):
		m.focus = paneContent
		return m, nil

	case key.Matches(msg, k.CmdClear):
		m.cmdClear()
		return m, nil

	case key.Matches(msg, k.CmdRun):
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
		if isInteractiveCmd(trimmed) {
			c := exec.Command("sh", "-c", trimmed)
			c.Dir = m.workDir
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return interactiveExecDoneMsg{err: err}
			})
		}
		return m, startStreamCmd(trimmed, m.workDir)

	case key.Matches(msg, k.Backspace):
		m.cmdDeleteBackward()
		m.cmdUpdateViewport()

	case key.Matches(msg, k.Delete):
		m.cmdDeleteForward()
		m.cmdUpdateViewport()

	case key.Matches(msg, k.Left):
		if m.cmdCursorPos > 0 {
			m.cmdCursorPos--
		}
		m.cmdUpdateViewport()

	case key.Matches(msg, k.Right):
		if m.cmdCursorPos < len([]rune(m.cmdCurrentLine)) {
			m.cmdCursorPos++
		}
		m.cmdUpdateViewport()

	case key.Matches(msg, k.Home):
		m.cmdCursorPos = 0
		m.cmdUpdateViewport()

	case key.Matches(msg, k.End):
		m.cmdCursorPos = len([]rune(m.cmdCurrentLine))
		m.cmdUpdateViewport()

	case key.Matches(msg, k.WordLeft):
		m.cmdCursorPos = parser.PrevWord(m.cmdCurrentLine, m.cmdCursorPos)
		m.cmdUpdateViewport()

	case key.Matches(msg, k.WordRight):
		m.cmdCursorPos = parser.NextWord(m.cmdCurrentLine, m.cmdCursorPos)
		m.cmdUpdateViewport()

	case key.Matches(msg, k.WordDelete):
		newPos := parser.PrevWord(m.cmdCurrentLine, m.cmdCursorPos)
		runes := []rune(m.cmdCurrentLine)
		m.cmdCurrentLine = string(append(runes[:newPos], runes[m.cmdCursorPos:]...))
		m.cmdCursorPos = newPos
		m.cmdUpdateViewport()

	case key.Matches(msg, k.CmdHistoryUp):
		if m.cmdHistIdx > 0 {
			m.cmdHistIdx--
			m.cmdCurrentLine = m.cmdHistory[m.cmdHistIdx]
			m.cmdCursorPos = len([]rune(m.cmdCurrentLine))
		}
		m.cmdUpdateViewport()

	case key.Matches(msg, k.CmdHistoryDown):
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

	case key.Matches(msg, k.ResizeCmdSmaller):
		m.cmdHeight = max(3, m.cmdHeight-1)
		m.vp.Height = m.bodyH() - 2
		m.cmdVp.Height = max(1, m.cmdH()-2)
		m.cmdUpdateViewport()

	case key.Matches(msg, k.ResizeCmdLarger):
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
