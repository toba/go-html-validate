package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ButtonName checks that buttons have accessible names.
type ButtonName struct{}

func (r *ButtonName) Name() string { return RuleButtonName }

func (r *ButtonName) Description() string {
	return "buttons must have text content or aria-label for accessibility"
}

func (r *ButtonName) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("button") {
			return true
		}

		if !HasAccessibleName(n) {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "button element missing accessible name",
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
