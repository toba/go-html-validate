package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// WcagH71 checks that fieldset elements contain a legend.
// Per WCAG H71: Providing a description for groups of form controls using fieldset and legend elements.
type WcagH71 struct{}

// Name returns the rule identifier.
func (r *WcagH71) Name() string { return RuleWcagH71 }

// Description returns what this rule checks.
func (r *WcagH71) Description() string {
	return "fieldset elements must contain a legend element"
}

// Check examines the document for fieldsets without legend.
func (r *WcagH71) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("fieldset") {
			return true
		}

		// Check if fieldset has a legend child
		hasLegend := false
		for _, child := range n.Children {
			if child.Type == html.ElementNode && child.IsElement("legend") {
				hasLegend = true
				break
			}
		}

		if !hasLegend {
			results = append(results, Result{
				Rule:     RuleWcagH71,
				Message:  "fieldset element must contain a legend element",
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
