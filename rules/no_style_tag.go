package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoStyleTag discourages use of inline <style> tags.
type NoStyleTag struct{}

// Name returns the rule identifier.
func (r *NoStyleTag) Name() string { return RuleNoStyleTag }

// Description returns what this rule checks.
func (r *NoStyleTag) Description() string {
	return "inline <style> tags should be avoided; use external stylesheets"
}

// Check examines the document for <style> tags.
func (r *NoStyleTag) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if strings.ToLower(n.Data) == "style" {
			results = append(results, Result{
				Rule:     RuleNoStyleTag,
				Message:  "inline <style> tags are discouraged; use external stylesheets",
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
