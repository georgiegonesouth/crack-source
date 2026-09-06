package main

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

type appMode int

const (
	modeNotebook appMode = iota
	modeEditor
)
