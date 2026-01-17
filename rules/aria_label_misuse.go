package rules

import (
	"fmt"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// AriaLabelMisuse checks for aria-label on elements that don't support it.
type AriaLabelMisuse struct{}

func (r *AriaLabelMisuse) Name() string { return RuleAriaLabelMisuse }

func (r *AriaLabelMisuse) Description() string {
	return "aria-label/aria-labelledby only allowed on labelable elements"
}

func (r *AriaLabelMisuse) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		ariaLabel := n.GetAttr("aria-label")
		ariaLabelledby := n.GetAttr("aria-labelledby")

		if ariaLabel == "" && ariaLabelledby == "" {
			return true
		}

		tagName := strings.ToLower(n.Data)

		// Check if element allows aria-label
		if AriaLabelableElements[tagName] {
			return true
		}

		// Elements with role attribute can use aria-label
		if n.GetAttr("role") != "" {
			return true
		}

		// Elements with tabindex can use aria-label (they're interactive)
		if n.GetAttr("tabindex") != "" {
			return true
		}

		attr := "aria-label"
		if ariaLabelledby != "" {
			attr = "aria-labelledby"
		}

		results = append(results, Result{
			Rule:     r.Name(),
			Message:  fmt.Sprintf("%s not allowed on <%s> without role attribute", attr, tagName),
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		})

		return true
	})

	return results
}
