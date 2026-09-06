package main

import (
	"os"
	"os/exec"

	"crack-source/internal/manifest"
	"crack-source/internal/parser"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width, height int
	cfg           tuiConfig
	keys          KeyMap
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

	localTokenKeys       []string
	localTokenValues     map[string]string
	localTokenFocus      int
	localTokenInput      textinput.Model
	localTokenEditActive bool
	localTokensInBlock   []string

	codeEditMode bool
	editText     string
	editCursor   int

	tokenInputs [7]textinput.Model
	tokenFocus  int

	cmdScrollback   string
	cmdCurrentLine  string
	cmdCursorPos    int
	cmdVp           viewport.Model
	cmdPromptAtTop  bool
	cmdClearLine    int
	cmdHistory      []string
	cmdHistIdx      int
	cmdHeight       int

	focus pane

	runningProc  *exec.Cmd
	ctrlCPending bool

	mode    appMode
	editor  editorState
	workDir string

	searchInput      textinput.Model
	searchCursor     int
	searchNavMode    bool
	sidebarOffset    int
	filtered         []manifest.Entry
	sidebarCollapsed map[int]bool
	contentIndex     map[string]string
}

func newModel(cfg tuiConfig) Model {
	files := manifest.LoadFiles(cfg.ManifestPath)
	p := loadProfile()
	cwd, _ := os.Getwd()

	var tis [7]textinput.Model
	for i, k := range tokenKeys {
		ti := textinput.New()
		ti.Placeholder = k
		ti.Width = 14
		ti.SetValue(p.Tokens[k])
		tis[i] = ti
	}

	search := textinput.New()
	search.Placeholder = "search..."
	search.Width = sidebarWidth - 4

	localTI := textinput.New()
	localTI.Width = 20

	return Model{
		cfg:              cfg,
		keys:             DefaultKeyMap,
		profile:          p,
		files:            files,
		sidebarCollapsed: initSidebarCollapsed(files),
		tokenInputs:      tis,
		searchInput:      search,
		localTokenInput:  localTI,
		localTokenValues: map[string]string{},
		cmdHeight:        10,
		focus:            paneSidebar,
		workDir:          cwd,
	}
}

func (m Model) Init() tea.Cmd {
	return buildContentIndex(m.cfg.RepoRoot, m.files)
}
