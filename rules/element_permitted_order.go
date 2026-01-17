package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// ElementPermittedOrder checks that elements appear in correct order.
type ElementPermittedOrder struct{}

// Name returns the rule identifier.
func (r *ElementPermittedOrder) Name() string { return RuleElementPermittedOrder }

// Description returns what this rule checks.
func (r *ElementPermittedOrder) Description() string {
	return "elements must appear in correct order"
}

// Check examines the document for incorrectly ordered elements.
func (r *ElementPermittedOrder) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		switch tag {
		case "html":
			results = append(results, r.checkHTMLOrder(n, doc)...)
		case "table":
			results = append(results, r.checkTableOrder(n, doc)...)
		case "details":
			results = append(results, r.checkDetailsOrder(n, doc)...)
		case "fieldset":
			results = append(results, r.checkFieldsetOrder(n, doc)...)
		}

		return true
	})

	return results
}

// checkHTMLOrder verifies head comes before body.
func (r *ElementPermittedOrder) checkHTMLOrder(n *parser.Node, doc *parser.Document) []Result {
	var results []Result
	var headSeen, bodySeen bool
	var bodyNode *parser.Node

	for _, child := range n.Children {
		if child.Type != html.ElementNode {
			continue
		}
		childTag := strings.ToLower(child.Data)

		if childTag == "body" {
			bodySeen = true
			bodyNode = child
		}
		if childTag == "head" {
			headSeen = true
			if bodySeen {
				results = append(results, Result{
					Rule:     RuleElementPermittedOrder,
					Message:  "<head> must come before <body>",
					Filename: doc.Filename,
					Line:     child.Line,
					Col:      child.Col,
					Severity: Error,
				})
			}
		}
	}

	// Also check if body comes before head (report on body)
	if bodySeen && !headSeen {
		// head might come later, which we already caught above
		_ = bodyNode // avoid unused warning
	}

	return results
}

// checkTableOrder verifies caption comes first, then colgroup, then thead/tbody/tfoot.
func (r *ElementPermittedOrder) checkTableOrder(n *parser.Node, doc *parser.Document) []Result {
	var results []Result

	// Track what we've seen
	var captionSeen, colgroupSeen, theadSeen, tbodySeen, trSeen bool

	for _, child := range n.Children {
		if child.Type != html.ElementNode {
			continue
		}
		childTag := strings.ToLower(child.Data)

		switch childTag {
		case "caption":
			captionSeen = true
			if colgroupSeen || theadSeen || tbodySeen || trSeen {
				results = append(results, Result{
					Rule:     RuleElementPermittedOrder,
					Message:  "<caption> must be first child of <table>",
					Filename: doc.Filename,
					Line:     child.Line,
					Col:      child.Col,
					Severity: Error,
				})
			}
		case "colgroup":
			colgroupSeen = true
			if theadSeen || tbodySeen || trSeen {
				results = append(results, Result{
					Rule:     RuleElementPermittedOrder,
					Message:  "<colgroup> must come before <thead>, <tbody>, <tfoot>, and <tr>",
					Filename: doc.Filename,
					Line:     child.Line,
					Col:      child.Col,
					Severity: Error,
				})
			}
		case "thead":
			theadSeen = true
			if tbodySeen || trSeen {
				results = append(results, Result{
					Rule:     RuleElementPermittedOrder,
					Message:  "<thead> must come before <tbody> and <tr>",
					Filename: doc.Filename,
					Line:     child.Line,
					Col:      child.Col,
					Severity: Error,
				})
			}
		case "tbody", "tfoot":
			tbodySeen = true
		case "tr":
			trSeen = true
		}
	}

	_ = captionSeen // avoid unused warning

	return results
}

// checkDetailsOrder verifies summary is first child.
func (r *ElementPermittedOrder) checkDetailsOrder(n *parser.Node, doc *parser.Document) []Result {
	var results []Result

	var otherSeen bool
	for _, child := range n.Children {
		if child.Type != html.ElementNode {
			continue
		}
		childTag := strings.ToLower(child.Data)

		if childTag == "summary" {
			if otherSeen {
				results = append(results, Result{
					Rule:     RuleElementPermittedOrder,
					Message:  "<summary> must be first child of <details>",
					Filename: doc.Filename,
					Line:     child.Line,
					Col:      child.Col,
					Severity: Error,
				})
			}
			break // Only check first summary
		}
		otherSeen = true
	}

	return results
}

// checkFieldsetOrder verifies legend is first child.
func (r *ElementPermittedOrder) checkFieldsetOrder(n *parser.Node, doc *parser.Document) []Result {
	var results []Result

	var otherSeen bool
	for _, child := range n.Children {
		if child.Type != html.ElementNode {
			continue
		}
		childTag := strings.ToLower(child.Data)

		if childTag == "legend" {
			if otherSeen {
				results = append(results, Result{
					Rule:     RuleElementPermittedOrder,
					Message:  "<legend> must be first child of <fieldset>",
					Filename: doc.Filename,
					Line:     child.Line,
					Col:      child.Col,
					Severity: Error,
				})
			}
			break // Only check first legend
		}
		otherSeen = true
	}

	return results
}
