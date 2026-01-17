package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// InputAttributesConfig holds htmx configuration for the InputAttributes rule.
type InputAttributesConfig struct {
	// HTMXEnabled allows htmx attributes on input elements.
	HTMXEnabled bool
	// HTMXVersion specifies which htmx version to validate ("2" or "4").
	HTMXVersion string
}

// InputAttributes checks that input elements only have type-appropriate attributes.
type InputAttributes struct {
	config InputAttributesConfig
}

// Configure sets the htmx configuration for this rule.
func (r *InputAttributes) Configure(htmxEnabled bool, htmxVersion string) {
	r.config.HTMXEnabled = htmxEnabled
	r.config.HTMXVersion = htmxVersion
}

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

			// Handle htmx attributes
			if IsHTMXAttribute(attrName) {
				if !r.config.HTMXEnabled {
					results = append(results, Result{
						Rule:     RuleInputAttributes,
						Message:  "htmx attribute '" + attrName + "' used but htmx not enabled in config",
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Warning,
					})
					continue
				}

				// Validate against htmx version
				version := r.config.HTMXVersion
				if version == "" {
					version = "2"
				}
				valid, deprecated, v4Only := ValidateHTMXAttribute(attrName, version)

				if !valid {
					if v4Only {
						results = append(results, Result{
							Rule:     RuleInputAttributes,
							Message:  "htmx attribute '" + attrName + "' is only available in htmx 4",
							Filename: doc.Filename,
							Line:     n.Line,
							Col:      n.Col,
							Severity: Warning,
						})
					} else {
						results = append(results, Result{
							Rule:     RuleInputAttributes,
							Message:  "unknown htmx attribute '" + attrName + "'",
							Filename: doc.Filename,
							Line:     n.Line,
							Col:      n.Col,
							Severity: Warning,
						})
					}
				} else if deprecated {
					results = append(results, Result{
						Rule:     RuleInputAttributes,
						Message:  "htmx attribute '" + attrName + "' is deprecated in htmx 4",
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Warning,
					})
				}
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
