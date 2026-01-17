package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// AriaHiddenBody checks that aria-hidden is not set on the body element.
type AriaHiddenBody struct{}

func (r *AriaHiddenBody) Name() string { return RuleAriaHiddenBody }

func (r *AriaHiddenBody) Description() string {
	return "aria-hidden must not be set on body element"
}

func (r *AriaHiddenBody) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("body") {
			return true
		}

		if n.GetAttr("aria-hidden") == "true" {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "aria-hidden on body hides entire page from assistive technology",
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
