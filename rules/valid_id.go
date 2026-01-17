package rules

import (
	"strings"
	"unicode"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ValidID ensures ID attributes are well-formed.
type ValidID struct{}

func (r *ValidID) Name() string { return RuleValidID }

func (r *ValidID) Description() string {
	return "ID attributes must be non-empty and not contain whitespace"
}

func (r *ValidID) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if !n.HasAttr("id") {
			return true
		}

		id := n.GetAttr("id")

		// Check for empty ID
		if id == "" {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "id attribute must not be empty; provide a unique identifier or remove the attribute",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return true
		}

		// Check for whitespace
		if strings.ContainsAny(id, " \t\n\r") {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "id \"" + id + "\" contains whitespace; use hyphens or underscores instead of spaces",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return true
		}

		// Warn if starts with a digit (CSS selector issues)
		if id != "" && unicode.IsDigit(rune(id[0])) {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "id \"" + id + "\" starts with digit; prefix with letter to avoid CSS selector issues",
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
