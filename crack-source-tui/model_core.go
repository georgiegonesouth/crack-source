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
	// terminal dimensions and app-wide config
	width, height int
	cfg           tuiConfig
	keys          KeyMap
	profile       profile

	// sidebar: manifest entries and navigation state
	files  []manifest.Entry
	cursor int

	// content pane: rendered document and scroll position
	vp        viewport.Model
	rawFile   string
	doc       parser.ParsedDoc
	tipsDoc   parser.ParsedDoc
	collapsed map[string]bool

	// content pane: heading and code block navigation
	headingIdx    int
	h2Lines       []int
	codeIdx       int
	codeNavActive bool
	visibleCodes  []parser.Block

	// per-page token bar: inputs for page-local placeholders
	localTokenKeys       []string
	localTokenValues     map[string]string
	localTokenFocus      int
	localTokenInput      textinput.Model
	localTokenEditActive bool
	localTokensInBlock   []string

	// inline code block editing
	codeEditMode bool
	editText     string
	editCursor   int

	// global token bar: TARGET_IP, PORT, LHOST, LPORT, USER, PASSWORD, DOMAIN
	tokenInputs [7]textinput.Model
	tokenFocus  int

	// cmd pane: shell state, scrollback, and history
	cmdScrollback  string
	cmdCurrentLine string
	cmdCursorPos   int
	cmdVp          viewport.Model
	cmdPromptAtTop bool
	cmdClearLine   int
	cmdHistory     []string
	cmdHistIdx     int
	cmdHeight      int

	// which pane has keyboard focus
	focus pane

	// background process tracking (long-running shell commands)
	runningProc  *exec.Cmd
	procSiginted bool
	ctrlCPending bool

	// mouse text selection state
	selActive    bool
	selAX, selAY int
	selEX, selEY int
	selX0, selX1 int

	// app mode, editor state, and shared working directory
	mode    appMode
	editor  editorState
	workDir string

	// sidebar search: input, results, and collapsed folder state
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
