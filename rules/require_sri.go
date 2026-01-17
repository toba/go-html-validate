package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// RequireSRI checks that external scripts and stylesheets have integrity attributes.
type RequireSRI struct{}

func (r *RequireSRI) Name() string { return RuleRequireSRI }

func (r *RequireSRI) Description() string {
	return "external resources should have subresource integrity (integrity attribute)"
}

func (r *RequireSRI) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tagName := strings.ToLower(n.Data)

		switch tagName {
		case "script":
			src := n.GetAttr("src")
			if src == "" {
				return true // inline script
			}
			if !isExternalURL(src) {
				return true // local resource
			}
			if n.GetAttr("integrity") == "" {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  "external script missing integrity attribute for SRI",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Warning,
				})
			}

		case "link":
			rel := strings.ToLower(n.GetAttr("rel"))
			// Only check stylesheets and preloads
			if rel != "stylesheet" && rel != "preload" && rel != "modulepreload" {
				return true
			}
			href := n.GetAttr("href")
			if href == "" {
				return true
			}
			if !isExternalURL(href) {
				return true // local resource
			}
			if n.GetAttr("integrity") == "" {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  "external stylesheet missing integrity attribute for SRI",
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

// isExternalURL checks if a URL points to an external resource.
func isExternalURL(url string) bool {
	url = strings.TrimSpace(url)
	// Absolute URLs with protocol
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return true
	}
	// Protocol-relative URLs
	if strings.HasPrefix(url, "//") {
		return true
	}
	return false
}
