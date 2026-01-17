package rules

import (
	"slices"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ElementPermittedParent checks that elements have valid parent elements.
type ElementPermittedParent struct{}

// Name returns the rule identifier.
func (r *ElementPermittedParent) Name() string { return RuleElementPermittedParent }

// Description returns what this rule checks.
func (r *ElementPermittedParent) Description() string {
	return "elements must have permitted parent elements"
}

// Check examines the document for elements with invalid parents.
func (r *ElementPermittedParent) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Get element spec
		spec, hasSpec := ElementSpecs[tag]
		if !hasSpec || len(spec.PermittedParents) == 0 {
			return true
		}

		// Skip if no parent (document root)
		if n.Parent == nil || n.Parent.Type != html.ElementNode {
			return true
		}

		parentTag := strings.ToLower(n.Parent.Data)

		// Check if parent is permitted
		if slices.Contains(spec.PermittedParents, parentTag) {
			return true
		}

		// Parent not permitted
		parentList := strings.Join(spec.PermittedParents, ", ")
		results = append(results, Result{
			Rule:     RuleElementPermittedParent,
			Message:  "<" + tag + "> must be child of " + parentList + ", not <" + parentTag + ">",
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		})

		return true
	})

	return results
}
