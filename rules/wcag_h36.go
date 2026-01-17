package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// WcagH36 checks that input type="image" has alt text.
// Per WCAG H36: Using alt attributes on images used as submit buttons.
type WcagH36 struct{}

// Name returns the rule identifier.
func (r *WcagH36) Name() string { return RuleWcagH36 }

// Description returns what this rule checks.
func (r *WcagH36) Description() string {
	return "input type=\"image\" must have alt attribute describing the action"
}

// Check examines the document for image inputs missing alt text.
func (r *WcagH36) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("input") {
			return true
		}

		// Check if this is an image input
		inputType := strings.ToLower(n.GetAttr("type"))
		if inputType != "image" {
			return true
		}

		// Check for alt attribute
		if !n.HasAttr("alt") {
			results = append(results, Result{
				Rule:     RuleWcagH36,
				Message:  "input type=\"image\" must have alt attribute",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return true
		}

		// Check for empty alt
		alt := n.GetAttr("alt")
		if alt == "" {
			results = append(results, Result{
				Rule:     RuleWcagH36,
				Message:  "input type=\"image\" has empty alt attribute",
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
