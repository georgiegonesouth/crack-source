package main

import (
	"strings"

	"crack-source/internal/parser"
)

func (m Model) bodyH() int {
	if m.mode == modeEditor {
		return m.height - 2 - m.cmdH()
	}
	return m.height - 2 - m.tokenBarRows() - m.cmdH()
}
func (m Model) cmdH() int     { return m.cmdHeight }
func (m Model) hasTips() bool { return len(m.tipsDoc.Blocks) > 0 }
func (m Model) contentW() int {
	w := m.width - sidebarWidth - 2
	if m.hasTips() {
		w -= tipsWidth + 1
	}
	return w
}

// focusBounds returns the inclusive inner screen rect (x0,y0,x1,y1) of the
// focused pane, excluding borders and padding so selection never captures │/─.
func (m Model) focusBounds() (x0, y0, x1, y1 int) {
	bodyTop := 1
	if m.mode == modeNotebook {
		bodyTop += m.tokenBarRows()
	}
	bodyBot := bodyTop + m.bodyH() + 1
	cmdTop := m.height - m.cmdH()
	switch m.focus {
	case paneTokens:
		return 0, 1, m.width - 1, bodyTop - 1
	case paneSidebar:
		// border on all sides; inner starts 1 inside each border
		return 1, bodyTop + 1, sidebarWidth, bodyBot - 1
	case paneContent:
		// pane starts at x=sidebarWidth+2; inner is 1 inside left/right borders
		return sidebarWidth + 3, bodyTop + 1, sidebarWidth + 2 + m.contentW(), bodyBot - 1
	case paneCmd:
		// border + Padding(0,1) on left/right; inner starts at x=2
		return 2, cmdTop + 1, m.width - 3, m.height - 2
	}
	return 0, 0, m.width - 1, m.height - 1
}

// syncContentVpHeight keeps m.vp.Height in sync with the actual lines available
// for scrolling content (bodyH minus border, pinned H1 title, and status line).
// renderContent is a value receiver so its m.vp.Height assignment is discarded;
// call this (pointer receiver) before any vp.LineUp/LineDown/SetContent.
func (m *Model) syncContentVpHeight() {
	if m.mode != modeNotebook {
		return
	}
	h := m.bodyH()
	vpH := h - 2
	w := m.contentW()
	for _, b := range m.doc.Blocks {
		if b.Kind == parser.BlockHeading && b.Level == 1 {
			vpH -= strings.Count(renderHeading(b, w-2, "", false), "\n")
			break
		}
	}
	if m.focus == paneContent {
		hs := headings(m.doc)
		if len(hs) > 0 || m.codeNavActive || m.codeEditMode || m.localTokenEditActive {
			vpH--
		}
	}
	m.vp.Height = max(1, vpH)
}

func (m *Model) scrollToCode() {
	if m.codeIdx < len(m.codeLines) {
		m.vp.SetYOffset(m.codeLines[m.codeIdx] -2)
	}
}

func (m *Model) scrollToHeading() {
	if m.headingIdx < len(m.h2Lines) {
		m.vp.SetYOffset(m.h2Lines[m.headingIdx])
	}
}

