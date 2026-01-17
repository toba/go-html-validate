package rules

import (
	"fmt"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// MaxTitleLength is the recommended maximum title length for SEO.
const MaxTitleLength = 70

// LongTitle checks that title elements don't exceed recommended length.
type LongTitle struct{}

func (r *LongTitle) Name() string { return RuleLongTitle }

func (r *LongTitle) Description() string {
	return "title element should not exceed 70 characters for SEO"
}

func (r *LongTitle) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		if strings.ToLower(n.Data) != "title" {
			return true
		}

		text := strings.TrimSpace(n.TextContent())
		if len(text) > MaxTitleLength {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  fmt.Sprintf("title text is %d characters, should be at most %d", len(text), MaxTitleLength),
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
