package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// Deprecated checks for deprecated HTML elements.
type Deprecated struct{}

// Name returns the rule identifier.
func (r *Deprecated) Name() string { return RuleDeprecated }

// Description returns what this rule checks.
func (r *Deprecated) Description() string {
	return "deprecated HTML elements should not be used"
}

// Check examines the document for deprecated elements.
func (r *Deprecated) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Check if element is deprecated
		if suggestion, deprecated := DeprecatedElements[tag]; deprecated {
			results = append(results, Result{
				Rule:     RuleDeprecated,
				Message:  "element <" + tag + "> is deprecated; " + suggestion,
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
