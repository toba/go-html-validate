package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// InputAttributes checks that input elements only have type-appropriate attributes.
type InputAttributes struct{}

// Name returns the rule identifier.
func (r *InputAttributes) Name() string { return RuleInputAttributes }

// Description returns what this rule checks.
func (r *InputAttributes) Description() string {
	return "input attributes must be appropriate for input type"
}

// Check examines the document for invalid input attribute combinations.
func (r *InputAttributes) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if strings.ToLower(n.Data) != "input" {
			return true
		}

		// Get input type (defaults to "text")
		inputType := strings.ToLower(n.GetAttr("type"))
		if inputType == "" {
			inputType = "text"
		}

		// Get allowed attributes for this type
		allowedAttrs, hasSpec := InputTypeAttributes[inputType]
		if !hasSpec {
			// Unknown input type, skip (attribute-allowed-values handles this)
			return true
		}

		// Check each attribute
		for _, attr := range n.Attr {
			attrName := strings.ToLower(attr.Key)

			// Skip common attributes valid on all inputs
			if isCommonInputAttr(attrName) {
				continue
			}

			// Skip global attributes, data-*, aria-*, etc.
			if isGlobalAttribute(attrName) || strings.HasPrefix(attrName, "data-") || strings.HasPrefix(attrName, "aria-") {
				continue
			}

			// Skip event handlers
			if strings.HasPrefix(attrName, "on") {
				continue
			}

			// Check if attribute is valid for this input type
			if !allowedAttrs[attrName] {
				results = append(results, Result{
					Rule:     RuleInputAttributes,
					Message:  "attribute '" + attrName + "' not valid for input type=\"" + inputType + "\"",
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

// isCommonInputAttr returns true for attributes valid on all input types.
func isCommonInputAttr(attr string) bool {
	commonAttrs := map[string]bool{
		"type":      true,
		"name":      true,
		"value":     true,
		"disabled":  true,
		"form":      true,
		"autofocus": true,
	}
	return commonAttrs[attr]
}
