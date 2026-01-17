package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// AllowedLinks checks that link hrefs are valid.
type AllowedLinks struct{}

// Name returns the rule identifier.
func (r *AllowedLinks) Name() string { return RuleAllowedLinks }

// Description returns what this rule checks.
func (r *AllowedLinks) Description() string {
	return "links must have valid href values"
}

// Check examines the document for problematic link hrefs.
func (r *AllowedLinks) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)
		if tag != "a" && tag != "area" {
			return true
		}

		href := n.GetAttr("href")

		// Skip template expressions
		if IsTemplateExpr(href) {
			return true
		}

		// Check for javascript: protocol (security risk)
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(href)), "javascript:") {
			results = append(results, Result{
				Rule:     RuleAllowedLinks,
				Message:  "javascript: URLs are not allowed",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return true
		}

		// Check for vbscript: protocol (security risk)
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(href)), "vbscript:") {
			results = append(results, Result{
				Rule:     RuleAllowedLinks,
				Message:  "vbscript: URLs are not allowed",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return true
		}

		// Check for data: URLs in links (potential security risk)
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(href)), "data:") {
			results = append(results, Result{
				Rule:     RuleAllowedLinks,
				Message:  "data: URLs in links are discouraged",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
			return true
		}

		// Check for empty href (common mistake)
		if href == "" && n.HasAttr("href") {
			results = append(results, Result{
				Rule:     RuleAllowedLinks,
				Message:  "empty href is often a mistake; use # for placeholder links",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Info,
			})
			return true
		}

		return true
	})

	return results
}
