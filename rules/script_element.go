package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ScriptElement checks script element constraints.
type ScriptElement struct{}

// Name returns the rule identifier.
func (r *ScriptElement) Name() string { return RuleScriptElement }

// Description returns what this rule checks.
func (r *ScriptElement) Description() string {
	return "script elements must follow HTML5 constraints"
}

// Check examines the document for script element issues.
func (r *ScriptElement) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if strings.ToLower(n.Data) != "script" {
			return true
		}

		hasSrc := n.HasAttr("src")
		hasAsync := n.HasAttr("async")
		hasDefer := n.HasAttr("defer")
		hasNomodule := n.HasAttr("nomodule")
		scriptType := strings.ToLower(n.GetAttr("type"))

		// async and defer are mutually exclusive for classic scripts
		// For module scripts, async is allowed but defer is ignored
		if hasAsync && hasDefer && scriptType != "module" {
			results = append(results, Result{
				Rule:     RuleScriptElement,
				Message:  "script should not have both async and defer",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}

		// async and defer only make sense with src
		if !hasSrc {
			if hasAsync {
				results = append(results, Result{
					Rule:     RuleScriptElement,
					Message:  "async attribute requires src attribute",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
			if hasDefer {
				results = append(results, Result{
					Rule:     RuleScriptElement,
					Message:  "defer attribute requires src attribute",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		}

		// nomodule only makes sense for classic scripts
		if hasNomodule && scriptType == "module" {
			results = append(results, Result{
				Rule:     RuleScriptElement,
				Message:  "nomodule attribute should not be on module scripts",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}

		// Script with src should not have inline content
		if hasSrc {
			// Check for non-whitespace content
			for _, child := range n.Children {
				if child.Type == html.TextNode && strings.TrimSpace(child.Data) != "" {
					results = append(results, Result{
						Rule:     RuleScriptElement,
						Message:  "script with src should not have inline content",
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Warning,
					})
					break
				}
			}
		}

		// Check type value if present
		if scriptType != "" && !ValidScriptTypes[scriptType] {
			// Unknown type acts as data block (valid but unusual)
			// Only warn for non-MIME type values that might be mistakes
			if !strings.Contains(scriptType, "/") {
				results = append(results, Result{
					Rule:     RuleScriptElement,
					Message:  "unknown script type: " + scriptType,
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Info,
				})
			}
		}

		return true
	})

	return results
}
