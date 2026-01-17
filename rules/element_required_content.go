package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ElementRequiredContent checks that elements have their required children.
type ElementRequiredContent struct{}

// Name returns the rule identifier.
func (r *ElementRequiredContent) Name() string { return RuleElementRequiredContent }

// Description returns what this rule checks.
func (r *ElementRequiredContent) Description() string {
	return "elements must have required child elements"
}

// Check examines the document for elements missing required children.
func (r *ElementRequiredContent) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Get element spec
		spec, hasSpec := ElementSpecs[tag]
		if !hasSpec || len(spec.RequiredChildren) == 0 {
			return true
		}

		// Build set of child element tags
		childTags := make(map[string]bool)
		for _, child := range n.Children {
			if child.Type == html.ElementNode {
				childTags[strings.ToLower(child.Data)] = true
			}
		}

		// Check for each required child
		for _, required := range spec.RequiredChildren {
			if !childTags[required] {
				results = append(results, Result{
					Rule:     RuleElementRequiredContent,
					Message:  "<" + tag + "> requires child element: <" + required + ">",
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
