package rules

import (
	"fmt"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// PreferButton checks for input elements that should be buttons.
type PreferButton struct{}

func (r *PreferButton) Name() string { return RulePreferButton }

func (r *PreferButton) Description() string {
	return "prefer <button> over <input type=\"button|submit|reset\">"
}

// buttonInputTypes are input types that should use <button> instead.
var buttonInputTypes = map[string]bool{
	"button": true,
	"submit": true,
	"reset":  true,
	"image":  true,
}

func (r *PreferButton) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if strings.ToLower(n.Data) != "input" {
			return true
		}

		inputType := strings.ToLower(n.GetAttr("type"))
		if buttonInputTypes[inputType] {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  fmt.Sprintf("prefer <button> over <input type=%q>", inputType),
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
