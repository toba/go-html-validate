package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoInlineStyle checks for inline style attributes.
type NoInlineStyle struct{}

func (r *NoInlineStyle) Name() string { return RuleNoInlineStyle }

func (r *NoInlineStyle) Description() string {
	return "avoid inline styles; use classes with separate stylesheets"
}

func (r *NoInlineStyle) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if n.GetAttr("style") != "" {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "avoid inline style attribute; use CSS classes instead",
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
