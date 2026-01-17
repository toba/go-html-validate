package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// HiddenFocusable checks that aria-hidden elements don't contain focusable content.
type HiddenFocusable struct{}

func (r *HiddenFocusable) Name() string { return RuleHiddenFocusable }

func (r *HiddenFocusable) Description() string {
	return "focusable elements must not be inside aria-hidden containers"
}

func (r *HiddenFocusable) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		// Check if this element has aria-hidden="true"
		if n.GetAttr("aria-hidden") != "true" {
			return true
		}

		// Find focusable descendants
		r.findFocusable(n, doc.Filename, &results)

		// Don't recurse into children from Walk since we handled them
		return true
	})

	return results
}

// findFocusable recursively finds focusable elements within a node.
func (r *HiddenFocusable) findFocusable(n *parser.Node, filename string, results *[]Result) {
	for _, child := range n.Children {
		if child.Type != html.ElementNode {
			r.findFocusable(child, filename, results)
			continue
		}

		if r.isFocusable(child) {
			*results = append(*results, Result{
				Rule:     r.Name(),
				Message:  child.Data + " is focusable but inside aria-hidden container",
				Filename: filename,
				Line:     child.Line,
				Col:      child.Col,
				Severity: Error,
			})
		}

		r.findFocusable(child, filename, results)
	}
}

// isFocusable checks if an element is natively or explicitly focusable.
func (r *HiddenFocusable) isFocusable(n *parser.Node) bool {
	// Check tabindex
	if n.HasAttr("tabindex") {
		tabindex := n.GetAttr("tabindex")
		// tabindex="-1" is still focusable via script
		if tabindex != "" {
			return true
		}
	}

	// Natively focusable elements
	switch n.Data {
	case "a":
		return n.HasAttr("href")
	case "button":
		return true
	case "input":
		return n.GetAttr("type") != "hidden"
	case "select", "textarea":
		return true
	case "area":
		return n.HasAttr("href")
	}

	// contenteditable
	if n.GetAttr("contenteditable") == "true" {
		return true
	}

	return false
}
