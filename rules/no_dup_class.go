package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoDupClass checks for duplicate class names within a class attribute.
type NoDupClass struct{}

// Name returns the rule identifier.
func (r *NoDupClass) Name() string { return RuleNoDupClass }

// Description returns what this rule checks.
func (r *NoDupClass) Description() string {
	return "elements should not have duplicate class names"
}

// Check examines the document for duplicate class names.
func (r *NoDupClass) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		classAttr := n.GetAttr("class")
		if classAttr == "" || classAttr == TemplateExprPlaceholder {
			return true
		}

		// Split class names and check for duplicates
		classes := strings.Fields(classAttr)
		seen := make(map[string]bool)
		for _, class := range classes {
			// Skip template placeholders
			if class == TemplateExprPlaceholder {
				continue
			}
			if seen[class] {
				results = append(results, Result{
					Rule:     RuleNoDupClass,
					Message:  "duplicate class name: " + class,
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Warning,
				})
			}
			seen[class] = true
		}

		return true
	})

	return results
}
