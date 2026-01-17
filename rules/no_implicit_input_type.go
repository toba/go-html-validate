package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoImplicitInputType suggests explicit type attribute on input elements.
type NoImplicitInputType struct{}

// Name returns the rule identifier.
func (r *NoImplicitInputType) Name() string { return RuleNoImplicitInputType }

// Description returns what this rule checks.
func (r *NoImplicitInputType) Description() string {
	return "input elements should have explicit type attribute"
}

// Check examines the document for inputs without type.
func (r *NoImplicitInputType) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if strings.ToLower(n.Data) != "input" {
			return true
		}

		if !n.HasAttr("type") {
			results = append(results, Result{
				Rule:     RuleNoImplicitInputType,
				Message:  "input should have explicit type attribute (defaults to \"text\")",
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
