package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoMissingReferences checks that ID references point to existing elements.
type NoMissingReferences struct{}

// Name returns the rule identifier.
func (r *NoMissingReferences) Name() string { return RuleNoMissingReferences }

// Description returns what this rule checks.
func (r *NoMissingReferences) Description() string {
	return "ID references must point to existing elements"
}

// Check examines the document for broken ID references.
func (r *NoMissingReferences) Check(doc *parser.Document) []Result {
	var results []Result

	// First pass: collect all IDs
	ids := make(map[string]bool)
	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}
		if id := n.GetAttr("id"); id != "" {
			ids[id] = true
		}
		return true
	})

	// Second pass: check references
	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		// Check for attribute
		if forID := n.GetAttr("for"); forID != "" {
			if !ids[forID] && !IsTemplateExpr(forID) {
				results = append(results, Result{
					Rule:     RuleNoMissingReferences,
					Message:  "for=\"" + forID + "\" references non-existent id",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		}

		// Check aria-labelledby (space-separated list)
		if labelledby := n.GetAttr("aria-labelledby"); labelledby != "" {
			for id := range strings.FieldsSeq(labelledby) {
				if !ids[id] && !IsTemplateExpr(id) {
					results = append(results, Result{
						Rule:     RuleNoMissingReferences,
						Message:  "aria-labelledby references non-existent id: " + id,
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Error,
					})
				}
			}
		}

		// Check aria-describedby (space-separated list)
		if describedby := n.GetAttr("aria-describedby"); describedby != "" {
			for id := range strings.FieldsSeq(describedby) {
				if !ids[id] && !IsTemplateExpr(id) {
					results = append(results, Result{
						Rule:     RuleNoMissingReferences,
						Message:  "aria-describedby references non-existent id: " + id,
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Error,
					})
				}
			}
		}

		// Check aria-controls (space-separated list)
		if controls := n.GetAttr("aria-controls"); controls != "" {
			for id := range strings.FieldsSeq(controls) {
				if !ids[id] && !IsTemplateExpr(id) {
					results = append(results, Result{
						Rule:     RuleNoMissingReferences,
						Message:  "aria-controls references non-existent id: " + id,
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Warning,
					})
				}
			}
		}

		// Check aria-owns (space-separated list)
		if owns := n.GetAttr("aria-owns"); owns != "" {
			for id := range strings.FieldsSeq(owns) {
				if !ids[id] && !IsTemplateExpr(id) {
					results = append(results, Result{
						Rule:     RuleNoMissingReferences,
						Message:  "aria-owns references non-existent id: " + id,
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Warning,
					})
				}
			}
		}

		// Check list attribute on input
		if list := n.GetAttr("list"); list != "" {
			if !ids[list] && !IsTemplateExpr(list) {
				results = append(results, Result{
					Rule:     RuleNoMissingReferences,
					Message:  "list=\"" + list + "\" references non-existent datalist",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		}

		// Check headers attribute on td/th (space-separated list)
		if headers := n.GetAttr("headers"); headers != "" {
			for id := range strings.FieldsSeq(headers) {
				if !ids[id] && !IsTemplateExpr(id) {
					results = append(results, Result{
						Rule:     RuleNoMissingReferences,
						Message:  "headers references non-existent id: " + id,
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Error,
					})
				}
			}
		}

		// Check usemap attribute (starts with #)
		if usemap := n.GetAttr("usemap"); usemap != "" && strings.HasPrefix(usemap, "#") {
			mapName := usemap[1:] // Remove #
			if !ids[mapName] && !IsTemplateExpr(mapName) {
				// usemap references name attribute, not id, but often they match
				// This is a simplified check
				results = append(results, Result{
					Rule:     RuleNoMissingReferences,
					Message:  "usemap=\"" + usemap + "\" may reference non-existent map",
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
