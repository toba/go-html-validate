package rules

import (
	"fmt"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// NoRedundantRole checks for explicit roles that match implicit roles.
type NoRedundantRole struct{}

func (r *NoRedundantRole) Name() string { return RuleNoRedundantRole }

func (r *NoRedundantRole) Description() string {
	return "element should not have role matching its implicit role"
}

func (r *NoRedundantRole) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		role := strings.ToLower(n.GetAttr("role"))
		if role == "" {
			return true
		}

		tagName := strings.ToLower(n.Data)
		implicitRole := GetImplicitRole(tagName, n)

		if role == implicitRole {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  fmt.Sprintf("<%s> has implicit role %q, explicit role is redundant", tagName, role),
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
