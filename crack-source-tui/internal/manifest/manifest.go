package manifest

import (
	"encoding/json"
	"os"
	"strings"
)

type manifestEntry struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Path     string          `json:"path"`
	Children []manifestEntry `json:"children"`
}

type Entry struct {
	Label   string
	Path    string
	Depth   int
	OrigIdx int
}

func LoadFiles(manifestPath string) []Entry {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil
	}
	var manifest struct {
		Tree []manifestEntry `json:"tree"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil
	}
	var out []Entry
	var walk func(items []manifestEntry, depth int)
	walk = func(items []manifestEntry, depth int) {
		for _, item := range items {
			if item.Type == "dir" {
				out = append(out, Entry{Label: strings.ToUpper(item.Name), Path: "", Depth: depth, OrigIdx: len(out)})
				walk(item.Children, depth+1)
			} else {
				label := strings.TrimSuffix(strings.TrimSuffix(item.Name, ".md"), ".pdf")
				out = append(out, Entry{Label: label, Path: item.Path, Depth: depth, OrigIdx: len(out)})
			}
		}
	}
	walk(manifest.Tree, 0)
	return out
}
