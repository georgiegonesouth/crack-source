package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func loadConfig(binaryDir string) (tuiConfig, error) {
	configPath := filepath.Join(binaryDir, "tui.config")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return tuiConfig{}, fmt.Errorf("tui.config not found at %s: %w", configPath, err)
	}
	cfg := tuiConfig{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(k) {
		case "MANIFEST_PATH":
			cfg.ManifestPath = strings.TrimSpace(v)
		case "REPO_ROOT":
			cfg.RepoRoot = strings.TrimSpace(v)
		}
	}
	if cfg.ManifestPath == "" {
		return tuiConfig{}, fmt.Errorf("tui.config: MANIFEST_PATH not set")
	}
	if cfg.RepoRoot == "" {
		return tuiConfig{}, fmt.Errorf("tui.config: REPO_ROOT not set")
	}
	return cfg, nil
}

func keybindingsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "crack-source", "keybindings.yaml")
}

func main() {
	binaryDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg, err := loadConfig(binaryDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	km, err := LoadKeyMap(keybindingsPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: keybindings: %v — using defaults\n", err)
		km = DefaultKeyMap
	}

	p := tea.NewProgram(newModel(cfg, km), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
