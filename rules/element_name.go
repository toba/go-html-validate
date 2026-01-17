package rules

import (
	"strings"
	"unicode"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ElementName checks that element names are valid.
type ElementName struct{}

// Name returns the rule identifier.
func (r *ElementName) Name() string { return RuleElementName }

// Description returns what this rule checks.
func (r *ElementName) Description() string {
	return "element names must be valid HTML element names or valid custom element names"
}

// Check examines the document for invalid element names.
func (r *ElementName) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tagName := strings.ToLower(n.Data)

		// Skip template placeholders
		if tagName == "tmpl" || tagName == "" {
			return true
		}

		// Check if it's a known valid element
		if ValidElements[tagName] {
			return true
		}

		// Check if it's a valid custom element
		if IsCustomElement(tagName) {
			return true
		}

		// Check for common typos or invalid characters
		if !isValidElementName(tagName) {
			results = append(results, Result{
				Rule:     RuleElementName,
				Message:  "invalid element name: " + tagName,
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return true
		}

		// Unknown element (not in valid list, not custom, but syntactically valid)
		results = append(results, Result{
			Rule:     RuleElementName,
			Message:  "unknown element: " + tagName,
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Warning,
		})

		return true
	})

	return results
}

// isValidElementName checks if a name follows element naming rules.
func isValidElementName(name string) bool {
	if name == "" {
		return false
	}
	// Must start with ASCII letter
	first := rune(name[0])
	if !unicode.IsLetter(first) {
		return false
	}
	// Rest must be alphanumeric or hyphen
	for _, r := range name[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			return false
		}
	}
	return true
}
