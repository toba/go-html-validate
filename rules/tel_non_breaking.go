package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// TelNonBreaking checks that tel: links use non-breaking formatting.
// Phone numbers should use &nbsp; or CSS to prevent awkward line breaks.
type TelNonBreaking struct{}

// Name returns the rule identifier.
func (r *TelNonBreaking) Name() string { return RuleTelNonBreaking }

// Description returns what this rule checks.
func (r *TelNonBreaking) Description() string {
	return "tel: links should use non-breaking spaces to prevent awkward line breaks"
}

// Check examines the document for tel: links with breaking spaces.
func (r *TelNonBreaking) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("a") {
			return true
		}

		// Check if this is a tel: link
		href := n.GetAttr("href")
		if !strings.HasPrefix(strings.ToLower(href), "tel:") {
			return true
		}

		// Check text content for regular spaces
		text := n.TextContent()
		if text == "" || text == TemplateExprPlaceholder {
			return true
		}

		// Look for regular spaces in the phone number text
		// Non-breaking space is \u00A0, regular space is \u0020
		if strings.Contains(text, " ") {
			results = append(results, Result{
				Rule:     RuleTelNonBreaking,
				Message:  "tel: link text contains regular spaces; use &nbsp; or CSS white-space: nowrap",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Info,
			})
		}

		return true
	})

	return results
}
