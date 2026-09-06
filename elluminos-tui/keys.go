package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"gopkg.in/yaml.v3"
)

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

// DefaultKeyMap preserves every hardcoded keybinding that existed before this
// config system was introduced. Any LoadKeyMap call with a missing file returns
// this value unchanged.
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

// KeyConfig is the user-facing YAML shape: snake_case action name → key strings.
type KeyConfig map[string][]string

// LoadKeyMap reads keybindings from path. If the file does not exist, it
// returns DefaultKeyMap with no error. Any action not present in the file keeps
// its default binding (partial overrides are fully supported). Unknown action
// names and empty key lists produce stderr warnings but never return an error.
func LoadKeyMap(path string) (KeyMap, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DefaultKeyMap, nil
	}
	if err != nil {
		return DefaultKeyMap, err
	}
	var cfg KeyConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultKeyMap, fmt.Errorf("keybindings.yaml: %w", err)
	}

	km := DefaultKeyMap
	for action, keys := range cfg {
		if len(keys) == 0 {
			fmt.Fprintf(os.Stderr, "warning: keybinding %q has no keys, using default\n", action)
			continue
		}
		b := key.NewBinding(key.WithKeys(keys...))
		switch action {
		case "quit":
			km.Quit = b
		case "mode_toggle":
			km.ModeToggle = b
		case "focus_cmd":
			km.FocusCmd = b
		case "focus_tokens":
			km.FocusTokens = b
		case "focus_search":
			km.FocusSearch = b
		case "focus_sidebar":
			km.FocusSidebar = b
		case "focus_content":
			km.FocusContent = b
		case "resize_cmd_snap_min":
			km.ResizeCmdSnapMin = b
		case "resize_cmd_snap_max":
			km.ResizeCmdSnapMax = b
		case "resize_cmd_smaller":
			km.ResizeCmdSmaller = b
		case "resize_cmd_larger":
			km.ResizeCmdLarger = b
		case "escape":
			km.Escape = b
		case "enter":
			km.Enter = b
		case "tab":
			km.Tab = b
		case "shift_tab":
			km.ShiftTab = b
		case "backspace":
			km.Backspace = b
		case "delete":
			km.Delete = b
		case "left":
			km.Left = b
		case "right":
			km.Right = b
		case "up":
			km.Up = b
		case "down":
			km.Down = b
		case "home":
			km.Home = b
		case "end":
			km.End = b
		case "word_left":
			km.WordLeft = b
		case "word_right":
			km.WordRight = b
		case "word_delete":
			km.WordDelete = b
		case "sidebar_down":
			km.SidebarDown = b
		case "sidebar_up":
			km.SidebarUp = b
		case "sidebar_jump_down":
			km.SidebarJumpDown = b
		case "sidebar_jump_up":
			km.SidebarJumpUp = b
		case "content_next_section":
			km.ContentNextSection = b
		case "content_prev_section":
			km.ContentPrevSection = b
		case "content_scroll_up":
			km.ContentScrollUp = b
		case "content_scroll_down":
			km.ContentScrollDown = b
		case "content_next_code":
			km.ContentNextCode = b
		case "content_prev_code":
			km.ContentPrevCode = b
		case "content_run_code":
			km.ContentRunCode = b
		case "content_use_token":
			km.ContentUseToken = b
		case "cmd_run":
			km.CmdRun = b
		case "cmd_clear":
			km.CmdClear = b
		case "cmd_history_up":
			km.CmdHistoryUp = b
		case "cmd_history_down":
			km.CmdHistoryDown = b
		case "cmd_copy_line":
			km.CmdCopyLine = b
		case "cmd_paste":
			km.CmdPaste = b
		case "editor_sidebar_up":
			km.EditorSidebarUp = b
		case "editor_sidebar_down":
			km.EditorSidebarDown = b
		case "editor_sidebar_open":
			km.EditorSidebarOpen = b
		case "editor_save":
			km.EditorSave = b
		default:
			fmt.Fprintf(os.Stderr, "warning: unknown keybinding %q, ignoring\n", action)
		}
	}

	validateKeyMap(km)
	return km, nil
}

// validateKeyMap warns to stderr about any binding left with no keys.
// It never returns an error — an empty binding simply never matches.
func validateKeyMap(km KeyMap) {
	check := []struct {
		name string
		b    key.Binding
	}{
		{"quit", km.Quit},
		{"mode_toggle", km.ModeToggle},
		{"focus_cmd", km.FocusCmd},
		{"focus_tokens", km.FocusTokens},
		{"focus_search", km.FocusSearch},
		{"focus_sidebar", km.FocusSidebar},
		{"focus_content", km.FocusContent},
		{"resize_cmd_snap_min", km.ResizeCmdSnapMin},
		{"resize_cmd_snap_max", km.ResizeCmdSnapMax},
		{"resize_cmd_smaller", km.ResizeCmdSmaller},
		{"resize_cmd_larger", km.ResizeCmdLarger},
		{"escape", km.Escape},
		{"enter", km.Enter},
		{"tab", km.Tab},
		{"shift_tab", km.ShiftTab},
		{"backspace", km.Backspace},
		{"delete", km.Delete},
		{"left", km.Left},
		{"right", km.Right},
		{"up", km.Up},
		{"down", km.Down},
		{"home", km.Home},
		{"end", km.End},
		{"word_left", km.WordLeft},
		{"word_right", km.WordRight},
		{"word_delete", km.WordDelete},
		{"sidebar_down", km.SidebarDown},
		{"sidebar_up", km.SidebarUp},
		{"sidebar_jump_down", km.SidebarJumpDown},
		{"sidebar_jump_up", km.SidebarJumpUp},
		{"content_next_section", km.ContentNextSection},
		{"content_prev_section", km.ContentPrevSection},
		{"content_scroll_up", km.ContentScrollUp},
		{"content_scroll_down", km.ContentScrollDown},
		{"content_next_code", km.ContentNextCode},
		{"content_prev_code", km.ContentPrevCode},
		{"content_run_code", km.ContentRunCode},
		{"content_use_token", km.ContentUseToken},
		{"cmd_run", km.CmdRun},
		{"cmd_clear", km.CmdClear},
		{"cmd_history_up", km.CmdHistoryUp},
		{"cmd_history_down", km.CmdHistoryDown},
		{"cmd_copy_line", km.CmdCopyLine},
		{"cmd_paste", km.CmdPaste},
		{"editor_sidebar_up", km.EditorSidebarUp},
		{"editor_sidebar_down", km.EditorSidebarDown},
		{"editor_sidebar_open", km.EditorSidebarOpen},
		{"editor_save", km.EditorSave},
	}
	for _, c := range check {
		if len(c.b.Keys()) == 0 {
			fmt.Fprintf(os.Stderr, "warning: keybinding %q has no keys assigned\n", c.name)
		}
	}
}
