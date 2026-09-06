package main

import (
	"os"
	"path/filepath"
	"strings"
)

type dirEntry struct {
	name  string
	isDir bool
}

type editorState struct {
	// directory browser
	dir        string
	dirEntries []dirEntry
	dirCursor  int
	dirOffset  int

	// file buffer
	lines    []string
	filePath string
	dirty    bool

	// cursor
	line int
	col  int

	// display scroll
	scrollTop int
	hScroll   int
}

func newEditorState(dir string) editorState {
	e := editorState{}
	if dir == "" {
		dir, _ = os.Getwd()
	}
	e.loadDir(dir)
	return e
}

func (e *editorState) loadDir(path string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return
	}
	e.dir = abs
	e.dirEntries = nil
	if abs != "/" {
		e.dirEntries = append(e.dirEntries, dirEntry{name: "..", isDir: true})
	}
	for _, de := range entries {
		e.dirEntries = append(e.dirEntries, dirEntry{name: de.Name(), isDir: de.IsDir()})
	}
	e.dirCursor = 0
	e.dirOffset = 0
}

func (e *editorState) reloadDir() {
	if e.dir == "" {
		return
	}
	name := ""
	if e.dirCursor < len(e.dirEntries) {
		name = e.dirEntries[e.dirCursor].name
	}
	e.loadDir(e.dir)
	if name != "" {
		for i, de := range e.dirEntries {
			if de.name == name {
				e.dirCursor = i
				return
			}
		}
	}
}

func (e *editorState) openFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	e.filePath = path
	e.lines = strings.Split(string(data), "\n")
	if len(e.lines) == 0 {
		e.lines = []string{""}
	}
	e.line = 0
	e.col = 0
	e.scrollTop = 0
	e.hScroll = 0
	e.dirty = false
}

func (e *editorState) saveFile() {
	if e.filePath == "" {
		return
	}
	_ = os.WriteFile(e.filePath, []byte(strings.Join(e.lines, "\n")), 0644)
	e.dirty = false
}

func (e *editorState) moveUp() {
	if e.line > 0 {
		e.line--
		e.clampCol()
	}
}

func (e *editorState) moveDown() {
	if e.line < len(e.lines)-1 {
		e.line++
		e.clampCol()
	}
}

func (e *editorState) moveLeft() {
	if e.col > 0 {
		e.col--
	} else if e.line > 0 {
		e.line--
		e.col = len([]rune(e.lines[e.line]))
	}
}

func (e *editorState) moveRight() {
	lineRunes := []rune(e.lines[e.line])
	if e.col < len(lineRunes) {
		e.col++
	} else if e.line < len(e.lines)-1 {
		e.line++
		e.col = 0
	}
}

func (e *editorState) moveLineStart() {
	e.col = 0
}

func (e *editorState) moveLineEnd() {
	if len(e.lines) > 0 {
		e.col = len([]rune(e.lines[e.line]))
	}
}

func (e *editorState) clampCol() {
	if len(e.lines) == 0 {
		e.col = 0
		return
	}
	lineLen := len([]rune(e.lines[e.line]))
	if e.col > lineLen {
		e.col = lineLen
	}
}

func (e *editorState) insertRune(r rune) {
	if len(e.lines) == 0 {
		e.lines = []string{""}
	}
	runes := []rune(e.lines[e.line])
	col := e.col
	if col > len(runes) {
		col = len(runes)
	}
	newRunes := make([]rune, len(runes)+1)
	copy(newRunes, runes[:col])
	newRunes[col] = r
	copy(newRunes[col+1:], runes[col:])
	e.lines[e.line] = string(newRunes)
	e.col++
	e.dirty = true
}

func (e *editorState) insertNewline() {
	if len(e.lines) == 0 {
		e.lines = []string{"", ""}
		e.line = 1
		e.col = 0
		e.dirty = true
		return
	}
	runes := []rune(e.lines[e.line])
	col := e.col
	if col > len(runes) {
		col = len(runes)
	}
	before := string(runes[:col])
	after := string(runes[col:])
	e.lines[e.line] = before
	newLines := make([]string, len(e.lines)+1)
	copy(newLines, e.lines[:e.line+1])
	newLines[e.line+1] = after
	copy(newLines[e.line+2:], e.lines[e.line+1:])
	e.lines = newLines
	e.line++
	e.col = 0
	e.dirty = true
}

func (e *editorState) deleteBackward() {
	if len(e.lines) == 0 {
		return
	}
	if e.col > 0 {
		runes := []rune(e.lines[e.line])
		col := e.col
		e.lines[e.line] = string(append(runes[:col-1], runes[col:]...))
		e.col--
		e.dirty = true
	} else if e.line > 0 {
		prevLen := len([]rune(e.lines[e.line-1]))
		e.lines[e.line-1] += e.lines[e.line]
		e.lines = append(e.lines[:e.line], e.lines[e.line+1:]...)
		e.line--
		e.col = prevLen
		e.dirty = true
	}
}

func (e *editorState) deleteForward() {
	if len(e.lines) == 0 {
		return
	}
	runes := []rune(e.lines[e.line])
	if e.col < len(runes) {
		e.lines[e.line] = string(append(runes[:e.col], runes[e.col+1:]...))
		e.dirty = true
	} else if e.line < len(e.lines)-1 {
		e.lines[e.line] += e.lines[e.line+1]
		e.lines = append(e.lines[:e.line+1], e.lines[e.line+2:]...)
		e.dirty = true
	}
}

func (e *editorState) updateScroll(visibleLines, visibleCols int) {
	if visibleLines < 1 {
		visibleLines = 1
	}
	if visibleCols < 1 {
		visibleCols = 1
	}
	if e.line < e.scrollTop {
		e.scrollTop = e.line
	}
	if e.line >= e.scrollTop+visibleLines {
		e.scrollTop = e.line - visibleLines + 1
	}
	if e.col < e.hScroll {
		e.hScroll = e.col
	}
	if e.col >= e.hScroll+visibleCols {
		e.hScroll = e.col - visibleCols + 1
	}
}
