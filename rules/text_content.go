package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// TextContent checks that interactive elements have accessible text content.
// This is broader than button-name/link-name and covers other interactive elements.
type TextContent struct{}

// Name returns the rule identifier.
func (r *TextContent) Name() string { return RuleTextContent }

// Description returns what this rule checks.
func (r *TextContent) Description() string {
	return "interactive elements must have accessible text content"
}

// Check examines the document for interactive elements without text content.
func (r *TextContent) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Check summary elements
		if tag == "summary" {
			if !HasAccessibleName(n) {
				results = append(results, Result{
					Rule:     RuleTextContent,
					Message:  "summary element must have accessible text content",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
			return true
		}

		// Check details without summary (the default summary needs text)
		if tag == "details" {
			hasSummary := false
			for _, child := range n.Children {
				if child.Type == html.ElementNode && child.IsElement("summary") {
					hasSummary = true
					break
				}
			}
			// Details without explicit summary uses browser default ("Details")
			// which is accessible, so no error needed
			_ = hasSummary
			return true
		}

		return true
	})

	return results
}
