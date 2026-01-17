package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ElementPermittedContent checks that elements contain only permitted children.
type ElementPermittedContent struct{}

// Name returns the rule identifier.
func (r *ElementPermittedContent) Name() string { return RuleElementPermittedContent }

// Description returns what this rule checks.
func (r *ElementPermittedContent) Description() string {
	return "elements must contain only permitted child elements"
}

// Check examines the document for elements with invalid children.
func (r *ElementPermittedContent) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Get element spec
		spec, hasSpec := ElementSpecs[tag]
		if !hasSpec {
			return true
		}

		// Check for forbidden descendants
		if len(spec.ForbiddenContent) > 0 {
			for _, child := range n.Children {
				if child.Type != html.ElementNode {
					continue
				}
				childTag := strings.ToLower(child.Data)
				for _, forbidden := range spec.ForbiddenContent {
					if childTag == forbidden {
						results = append(results, Result{
							Rule:     RuleElementPermittedContent,
							Message:  "<" + tag + "> must not contain <" + forbidden + ">",
							Filename: doc.Filename,
							Line:     child.Line,
							Col:      child.Col,
							Severity: Error,
						})
					}
				}
			}
		}

		// Check permitted content if specified
		if len(spec.PermittedContent) == 0 {
			return true
		}

		// Build permitted set for fast lookup
		permitted := make(map[string]bool, len(spec.PermittedContent))
		for _, p := range spec.PermittedContent {
			permitted[p] = true
		}

		// Check each child element
		for _, child := range n.Children {
			if child.Type != html.ElementNode {
				continue
			}

			childTag := strings.ToLower(child.Data)

			// Skip custom elements (contain hyphen)
			if IsCustomElement(childTag) {
				continue
			}

			if !permitted[childTag] {
				results = append(results, Result{
					Rule:     RuleElementPermittedContent,
					Message:  "<" + childTag + "> is not permitted as child of <" + tag + ">",
					Filename: doc.Filename,
					Line:     child.Line,
					Col:      child.Col,
					Severity: Error,
				})
			}
		}

		return true
	})

	return results
}
