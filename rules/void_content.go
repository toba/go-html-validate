package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// VoidContent checks that void elements have no children.
type VoidContent struct{}

// Name returns the rule identifier.
func (r *VoidContent) Name() string { return RuleVoidContent }

// Description returns what this rule checks.
func (r *VoidContent) Description() string {
	return "void elements must not have content"
}

// Check examines the document for void elements with children.
func (r *VoidContent) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Check if element is void
		if !VoidElements[tag] {
			return true
		}

		// Check for child elements (text nodes are also invalid)
		for _, child := range n.Children {
			if child.Type == html.ElementNode {
				results = append(results, Result{
					Rule:     RuleVoidContent,
					Message:  "void element <" + tag + "> must not have child elements",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
				break
			}
			// Check for non-whitespace text content
			if child.Type == html.TextNode && strings.TrimSpace(child.Data) != "" {
				results = append(results, Result{
					Rule:     RuleVoidContent,
					Message:  "void element <" + tag + "> must not have text content",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
				break
			}
		}

		return true
	})

	return results
}
