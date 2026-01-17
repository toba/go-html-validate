package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// MapDupName checks for duplicate area names within a map element.
type MapDupName struct{}

// Name returns the rule identifier.
func (r *MapDupName) Name() string { return RuleMapDupName }

// Description returns what this rule checks.
func (r *MapDupName) Description() string {
	return "area elements within a map should have unique names"
}

// Check examines the document for duplicate area names within maps.
func (r *MapDupName) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("map") {
			return true
		}

		// Collect area names
		names := make(map[string]*parser.Node)
		for _, child := range n.Children {
			if child.Type == html.ElementNode && child.IsElement("area") {
				name := child.GetAttr("name")
				if name == "" || name == TemplateExprPlaceholder {
					continue
				}
				if first, exists := names[name]; exists {
					results = append(results, Result{
						Rule:     RuleMapDupName,
						Message:  "duplicate area name: " + name,
						Filename: doc.Filename,
						Line:     child.Line,
						Col:      child.Col,
						Severity: Warning,
					})
					_ = first // First occurrence tracked but not reported
				} else {
					names[name] = child
				}
			}
		}

		return true
	})

	return results
}
