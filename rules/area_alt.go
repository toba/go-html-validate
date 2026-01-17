package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// AreaAlt checks that <area> elements have alt text.
// Per WCAG 2.0 H24: Providing text alternatives for the area elements of image maps.
type AreaAlt struct{}

// Name returns the rule identifier.
func (r *AreaAlt) Name() string { return RuleAreaAlt }

// Description returns what this rule checks.
func (r *AreaAlt) Description() string {
	return "area elements must have alt text describing the link destination"
}

// Check examines the document for area elements missing alt text.
func (r *AreaAlt) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("area") {
			return true
		}

		// Only areas with href need alt (non-link areas don't need it)
		if !n.HasAttr("href") {
			return true
		}

		// Check for alt attribute
		if !n.HasAttr("alt") {
			results = append(results, Result{
				Rule:     RuleAreaAlt,
				Message:  "area element with href must have alt attribute",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return true
		}

		// Check for empty alt (allowed only if another area has same href with alt)
		alt := n.GetAttr("alt")
		if alt == "" || alt == TemplateExprPlaceholder {
			// Empty alt is a warning - may be intentional for redundant areas
			results = append(results, Result{
				Rule:     RuleAreaAlt,
				Message:  "area element has empty alt attribute",
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
