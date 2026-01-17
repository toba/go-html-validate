package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// roleToNative maps ARIA roles to their native HTML element equivalents.
var roleToNative = map[string]string{
	"main":          "<main>",
	"navigation":    "<nav>",
	"banner":        "<header>",
	"contentinfo":   "<footer>",
	"complementary": "<aside>",
	"article":       "<article>",
	"region":        "<section>",
	"form":          "<form>",
	"search":        "<search>",
}

// PreferNativeElement checks for ARIA roles that have native HTML equivalents.
type PreferNativeElement struct{}

func (r *PreferNativeElement) Name() string { return RulePreferNativeElement }

func (r *PreferNativeElement) Description() string {
	return "prefer native HTML elements over ARIA roles"
}

func (r *PreferNativeElement) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		role := n.GetAttr("role")
		if role == "" {
			return true
		}

		native, ok := roleToNative[role]
		if !ok {
			return true
		}

		results = append(results, Result{
			Rule:     r.Name(),
			Message:  "use " + native + " element instead of <" + n.Data + " role=\"" + role + "\">; native elements have better accessibility support",
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Warning,
		})

		return true
	})

	return results
}
