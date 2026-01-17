package rules

import (
	"fmt"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// UniqueLandmark checks that duplicate landmark types have unique names.
type UniqueLandmark struct{}

func (r *UniqueLandmark) Name() string { return RuleUniqueLandmark }

func (r *UniqueLandmark) Description() string {
	return "multiple landmarks of same type must have unique accessible names"
}

// landmarkInfo holds information about a landmark element.
type landmarkInfo struct {
	node *parser.Node
	name string // accessible name from aria-label/aria-labelledby
}

func (r *UniqueLandmark) Check(doc *parser.Document) []Result {
	var results []Result

	// Collect landmarks by role
	landmarks := make(map[string][]landmarkInfo)

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tagName := strings.ToLower(n.Data)
		role := LandmarkElements[tagName]

		// Check for explicit role attribute
		if explicitRole := n.GetAttr("role"); explicitRole != "" {
			explicitRole = strings.ToLower(explicitRole)
			// role="presentation" or role="none" removes landmark semantics
			if explicitRole == "presentation" || explicitRole == "none" {
				return true
			}
			role = explicitRole
		}

		if role == "" {
			return true
		}

		// Skip form/section without accessible name (they're not landmarks)
		name := GetAccessibleName(n)
		if (tagName == "form" || tagName == "section") && name == "" {
			return true
		}

		landmarks[role] = append(landmarks[role], landmarkInfo{
			node: n,
			name: name,
		})

		return true
	})

	// Check for duplicate landmarks without unique names
	for role, infos := range landmarks {
		if len(infos) <= 1 {
			continue
		}

		// Check for duplicates or missing names
		seenNames := make(map[string]*parser.Node)
		for _, info := range infos {
			if info.name == "" {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  fmt.Sprintf("multiple %q landmarks, this one needs aria-label to distinguish it", role),
					Filename: doc.Filename,
					Line:     info.node.Line,
					Col:      info.node.Col,
					Severity: Warning,
				})
				continue
			}

			normalizedName := strings.ToLower(strings.TrimSpace(info.name))
			if prev, exists := seenNames[normalizedName]; exists {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  fmt.Sprintf("duplicate %q landmark name %q (first at line %d)", role, info.name, prev.Line),
					Filename: doc.Filename,
					Line:     info.node.Line,
					Col:      info.node.Col,
					Severity: Warning,
				})
			} else {
				seenNames[normalizedName] = info.node
			}
		}
	}

	return results
}
