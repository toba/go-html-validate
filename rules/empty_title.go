package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// EmptyTitle checks that title elements have text content.
type EmptyTitle struct{}

func (r *EmptyTitle) Name() string { return RuleEmptyTitle }

func (r *EmptyTitle) Description() string {
	return "<title> element must have text content"
}

func (r *EmptyTitle) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if strings.ToLower(n.Data) != "title" {
			return true
		}

		text := strings.TrimSpace(n.TextContent())
		if text == "" {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "<title> cannot be empty, must have text content",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		}

		return true
	})

	return results
}
