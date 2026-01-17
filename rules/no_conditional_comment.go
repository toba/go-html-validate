package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoConditionalComment checks for IE conditional comments.
type NoConditionalComment struct{}

// Name returns the rule identifier.
func (r *NoConditionalComment) Name() string { return RuleNoConditionalComment }

// Description returns what this rule checks.
func (r *NoConditionalComment) Description() string {
	return "IE conditional comments should not be used"
}

// Check examines the document for IE conditional comments.
func (r *NoConditionalComment) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		// Check for HTML comments containing IE conditionals
		if n.Type != html.CommentNode {
			return true
		}

		comment := n.Data

		// Check for common IE conditional patterns
		// <!--[if IE]>, <!--[if lt IE 9]>, etc.
		if strings.Contains(comment, "[if ") && strings.Contains(comment, "]>") {
			results = append(results, Result{
				Rule:     RuleNoConditionalComment,
				Message:  "IE conditional comments are deprecated and not supported in modern browsers",
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
