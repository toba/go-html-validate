package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// SVGFocusable checks that SVGs inside interactive elements have focusable="false".
type SVGFocusable struct{}

func (r *SVGFocusable) Name() string { return RuleSVGFocusable }

func (r *SVGFocusable) Description() string {
	return "SVGs inside interactive elements should have focusable=\"false\""
}

func (r *SVGFocusable) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		// Check if this is an interactive element
		if !r.isInteractive(n) {
			return true
		}

		// Look for SVG children
		for _, child := range n.Children {
			r.checkSVG(child, doc.Filename, &results)
		}

		return true
	})

	return results
}

func (r *SVGFocusable) isInteractive(n *parser.Node) bool {
	switch n.Data {
	case "a":
		return n.HasAttr("href")
	case "button":
		return true
	default:
		return false
	}
}

func (r *SVGFocusable) checkSVG(n *parser.Node, filename string, results *[]Result) {
	if n.Type != html.ElementNode {
		return
	}

	if n.IsElement("svg") {
		// Check if focusable is explicitly set to "false"
		focusable := n.GetAttr("focusable")
		if focusable != "false" {
			*results = append(*results, Result{
				Rule:     r.Name(),
				Message:  "SVG inside interactive element should have focusable=\"false\"",
				Filename: filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}
		return
	}

	// Recurse into children
	for _, child := range n.Children {
		r.checkSVG(child, filename, results)
	}
}
