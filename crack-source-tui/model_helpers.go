package main

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

func (m *Model) scrollToHeading() {
	if m.headingIdx < len(m.h2Lines) {
		m.vp.SetYOffset(m.h2Lines[m.headingIdx])
	}
}

