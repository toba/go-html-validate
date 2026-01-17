package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// MultipleLabeledControls checks that labels don't reference multiple controls.
type MultipleLabeledControls struct{}

func (r *MultipleLabeledControls) Name() string { return RuleMultipleLabeledControls }

func (r *MultipleLabeledControls) Description() string {
	return "label element should only be associated with one control"
}

func (r *MultipleLabeledControls) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if strings.ToLower(n.Data) != "label" {
			return true
		}

		controlCount := countLabeledControls(n)
		if controlCount > 1 {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "label is associated with multiple controls",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		}

		return true
	})

	return results
}

// countLabeledControls counts controls associated with a label.
func countLabeledControls(label *parser.Node) int {
	count := 0
	nestedIDs := make(map[string]bool)

	// Count labelable elements nested inside the label and track their IDs
	var countNested func(n *parser.Node)
	countNested = func(n *parser.Node) {
		for _, child := range n.Children {
			if child.Type == html.ElementNode {
				tagName := strings.ToLower(child.Data)
				if LabelableElements[tagName] {
					// Skip hidden inputs
					if tagName == "input" && strings.ToLower(child.GetAttr("type")) == "hidden" {
						countNested(child)
						continue
					}
					count++
					// Track ID for deduplication with "for" attribute
					if id := child.GetAttr("id"); id != "" {
						nestedIDs[id] = true
					}
				}
				countNested(child)
			}
		}
	}
	countNested(label)

	// Check for "for" attribute pointing to another element
	forAttr := label.GetAttr("for")
	if forAttr != "" {
		// Only count if "for" points to an element NOT already nested
		if !nestedIDs[forAttr] {
			count++
		}
	}

	return count
}
