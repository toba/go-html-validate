package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// InputLabel checks that form inputs have associated labels.
type InputLabel struct{}

func (r *InputLabel) Name() string { return RuleInputLabel }

func (r *InputLabel) Description() string {
	return "form inputs must have associated label, aria-label, or aria-labelledby"
}

func (r *InputLabel) Check(doc *parser.Document) []Result {
	var results []Result

	// Collect all label for= values
	labelFor := make(map[string]bool)
	doc.Walk(func(n *parser.Node) bool {
		if n.IsElement("label") {
			if forAttr := n.GetAttr("for"); forAttr != "" {
				labelFor[forAttr] = true
			}
		}
		return true
	})

	// Check inputs
	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}
		if !n.IsElement("input") && !n.IsElement("select") && !n.IsElement("textarea") {
			return true
		}

		// Hidden and submit inputs don't need labels
		if n.IsElement("input") {
			inputType := n.GetAttr("type")
			if inputType == "hidden" || inputType == "submit" || inputType == "button" || inputType == "reset" || inputType == "image" {
				return true
			}
		}

		// Check for accessible label
		hasLabel := n.HasAttr("aria-label") && n.GetAttr("aria-label") != ""

		// Check aria-label

		// Check aria-labelledby
		if n.HasAttr("aria-labelledby") && n.GetAttr("aria-labelledby") != "" {
			hasLabel = true
		}

		// Check for associated label via id
		if id := n.GetAttr("id"); id != "" {
			if labelFor[id] {
				hasLabel = true
			}
		}

		// Check if input is inside a label
		for p := n.Parent; p != nil; p = p.Parent {
			if p.IsElement("label") {
				hasLabel = true
				break
			}
		}

		// Check title attribute as fallback
		if n.HasAttr("title") && n.GetAttr("title") != "" {
			hasLabel = true
		}

		if !hasLabel {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  n.Data + " element missing accessible label",
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
