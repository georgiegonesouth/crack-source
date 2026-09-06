package main

import "github.com/charmbracelet/bubbles/key"

// KeyMap holds one key.Binding per distinct action across all panes.
// The same binding may be reused in multiple panes when the same physical key
// triggers semantically different behaviour — context routing in Update handles
// that disambiguation. Editor sidebar navigation uses separate bindings from
// notebook sidebar because j/k directions are inverted between the two modes.
type KeyMap struct {
	// Application-level
	Quit       key.Binding
	ModeToggle key.Binding

	// Pane focus
	FocusCmd     key.Binding
	FocusTokens  key.Binding
	FocusSearch  key.Binding
	FocusSidebar key.Binding
	FocusContent key.Binding

	// Cmd pane height
	ResizeCmdSnapMin key.Binding
	ResizeCmdSnapMax key.Binding
	ResizeCmdSmaller key.Binding
	ResizeCmdLarger  key.Binding

	// Shared primitives used across panes for text editing / cursor movement
	Escape     key.Binding
	Enter      key.Binding
	Tab        key.Binding
	ShiftTab   key.Binding
	Backspace  key.Binding
	Delete     key.Binding
	Left       key.Binding
	Right      key.Binding
	Up         key.Binding
	Down       key.Binding
	Home       key.Binding
	End        key.Binding
	WordLeft   key.Binding
	WordRight  key.Binding
	WordDelete key.Binding

	// Notebook sidebar navigation (also shared with search pane nav mode)
	SidebarDown     key.Binding
	SidebarUp       key.Binding
	SidebarJumpDown key.Binding
	SidebarJumpUp   key.Binding

	// Notebook content pane
	ContentNextSection key.Binding
	ContentPrevSection key.Binding
	ContentScrollUp    key.Binding
	ContentScrollDown  key.Binding
	ContentNextCode    key.Binding
	ContentPrevCode    key.Binding
	ContentRunCode     key.Binding
	ContentUseToken    key.Binding

	// Cmd pane actions (separate from primitives so users can rebind independently)
	CmdRun         key.Binding
	CmdClear       key.Binding
	CmdHistoryUp   key.Binding
	CmdHistoryDown key.Binding
	CmdCopyLine    key.Binding
	CmdPaste       key.Binding

	// Editor mode — sidebar navigation (j/k inverted vs notebook sidebar)
	EditorSidebarUp   key.Binding
	EditorSidebarDown key.Binding
	EditorSidebarOpen key.Binding
	EditorSave        key.Binding
}

var DefaultKeyMap = KeyMap{
	Quit:       key.NewBinding(key.WithKeys("ctrl+c", "q")),
	ModeToggle: key.NewBinding(key.WithKeys("ctrl+e")),

	FocusCmd:     key.NewBinding(key.WithKeys("C")),
	FocusTokens:  key.NewBinding(key.WithKeys("t")),
	FocusSearch:  key.NewBinding(key.WithKeys("'")),
	FocusSidebar: key.NewBinding(key.WithKeys("h")),
	FocusContent: key.NewBinding(key.WithKeys("tab", "l")),

	ResizeCmdSnapMin: key.NewBinding(key.WithKeys("ctrl+shift+down")),
	ResizeCmdSnapMax: key.NewBinding(key.WithKeys("ctrl+shift+up")),
	ResizeCmdSmaller: key.NewBinding(key.WithKeys("shift+down")),
	ResizeCmdLarger:  key.NewBinding(key.WithKeys("shift+up")),

	Escape:     key.NewBinding(key.WithKeys("escape", "esc")),
	Enter:      key.NewBinding(key.WithKeys("enter")),
	Tab:        key.NewBinding(key.WithKeys("tab")),
	ShiftTab:   key.NewBinding(key.WithKeys("shift+tab")),
	Backspace:  key.NewBinding(key.WithKeys("backspace", "ctrl+h")),
	Delete:     key.NewBinding(key.WithKeys("delete")),
	Left:       key.NewBinding(key.WithKeys("left")),
	Right:      key.NewBinding(key.WithKeys("right")),
	Up:         key.NewBinding(key.WithKeys("up")),
	Down:       key.NewBinding(key.WithKeys("down")),
	Home:       key.NewBinding(key.WithKeys("home", "ctrl+a")),
	End:        key.NewBinding(key.WithKeys("end")),
	WordLeft:   key.NewBinding(key.WithKeys("alt+left")),
	WordRight:  key.NewBinding(key.WithKeys("alt+right")),
	WordDelete: key.NewBinding(key.WithKeys("ctrl+w")),

	SidebarDown:     key.NewBinding(key.WithKeys("k", "down")),
	SidebarUp:       key.NewBinding(key.WithKeys("j", "up")),
	SidebarJumpDown: key.NewBinding(key.WithKeys("K")),
	SidebarJumpUp:   key.NewBinding(key.WithKeys("J")),

	ContentNextSection: key.NewBinding(key.WithKeys("right", "k")),
	ContentPrevSection: key.NewBinding(key.WithKeys("left", "j")),
	ContentScrollUp:    key.NewBinding(key.WithKeys("l", "up")),
	ContentScrollDown:  key.NewBinding(key.WithKeys("ö", "down")),
	ContentNextCode:    key.NewBinding(key.WithKeys("f")),
	ContentPrevCode:    key.NewBinding(key.WithKeys("d")),
	ContentRunCode:     key.NewBinding(key.WithKeys("x")),
	ContentUseToken:    key.NewBinding(key.WithKeys("u")),

	CmdRun:         key.NewBinding(key.WithKeys("enter")),
	CmdClear:       key.NewBinding(key.WithKeys("ctrl+l")),
	CmdHistoryUp:   key.NewBinding(key.WithKeys("up")),
	CmdHistoryDown: key.NewBinding(key.WithKeys("down")),
	CmdCopyLine:    key.NewBinding(key.WithKeys("ctrl+shift+c")),
	CmdPaste:       key.NewBinding(key.WithKeys("ctrl+shift+v")),

	EditorSidebarUp:   key.NewBinding(key.WithKeys("up", "j")),
	EditorSidebarDown: key.NewBinding(key.WithKeys("down", "k")),
	EditorSidebarOpen: key.NewBinding(key.WithKeys("enter", "l")),
	EditorSave:        key.NewBinding(key.WithKeys("ctrl+s")),
}
