package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ImgAlt checks that all img elements have alt attributes.
type ImgAlt struct{}

func (r *ImgAlt) Name() string { return RuleImgAlt }

func (r *ImgAlt) Description() string {
	return "images must have alt attribute for accessibility"
}

func (r *ImgAlt) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type == html.ElementNode && n.IsElement("img") {
			if !n.HasAttr("alt") {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  "img element missing alt attribute",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		}
		return true
	})

	return results
}
