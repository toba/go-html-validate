package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// RedundantAriaLabel checks for aria-label that duplicates visible text.
type RedundantAriaLabel struct{}

func (r *RedundantAriaLabel) Name() string { return RuleRedundantAriaLabel }

func (r *RedundantAriaLabel) Description() string {
	return "aria-label should not duplicate visible text content"
}

func (r *RedundantAriaLabel) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		ariaLabel := n.GetAttr("aria-label")
		if ariaLabel == "" {
			return true
		}

		// Get visible text content
		textContent := n.TextContent()

		// Normalize both for comparison
		normalizedAriaLabel := NormalizeText(ariaLabel)
		normalizedText := NormalizeText(textContent)

		// Check if they're the same (ignoring case and whitespace)
		if normalizedAriaLabel != "" && normalizedText != "" &&
			normalizedAriaLabel == normalizedText {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "aria-label duplicates visible text content",
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
