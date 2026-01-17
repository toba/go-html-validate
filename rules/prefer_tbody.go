package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// PreferTbody suggests using explicit tbody in tables.
type PreferTbody struct{}

// Name returns the rule identifier.
func (r *PreferTbody) Name() string { return RulePreferTbody }

// Description returns what this rule checks.
func (r *PreferTbody) Description() string {
	return "tables should use explicit <tbody> element"
}

// Check examines the document for tables without explicit tbody.
func (r *PreferTbody) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if strings.ToLower(n.Data) != "table" {
			return true
		}

		// Check if table has tbody
		hasTbody := false
		hasTr := false

		for _, child := range n.Children {
			if child.Type != html.ElementNode {
				continue
			}
			childTag := strings.ToLower(child.Data)
			if childTag == "tbody" {
				hasTbody = true
			}
			if childTag == "tr" {
				hasTr = true
			}
		}

		// Only report if table has direct tr children without tbody
		if hasTr && !hasTbody {
			results = append(results, Result{
				Rule:     RulePreferTbody,
				Message:  "table should use explicit <tbody> element",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Info,
			})
		}

		return true
	})

	return results
}
