package main

import (
	"pt-tui/internal/parser"

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

	case tea.KeyMsg:
		if msg.String() == "ctrl+e" {
			return handleModeToggle(m)
		}
		if m.mode == modeEditor && (m.focus == paneSidebar || m.focus == paneContent) {
			return handleEditorKey(m, msg)
		}

		switch msg.String() {
		case "ctrl+shift+down":
			m.cmdHeight = 10
			m.vp.Height = m.bodyH() - 2
			m.cmdVp.Height = max(1, m.cmdH()-2)
		case "ctrl+shift+up":
			m.cmdHeight = m.height * 2 / 3
			m.vp.Height = m.bodyH() - 2
			m.cmdVp.Height = max(1, m.cmdH()-2)
		}

		switch m.focus {
		case paneSearch:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "escape", "esc":
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				m.filtered = nil
				m.sidebarOffset = 0
				m.searchNavMode = false
				m.focus = paneSidebar
			case "enter":
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
			case "down":
				list := m.searchList()
				next := m.searchCursor + 1
				for next < len(list) && list[next].Path == "" {
					next++
				}
				if next < len(list) {
					m.searchCursor = next
				}
				m.sidebarOffset = clampOffset(m.searchCursor, m.sidebarOffset, m.bodyH()-3)
			case "up":
				list := m.searchList()
				prev := m.searchCursor - 1
				for prev >= 0 && list[prev].Path == "" {
					prev--
				}
				if prev >= 0 {
					m.searchCursor = prev
				}
				m.sidebarOffset = clampOffset(m.searchCursor, m.sidebarOffset, m.bodyH()-3)
			case "k", "j":
				if m.searchNavMode {
					list := m.searchList()
					if msg.String() == "k" {
						next := m.searchCursor + 1
						for next < len(list) && list[next].Path == "" {
							next++
						}
						if next < len(list) {
							m.searchCursor = next
						}
					} else {
						prev := m.searchCursor - 1
						for prev >= 0 && list[prev].Path == "" {
							prev--
						}
						if prev >= 0 {
							m.searchCursor = prev
						}
					}
					m.sidebarOffset = clampOffset(m.searchCursor, m.sidebarOffset, m.bodyH()-3)
				} else {
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
			case "'":
				if m.searchNavMode {
					m.searchNavMode = false
					cmds = append(cmds, m.searchInput.Focus())
				}
			default:
				if len(msg.Runes) == 1 && (msg.Runes[0] == 'K' || msg.Runes[0] == 'J') {
					isK := msg.Runes[0] == 'K'
					if m.searchNavMode {
						list := m.searchList()
						if isK {
							next := m.searchCursor + 1
							for next < len(list) && list[next].Path != "" {
								next++
							}
							if next < len(list) {
								m.searchCursor = next
							}
						} else {
							prev := m.searchCursor - 1
							for prev >= 0 && list[prev].Path != "" {
								prev--
							}
							if prev >= 0 {
								m.searchCursor = prev
							}
						}
						m.sidebarOffset = clampOffset(m.searchCursor, m.sidebarOffset, m.bodyH()-3)
						break
					}
				}
				if m.searchNavMode && len(msg.Runes) > 0 {
					m.searchNavMode = false
					cmds = append(cmds, m.searchInput.Focus())
				}
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

		case paneSidebar:
			tree := m.collapsibleTree()
			switch msg.String() {
			case "enter":
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
			case "'":
				m.focus = paneSearch
				m.searchInput.SetValue("")
				m.filtered = nil
				m.searchCursor = 0
				m.searchNavMode = false
				cmds = append(cmds, m.searchInput.Focus())
			case "ctrl+c", "q":
				return m, tea.Quit
			case "k", "down":
				m.cursor = min(m.cursor+1, len(tree)-1)
				m.sidebarOffset = clampOffset(m.cursor, m.sidebarOffset, m.bodyH()-2)
			case "j", "up":
				m.cursor = max(m.cursor-1, 0)
				m.sidebarOffset = clampOffset(m.cursor, m.sidebarOffset, m.bodyH()-2)
			case "tab", "l":
				m.focus = paneContent
			case "t":
				m.focus = paneTokens
				m.tokenInputs[m.tokenFocus].Focus()
			case "C":
				m.focus = paneCmd
			default:
				if len(msg.Runes) == 1 && (msg.Runes[0] == 'K' || msg.Runes[0] == 'J') {
					isK := msg.Runes[0] == 'K'
					if isK {
						next := m.cursor + 1
						for next < len(tree) && tree[next].Path != "" {
							next++
						}
						if next < len(tree) {
							m.cursor = next
						}
					} else {
						prev := m.cursor - 1
						for prev >= 0 && tree[prev].Path != "" {
							prev--
						}
						if prev >= 0 {
							m.cursor = prev
						}
					}
					m.sidebarOffset = clampOffset(m.cursor, m.sidebarOffset, m.bodyH()-2)
				}
			}

		case paneContent:
			if m.codeEditMode {
				switch msg.String() {
				case "ctrl+c":
					return m, tea.Quit
				case "enter":
					m.codeEditMode = false
				case "escape", "esc":
					m.codeEditMode = false
					m.editText = ""
					m.editCursor = 0
				case "left":
					if m.editCursor > 0 {
						m.editCursor--
					}
				case "right":
					if m.editCursor < len([]rune(m.editText)) {
						m.editCursor++
					}
				case "alt+left":
					m.editCursor = parser.PrevWord(m.editText, m.editCursor)
				case "alt+right":
					m.editCursor = parser.NextWord(m.editText, m.editCursor)
				case "ctrl+a":
					m.editCursor = 0
				case "ctrl+e":
					m.editCursor = len([]rune(m.editText))
				case "backspace":
					if m.editCursor > 0 {
						r := []rune(m.editText)
						m.editText = string(r[:m.editCursor-1]) + string(r[m.editCursor:])
						m.editCursor--
					}
				case "delete":
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
				switch msg.String() {
				case "ctrl+c":
					return m, tea.Quit
				case "escape", "esc":
					m.localTokenInput.Blur()
					m.localTokenEditActive = false
					cmds = append(cmds, m.rerenderFile())
				case "tab":
					m.localTokenValues[m.localTokensInBlock[m.localTokenFocus]] = m.localTokenInput.Value()
					m.localTokenFocus = (m.localTokenFocus + 1) % len(m.localTokensInBlock)
					m.localTokenInput.SetValue(m.localTokenValues[m.localTokensInBlock[m.localTokenFocus]])
					cmds = append(cmds, m.rerenderFile())
				case "shift+tab":
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

			switch msg.String() {
			case "escape", "esc":
				if m.codeNavActive {
					m.codeNavActive = false
					cmds = append(cmds, m.rerenderFile())
				} else {
					m.focus = paneSidebar
				}
			case "h":
				m.focus = paneSidebar
			case "'":
				m.focus = paneSearch
				m.searchInput.SetValue("")
				m.filtered = nil
				m.searchCursor = 0
				m.searchNavMode = false
				cmds = append(cmds, m.searchInput.Focus())
			case "ctrl+c", "q":
				return m, tea.Quit
			case "t":
				m.focus = paneTokens
				m.tokenInputs[m.tokenFocus].Focus()
			case "C":
				m.focus = paneCmd
			case "right", "k":
				hs := headings(m.doc)
				if m.headingIdx < len(hs)-1 {
					m.headingIdx++
					m.codeNavActive = false
					m.codeEditMode = false
					m.editText = ""
					m.editCursor = 0
					cmds = append(cmds, m.rerenderFile())
				}
			case "left", "j":
				if m.headingIdx > 0 {
					m.headingIdx--
					m.codeNavActive = false
					m.codeEditMode = false
					m.editText = ""
					m.editCursor = 0
					cmds = append(cmds, m.rerenderFile())
				}
			case "enter":
				if m.codeNavActive && m.codeIdx < len(m.visibleCodes) {
					m.codeEditMode = true
					m.editText = parser.StripANSI(m.visibleCodes[m.codeIdx].Text)
					m.editCursor = len([]rune(m.editText))
					cmds = append(cmds, m.rerenderFile())
				} else {
					hs := headings(m.doc)
					if len(hs) > 0 && m.headingIdx < len(hs) {
						key := hs[m.headingIdx].Text
						if m.collapsed == nil {
							m.collapsed = map[string]bool{}
						}
						m.collapsed[key] = !m.collapsed[key]
						cmds = append(cmds, m.rerenderFile())
					}
				}
			case "f":
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
			case "d":
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
			case "u":
				if m.codeNavActive && m.codeIdx < len(m.visibleCodes) {
					b := m.visibleCodes[m.codeIdx]
					if b.TokenOption != "" {
						val := parser.ExtractTokenOptionValue(parser.StripANSI(b.Text))
						m.localTokenValues[b.TokenOption] = val
						cmds = append(cmds, m.rerenderFile())
					}
				}
			case "x":
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
			case "tab":
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
			case "l", "up":
				m.vp.LineUp(1)
			case "ö", "down":
				m.vp.LineDown(1)
			default:
				var cmd tea.Cmd
				m.vp, cmd = m.vp.Update(msg)
				cmds = append(cmds, cmd)
			}

		case paneTokens:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "escape", "esc":
				m.tokenInputs[m.tokenFocus].Blur()
				m.focus = paneContent
				m.commitTokens()
				cmds = append(cmds, m.rerenderFile())
			case "tab", "right":
				m.tokenInputs[m.tokenFocus].Blur()
				m.tokenFocus = (m.tokenFocus + 1) % len(tokenKeys)
				m.tokenInputs[m.tokenFocus].Focus()
			case "shift+tab", "left":
				m.tokenInputs[m.tokenFocus].Blur()
				m.tokenFocus = (m.tokenFocus + len(tokenKeys) - 1) % len(tokenKeys)
				m.tokenInputs[m.tokenFocus].Focus()
			case "enter":
				m.tokenInputs[m.tokenFocus].Blur()
				m.focus = paneContent
				m.commitTokens()
				cmds = append(cmds, m.rerenderFile())
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
