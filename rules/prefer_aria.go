package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// PreferAria checks for custom data attributes that should use ARIA equivalents.
type PreferAria struct{}

func (r *PreferAria) Name() string { return RulePreferAria }

func (r *PreferAria) Description() string {
	return "prefer ARIA attributes over custom data-* attributes for accessibility semantics"
}

// ariaEquivalents maps data-* patterns to their ARIA equivalents.
var ariaEquivalents = map[string]string{
	"data-label":        "aria-label",
	"data-description":  "aria-describedby (with referenced element)",
	"data-expanded":     "aria-expanded",
	"data-selected":     "aria-selected",
	"data-checked":      "aria-checked",
	"data-disabled":     "aria-disabled",
	"data-hidden":       "aria-hidden",
	"data-pressed":      "aria-pressed",
	"data-current":      "aria-current",
	"data-busy":         "aria-busy",
	"data-live":         "aria-live",
	"data-role":         "role",
	"data-required":     "aria-required",
	"data-invalid":      "aria-invalid",
	"data-readonly":     "aria-readonly",
	"data-sort":         "aria-sort",
	"data-level":        "aria-level",
	"data-posinset":     "aria-posinset",
	"data-setsize":      "aria-setsize",
	"data-valuenow":     "aria-valuenow",
	"data-valuemin":     "aria-valuemin",
	"data-valuemax":     "aria-valuemax",
	"data-valuetext":    "aria-valuetext",
	"data-orientation":  "aria-orientation",
	"data-autocomplete": "aria-autocomplete",
	"data-haspopup":     "aria-haspopup",
	"data-modal":        "aria-modal",
	"data-multiline":    "aria-multiline",
	"data-placeholder":  "aria-placeholder",
}

func (r *PreferAria) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || n.Node == nil {
			return true
		}

		for _, attr := range n.Attr {
			// Check if this is a data-* attribute with an ARIA equivalent
			if !strings.HasPrefix(attr.Key, "data-") {
				continue
			}

			// Check exact matches
			if ariaAttr, ok := ariaEquivalents[attr.Key]; ok {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  attr.Key + " should use " + ariaAttr + " instead",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Warning,
				})
				continue
			}

			// Check for patterns like data-is-expanded, data-is-selected
			lowerKey := strings.ToLower(attr.Key)
			for dataAttr, ariaAttr := range ariaEquivalents {
				// Check variations: data-is-expanded, data-isexpanded
				baseName := strings.TrimPrefix(dataAttr, "data-")
				if strings.Contains(lowerKey, baseName) {
					results = append(results, Result{
						Rule:     r.Name(),
						Message:  attr.Key + " appears to indicate state; consider using " + ariaAttr,
						Filename: doc.Filename,
						Line:     n.Line,
						Col:      n.Col,
						Severity: Info,
					})
					break
				}
			}
		}

		return true
	})

	return results
}
