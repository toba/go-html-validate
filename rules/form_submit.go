package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// FormSubmit checks that forms have submit buttons (WCAG H32).
type FormSubmit struct{}

func (r *FormSubmit) Name() string { return RuleFormSubmit }

func (r *FormSubmit) Description() string {
	return "forms must have a submit button (WCAG H32)"
}

func (r *FormSubmit) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("form") {
			return true
		}

		// Skip forms that use HTMX for submission (they submit via JavaScript)
		if isHTMXForm(n) {
			return true
		}

		// Check if form has a submit button
		if !HasDescendant(n, isSubmitButton) {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "form has no submit button",
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

// isHTMXForm checks if a form uses HTMX for submission via hx-* attributes.
func isHTMXForm(n *parser.Node) bool {
	for _, attr := range n.Attr {
		if len(attr.Key) > 3 && attr.Key[:3] == "hx-" {
			return true
		}
	}
	// Also check if any descendant element triggers form submission via HTMX
	return HasDescendant(n, func(child *parser.Node) bool {
		if child.Type != html.ElementNode {
			return false
		}
		hasHTMXRequest := false
		hasHTMXTrigger := false
		for _, attr := range child.Attr {
			// Check for hx-post, hx-put, hx-delete, hx-patch, hx-get
			if attr.Key == "hx-post" || attr.Key == "hx-put" ||
				attr.Key == "hx-delete" || attr.Key == "hx-patch" ||
				attr.Key == "hx-get" {
				hasHTMXRequest = true
			}
			// Check for hx-trigger which indicates user-initiated action
			if attr.Key == "hx-trigger" {
				hasHTMXTrigger = true
			}
		}
		// hx-get with hx-trigger indicates form-like behavior
		return hasHTMXRequest && hasHTMXTrigger
	})
}

// isSubmitButton checks if a node is a form submit mechanism.
func isSubmitButton(n *parser.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	// <button> without type or type="submit"
	if n.IsElement("button") {
		btnType := n.GetAttr("type")
		return btnType == "" || btnType == "submit"
	}
	// <input type="submit"> or <input type="image">
	if n.IsElement("input") {
		inputType := n.GetAttr("type")
		return inputType == "submit" || inputType == "image"
	}
	return false
}
