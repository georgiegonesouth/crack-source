package main

import (
	"strings"

	"crack-source/internal/parser"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	styleH1 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).
		BorderBottom(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("12"))
	styleH2        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	styleH2Active  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("12"))
	styleH3        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	styleH4        = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	styleLang             = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleCode             = lipgloss.NewStyle().Background(lipgloss.Color("235")).Padding(0, 1)
	styleCodeActive       = lipgloss.NewStyle().Background(lipgloss.Color("235")).Padding(0, 1).
					Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("205"))
	styleCodeOption       = lipgloss.NewStyle().Background(lipgloss.Color("235")).Padding(0, 1).
					Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("3"))
	styleCodeOptionActive = lipgloss.NewStyle().Background(lipgloss.Color("235")).Padding(0, 1).
					Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("205"))
	styleOptionLabel      = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	styleToken            = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	styleTokenActive      = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("205")).Bold(true)
)

func splitTips(doc parser.ParsedDoc) (main parser.ParsedDoc, tips parser.ParsedDoc) {
	inTips := false
	for _, b := range doc.Blocks {
		if b.Kind == parser.BlockHeading && b.Level == 2 {
			if strings.EqualFold(b.Text, "tips") {
				inTips = true
				continue
			}
			inTips = false
		}
		if inTips {
			tips.Blocks = append(tips.Blocks, b)
		} else {
			main.Blocks = append(main.Blocks, b)
		}
	}
	return
}

func headings(doc parser.ParsedDoc) []parser.Block {
	var out []parser.Block
	for _, b := range doc.Blocks {
		if b.Kind == parser.BlockHeading && b.Level == 2 {
			out = append(out, b)
		}
	}
	return out
}

func renderBlocks(doc parser.ParsedDoc, width int, collapsed map[string]bool, activeHeading string, selectedCode int, editOverride string, skipH1 bool) (string, []int, []int, []parser.Block) {
	if width <= 4 {
		width = 80
	}
	var sb strings.Builder
	var h2Lines []int
	var codeLines []int
	var visibleCodes []parser.Block
	skipBelow := 0
	lineCount := 0
	h1Skipped := false

	for _, b := range doc.Blocks {
		if skipH1 && !h1Skipped && b.Kind == parser.BlockHeading && b.Level == 1 {
			h1Skipped = true
			continue
		}
		if b.Kind == parser.BlockHeading {
			if skipBelow > 0 && b.Level <= skipBelow {
				skipBelow = 0
			}
			if skipBelow > 0 {
				continue
			}
			active := b.Level == 2 && b.Text == activeHeading
			isCollapsed := b.Level == 2 && collapsed[b.Text]
			indicator := ""
			if b.Level == 2 {
				h2Lines = append(h2Lines, lineCount)
				if isCollapsed {
					indicator = " [+]"
					skipBelow = 2
				} else {
					indicator = " [-]"
				}
			}
			chunk := renderHeading(b, width, indicator, active) + "\n"
			lineCount += strings.Count(chunk, "\n")
			sb.WriteString(chunk)
			continue
		}
		if skipBelow > 0 {
			continue
		}
		if b.Kind == parser.BlockCode {
			idx := len(visibleCodes)
			visibleCodes = append(visibleCodes, b)
			codeLines = append(codeLines, lineCount)
			display := b
			if idx == selectedCode && editOverride != "" {
				display.Text = editOverride
			}
			chunk := renderCodeBlock(display, width, idx == selectedCode) + "\n"
			lineCount += strings.Count(chunk, "\n")
			sb.WriteString(chunk)
			continue
		}
		chunk := renderBlock(b, width) + "\n"
		lineCount += strings.Count(chunk, "\n")
		sb.WriteString(chunk)
	}
	return sb.String(), h2Lines, codeLines, visibleCodes 
}

func renderHeading(b parser.Block, width int, indicator string, active bool) string {
	text := b.Text + indicator
	switch b.Level {
	case 1:
		return styleH1.Width(width - 2).Render(text) + "\n"
	case 2:
		if active {
			return "\n" + styleH2Active.Render("  "+text+"  ") + "\n"
		}
		return "\n" + styleH2.Render("## "+text) + "\n"
	case 3:
		return "\n" + styleH3.Render("### "+text) + "\n"
	default:
		return "\n" + styleH4.Render(strings.Repeat("#", b.Level)+" "+text) + "\n"
	}
}

func renderCodeBlock(b parser.Block, width int, active bool) string {
	if b.TokenOption != "" {
		header := styleOptionLabel.Render("→ "+b.TokenOption+"  [u]") + "\n"
		var code string
		if active {
			code = styleCodeOptionActive.Width(width - 6).Render(b.Text)
		} else {
			code = styleCodeOption.Width(width - 6).Render(b.Text)
		}
		return header + code
	}
	header := ""
	if b.Lang != "" {
		header = styleLang.Render(b.Lang) + "\n"
	}
	var code string
	if active {
		code = styleCodeActive.Width(width - 6).Render(b.Text)
	} else {
		code = styleCode.Width(width - 4).Render(b.Text)
	}
	return header + code
}

func renderTable(rows [][]string, width int) string {
	if len(rows) == 0 {
		return ""
	}
	cols := len(rows[0])
	colW := make([]int, cols)
	for _, row := range rows {
		for i, cell := range row {
			if i < cols && len(cell) > colW[i] {
				colW[i] = len(cell)
			}
		}
	}
	// truncate columns if total width would exceed pane width
	// total = sum(colW) + (cols-1)*3 (for " │ ") + 2 padding
	for {
		total := 2
		for _, w := range colW {
			total += w + 3
		}
		if total <= width || cols == 0 {
			break
		}
		// trim widest column by 1
		max, idx := 0, 0
		for i, w := range colW {
			if w > max {
				max, idx = w, i
			}
		}
		if max == 0 {
			break
		}
		colW[idx]--
	}
	sep := styleDim.Render("│")
	renderRow := func(row []string) string {
		var sb strings.Builder
		for i := 0; i < cols; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			if len(cell) > colW[i] {
				cell = cell[:colW[i]]
			}
			sb.WriteString(" ")
			sb.WriteString(cell)
			sb.WriteString(strings.Repeat(" ", colW[i]-len(cell)))
			sb.WriteString(" ")
			if i < cols-1 {
				sb.WriteString(sep)
			}
		}
		return sb.String()
	}
	var sb strings.Builder
	sb.WriteString(renderRow(rows[0]) + "\n")
	// separator
	parts := make([]string, cols)
	for i, w := range colW {
		parts[i] = strings.Repeat("─", w+2)
	}
	sb.WriteString(styleDim.Render(strings.Join(parts, "┼")) + "\n")
	for _, row := range rows[1:] {
		sb.WriteString(renderRow(row) + "\n")
	}
	return sb.String()
}

func renderBlock(b parser.Block, width int) string {
	switch b.Kind {
	case parser.BlockHeading:
		return renderHeading(b, width, "", false)
	case parser.BlockCode:
		return renderCodeBlock(b, width, false)
	case parser.BlockRule:
		return styleDim.Render(strings.Repeat("─", width-2)) + "\n"
	case parser.BlockText:
		wrapped := wordwrap.String(b.Text, width-2)
		return wrapped + "\n"
	case parser.BlockTable:
		return renderTable(b.Rows, width)
	}
	return ""
}
