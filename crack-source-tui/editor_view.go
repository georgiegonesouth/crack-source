package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleEditorCursor  = lipgloss.NewStyle().Reverse(true)
	styleEditorLineNum = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleEditorDir     = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	styleEditorDirty   = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
)

func (m Model) renderEditorSidebar() string {
	e := m.editor
	h := m.bodyH()
	w := sidebarWidth

	var sb strings.Builder

	dirLabel := e.dir
	if len([]rune(dirLabel)) > w-2 {
		runes := []rune(dirLabel)
		dirLabel = "…" + string(runes[len(runes)-(w-3):])
	}
	sb.WriteString(styleEditorDir.Render(dirLabel))
	sb.WriteByte('\n')

	maxLines := h - 3
	start := e.dirOffset
	end := min(start+maxLines, len(e.dirEntries))
	for i := start; i < end; i++ {
		de := e.dirEntries[i]
		label := de.name
		if de.isDir {
			label += "/"
		}
		labelRunes := []rune(label)
		if len(labelRunes) > w-2 {
			label = string(labelRunes[:w-2])
		}
		if i == e.dirCursor {
			sb.WriteString(styleSelected.Width(w - 2).Render(label))
		} else if de.isDir {
			sb.WriteString(styleEditorDir.Render(label))
		} else {
			sb.WriteString(label)
		}
		sb.WriteByte('\n')
	}

	borderColor := lipgloss.Color("8")
	if m.focus == paneSidebar {
		borderColor = lipgloss.Color("12")
	}
	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Render(sb.String())
}

func (m Model) renderEditorContent() string {
	e := m.editor
	h := m.bodyH()
	w := m.contentW()

	borderColor := lipgloss.Color("8")
	if m.focus == paneContent {
		borderColor = lipgloss.Color("12")
	}
	border := lipgloss.NewStyle().
		Width(w).
		Height(h).
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor)

	if e.filePath == "" {
		return border.Render(styleDim.Render("open a file from the sidebar"))
	}

	name := filepath.Base(e.filePath)
	dirtyMark := ""
	if e.dirty {
		dirtyMark = styleEditorDirty.Render(" ●")
	}
	status := styleH2.Render(name) + dirtyMark +
		styleDim.Render(fmt.Sprintf("  %d:%d  [ctrl+s] save  [esc] sidebar", e.line+1, e.col+1))

	vpHeight := h - 2 - 1
	if vpHeight < 1 {
		vpHeight = 1
	}

	lineNumWidth := len(fmt.Sprintf("%d", len(e.lines)))
	gutter := lineNumWidth + 3
	innerW := w - 2 - gutter
	if innerW < 4 {
		innerW = 4
	}

	var sb strings.Builder
	for i := e.scrollTop; i < e.scrollTop+vpHeight && i < len(e.lines); i++ {
		lineNum := styleEditorLineNum.Render(fmt.Sprintf("%*d │ ", lineNumWidth, i+1))
		runes := []rune(e.lines[i])

		if i == e.line {
			col := e.col
			if col > len(runes) {
				col = len(runes)
			}
			hStart := e.hScroll
			var visible []rune
			if hStart < len(runes) {
				visible = runes[hStart:]
			}
			displayCol := col - hStart
			if displayCol < 0 {
				displayCol = 0
			}

			var rendered string
			if displayCol < len(visible) {
				before := string(visible[:displayCol])
				cur := styleEditorCursor.Render(string(visible[displayCol : displayCol+1]))
				end := len(visible)
				if end > innerW {
					end = innerW
				}
				rest := ""
				if displayCol+1 < end {
					rest = string(visible[displayCol+1 : end])
				}
				rendered = before + cur + rest
			} else {
				before := string(visible)
				beforeRunes := []rune(before)
				if len(beforeRunes) > innerW-1 {
					before = string(beforeRunes[:innerW-1])
				}
				rendered = before + styleEditorCursor.Render(" ")
			}
			sb.WriteString(lineNum)
			sb.WriteString(rendered)
		} else {
			hStart := e.hScroll
			var visible []rune
			if hStart < len(runes) {
				visible = runes[hStart:]
			}
			if len(visible) > innerW {
				visible = visible[:innerW]
			}
			sb.WriteString(lineNum)
			sb.WriteString(string(visible))
		}
		sb.WriteByte('\n')
	}

	return border.Render(status + "\n" + sb.String())
}
