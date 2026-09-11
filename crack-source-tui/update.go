package main

import (
	"strings"
	"syscall"

	"crack-source/internal/parser"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.vp.Width = m.contentW() - 2
		m.vp.Height = m.bodyH() - 2
		m.cmdVp.Width = m.width - 4
		m.cmdVp.Height = max(1, m.cmdH()-2)
		m.cmdUpdateViewport()
		if m.rawFile != "" {
			cmds = append(cmds, m.rerenderFile())
		}
		if m.mode == modeEditor {
			m.editorScrollUpdate()
		}

	case fileLoadedMsg:
		m.rawFile = msg.raw
		m.doc = msg.doc
		m.tipsDoc = msg.tipsDoc
		m.collapsed = msg.collapsed
		m.h2Lines = msg.h2Lines
		m.visibleCodes = msg.visibleCodes
		m.headingIdx = 0
		m.codeIdx = 0
		m.codeNavActive = false
		m.codeEditMode = false
		m.editText = ""
		m.editCursor = 0
		m.localTokenKeys = parser.ExtractLocalTokens(msg.raw, tokenKeys)
		m.localTokenValues = map[string]string{}
		m.localTokenFocus = 0
		m.localTokenEditActive = false
		m.localTokenInput.Blur()
		m.vp.GotoTop()
		m.vp.SetContent(msg.rendered)
		m.scrollToHeading()

	case contentIndexedMsg:
		m.contentIndex = map[string]string(msg)

	case renderedMsg:
		m.h2Lines = msg.h2Lines
		m.visibleCodes = msg.visibleCodes
		m.vp.SetContent(msg.content)
		m.scrollToHeading()

	case tea.MouseMsg:
		switch msg.Action {
		case tea.MouseActionPress:
			if msg.Button == tea.MouseButtonLeft && msg.Shift {
				x0, y0, x1, y1 := m.focusBounds()
				cx := max(x0, min(msg.X, x1))
				cy := max(y0, min(msg.Y, y1))
				m.selActive = true
				m.selAX, m.selAY = cx, cy
				m.selEX, m.selEY = cx, cy
				m.selX0, m.selX1 = x0, x1
			} else if msg.Button == tea.MouseButtonLeft {
				m.selActive = false
				bodyTop := 1
				if m.mode == modeNotebook {
					bodyTop += m.tokenBarRows()
				}
				cmdTop := m.height - m.cmdH()
				switch {
				case msg.Y >= cmdTop:
					m.focus = paneCmd
				case msg.Y < bodyTop && m.mode == modeNotebook && msg.Y >= 1:
					m.tokenInputs[m.tokenFocus].Focus()
					m.focus = paneTokens
				case msg.X < sidebarWidth+2:
					m.searchInput.Blur()
					m.searchNavMode = false
					m.focus = paneSidebar
				default:
					m.focus = paneContent
				}
			}
		case tea.MouseActionMotion:
			if m.selActive {
				x0, y0, x1, y1 := m.focusBounds()
				m.selEX = max(x0, min(msg.X, x1))
				m.selEY = max(y0, min(msg.Y, y1))
			}
		case tea.MouseActionRelease:
			if m.selActive && msg.Button == tea.MouseButtonLeft {
				text := extractSelection(m.buildFrame(), m.selAX, m.selAY, m.selEX, m.selEY, m.selX0, m.selX1)
				if text != "" {
					_ = clipboard.WriteAll(text)
				}
				m.selActive = false
			}
		}
		if m.focus == paneCmd {
			if msg.Button == tea.MouseButtonWheelUp {
				m.cmdVp.LineUp(3)
			} else if msg.Button == tea.MouseButtonWheelDown {
				m.cmdVp.LineDown(3)
			}
		}

	case cmdOutputMsg:
		m.cmdScrollback += string(msg) + "\n"
		m.cmdPromptAtTop = false
		m.cmdUpdateViewport()
		if m.mode == modeEditor {
			m.editor.reloadDir()
		}

	case cmdStartedMsg:
		m.runningProc = msg.proc
		return m, func() tea.Msg { return readChunk(msg.r, msg.proc) }

	case cmdChunkMsg:
		m.cmdScrollback += msg.text
		m.cmdPromptAtTop = false
		m.cmdUpdateViewport()
		return m, func() tea.Msg { return readChunk(msg.r, msg.proc) }

	case cmdDoneMsg:
		m.runningProc = nil
		m.procSiginted = false
		if msg.text != "" {
			m.cmdScrollback += msg.text
		}
		if msg.err != nil {
			m.cmdScrollback += "[exit: " + msg.err.Error() + "]\n"
		}
		m.cmdUpdateViewport()
		if m.mode == modeEditor {
			m.editor.reloadDir()
		}

	case ctrlCTimeoutMsg:
		m.ctrlCPending = false

	case interactiveExecDoneMsg:
		if msg.err != nil {
			m.cmdScrollback += "[exit: " + msg.err.Error() + "]\n"
		}
		m.cmdUpdateViewport()
		if m.mode == modeEditor {
			m.editor.reloadDir()
		}

	case tea.KeyMsg:
		// Bracketed paste: bubbletea v1.x delivers it as KeyMsg with Paste=true.
		// Route it through the same logic as ctrl+shift+v.
		if msg.Paste {
			text := string(msg.Runes)
			noNL := strings.NewReplacer("\r\n", "", "\r", "", "\n", "").Replace(text)
			switch m.focus {
			case paneCmd:
				for _, r := range []rune(noNL) {
					m.cmdInsertRune(r)
				}
				m.cmdUpdateViewport()
			case paneSearch:
				if !m.searchNavMode {
					m.searchInput.SetValue(m.searchInput.Value() + noNL)
					m.filtered = m.filterEntries(m.searchInput.Value())
				}
			case paneTokens:
				m.tokenInputs[m.tokenFocus].SetValue(m.tokenInputs[m.tokenFocus].Value() + noNL)
				m.commitTokens()
				return m, m.rerenderFile()
			case paneContent:
				switch {
				case m.localTokenEditActive:
					m.localTokenInput.SetValue(m.localTokenInput.Value() + noNL)
					if len(m.localTokensInBlock) > 0 {
						m.localTokenValues[m.localTokensInBlock[m.localTokenFocus]] = m.localTokenInput.Value()
					}
					return m, m.rerenderFile()
				case m.codeEditMode:
					runes := []rune(m.editText)
					ins := []rune(noNL)
					m.editText = string(runes[:m.editCursor]) + string(ins) + string(runes[m.editCursor:])
					m.editCursor += len(ins)
					return m, m.rerenderFile()
				case m.mode == modeEditor:
					for _, r := range []rune(text) {
						if r == '\n' {
							m.editor.insertNewline()
						} else if r != '\r' {
							m.editor.insertRune(r)
						}
					}
					m.editorScrollUpdate()
				}
			}
			return m, nil
		}

		if msg.Type == tea.KeyCtrlC {
			if m.runningProc != nil {
				pgid := m.runningProc.Process.Pid
				if m.procSiginted {
					syscall.Kill(-pgid, syscall.SIGKILL)
				} else {
					syscall.Kill(-pgid, syscall.SIGINT)
					m.procSiginted = true
				}
				return m, nil
			}
			if m.ctrlCPending {
				return m, tea.Quit
			}
			m.ctrlCPending = true
			m.cmdScrollback += "\n[press Ctrl+C again to exit]\n"
			m.cmdUpdateViewport()
			return m, ctrlCTimeout()
		}
		m.ctrlCPending = false

		if key.Matches(msg, m.keys.CmdCopyLine) {
			var text string
			switch m.focus {
			case paneCmd:
				text = m.cmdCurrentLine
			case paneContent:
				if m.mode == modeEditor && m.editor.line < len(m.editor.lines) {
					text = m.editor.lines[m.editor.line]
				} else if m.codeNavActive && m.codeIdx < len(m.visibleCodes) {
					text = m.editText
					if text == "" {
						text = parser.StripANSI(m.visibleCodes[m.codeIdx].Text)
					}
				}
			}
			if text != "" {
				_ = clipboard.WriteAll(text)
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.CmdPaste) {
			text, err := clipboard.ReadAll()
			if err != nil || text == "" {
				return m, nil
			}
			noNL := strings.NewReplacer("\r\n", "", "\r", "", "\n", "").Replace(text)
			switch m.focus {
			case paneCmd:
				for _, r := range []rune(noNL) {
					m.cmdInsertRune(r)
				}
				m.cmdUpdateViewport()
			case paneSearch:
				if !m.searchNavMode {
					m.searchInput.SetValue(m.searchInput.Value() + noNL)
					m.filtered = m.filterEntries(m.searchInput.Value())
				}
			case paneTokens:
				m.tokenInputs[m.tokenFocus].SetValue(m.tokenInputs[m.tokenFocus].Value() + noNL)
				m.commitTokens()
				return m, m.rerenderFile()
			case paneContent:
				switch {
				case m.localTokenEditActive:
					m.localTokenInput.SetValue(m.localTokenInput.Value() + noNL)
					if len(m.localTokensInBlock) > 0 {
						m.localTokenValues[m.localTokensInBlock[m.localTokenFocus]] = m.localTokenInput.Value()
					}
					return m, m.rerenderFile()
				case m.codeEditMode:
					runes := []rune(m.editText)
					ins := []rune(noNL)
					m.editText = string(runes[:m.editCursor]) + string(ins) + string(runes[m.editCursor:])
					m.editCursor += len(ins)
					return m, m.rerenderFile()
				case m.mode == modeEditor:
					for _, r := range []rune(text) {
						if r == '\n' {
							m.editor.insertNewline()
						} else if r != '\r' {
							m.editor.insertRune(r)
						}
					}
					m.editorScrollUpdate()
				}
			}
			return m, nil
		}

		// ModeToggle fires globally before any pane routing so it always works,
		// including when paneContent is in codeEditMode.
		if key.Matches(msg, m.keys.ModeToggle) {
			return handleModeToggle(m)
		}
		if m.mode == modeEditor && (m.focus == paneSidebar || m.focus == paneContent) {
			return handleEditorKey(m, msg)
		}

		switch {
		case key.Matches(msg, m.keys.ResizeCmdSnapMin):
			m.cmdHeight = 10
			m.vp.Height = m.bodyH() - 2
			m.cmdVp.Height = max(1, m.cmdH()-2)
		case key.Matches(msg, m.keys.ResizeCmdSnapMax):
			m.cmdHeight = m.height * 2 / 3
			m.vp.Height = m.bodyH() - 2
			m.cmdVp.Height = max(1, m.cmdH()-2)
		}

		switch m.focus {
		case paneSearch:
			switch {
			case key.Matches(msg, m.keys.Escape):
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				m.filtered = nil
				m.sidebarOffset = 0
				m.searchNavMode = false
				m.focus = paneSidebar

			case key.Matches(msg, m.keys.Enter):
				if !m.searchNavMode {
					m.searchNavMode = true
					m.searchInput.Blur()
				} else {
					list := m.searchList()
					if m.searchCursor < len(list) {
						sel := list[m.searchCursor]
						if sel.Path != "" {
							cmds = append(cmds, m.openFile(sel.Path))
							m.searchNavMode = false
							m.focus = paneContent
						}
					}
				}

			default:
				if m.searchNavMode {
					// In nav mode, handle navigation keys; anything else exits nav mode.
					switch {
					case key.Matches(msg, m.keys.Quit):
						return m, tea.Quit
					case key.Matches(msg, m.keys.SidebarDown):
						list := m.searchList()
						next := m.searchCursor + 1
						for next < len(list) && list[next].Path == "" {
							next++
						}
						if next < len(list) {
							m.searchCursor = next
						}
						m.sidebarOffset = clampOffset(m.searchCursor, m.sidebarOffset, m.bodyH()-3)

					case key.Matches(msg, m.keys.SidebarUp):
						list := m.searchList()
						prev := m.searchCursor - 1
						for prev >= 0 && list[prev].Path == "" {
							prev--
						}
						if prev >= 0 {
							m.searchCursor = prev
						}
						m.sidebarOffset = clampOffset(m.searchCursor, m.sidebarOffset, m.bodyH()-3)

					case key.Matches(msg, m.keys.SidebarJumpDown):
						list := m.searchList()
						next := m.searchCursor + 1
						for next < len(list) && list[next].Path != "" {
							next++
						}
						if next < len(list) {
							m.searchCursor = next
						}
						m.sidebarOffset = clampOffset(m.searchCursor, m.sidebarOffset, m.bodyH()-3)

					case key.Matches(msg, m.keys.SidebarJumpUp):
						list := m.searchList()
						prev := m.searchCursor - 1
						for prev >= 0 && list[prev].Path != "" {
							prev--
						}
						if prev >= 0 {
							m.searchCursor = prev
						}
						m.sidebarOffset = clampOffset(m.searchCursor, m.sidebarOffset, m.bodyH()-3)

					case key.Matches(msg, m.keys.FocusSearch):
						// Toggle back to input mode.
						m.searchNavMode = false
						cmds = append(cmds, m.searchInput.Focus())

					default:
						// Any other key exits nav mode and is forwarded to the input.
						m.searchNavMode = false
						cmds = append(cmds, m.searchInput.Focus())
						var cmd tea.Cmd
						m.searchInput, cmd = m.searchInput.Update(msg)
						cmds = append(cmds, cmd)
						m.filtered = m.filterEntries(m.searchInput.Value())
						m.searchCursor = 0
						list := m.searchList()
						for m.searchCursor < len(list) && list[m.searchCursor].Path == "" {
							m.searchCursor++
						}
					}
				} else {
					// Non-nav mode: all keys feed the search input.
					var cmd tea.Cmd
					m.searchInput, cmd = m.searchInput.Update(msg)
					cmds = append(cmds, cmd)
					m.filtered = m.filterEntries(m.searchInput.Value())
					m.searchCursor = 0
					list := m.searchList()
					for m.searchCursor < len(list) && list[m.searchCursor].Path == "" {
						m.searchCursor++
					}
				}
			}

		case paneSidebar:
			tree := m.collapsibleTree()
			switch {
			case key.Matches(msg, m.keys.Enter):
				if m.cursor < len(tree) {
					sel := tree[m.cursor]
					if sel.Path == "" {
						m.sidebarCollapsed[sel.OrigIdx] = !m.sidebarCollapsed[sel.OrigIdx]
						newTree := m.collapsibleTree()
						if m.cursor >= len(newTree) {
							m.cursor = max(0, len(newTree)-1)
						}
					} else {
						cmds = append(cmds, m.openFile(sel.Path))
						m.focus = paneContent
					}
				}
			case key.Matches(msg, m.keys.FocusSearch):
				m.focus = paneSearch
				m.searchInput.SetValue("")
				m.filtered = nil
				m.searchCursor = 0
				m.searchNavMode = false
				cmds = append(cmds, m.searchInput.Focus())
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.SidebarDown):
				m.cursor = min(m.cursor+1, len(tree)-1)
				m.sidebarOffset = clampOffset(m.cursor, m.sidebarOffset, m.bodyH()-2)
			case key.Matches(msg, m.keys.SidebarUp):
				m.cursor = max(m.cursor-1, 0)
				m.sidebarOffset = clampOffset(m.cursor, m.sidebarOffset, m.bodyH()-2)
			case key.Matches(msg, m.keys.FocusContent):
				m.focus = paneContent
			case key.Matches(msg, m.keys.FocusTokens):
				m.focus = paneTokens
				m.tokenInputs[m.tokenFocus].Focus()
			case key.Matches(msg, m.keys.FocusCmd):
				m.focus = paneCmd
			case key.Matches(msg, m.keys.SidebarJumpDown):
				next := m.cursor + 1
				for next < len(tree) && tree[next].Path != "" {
					next++
				}
				if next < len(tree) {
					m.cursor = next
				}
				m.sidebarOffset = clampOffset(m.cursor, m.sidebarOffset, m.bodyH()-2)
			case key.Matches(msg, m.keys.SidebarJumpUp):
				prev := m.cursor - 1
				for prev >= 0 && tree[prev].Path != "" {
					prev--
				}
				if prev >= 0 {
					m.cursor = prev
				}
				m.sidebarOffset = clampOffset(m.cursor, m.sidebarOffset, m.bodyH()-2)
			}

		case paneContent:
			if m.codeEditMode {
				switch {
				case key.Matches(msg, m.keys.Enter):
					m.codeEditMode = false
				case key.Matches(msg, m.keys.Escape):
					m.codeEditMode = false
					m.editText = ""
					m.editCursor = 0
				case key.Matches(msg, m.keys.Left):
					if m.editCursor > 0 {
						m.editCursor--
					}
				case key.Matches(msg, m.keys.Right):
					if m.editCursor < len([]rune(m.editText)) {
						m.editCursor++
					}
				case key.Matches(msg, m.keys.WordLeft):
					m.editCursor = parser.PrevWord(m.editText, m.editCursor)
				case key.Matches(msg, m.keys.WordRight):
					m.editCursor = parser.NextWord(m.editText, m.editCursor)
				case key.Matches(msg, m.keys.Home):
					m.editCursor = 0
				case key.Matches(msg, m.keys.End):
					m.editCursor = len([]rune(m.editText))
				case key.Matches(msg, m.keys.Backspace):
					if m.editCursor > 0 {
						r := []rune(m.editText)
						m.editText = string(r[:m.editCursor-1]) + string(r[m.editCursor:])
						m.editCursor--
					}
				case key.Matches(msg, m.keys.Delete):
					r := []rune(m.editText)
					if m.editCursor < len(r) {
						m.editText = string(r[:m.editCursor]) + string(r[m.editCursor+1:])
					}
				default:
					if len(msg.Runes) > 0 {
						r := []rune(m.editText)
						ins := string(msg.Runes)
						m.editText = string(r[:m.editCursor]) + ins + string(r[m.editCursor:])
						m.editCursor += len(msg.Runes)
					}
				}
				cmds = append(cmds, m.rerenderFile())
				break
			}

			if m.localTokenEditActive {
				switch {
				case key.Matches(msg, m.keys.Escape):
					m.localTokenInput.Blur()
					m.localTokenEditActive = false
					cmds = append(cmds, m.rerenderFile())
				case key.Matches(msg, m.keys.Tab):
					m.localTokenValues[m.localTokensInBlock[m.localTokenFocus]] = m.localTokenInput.Value()
					m.localTokenFocus = (m.localTokenFocus + 1) % len(m.localTokensInBlock)
					m.localTokenInput.SetValue(m.localTokenValues[m.localTokensInBlock[m.localTokenFocus]])
					cmds = append(cmds, m.rerenderFile())
				case key.Matches(msg, m.keys.ShiftTab):
					m.localTokenValues[m.localTokensInBlock[m.localTokenFocus]] = m.localTokenInput.Value()
					m.localTokenFocus = (m.localTokenFocus + len(m.localTokensInBlock) - 1) % len(m.localTokensInBlock)
					m.localTokenInput.SetValue(m.localTokenValues[m.localTokensInBlock[m.localTokenFocus]])
					cmds = append(cmds, m.rerenderFile())
				default:
					var cmd tea.Cmd
					m.localTokenInput, cmd = m.localTokenInput.Update(msg)
					cmds = append(cmds, cmd)
					if len(m.localTokensInBlock) > 0 {
						m.localTokenValues[m.localTokensInBlock[m.localTokenFocus]] = m.localTokenInput.Value()
						cmds = append(cmds, m.rerenderFile())
					}
				}
				break
			}

			switch {
			case key.Matches(msg, m.keys.Escape):
				if m.codeNavActive {
					m.codeNavActive = false
					cmds = append(cmds, m.rerenderFile())
				} else {
					m.focus = paneSidebar
				}
			case key.Matches(msg, m.keys.FocusSidebar):
				m.focus = paneSidebar
			case key.Matches(msg, m.keys.FocusSearch):
				m.focus = paneSearch
				m.searchInput.SetValue("")
				m.filtered = nil
				m.searchCursor = 0
				m.searchNavMode = false
				cmds = append(cmds, m.searchInput.Focus())
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.FocusTokens):
				m.focus = paneTokens
				m.tokenInputs[m.tokenFocus].Focus()
			case key.Matches(msg, m.keys.FocusCmd):
				m.focus = paneCmd
			case key.Matches(msg, m.keys.ContentNextSection):
				hs := headings(m.doc)
				if m.headingIdx < len(hs)-1 {
					m.headingIdx++
					m.codeNavActive = false
					m.codeEditMode = false
					m.editText = ""
					m.editCursor = 0
					cmds = append(cmds, m.rerenderFile())
				}
			case key.Matches(msg, m.keys.ContentPrevSection):
				if m.headingIdx > 0 {
					m.headingIdx--
					m.codeNavActive = false
					m.codeEditMode = false
					m.editText = ""
					m.editCursor = 0
					cmds = append(cmds, m.rerenderFile())
				}
			case key.Matches(msg, m.keys.Enter):
				if m.codeNavActive && m.codeIdx < len(m.visibleCodes) {
					m.codeEditMode = true
					m.editText = parser.StripANSI(m.visibleCodes[m.codeIdx].Text)
					m.editCursor = len([]rune(m.editText))
					cmds = append(cmds, m.rerenderFile())
				} else {
					hs := headings(m.doc)
					if len(hs) > 0 && m.headingIdx < len(hs) {
						k := hs[m.headingIdx].Text
						if m.collapsed == nil {
							m.collapsed = map[string]bool{}
						}
						m.collapsed[k] = !m.collapsed[k]
						cmds = append(cmds, m.rerenderFile())
					}
				}
			case key.Matches(msg, m.keys.ContentNextCode):
				if len(m.visibleCodes) > 0 {
					next := min(m.codeIdx+1, len(m.visibleCodes)-1)
					if next != m.codeIdx {
						m.editText = ""
						m.editCursor = 0
						m.codeEditMode = false
					}
					m.codeNavActive = true
					m.codeIdx = next
					cmds = append(cmds, m.rerenderFile())
				}
			case key.Matches(msg, m.keys.ContentPrevCode):
				if len(m.visibleCodes) > 0 {
					prev := max(m.codeIdx-1, 0)
					if prev != m.codeIdx {
						m.editText = ""
						m.editCursor = 0
						m.codeEditMode = false
					}
					m.codeNavActive = true
					m.codeIdx = prev
					cmds = append(cmds, m.rerenderFile())
				}
			case key.Matches(msg, m.keys.ContentUseToken):
				if m.codeNavActive && m.codeIdx < len(m.visibleCodes) {
					b := m.visibleCodes[m.codeIdx]
					if b.TokenOption != "" {
						val := parser.ExtractTokenOptionValue(parser.StripANSI(b.Text))
						m.localTokenValues[b.TokenOption] = val
						cmds = append(cmds, m.rerenderFile())
					}
				}
			case key.Matches(msg, m.keys.ContentRunCode):
				if m.codeNavActive && m.codeIdx < len(m.visibleCodes) {
					execText := m.editText
					if execText == "" {
						execText = parser.StripANSI(m.visibleCodes[m.codeIdx].Text)
					}
					m.cmdCurrentLine = execText
					m.cmdCursorPos = len([]rune(execText))
					m.cmdPromptAtTop = false
					m.cmdUpdateViewport()
					m.focus = paneCmd
				}
			case key.Matches(msg, m.keys.Tab):
				if m.codeNavActive && m.codeIdx < len(m.visibleCodes) {
					blockTokens := parser.ExtractLocalTokens(parser.StripANSI(m.visibleCodes[m.codeIdx].Text), tokenKeys)
					if len(blockTokens) > 0 {
						m.localTokensInBlock = blockTokens
						m.localTokenEditActive = true
						m.localTokenFocus = 0
						m.localTokenInput.SetValue(m.localTokenValues[blockTokens[0]])
						cmds = append(cmds, m.localTokenInput.Focus())
						cmds = append(cmds, m.rerenderFile())
					}
				}
			case key.Matches(msg, m.keys.ContentScrollUp):
				m.vp.LineUp(1)
			case key.Matches(msg, m.keys.ContentScrollDown):
				m.vp.LineDown(1)
			default:
				var cmd tea.Cmd
				m.vp, cmd = m.vp.Update(msg)
				cmds = append(cmds, cmd)
			}

		case paneTokens:
			switch {
			case key.Matches(msg, m.keys.Escape), key.Matches(msg, m.keys.Enter):
				m.tokenInputs[m.tokenFocus].Blur()
				m.focus = paneContent
				m.commitTokens()
				cmds = append(cmds, m.rerenderFile())
			case key.Matches(msg, m.keys.Tab), key.Matches(msg, m.keys.Right):
				m.tokenInputs[m.tokenFocus].Blur()
				m.tokenFocus = (m.tokenFocus + 1) % len(tokenKeys)
				m.tokenInputs[m.tokenFocus].Focus()
			case key.Matches(msg, m.keys.ShiftTab), key.Matches(msg, m.keys.Left):
				m.tokenInputs[m.tokenFocus].Blur()
				m.tokenFocus = (m.tokenFocus + len(tokenKeys) - 1) % len(tokenKeys)
				m.tokenInputs[m.tokenFocus].Focus()
			default:
				var cmd tea.Cmd
				m.tokenInputs[m.tokenFocus], cmd = m.tokenInputs[m.tokenFocus].Update(msg)
				cmds = append(cmds, cmd)
				cmds = append(cmds, m.rerenderFile())
			}

		case paneCmd:
			return handleCmdKey(m, msg)
		}
	}

	return m, tea.Batch(cmds...)
}
