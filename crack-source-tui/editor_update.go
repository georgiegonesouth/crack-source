package main

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// editorScrollUpdate recomputes scrollTop/hScroll after cursor moves.
// Called from Update (pointer receiver so mutations persist on return).
func (m *Model) editorScrollUpdate() {
	if len(m.editor.lines) == 0 {
		return
	}
	h := m.bodyH()
	w := m.contentW()
	vpHeight := h - 2 - 1
	if vpHeight < 1 {
		vpHeight = 1
	}
	lineNumWidth := len(fmt.Sprintf("%d", len(m.editor.lines)))
	gutter := lineNumWidth + 3
	innerW := w - 2 - gutter
	if innerW < 4 {
		innerW = 4
	}
	m.editor.updateScroll(vpHeight, innerW)
}

func handleModeToggle(m Model) (tea.Model, tea.Cmd) {
	if m.mode == modeNotebook {
		m.mode = modeEditor
		if m.editor.dir == "" {
			m.editor = newEditorState(m.workDir)
		}
		m.focus = paneSidebar
	} else {
		m.mode = modeNotebook
		m.focus = paneSidebar
	}
	return m, nil
}

func handleEditorKey(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.focus {
	case paneSidebar:
		return handleEditorSidebarKey(m, msg)
	case paneContent:
		return handleEditorContentKey(m, msg)
	}
	return m, nil
}

func handleEditorSidebarKey(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := m.keys
	e := &m.editor

	switch {
	case key.Matches(msg, k.ResizeCmdSnapMin):
		m.cmdHeight = 10
		m.vp.Height = m.bodyH() - 2
	case key.Matches(msg, k.ResizeCmdSnapMax):
		m.cmdHeight = m.height * 2 / 3
		m.vp.Height = m.bodyH() - 2
	case key.Matches(msg, k.Quit):
		return m, tea.Quit
	case key.Matches(msg, k.ModeToggle):
		return handleModeToggle(m)
	case key.Matches(msg, k.EditorSidebarUp):
		if e.dirCursor > 0 {
			e.dirCursor--
			e.dirOffset = clampOffset(e.dirCursor, e.dirOffset, m.bodyH()-3)
		}
	case key.Matches(msg, k.EditorSidebarDown):
		if e.dirCursor < len(e.dirEntries)-1 {
			e.dirCursor++
			e.dirOffset = clampOffset(e.dirCursor, e.dirOffset, m.bodyH()-3)
		}
	case key.Matches(msg, k.EditorSidebarOpen):
		if e.dirCursor < len(e.dirEntries) {
			de := e.dirEntries[e.dirCursor]
			if de.isDir {
				e.loadDir(filepath.Join(e.dir, de.name))
				m.workDir = m.editor.dir
			} else {
				e.openFile(filepath.Join(e.dir, de.name))
				m.focus = paneContent
				m.editorScrollUpdate()
			}
		}
	case key.Matches(msg, k.FocusContent):
		if e.filePath != "" {
			m.focus = paneContent
		}
	case key.Matches(msg, k.FocusCmd):
		m.focus = paneCmd
	}
	return m, nil
}

func handleEditorContentKey(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	k := m.keys
	e := &m.editor

	switch {
	case key.Matches(msg, k.ResizeCmdSnapMin):
		m.cmdHeight = 10
		m.vp.Height = m.bodyH() - 2
	case key.Matches(msg, k.ResizeCmdSnapMax):
		m.cmdHeight = m.height * 2 / 3
		m.vp.Height = m.bodyH() - 2
	case key.Matches(msg, k.Quit):
		return m, tea.Quit
	case key.Matches(msg, k.EditorSave):
		e.saveFile()
	case key.Matches(msg, k.ModeToggle):
		return handleModeToggle(m)
	case key.Matches(msg, k.Escape):
		m.focus = paneSidebar
	case key.Matches(msg, k.Up):
		e.moveUp()
		m.editorScrollUpdate()
	case key.Matches(msg, k.Down):
		e.moveDown()
		m.editorScrollUpdate()
	case key.Matches(msg, k.Left):
		e.moveLeft()
		m.editorScrollUpdate()
	case key.Matches(msg, k.Right):
		e.moveRight()
		m.editorScrollUpdate()
	case key.Matches(msg, k.Home):
		e.moveLineStart()
		m.editorScrollUpdate()
	case key.Matches(msg, k.End):
		e.moveLineEnd()
		m.editorScrollUpdate()
	case key.Matches(msg, k.Enter):
		e.insertNewline()
		m.editorScrollUpdate()
	case key.Matches(msg, k.Backspace):
		e.deleteBackward()
		m.editorScrollUpdate()
	case key.Matches(msg, k.Delete):
		e.deleteForward()
		m.editorScrollUpdate()
	case key.Matches(msg, k.Tab):
		for i := 0; i < 4; i++ {
			e.insertRune(' ')
		}
		m.editorScrollUpdate()
	case key.Matches(msg, k.FocusCmd):
		m.focus = paneCmd
	default:
		if len(msg.Runes) == 1 {
			e.insertRune(msg.Runes[0])
			m.editorScrollUpdate()
		}
	}
	return m, nil
}
