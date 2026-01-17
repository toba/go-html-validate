package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ValidAutocomplete checks that autocomplete attributes have valid values.
type ValidAutocomplete struct{}

// Name returns the rule identifier.
func (r *ValidAutocomplete) Name() string { return RuleValidAutocomplete }

// Description returns what this rule checks.
func (r *ValidAutocomplete) Description() string {
	return "autocomplete attribute must have valid token values"
}

// Check examines the document for invalid autocomplete values.
func (r *ValidAutocomplete) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		// Check elements that support autocomplete
		tag := strings.ToLower(n.Data)
		switch tag {
		case "input", "select", "textarea", "form":
			// These support autocomplete
		default:
			return true
		}

		autocomplete := n.GetAttr("autocomplete")
		if autocomplete == "" || autocomplete == TemplateExprPlaceholder {
			return true
		}

		// Parse autocomplete tokens
		for token := range strings.FieldsSeq(autocomplete) {
			token = strings.ToLower(token)

			// Skip template placeholders
			if token == "tmpl" {
				continue
			}

			// Check for section- prefix
			if strings.HasPrefix(token, "section-") {
				// Section tokens are valid as long as they have content after prefix
				if len(token) > len("section-") {
					continue
				}
			}

			// Check for billing/shipping prefix
			if token == "shipping" || token == "billing" {
				continue
			}

			// Check if it's a valid autocomplete token
			if !AutocompleteTokens[token] {
				results = append(results, Result{
					Rule:     RuleValidAutocomplete,
					Message:  "invalid autocomplete token: " + token,
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Warning,
				})
			}
		}

		return true
	})

	return results
}
