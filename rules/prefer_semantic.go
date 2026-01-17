package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// PreferSemantic checks for non-semantic elements used as interactive controls.
type PreferSemantic struct{}

func (r *PreferSemantic) Name() string { return RulePreferSemantic }

func (r *PreferSemantic) Description() string {
	return "prefer semantic elements (button, a) over div/span with click handlers"
}

func (r *PreferSemantic) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		// Check div and span elements
		if !n.IsElement("div") && !n.IsElement("span") {
			return true
		}

		// Check for interactive attributes that suggest it should be a button
		hasClickHandler := n.HasAttr("onclick") ||
			n.HasAttr("onkeydown") ||
			n.HasAttr("onkeyup") ||
			n.HasAttr("onkeypress")

		// Check for HTMX click attributes
		hasHTMXClick := n.HasAttr("hx-get") ||
			n.HasAttr("hx-post") ||
			n.HasAttr("hx-put") ||
			n.HasAttr("hx-delete") ||
			n.HasAttr("hx-patch")

		// Check for role="button" which indicates it should just be a button
		hasButtonRole := n.GetAttr("role") == "button"

		// Check for tabindex which suggests interactive intent
		hasTabindex := n.HasAttr("tabindex")

		if hasClickHandler || (hasHTMXClick && hasButtonRole) {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  n.Data + " with click handler should be a <button> element",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		} else if hasButtonRole && hasTabindex {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  n.Data + " with role=\"button\" should be a <button> element",
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
