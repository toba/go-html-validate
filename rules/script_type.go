package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ScriptType checks that script elements have valid type attributes.
type ScriptType struct{}

// Name returns the rule identifier.
func (r *ScriptType) Name() string { return RuleScriptType }

// Description returns what this rule checks.
func (r *ScriptType) Description() string {
	return "script type attribute must have a valid value"
}

// Check examines the document for script elements with invalid type.
func (r *ScriptType) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("script") {
			return true
		}

		// Type is optional (defaults to JavaScript)
		if !n.HasAttr("type") {
			return true
		}

		scriptType := strings.ToLower(strings.TrimSpace(n.GetAttr("type")))

		// Skip template values
		if scriptType == "tmpl" {
			return true
		}

		// Check if it's a valid type
		if !ValidScriptTypes[scriptType] {
			results = append(results, Result{
				Rule:     RuleScriptType,
				Message:  "invalid script type: " + n.GetAttr("type"),
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}

		return true
	})

	return results
}
