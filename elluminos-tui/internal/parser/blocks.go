package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type BlockKind int

const (
	BlockHeading BlockKind = iota
	BlockCode
	BlockText
	BlockRule
)

type Block struct {
	Kind        BlockKind
	Level       int
	Lang        string
	Text        string
	TokenOption string
}

type ParsedDoc struct {
	Blocks []Block
}

func ParseDoc(src string) ParsedDoc {
	md := goldmark.New()
	source := []byte(src)
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	var blocks []Block
	currentTokenSection := ""

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {

		case *ast.HTMLBlock:
			var buf bytes.Buffer
			for i := 0; i < node.Lines().Len(); i++ {
				line := node.Lines().At(i)
				buf.Write(line.Value(source))
			}
			if m := reTokenSection.FindStringSubmatch(strings.TrimSpace(buf.String())); m != nil {
				currentTokenSection = "<" + m[1] + ">"
			}
			return ast.WalkSkipChildren, nil

		case *ast.Heading:
			if node.Level <= 2 {
				currentTokenSection = ""
			}
			var buf bytes.Buffer
			collectText(node, source, &buf)
			blocks = append(blocks, Block{
				Kind:  BlockHeading,
				Level: node.Level,
				Text:  strings.TrimSpace(buf.String()),
			})
			return ast.WalkSkipChildren, nil

		case *ast.FencedCodeBlock:
			lang := ""
			if node.Info != nil {
				lang = strings.TrimSpace(string(node.Info.Segment.Value(source)))
			}
			var buf bytes.Buffer
			for i := 0; i < node.Lines().Len(); i++ {
				line := node.Lines().At(i)
				buf.Write(line.Value(source))
			}
			blocks = append(blocks, Block{
				Kind:        BlockCode,
				Lang:        lang,
				Text:        strings.TrimRight(buf.String(), "\n"),
				TokenOption: currentTokenSection,
			})
			return ast.WalkSkipChildren, nil

		case *ast.CodeBlock:
			var buf bytes.Buffer
			for i := 0; i < node.Lines().Len(); i++ {
				line := node.Lines().At(i)
				buf.Write(line.Value(source))
			}
			blocks = append(blocks, Block{
				Kind:        BlockCode,
				Text:        strings.TrimRight(buf.String(), "\n"),
				TokenOption: currentTokenSection,
			})
			return ast.WalkSkipChildren, nil

		case *ast.ThematicBreak:
			blocks = append(blocks, Block{Kind: BlockRule})
			return ast.WalkSkipChildren, nil

		case *ast.List:
			var buf bytes.Buffer
			collectList(node, source, &buf, 0)
			if t := strings.TrimSpace(buf.String()); t != "" {
				blocks = append(blocks, Block{Kind: BlockText, Text: t})
			}
			return ast.WalkSkipChildren, nil

		case *ast.Paragraph:
			var buf bytes.Buffer
			collectText(n, source, &buf)
			t := strings.TrimSpace(buf.String())
			if t != "" {
				blocks = append(blocks, Block{Kind: BlockText, Text: t})
			}
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})

	return ParsedDoc{Blocks: blocks}
}

func collectList(n *ast.List, source []byte, buf *bytes.Buffer, depth int) {
	bullet := "• "
	if !n.IsOrdered() == false {
		bullet = "• "
	}
	indent := strings.Repeat("  ", depth)
	num := 1
	for item := n.FirstChild(); item != nil; item = item.NextSibling() {
		if n.IsOrdered() {
			buf.WriteString(indent + fmt.Sprintf("%d. ", num))
			num++
		} else {
			buf.WriteString(indent + bullet)
		}
		for c := item.FirstChild(); c != nil; c = c.NextSibling() {
			if sub, ok := c.(*ast.List); ok {
				buf.WriteByte('\n')
				collectList(sub, source, buf, depth+1)
			} else {
				var tb bytes.Buffer
				collectText(c, source, &tb)
				buf.WriteString(strings.TrimSpace(tb.String()))
			}
		}
		buf.WriteByte('\n')
	}
}

func collectText(n ast.Node, source []byte, buf *bytes.Buffer) {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		switch node := c.(type) {
		case *ast.Text:
			buf.Write(node.Segment.Value(source))
			if node.HardLineBreak() || node.SoftLineBreak() {
				buf.WriteByte('\n')
			}
		case *ast.String:
			buf.Write(node.Value)
		case *ast.CodeSpan:
			buf.WriteByte('`')
			collectText(node, source, buf)
			buf.WriteByte('`')
		default:
			collectText(c, source, buf)
		}
	}
}
