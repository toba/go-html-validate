package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// MetaRefresh checks for auto-refresh meta tags.
type MetaRefresh struct{}

func (r *MetaRefresh) Name() string { return RuleMetaRefresh }

func (r *MetaRefresh) Description() string {
	return "meta refresh should not be used for auto-redirect (WCAG)"
}

func (r *MetaRefresh) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("meta") {
			return true
		}

		httpEquiv := strings.ToLower(n.GetAttr("http-equiv"))
		if httpEquiv != "refresh" {
			return true
		}

		content := n.GetAttr("content")
		if content == "" {
			return true
		}

		// Any refresh is problematic for accessibility
		// Immediate redirects (0 seconds) are less bad but still flagged
		results = append(results, Result{
			Rule:     r.Name(),
			Message:  "meta refresh causes automatic page change, disorienting users",
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		})

		return true
	})

	return results
}
