package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoDeprecatedAttr checks for deprecated HTML attributes.
type NoDeprecatedAttr struct{}

// Name returns the rule identifier.
func (r *NoDeprecatedAttr) Name() string { return RuleNoDeprecatedAttr }

// Description returns what this rule checks.
func (r *NoDeprecatedAttr) Description() string {
	return "deprecated HTML attributes should not be used"
}

// Check examines the document for deprecated attributes.
func (r *NoDeprecatedAttr) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		for _, attr := range n.Attr {
			attrName := strings.ToLower(attr.Key)

			// Check element-specific deprecated attributes
			if elemAttrs, ok := DeprecatedAttributes[tag]; ok {
				if suggestion, deprecated := elemAttrs[attrName]; deprecated {
					results = append(results, Result{
						Rule:     RuleNoDeprecatedAttr,
						Message:  "attribute \"" + attrName + "\" on <" + tag + "> is deprecated; " + suggestion,
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Warning,
					})
					continue
				}
			}

			// Check global deprecated attributes
			if globalAttrs, ok := DeprecatedAttributes[""]; ok {
				if suggestion, deprecated := globalAttrs[attrName]; deprecated {
					// Some attributes (width, height) are not deprecated on certain elements
					if (attrName == "width" || attrName == "height") && NonDeprecatedSizeAttrs[tag] {
						continue
					}
					results = append(results, Result{
						Rule:     RuleNoDeprecatedAttr,
						Message:  "attribute \"" + attrName + "\" is deprecated; " + suggestion,
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Warning,
					})
				}
			}
		}

		return true
	})

	return results
}
