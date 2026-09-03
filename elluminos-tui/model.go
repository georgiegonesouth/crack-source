package main

import (
	"pt-tui/internal/manifest"
	"pt-tui/internal/parser"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type tuiConfig struct {
	ManifestPath string
	RepoRoot     string
}

type pane int

const (
	paneSidebar pane = iota
	paneContent
	paneTokens
	paneCmd
	paneSearch
)

const sidebarWidth = 30
const tipsWidth = 28

var tokenKeys = [7]string{"TARGET_IP", "PORT", "LHOST", "LPORT", "USER", "PASSWORD", "DOMAIN"}

type Model struct {
	width, height int
	cfg           tuiConfig
	profile       profile
	files         []manifest.Entry
	cursor        int

	vp            viewport.Model
	rawFile       string
	doc           parser.ParsedDoc
	tipsDoc       parser.ParsedDoc
	collapsed     map[string]bool
	headingIdx    int
	h2Lines       []int
	codeIdx       int
	codeNavActive bool
	visibleCodes  []parser.Block

	// local token editing (active in code nav mode)
	localTokenKeys       []string
	localTokenValues     map[string]string
	localTokenFocus      int
	localTokenInput      textinput.Model
	localTokenEditActive bool
	localTokensInBlock   []string // tokens present in the currently selected block

	// manual code block editing
	codeEditMode bool
	editText     string
	editCursor   int

	tokenInputs [7]textinput.Model
	tokenFocus  int

	cmdInput   textinput.Model
	cmdHistory []string
	cmdHistIdx int
	cmdHeight  int
	output     string

	focus pane

	searchInput      textinput.Model
	searchCursor     int
	searchNavMode    bool // true after Enter commits query; j/k navigate instead of typing
	sidebarOffset    int
	filtered         []manifest.Entry
	sidebarCollapsed map[int]bool
	contentIndex     map[string]string
}

func newModel(cfg tuiConfig) Model {
	files := manifest.LoadFiles(cfg.ManifestPath)
	p := loadProfile("default")

	var tis [7]textinput.Model
	for i, k := range tokenKeys {
		ti := textinput.New()
		ti.Placeholder = k
		ti.Width = 14
		ti.SetValue(p.Tokens[k])
		tis[i] = ti
	}

	cmd := textinput.New()
	cmd.Placeholder = "command..."
	cmd.Width = 60

	search := textinput.New()
	search.Placeholder = "search..."
	search.Width = sidebarWidth - 4

	localTI := textinput.New()
	localTI.Width = 20

	return Model{
		cfg:              cfg,
		profile:          p,
		files:            files,
		sidebarCollapsed: initSidebarCollapsed(files),
		tokenInputs:      tis,
		cmdInput:         cmd,
		searchInput:      search,
		localTokenInput:  localTI,
		localTokenValues: map[string]string{},
		cmdHeight:        10,
		focus:            paneSidebar,
	}
}

func (m Model) Init() tea.Cmd {
	return buildContentIndex(m.cfg.RepoRoot, m.files)
}

func (m Model) bodyH() int { return m.height - 2 - m.tokenBarRows() - m.cmdH() }
func (m Model) cmdH() int  { return m.cmdHeight }
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
