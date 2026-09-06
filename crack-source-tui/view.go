package main

import (
	"strings"

	"crack-source/internal/parser"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleActive   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	styleDim      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleSelected = lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("15"))
	styleHeader   = lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("12")).Bold(true).Padding(0, 1)
)

func (m Model) buildFrame() string {
	if m.width == 0 {
		return "loading..."
	}

	var headerText string
	if m.mode == modeEditor {
		headerText = "PT Notes  editor  [ctrl+e] notebook  [C] cmd  [ctrl+s] save  [q] quit"
	} else {
		headerText = "PT Notes  [/] search  [t] tokens  [C] cmd  [ctrl+e] editor  [q] quit"
	}
	header := styleHeader.Width(m.width).Render(headerText)
	sidebar := m.renderSidebar()
	content := m.renderContent()
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
	if m.mode == modeNotebook && m.hasTips() {
		body = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content, m.renderTips())
	}
	cmdPane := m.renderCmd()

	if m.mode == modeEditor {
		return lipgloss.JoinVertical(lipgloss.Left, header, body, cmdPane)
	}
	tokenBar := m.renderTokenBar()
	return lipgloss.JoinVertical(lipgloss.Left, header, tokenBar, body, cmdPane)
}

func (m Model) View() string {
	frame := m.buildFrame()
	if m.selActive {
		return applySelectionHighlight(frame, m.selAX, m.selAY, m.selEX, m.selEY)
	}
	return frame
}

func selNorm(ax, ay, ex, ey int) (int, int, int, int) {
	if ay > ey || (ay == ey && ax > ex) {
		return ex, ey, ax, ay
	}
	return ax, ay, ex, ey
}

func applySelectionHighlight(frame string, ax, ay, ex, ey int) string {
	ax, ay, ex, ey = selNorm(ax, ay, ex, ey)
	lines := strings.Split(frame, "\n")
	hlStyle := lipgloss.NewStyle().Reverse(true)
	for y := ay; y <= ey && y < len(lines); y++ {
		plain := stripANSIString(lines[y])
		runes := []rune(plain)
		startX, endX := 0, len(runes)
		if y == ay {
			startX = ax
		}
		if y == ey {
			endX = ex + 1
		}
		if startX < 0 {
			startX = 0
		}
		if endX > len(runes) {
			endX = len(runes)
		}
		if startX > endX {
			startX = endX
		}
		lines[y] = string(runes[:startX]) + hlStyle.Render(string(runes[startX:endX])) + string(runes[endX:])
	}
	return strings.Join(lines, "\n")
}

func extractSelection(frame string, ax, ay, ex, ey int) string {
	ax, ay, ex, ey = selNorm(ax, ay, ex, ey)
	lines := strings.Split(frame, "\n")
	var parts []string
	for y := ay; y <= ey && y < len(lines); y++ {
		plain := stripANSIString(lines[y])
		runes := []rune(plain)
		startX, endX := 0, len(runes)
		if y == ay {
			startX = ax
		}
		if y == ey {
			endX = ex + 1
		}
		if startX < 0 {
			startX = 0
		}
		if endX > len(runes) {
			endX = len(runes)
		}
		if startX > endX {
			startX = endX
		}
		parts = append(parts, string(runes[startX:endX]))
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

// stripANSIString is a thin wrapper so view.go can strip ANSI without importing parser directly.
func stripANSIString(s string) string {
	return parser.StripANSI(s)
}

func (m Model) tokenBarRows() int {
	needed := 2
	for i, k := range tokenKeys {
		needed += len(k) + 1 + 1 + m.tokenInputs[i].Width
		if i < len(tokenKeys)-1 {
			needed += 2
		}
	}
	if needed <= m.width {
		return 1
	}
	return 2
}

func (m Model) renderTokenBar() string {
	parts := make([]string, len(tokenKeys))
	for i, ti := range m.tokenInputs {
		label := styleDim.Render(tokenKeys[i] + ":")
		active := m.focus == paneTokens && m.tokenFocus == i
		inp := ti.View()
		if active {
			inp = styleActive.Render(inp)
		}
		parts[i] = lipgloss.JoinHorizontal(lipgloss.Center, label, " ", inp)
	}
	style := lipgloss.NewStyle().Background(lipgloss.Color("236")).Width(m.width).Padding(0, 1)
	if m.tokenBarRows() == 1 {
		row := lipgloss.JoinHorizontal(lipgloss.Center,
			parts[0], "  ", parts[1], "  ", parts[2], "  ", parts[3],
			"  ", parts[4], "  ", parts[5], "  ", parts[6],
		)
		return style.Render(row)
	}
	row1 := lipgloss.JoinHorizontal(lipgloss.Center,
		parts[0], "  ", parts[1], "  ", parts[2], "  ", parts[3],
	)
	row2 := lipgloss.JoinHorizontal(lipgloss.Center,
		parts[4], "  ", parts[5], "  ", parts[6],
	)
	return lipgloss.JoinVertical(lipgloss.Left, style.Render(row1), style.Render(row2))
}

func (m Model) renderSidebar() string {
	if m.mode == modeEditor {
		return m.renderEditorSidebar()
	}
	h := m.bodyH()
	w := sidebarWidth

	var body string

	if m.focus == paneSearch {
		list := m.searchList()
		maxLines := h - 3
		var sb strings.Builder
		start := m.sidebarOffset
		end := min(start+maxLines, len(list))
		for i := start; i < end; i++ {
			f := list[i]
			indent := strings.Repeat("  ", f.Depth)
			label := indent + f.Label
			if f.Path == "" {
				if m.sidebarCollapsed[f.OrigIdx] {
					label += " ▶"
				} else {
					label += " ▼"
				}
			}
			if len(label) > w-2 {
				label = label[:w-2]
			}
			if i == m.searchCursor {
				sb.WriteString(styleSelected.Width(w - 2).Render(label))
				sb.WriteByte('\n')
			} else if f.Path == "" {
				sb.WriteString(styleDim.Render(label))
				sb.WriteByte('\n')
			} else {
				sb.WriteString(label)
				sb.WriteByte('\n')
			}
		}
		if len(list) == 0 {
			sb.WriteString(styleDim.Render("no results"))
			sb.WriteByte('\n')
		}
		body = lipgloss.JoinVertical(lipgloss.Left,
			styleActive.Render("/ ")+m.searchInput.View(),
			sb.String(),
		)
	} else {
		tree := m.collapsibleTree()
		maxLines := h - 2
		var sb strings.Builder
		start := m.sidebarOffset
		end := min(start+maxLines, len(tree))
		for i := start; i < end; i++ {
			f := tree[i]
			indent := strings.Repeat("  ", f.Depth)
			label := indent + f.Label
			if f.Path == "" {
				if m.sidebarCollapsed[f.OrigIdx] {
					label += " ▶"
				} else {
					label += " ▼"
				}
			}
			if len(label) > w-2 {
				label = label[:w-2]
			}
			if i == m.cursor {
				sb.WriteString(styleSelected.Width(w - 2).Render(label))
				sb.WriteByte('\n')
			} else if f.Path == "" {
				sb.WriteString(styleDim.Render(label))
				sb.WriteByte('\n')
			} else {
				sb.WriteString(label)
				sb.WriteByte('\n')
			}
		}
		body = sb.String()
	}

	borderColor := lipgloss.Color("8")
	if m.focus == paneSidebar || m.focus == paneSearch {
		borderColor = lipgloss.Color("12")
	}
	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Render(body)
}

func (m Model) renderContent() string {
	if m.mode == modeEditor {
		return m.renderEditorContent()
	}
	h := m.bodyH()
	w := m.contentW()
	borderColor := lipgloss.Color("8")
	if m.focus == paneContent {
		borderColor = lipgloss.Color("12")
	}

	pinnedTitle := ""
	pinnedTitleH := 0
	for _, b := range m.doc.Blocks {
		if b.Kind == parser.BlockHeading && b.Level == 1 {
			pinnedTitle = renderHeading(b, w-2, "", false)
			pinnedTitleH = strings.Count(pinnedTitle, "\n")
			break
		}
	}

	statusLine := ""
	vpHeight := h - 2 - pinnedTitleH
	if m.focus == paneContent {
		hs := headings(m.doc)
		showStatus := len(hs) > 0 || m.codeNavActive || m.codeEditMode || m.localTokenEditActive
		if showStatus {
			vpHeight--
			switch {
			case m.codeEditMode:
				statusLine = styleDim.Render("  ←→: cursor  alt+←→: word  ctrl+a/e: ends  enter: done  esc: discard")
			case m.localTokenEditActive && len(m.localTokensInBlock) > 0:
				tok := m.localTokensInBlock[m.localTokenFocus]
				statusLine = styleDim.Render(tok+": ") + m.localTokenInput.View() + styleDim.Render("  tab: next  esc: done")
			case m.codeNavActive:
				statusLine = styleDim.Render("  d/f: block  x: run  u: use  enter: edit  tab: token  esc: deselect")
			default:
				statusLine = styleDim.Render("  ← →: section  enter: expand/collapse")
			}
		}
	}

	m.vp.Width = w - 2
	m.vp.Height = max(1, vpHeight)

	inner := m.vp.View()
	if pinnedTitle != "" {
		if statusLine != "" {
			inner = statusLine + "\n" + pinnedTitle + inner
		} else {
			inner = pinnedTitle + inner
		}
	} else if statusLine != "" {
		inner = statusLine + "\n" + inner
	}

	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Render(inner)
}

func (m Model) renderTips() string {
	h := m.bodyH()
	rendered, _, _ := renderBlocks(m.tipsDoc, tipsWidth-2, nil, "", -1, "", false)
	return lipgloss.NewStyle().
		Width(tipsWidth).
		Height(h).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("3")).
		Render(styleDim.Render("tips") + "\n" + rendered)
}

func (m Model) renderCmd() string {
	h := m.cmdH()
	w := m.width
	borderColor := lipgloss.Color("8")
	if m.focus == paneCmd {
		borderColor = lipgloss.Color("12")
	}
	m.cmdVp.Width = w - 4
	m.cmdVp.Height = max(1, h-2)
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(w - 2).
		Height(h - 2).
		Render(m.cmdVp.View())
}
