package main

import (
	"os"
	"path/filepath"
	"strings"

	"crack-source/internal/manifest"

	tea "github.com/charmbracelet/bubbletea"
)

type contentIndexedMsg map[string]string

func (m Model) collapsibleTree() []manifest.Entry {
	var out []manifest.Entry
	skipBelow := -1
	for _, f := range m.files {
		if skipBelow >= 0 && f.Depth > skipBelow {
			continue
		}
		skipBelow = -1
		out = append(out, f)
		if f.Path == "" && m.sidebarCollapsed[f.OrigIdx] {
			skipBelow = f.Depth
		}
	}
	return out
}

func (m Model) searchList() []manifest.Entry {
	if m.searchInput.Value() == "" {
		return m.collapsibleTree()
	}
	return m.filtered
}

func initSidebarCollapsed(files []manifest.Entry) map[int]bool {
	m := map[int]bool{}
	for _, f := range files {
		if f.Path == "" {
			m[f.OrigIdx] = true
		}
	}
	return m
}

func buildContentIndex(repoRoot string, files []manifest.Entry) tea.Cmd {
	return func() tea.Msg {
		index := make(map[string]string)
		for _, f := range files {
			if f.Path == "" || !strings.HasSuffix(f.Path, ".md") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(repoRoot, f.Path))
			if err != nil {
				continue
			}
			index[f.Path] = strings.ToLower(string(data))
		}
		return contentIndexedMsg(index)
	}
}

func clampOffset(cursor, offset, maxLines int) int {
	if cursor < offset {
		return cursor
	}
	if cursor >= offset+maxLines {
		return cursor - maxLines + 1
	}
	return offset
}

func (m Model) filterEntries(q string) []manifest.Entry {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return nil
	}

	matchSet := map[int]bool{}
	for i, f := range m.files {
		if f.Path == "" {
			continue
		}
		if strings.Contains(strings.ToLower(f.Label), q) || strings.Contains(m.contentIndex[f.Path], q) {
			matchSet[i] = true
		}
	}

	include := map[int]bool{}
	for idx := range matchSet {
		include[idx] = true
		depth := m.files[idx].Depth
		for j := idx - 1; j >= 0 && depth > 0; j-- {
			e := m.files[j]
			if e.Path == "" && e.Depth < depth {
				include[j] = true
				depth = e.Depth
			}
		}
	}

	var out []manifest.Entry
	for i, f := range m.files {
		if include[i] {
			out = append(out, f)
		}
	}
	return out
}
