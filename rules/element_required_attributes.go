package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ElementRequiredAttributes checks that elements have their required attributes.
type ElementRequiredAttributes struct{}

// Name returns the rule identifier.
func (r *ElementRequiredAttributes) Name() string { return RuleElementRequiredAttributes }

// Description returns what this rule checks.
func (r *ElementRequiredAttributes) Description() string {
	return "elements must have required attributes"
}

// Check examines the document for elements missing required attributes.
func (r *ElementRequiredAttributes) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Get element spec
		spec, hasSpec := ElementSpecs[tag]
		if !hasSpec || len(spec.RequiredAttributes) == 0 {
			return true
		}

		// Check for each required attribute
		for _, attr := range spec.RequiredAttributes {
			if !n.HasAttr(attr) {
				results = append(results, Result{
					Rule:     RuleElementRequiredAttributes,
					Message:  "<" + tag + "> requires attribute: " + attr,
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		}

		return true
	})

	return results
}
