package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// FormDupName checks for duplicate name attributes within a form.
// Radio buttons and checkboxes with the same name are allowed.
type FormDupName struct{}

// Name returns the rule identifier.
func (r *FormDupName) Name() string { return RuleFormDupName }

// Description returns what this rule checks.
func (r *FormDupName) Description() string {
	return "form controls should have unique names (except radio/checkbox groups)"
}

// Check examines the document for duplicate names within forms.
func (r *FormDupName) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("form") {
			return true
		}

		// Collect all named controls in this form
		// Map name -> list of (element, type)
		type controlInfo struct {
			node      *parser.Node
			inputType string
		}
		names := make(map[string][]controlInfo)

		var collectNames func(node *parser.Node)
		collectNames = func(node *parser.Node) {
			for _, child := range node.Children {
				if child.Type != html.ElementNode {
					continue
				}

				// Don't recurse into nested forms
				if child.IsElement("form") {
					continue
				}

				// Check for named form controls
				tag := strings.ToLower(child.Data)
				name := child.GetAttr("name")
				if name != "" && name != TemplateExprPlaceholder {
					switch tag {
					case "input", "select", "textarea", "button", "output":
						inputType := strings.ToLower(child.GetAttr("type"))
						names[name] = append(names[name], controlInfo{
							node:      child,
							inputType: inputType,
						})
					}
				}

				collectNames(child)
			}
		}
		collectNames(n)

		// Check for duplicates (ignoring radio/checkbox)
		for name, controls := range names {
			if len(controls) <= 1 {
				continue
			}

			// Check if all are radio or checkbox (which can share names)
			allRadioCheckbox := true
			for _, ctrl := range controls {
				if ctrl.inputType != "radio" && ctrl.inputType != "checkbox" {
					allRadioCheckbox = false
					break
				}
			}

			if !allRadioCheckbox {
				// Report duplicate for non-radio/checkbox controls
				for i := 1; i < len(controls); i++ {
					ctrl := controls[i]
					results = append(results, Result{
						Rule:     RuleFormDupName,
						Message:  "duplicate form control name: " + name,
						Filename: doc.Filename,
						Line:     ctrl.node.Line,
						Col:      ctrl.node.Col,
						Severity: Warning,
					})
				}
			}
		}

		return true
	})

	return results
}
