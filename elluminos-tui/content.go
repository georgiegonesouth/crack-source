package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"pt-tui/internal/parser"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type fileLoadedMsg struct {
	raw          string
	doc          parser.ParsedDoc
	tipsDoc      parser.ParsedDoc
	collapsed    map[string]bool
	rendered     string
	h2Lines      []int
	visibleCodes []parser.Block
	tipsRendered string
}

type renderedMsg struct {
	content      string
	tips         string
	h2Lines      []int
	visibleCodes []parser.Block
}

type cmdOutputMsg string

type interactiveExecDoneMsg struct{ err error }

func (m Model) openFile(path string) tea.Cmd {
	repoRoot := m.cfg.RepoRoot
	tis := m.tokenInputs
	totalW := m.width
	return func() tea.Msg {
		full := filepath.Join(repoRoot, path)
		data, err := os.ReadFile(full)
		if err != nil {
			return fileLoadedMsg{rendered: "error reading file: " + err.Error()}
		}
		raw := string(data)
		raw = parser.StripFrontMatter(raw)
		raw = parser.StripDirectives(raw)
		text := substitute(raw, tis, nil, "")
		doc := parser.ParseDoc(text)
		mainDoc, tipsDoc := splitTips(doc)

		hasTips := len(tipsDoc.Blocks) > 0
		cw := totalW - sidebarWidth - 2
		if hasTips {
			cw -= tipsWidth + 1
		}
		glWidth := cw - 2

		collapsed := map[string]bool{}
		activeHeading := ""
		for i, h := range headings(mainDoc) {
			collapsed[h.Text] = true
			if i == 0 {
				activeHeading = h.Text
			}
		}

		rendered, h2Lines, visibleCodes := renderBlocks(mainDoc, glWidth, collapsed, activeHeading, -1, "", true)
		tipsRendered, _, _ := renderBlocks(tipsDoc, tipsWidth-2, nil, "", -1, "", false)
		return fileLoadedMsg{
			raw:          raw,
			doc:          mainDoc,
			tipsDoc:      tipsDoc,
			collapsed:    collapsed,
			rendered:     rendered,
			h2Lines:      h2Lines,
			visibleCodes: visibleCodes,
			tipsRendered: tipsRendered,
		}
	}
}

func (m Model) rerenderFile() tea.Cmd {
	if m.rawFile == "" {
		return nil
	}
	raw := m.rawFile
	collapsed := m.collapsed
	tis := m.tokenInputs
	localVals := m.localTokenValues
	highlightToken := ""
	if m.localTokenEditActive && m.localTokenFocus < len(m.localTokensInBlock) {
		highlightToken = m.localTokensInBlock[m.localTokenFocus]
	}
	glWidth := m.contentW() - 2
	tipsDoc := m.tipsDoc
	hs := headings(m.doc)
	activeHeading := ""
	if m.headingIdx < len(hs) {
		activeHeading = hs[m.headingIdx].Text
	}
	selectedCode := -1
	if m.codeNavActive {
		selectedCode = m.codeIdx
	}
	editOverride := ""
	if m.codeNavActive && m.editText != "" && m.codeIdx < len(m.visibleCodes) {
		if m.codeEditMode {
			editOverride = parser.InsertCursorMark(m.editText, m.editCursor)
		} else {
			editOverride = m.editText
		}
	}
	return func() tea.Msg {
		text := substitute(raw, tis, localVals, highlightToken)
		doc := parser.ParseDoc(text)
		mainDoc, _ := splitTips(doc)
		content, h2Lines, visibleCodes := renderBlocks(mainDoc, glWidth, collapsed, activeHeading, selectedCode, editOverride, true)
		tips, _, _ := renderBlocks(tipsDoc, tipsWidth-2, nil, "", -1, "", false)
		return renderedMsg{content: content, tips: tips, h2Lines: h2Lines, visibleCodes: visibleCodes}
	}
}

func runCmd(c, dir string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("sh", "-c", c)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			return cmdOutputMsg(string(out) + "\n[exit: " + err.Error() + "]")
		}
		return cmdOutputMsg(string(out))
	}
}

func isCdCmd(c string) bool {
	return c == "cd" || strings.HasPrefix(c, "cd ") || strings.HasPrefix(c, "cd\t")
}

func resolveCd(arg, workDir string) (string, bool) {
	switch arg {
	case "", "~":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", false
		}
		return home, true
	default:
		var target string
		if filepath.IsAbs(arg) {
			target = filepath.Clean(arg)
		} else {
			target = filepath.Clean(filepath.Join(workDir, arg))
		}
		info, err := os.Stat(target)
		if err != nil || !info.IsDir() {
			return "", false
		}
		return target, true
	}
}

func substitute(text string, tis [7]textinput.Model, localVals map[string]string, highlightToken string) string {
	for i, k := range tokenKeys {
		if v := tis[i].Value(); v != "" {
			text = strings.ReplaceAll(text, "<"+k+">", styleToken.Render(v))
		}
	}
	for tok, v := range localVals {
		if v != "" {
			text = strings.ReplaceAll(text, tok, styleToken.Render(v))
		}
	}
	text = parser.ReToken.ReplaceAllStringFunc(text, func(m string) string {
		if highlightToken != "" && m == highlightToken {
			return styleTokenActive.Render(m)
		}
		return styleToken.Render(m)
	})
	return text
}
