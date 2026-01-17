package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// LinkName checks that links have accessible names.
type LinkName struct{}

func (r *LinkName) Name() string { return RuleLinkName }

func (r *LinkName) Description() string {
	return "links must have text content or aria-label for accessibility"
}

func (r *LinkName) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("a") {
			return true
		}

		// Skip anchors without href (not really links)
		if !n.HasAttr("href") {
			return true
		}

		if !HasAccessibleName(n) {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "link element missing accessible name",
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
