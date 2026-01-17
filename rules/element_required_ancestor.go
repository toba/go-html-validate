package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ElementRequiredAncestor checks that elements have their required ancestors.
type ElementRequiredAncestor struct{}

// Name returns the rule identifier.
func (r *ElementRequiredAncestor) Name() string { return RuleElementRequiredAncestor }

// Description returns what this rule checks.
func (r *ElementRequiredAncestor) Description() string {
	return "elements must have required ancestor elements"
}

// Check examines the document for elements missing required ancestors.
func (r *ElementRequiredAncestor) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Check if element has required ancestors
		requiredAncestors, hasRequirement := RequiredAncestors[tag]
		if !hasRequirement {
			return true
		}

		// Check if any required ancestor is present
		if HasAncestor(n, requiredAncestors...) {
			return true
		}

		// For template fragments, skip errors for top-level orphaned elements.
		// These are meant to be included into parent templates that provide ancestors.
		if doc.IsTemplateFragment && isTopLevel(n) {
			return true
		}

		// Missing required ancestor
		ancestorList := strings.Join(requiredAncestors, ", ")
		results = append(results, Result{
			Rule:     RuleElementRequiredAncestor,
			Message:  "<" + tag + "> requires ancestor: " + ancestorList,
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		})

		return true
	})

	return results
}

// isTopLevel checks if a node is at or near the top level of the document.
// For template fragments, top-level elements may lack required ancestors
// because they're meant to be included in parent templates.
func isTopLevel(n *parser.Node) bool {
	// Check parent chain depth - if within 2-3 levels of root, consider top-level
	depth := 0
	for p := n.Parent; p != nil; p = p.Parent {
		depth++
		if depth > 3 {
			return false
		}
	}
	return true
}
