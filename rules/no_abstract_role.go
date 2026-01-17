package rules

import (
	"fmt"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoAbstractRole checks for abstract ARIA roles that shouldn't be used.
type NoAbstractRole struct{}

func (r *NoAbstractRole) Name() string { return RuleNoAbstractRole }

func (r *NoAbstractRole) Description() string {
	return "abstract ARIA roles must not be used in content"
}

func (r *NoAbstractRole) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		roleAttr := n.GetAttr("role")
		if roleAttr == "" {
			return true
		}

		// Role attribute can contain multiple space-separated roles
		for role := range strings.FieldsSeq(roleAttr) {
			role = strings.ToLower(role)
			if AbstractRoles[role] {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  fmt.Sprintf("abstract role %q must not be used", role),
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		}

		return true
	})

	return results
}
