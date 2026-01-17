package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoRedundantFor checks for label elements with for attribute when they wrap the control.
type NoRedundantFor struct{}

// Name returns the rule identifier.
func (r *NoRedundantFor) Name() string { return RuleNoRedundantFor }

// Description returns what this rule checks.
func (r *NoRedundantFor) Description() string {
	return "label for attribute is redundant when label wraps the control"
}

// Check examines the document for redundant for attributes on labels.
func (r *NoRedundantFor) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("label") {
			return true
		}

		forAttr := n.GetAttr("for")
		if forAttr == "" || forAttr == TemplateExprPlaceholder {
			return true
		}

		// Check if label contains a labelable element with matching id
		hasMatchingChild := false
		var checkChildren func(node *parser.Node)
		checkChildren = func(node *parser.Node) {
			for _, child := range node.Children {
				if child.Type == html.ElementNode {
					// Check if this is a labelable element with matching id
					if LabelableElements[child.Data] {
						childID := child.GetAttr("id")
						if childID == forAttr {
							hasMatchingChild = true
							return
						}
					}
					// Check nested children
					checkChildren(child)
				}
			}
		}
		checkChildren(n)

		if hasMatchingChild {
			results = append(results, Result{
				Rule:     RuleNoRedundantFor,
				Message:  "label for=\"" + forAttr + "\" is redundant when label wraps the control",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Info,
			})
		}

		return true
	})

	return results
}
