package rules

import (
	"fmt"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// HeadingContent checks that heading elements have text content.
type HeadingContent struct{}

func (r *HeadingContent) Name() string { return RuleHeadingContent }

func (r *HeadingContent) Description() string {
	return "heading elements (h1-h6) must have text content"
}

func (r *HeadingContent) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if !HeadingTags[strings.ToLower(n.Data)] {
			return true
		}

		if !HasAccessibleName(n) {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  fmt.Sprintf("<%s> element has no text content", n.Data),
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
