package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// MapIDName checks that map elements have matching id and name attributes.
type MapIDName struct{}

// Name returns the rule identifier.
func (r *MapIDName) Name() string { return RuleMapIDName }

// Description returns what this rule checks.
func (r *MapIDName) Description() string {
	return "map element id and name attributes should match for compatibility"
}

// Check examines the document for map elements with mismatched id and name.
func (r *MapIDName) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("map") {
			return true
		}

		id := n.GetAttr("id")
		name := n.GetAttr("name")

		// Skip template values
		if id == TemplateExprPlaceholder || name == TemplateExprPlaceholder {
			return true
		}

		// Map must have name attribute
		if name == "" {
			results = append(results, Result{
				Rule:     RuleMapIDName,
				Message:  "map element must have name attribute",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return true
		}

		// If both present, they should match
		if id != "" && id != name {
			results = append(results, Result{
				Rule:     RuleMapIDName,
				Message:  "map id=\"" + id + "\" and name=\"" + name + "\" should match",
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
