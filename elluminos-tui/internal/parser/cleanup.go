package parser

import (
	"regexp"
	"strings"
)

var (
	reFrontMatter  = regexp.MustCompile(`(?s)^---\n.*?\n---\n?`)
	reDirectives   = regexp.MustCompile(`(?s)<!--\s*token-(?:options|table|exclude):.*?-->`)
	ReToken        = regexp.MustCompile(`<[A-Za-z_][A-Za-z0-9_]*>`)
	reANSI         = regexp.MustCompile(`\x1b\[[0-9;]*[mKGHF]`)
	reLocalToken   = regexp.MustCompile(`<([a-z][a-z0-9_]*)>`)
	reTokenSection = regexp.MustCompile(`(?i)<!--\s*token-section:([A-Za-z_][A-Za-z0-9_]*)\s*-->`)
)

func StripFrontMatter(text string) string {
	return reFrontMatter.ReplaceAllString(text, "")
}

func StripDirectives(text string) string {
	return reDirectives.ReplaceAllString(text, "")
}

func StripANSI(s string) string {
	return strings.TrimSpace(reANSI.ReplaceAllString(s, ""))
}

func ExtractTokenOptionValue(text string) string {
	text = strings.TrimSpace(text)
	for _, prefix := range []string{"- ", "* ", "• "} {
		if strings.HasPrefix(text, prefix) {
			return strings.TrimSpace(text[len(prefix):])
		}
	}
	return text
}

func ExtractLocalTokens(raw string, globalKeys [7]string) []string {
	globals := make(map[string]bool, len(globalKeys))
	for _, k := range globalKeys {
		globals["<"+k+">"] = true
	}
	seen := map[string]bool{}
	var out []string
	for _, tok := range ReToken.FindAllString(raw, -1) {
		if !globals[tok] && !seen[tok] {
			seen[tok] = true
			out = append(out, tok)
		}
	}
	return out
}

func InsertCursorMark(text string, cursor int) string {
	runes := []rune(text)
	if cursor > len(runes) {
		cursor = len(runes)
	}
	return string(runes[:cursor]) + "█" + string(runes[cursor:])
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

func PrevWord(text string, cursor int) int {
	runes := []rune(text)
	if cursor == 0 {
		return 0
	}
	cursor--
	for cursor > 0 && !isWordChar(runes[cursor]) {
		cursor--
	}
	for cursor > 0 && isWordChar(runes[cursor-1]) {
		cursor--
	}
	return cursor
}

func NextWord(text string, cursor int) int {
	runes := []rune(text)
	n := len(runes)
	if cursor >= n {
		return n
	}
	for cursor < n && !isWordChar(runes[cursor]) {
		cursor++
	}
	for cursor < n && isWordChar(runes[cursor]) {
		cursor++
	}
	return cursor
}
