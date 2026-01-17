package rules

import (
	"fmt"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// HeadingLevel checks that heading levels don't skip (e.g., h1 to h3).
type HeadingLevel struct{}

func (r *HeadingLevel) Name() string { return RuleHeadingLevel }

func (r *HeadingLevel) Description() string {
	return "heading levels must not skip (h1 followed by h3 is invalid)"
}

func (r *HeadingLevel) Check(doc *parser.Document) []Result {
	var results []Result
	lastRank := 0

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		rank := HeadingRank(n.Data)
		if rank == 0 {
			return true
		}

		// Check for skipped levels
		// First heading can be any level (templates may be partials)
		// But subsequent headings must not skip more than one level
		if lastRank > 0 && rank > lastRank+1 {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  fmt.Sprintf("heading level skipped from h%d to h%d", lastRank, rank),
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}

		lastRank = rank
		return true
	})

	return results
}
