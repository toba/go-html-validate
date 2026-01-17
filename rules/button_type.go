package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ButtonType checks that buttons have explicit type attributes.
type ButtonType struct{}

func (r *ButtonType) Name() string { return RuleButtonType }

func (r *ButtonType) Description() string {
	return "buttons should have explicit type attribute (submit, button, or reset)"
}

func (r *ButtonType) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("button") {
			return true
		}

		// Check if button has type attribute
		if !n.HasAttr("type") {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "button missing type attribute (defaults to submit)",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}

		return true
	})

	return results
}
