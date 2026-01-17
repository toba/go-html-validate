package rules

import (
	"strconv"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// TabindexNoPositive checks that tabindex values are not positive.
type TabindexNoPositive struct{}

func (r *TabindexNoPositive) Name() string { return RuleTabindexNoPositive }

func (r *TabindexNoPositive) Description() string {
	return "tabindex should be 0 or -1, not positive (breaks natural tab order)"
}

func (r *TabindexNoPositive) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tabindex := n.GetAttr("tabindex")
		if tabindex == "" {
			return true
		}

		// Parse tabindex value
		val, err := strconv.Atoi(tabindex)
		if err != nil {
			// Non-numeric tabindex, skip (could be template variable)
			return true
		}

		if val > 0 {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "positive tabindex disrupts natural tab order; use 0 or -1",
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
