package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// Rule name constants for doctype rules.
const (
	RuleDoctypeHTML    = "doctype-html"
	RuleMissingDoctype = "missing-doctype"
)

// DoctypeHTML checks that DOCTYPE is the HTML5 doctype.
type DoctypeHTML struct{}

// Name returns the rule identifier.
func (r *DoctypeHTML) Name() string { return RuleDoctypeHTML }

// Description returns what this rule checks.
func (r *DoctypeHTML) Description() string {
	return "DOCTYPE must be html (HTML5)"
}

// Check examines the document for non-HTML5 doctypes.
func (r *DoctypeHTML) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.DoctypeNode {
			return true
		}

		// HTML5 doctype should be just "html" with no public/system identifiers
		if strings.ToLower(n.Data) != "html" {
			results = append(results, Result{
				Rule:     RuleDoctypeHTML,
				Message:  "DOCTYPE should be html (HTML5)",
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

// MissingDoctype checks that a DOCTYPE declaration is present.
type MissingDoctype struct{}

// Name returns the rule identifier.
func (r *MissingDoctype) Name() string { return RuleMissingDoctype }

// Description returns what this rule checks.
func (r *MissingDoctype) Description() string {
	return "document must have DOCTYPE declaration"
}

// Check examines the document for missing DOCTYPE.
func (r *MissingDoctype) Check(doc *parser.Document) []Result {
	// Only check full documents (not fragments)
	// A full document typically has <html> as root
	hasHTML := false
	hasDoctype := false

	doc.Walk(func(n *parser.Node) bool {
		if n.Type == html.DoctypeNode {
			hasDoctype = true
		}
		if n.Type == html.ElementNode && strings.ToLower(n.Data) == "html" {
			hasHTML = true
		}
		return true
	})

	// Only report if this looks like a full document (has <html>)
	if hasHTML && !hasDoctype {
		return []Result{{
			Rule:     RuleMissingDoctype,
			Message:  "document is missing DOCTYPE declaration",
			Filename: doc.Filename,
			Line:     1,
			Col:      1,
			Severity: Warning,
		}}
	}

	return nil
}
