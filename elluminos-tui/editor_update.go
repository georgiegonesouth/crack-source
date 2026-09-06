package main

import (
	"fmt"
	"path/filepath"

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
	switch msg.String() {
	case "ctrl+shift+down":
		m.cmdHeight = 10
		m.vp.Height = m.bodyH() - 2
	case "ctrl+shift+up":
		m.cmdHeight = m.height * 2 / 3
		m.vp.Height = m.bodyH() - 2
	}
	e := &m.editor
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "ctrl+e":
		return handleModeToggle(m)
	case "up", "j":
		if e.dirCursor > 0 {
			e.dirCursor--
			e.dirOffset = clampOffset(e.dirCursor, e.dirOffset, m.bodyH()-3)
		}
	case "down", "k":
		if e.dirCursor < len(e.dirEntries)-1 {
			e.dirCursor++
			e.dirOffset = clampOffset(e.dirCursor, e.dirOffset, m.bodyH()-3)
		}
	case "enter", "l":
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
	case "tab", "right":
		if e.filePath != "" {
			m.focus = paneContent
		}
	case "C":
		m.focus = paneCmd
		return m, nil
	}
	return m, nil
}

func handleEditorContentKey(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+shift+down":
		m.cmdHeight = 10
		m.vp.Height = m.bodyH() - 2
	case "ctrl+shift+up":
		m.cmdHeight = m.height * 2 / 3
		m.vp.Height = m.bodyH() - 2
	}
	e := &m.editor
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "ctrl+s":
		e.saveFile()
	case "ctrl+e":
		return handleModeToggle(m)
	case "escape", "esc":
		m.focus = paneSidebar
	case "up":
		e.moveUp()
		m.editorScrollUpdate()
	case "down":
		e.moveDown()
		m.editorScrollUpdate()
	case "left":
		e.moveLeft()
		m.editorScrollUpdate()
	case "right":
		e.moveRight()
		m.editorScrollUpdate()
	case "home", "ctrl+a":
		e.moveLineStart()
		m.editorScrollUpdate()
	case "end":
		e.moveLineEnd()
		m.editorScrollUpdate()
	case "enter":
		e.insertNewline()
		m.editorScrollUpdate()
	case "backspace", "ctrl+h":
		e.deleteBackward()
		m.editorScrollUpdate()
	case "delete":
		e.deleteForward()
		m.editorScrollUpdate()
	case "tab":
		for i := 0; i < 4; i++ {
			e.insertRune(' ')
		}
		m.editorScrollUpdate()
	case "C":
		m.focus = paneCmd
		return m, nil
	default:
		if len(msg.Runes) == 1 {
			e.insertRune(msg.Runes[0])
			m.editorScrollUpdate()
		}
	}
	return m, nil
}
