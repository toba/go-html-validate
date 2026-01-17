package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoDupAttr checks for duplicate attributes on an element.
// Note: The Go HTML parser already handles this by keeping the first occurrence,
// so this rule may not catch duplicates in the raw HTML. It validates the parsed result.
type NoDupAttr struct{}

// Name returns the rule identifier.
func (r *NoDupAttr) Name() string { return RuleNoDupAttr }

// Description returns what this rule checks.
func (r *NoDupAttr) Description() string {
	return "elements should not have duplicate attributes"
}

// Check examines the document for duplicate attributes.
func (r *NoDupAttr) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		// Track seen attribute names (case-insensitive)
		seen := make(map[string]bool)
		for _, attr := range n.Attr {
			key := strings.ToLower(attr.Key)
			if seen[key] {
				results = append(results, Result{
					Rule:     RuleNoDupAttr,
					Message:  "duplicate attribute: " + attr.Key,
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
			seen[key] = true
		}

		return true
	})

	return results
}
