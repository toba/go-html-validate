package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// RequireCSPNonce checks that inline scripts and styles have CSP nonce.
type RequireCSPNonce struct{}

// Name returns the rule identifier.
func (r *RequireCSPNonce) Name() string { return RuleRequireCSPNonce }

// Description returns what this rule checks.
func (r *RequireCSPNonce) Description() string {
	return "inline scripts and styles should have CSP nonce attribute"
}

// Check examines the document for inline scripts/styles without nonce.
func (r *RequireCSPNonce) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		switch tag {
		case "script":
			// External scripts with src don't need nonce (unless they also have inline content)
			if n.HasAttr("src") && !hasInlineContent(n) {
				return true
			}
			// Scripts with type that's not JavaScript don't need nonce
			scriptType := strings.ToLower(n.GetAttr("type"))
			if scriptType != "" && !isJavaScriptType(scriptType) {
				return true
			}
			// Check for nonce attribute
			if !n.HasAttr("nonce") {
				results = append(results, Result{
					Rule:     RuleRequireCSPNonce,
					Message:  "inline script should have nonce attribute for CSP",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Info, // Info level since not all sites use CSP
				})
			}

		case "style":
			// Check for nonce attribute
			if !n.HasAttr("nonce") {
				results = append(results, Result{
					Rule:     RuleRequireCSPNonce,
					Message:  "inline style should have nonce attribute for CSP",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Info, // Info level since not all sites use CSP
				})
			}
		}

		return true
	})

	return results
}

// hasInlineContent checks if element has non-whitespace text content.
func hasInlineContent(n *parser.Node) bool {
	for _, child := range n.Children {
		if child.Type == html.TextNode && strings.TrimSpace(child.Data) != "" {
			return true
		}
	}
	return false
}

// isJavaScriptType checks if script type is JavaScript.
func isJavaScriptType(scriptType string) bool {
	jsTypes := map[string]bool{
		"":                       true,
		"text/javascript":        true,
		"application/javascript": true,
		"text/ecmascript":        true,
		"application/ecmascript": true,
		"module":                 true,
	}
	return jsTypes[scriptType]
}
