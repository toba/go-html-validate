package rules

import (
	"fmt"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// DuplicateID checks that id attributes are unique within a document.
type DuplicateID struct{}

func (r *DuplicateID) Name() string { return RuleDuplicateID }

func (r *DuplicateID) Description() string {
	return "id attributes must be unique within a document"
}

type idLocation struct {
	line int
	col  int
}

func (r *DuplicateID) Check(doc *parser.Document) []Result {
	var results []Result
	seenIDs := make(map[string]idLocation)

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		id := n.GetAttr("id")
		if id == "" {
			return true
		}

		// Skip template placeholders
		if id == TemplateExprPlaceholder {
			return true
		}

		if first, exists := seenIDs[id]; exists {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  fmt.Sprintf("duplicate id %q (first defined at line %d)", id, first.line),
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		} else {
			seenIDs[id] = idLocation{line: n.Line, col: n.Col}
		}

		return true
	})

	return results
}
