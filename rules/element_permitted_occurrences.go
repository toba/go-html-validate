package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ElementPermittedOccurrences checks cardinality constraints (e.g., one title per head).
type ElementPermittedOccurrences struct{}

// Name returns the rule identifier.
func (r *ElementPermittedOccurrences) Name() string { return RuleElementPermittedOccurrences }

// Description returns what this rule checks.
func (r *ElementPermittedOccurrences) Description() string {
	return "elements must not exceed permitted occurrences"
}

// Check examines the document for elements appearing more than allowed.
func (r *ElementPermittedOccurrences) Check(doc *parser.Document) []Result {
	var results []Result

	// Track occurrences within each context
	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Check if this element type has uniqueness constraints
		contextTag, hasConstraint := UniqueElements[tag]
		if !hasConstraint {
			return true
		}

		// Find the context element
		context := AncestorWithTag(n, contextTag)
		if context == nil {
			return true
		}

		// Count occurrences of this tag within the context
		count := 0
		var firstOccurrence *parser.Node

		// Walk through context's children to find all occurrences
		var countOccurrences func(node *parser.Node)
		countOccurrences = func(node *parser.Node) {
			for _, child := range node.Children {
				if child.Type == html.ElementNode {
					if strings.ToLower(child.Data) == tag {
						count++
						if firstOccurrence == nil {
							firstOccurrence = child
						}
					}
					countOccurrences(child)
				}
			}
		}
		countOccurrences(context)

		// Report error only on second and subsequent occurrences
		if count > 1 && n != firstOccurrence {
			results = append(results, Result{
				Rule:     RuleElementPermittedOccurrences,
				Message:  "<" + tag + "> must only appear once within <" + contextTag + ">",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		}

		return true
	})

	return results
}
