package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoMultipleMain ensures only one visible <main> element per document.
type NoMultipleMain struct{}

func (r *NoMultipleMain) Name() string { return RuleNoMultipleMain }

func (r *NoMultipleMain) Description() string {
	return "only one visible <main> element allowed per document"
}

func (r *NoMultipleMain) Check(doc *parser.Document) []Result {
	var visibleMains []*parser.Node

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("main") {
			return true
		}

		// Check if element is hidden
		if n.HasAttr("hidden") {
			return true
		}

		visibleMains = append(visibleMains, n)
		return true
	})

	if len(visibleMains) <= 1 {
		return nil
	}

	// Report error on all but the first main element
	var results []Result
	for i := 1; i < len(visibleMains); i++ {
		n := visibleMains[i]
		results = append(results, Result{
			Rule:     r.Name(),
			Message:  "document has multiple visible <main> elements; remove this <main> or add hidden attribute",
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		})
	}

	return results
}
