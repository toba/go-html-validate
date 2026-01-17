package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// RequireLang ensures <html> elements have a lang attribute.
type RequireLang struct{}

func (r *RequireLang) Name() string { return RuleRequireLang }

func (r *RequireLang) Description() string {
	return "<html> element must have a lang attribute"
}

func (r *RequireLang) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("html") {
			return true
		}

		if !n.HasAttr("lang") {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "<html> element must have a lang attribute; add lang=\"en\" for English content",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		} else if n.GetAttr("lang") == "" {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "lang attribute must not be empty; use BCP 47 code like \"en\" or \"en-US\"",
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
