package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// WcagH63 checks that th elements have scope attribute.
// Per WCAG H63: Using the scope attribute to associate header cells and data cells.
type WcagH63 struct{}

// Name returns the rule identifier.
func (r *WcagH63) Name() string { return RuleWcagH63 }

// Description returns what this rule checks.
func (r *WcagH63) Description() string {
	return "th elements should have scope attribute for accessibility"
}

// Check examines the document for th elements without scope.
func (r *WcagH63) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("th") {
			return true
		}

		// Check for scope attribute
		if !n.HasAttr("scope") {
			results = append(results, Result{
				Rule:     RuleWcagH63,
				Message:  "th element should have scope attribute (col, row, colgroup, or rowgroup)",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
			return true
		}

		// Validate scope value
		scope := n.GetAttr("scope")
		if !ValidScopeValues[scope] && scope != TemplateExprPlaceholder {
			results = append(results, Result{
				Rule:     RuleWcagH63,
				Message:  "th scope attribute has invalid value: " + scope,
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
