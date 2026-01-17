package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// WcagH67 checks that decorative images have empty alt and no title.
// Per WCAG H67: Using null alt text and no title attribute for decorative images.
type WcagH67 struct{}

// Name returns the rule identifier.
func (r *WcagH67) Name() string { return RuleWcagH67 }

// Description returns what this rule checks.
func (r *WcagH67) Description() string {
	return "decorative images (alt=\"\") should not have title attribute"
}

// Check examines the document for images with empty alt that also have title.
func (r *WcagH67) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("img") {
			return true
		}

		// Check for decorative image pattern (empty alt)
		alt := n.GetAttr("alt")
		if alt != "" {
			// Not a decorative image
			return true
		}

		// Decorative images should not have title
		if n.HasAttr("title") {
			title := n.GetAttr("title")
			if title != "" && title != TemplateExprPlaceholder {
				results = append(results, Result{
					Rule:     RuleWcagH67,
					Message:  "decorative image (alt=\"\") should not have title attribute",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Warning,
				})
			}
		}

		// Decorative images should not have role="img"
		role := strings.ToLower(n.GetAttr("role"))
		if role == "img" {
			results = append(results, Result{
				Rule:     RuleWcagH67,
				Message:  "decorative image (alt=\"\") should have role=\"presentation\" or role=\"none\"",
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
