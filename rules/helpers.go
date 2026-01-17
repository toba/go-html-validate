package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// TemplateExprPlaceholder is the placeholder inserted for Go template expressions.
// The parser replaces {{.Field}} expressions with this value.
const TemplateExprPlaceholder = "TMPL"

// IsTemplateExpr returns true if value contains template placeholder.
func IsTemplateExpr(s string) bool {
	return strings.Contains(s, TemplateExprPlaceholder) || strings.Contains(s, "{{")
}

// Tag returns the lowercase tag name of an element node.
func Tag(n *parser.Node) string {
	return strings.ToLower(n.Data)
}

// TagEquals checks if node's tag matches (case-insensitive).
func TagEquals(n *parser.Node, tag string) bool {
	return strings.EqualFold(n.Data, tag)
}

// TagIn checks if node's tag matches any given tag.
func TagIn(n *parser.Node, tags ...string) bool {
	for _, t := range tags {
		if strings.EqualFold(n.Data, t) {
			return true
		}
	}
	return false
}

// NewResult creates a Result with standard location info.
func NewResult(rule, message string, n *parser.Node, doc *parser.Document, sev Severity) Result {
	return Result{
		Rule:     rule,
		Message:  message,
		Filename: doc.Filename,
		Line:     n.Line,
		Col:      n.Col,
		Severity: sev,
	}
}

// NormalizeText collapses whitespace and lowercases text for comparison.
// Removes template placeholders (TMPL) used during parsing.
func NormalizeText(s string) string {
	s = strings.ReplaceAll(s, TemplateExprPlaceholder, "")
	s = strings.Join(strings.Fields(s), " ")
	return strings.ToLower(s)
}

// HasAccessibleName checks if an element has an accessible name via:
// aria-label, aria-labelledby, title, text content, or child img/svg with alt.
func HasAccessibleName(n *parser.Node) bool {
	// Check aria-label
	if n.GetAttr("aria-label") != "" {
		return true
	}

	// Check aria-labelledby
	if n.HasAttr("aria-labelledby") {
		return true
	}

	// Check title attribute
	if n.GetAttr("title") != "" {
		return true
	}

	// Check text content
	// Note: "TMPL" is the placeholder for Go template expressions {{.Field}}.
	// We treat TMPL as valid content since it will be replaced at runtime.
	text := n.TextContent()
	if strings.TrimSpace(text) != "" {
		return true
	}

	// Check for child img with alt or svg with aria-label/title
	for _, child := range n.Children {
		if child.IsElement("img") && child.GetAttr("alt") != "" {
			return true
		}
		if child.IsElement("svg") {
			if child.HasAttr("aria-label") || child.HasAttr("title") {
				return true
			}
		}
	}

	return false
}

// GetAccessibleName returns the accessible name from aria-label or aria-labelledby.
// Returns empty string if no accessible name is set.
func GetAccessibleName(n *parser.Node) string {
	if label := n.GetAttr("aria-label"); label != "" {
		return label
	}
	// aria-labelledby resolution would require ID lookup; indicate presence
	if n.GetAttr("aria-labelledby") != "" {
		return "[referenced]"
	}
	return ""
}

// FindDescendant searches for a descendant matching predicate (depth-first).
// Returns the first match or nil if none found.
func FindDescendant(n *parser.Node, pred func(*parser.Node) bool) *parser.Node {
	for _, child := range n.Children {
		if pred(child) {
			return child
		}
		if found := FindDescendant(child, pred); found != nil {
			return found
		}
	}
	return nil
}

// HasDescendant returns true if any descendant matches predicate.
func HasDescendant(n *parser.Node, pred func(*parser.Node) bool) bool {
	return FindDescendant(n, pred) != nil
}

// IsElementNode returns true if node is an HTML element.
func IsElementNode(n *parser.Node) bool {
	return n.Type == html.ElementNode
}

// GetImplicitRole returns the implicit ARIA role for an element.
func GetImplicitRole(tagName string, n *parser.Node) string {
	switch tagName {
	case "a":
		if n.GetAttr("href") != "" {
			return "link"
		}
		return ""
	case "input":
		inputType := strings.ToLower(n.GetAttr("type"))
		if inputType == "" {
			inputType = "text"
		}
		return InputTypeRoles[inputType]
	case "img":
		return "img"
	}
	return ImplicitRoles[tagName]
}

// AncestorWithTag walks up the tree looking for an ancestor with the given tag.
// Returns the first matching ancestor or nil if none found.
func AncestorWithTag(n *parser.Node, tag string) *parser.Node {
	tag = strings.ToLower(tag)
	for p := n.Parent; p != nil; p = p.Parent {
		if p.Type == html.ElementNode && strings.ToLower(p.Data) == tag {
			return p
		}
	}
	return nil
}

// HasAncestor returns true if the node has an ancestor matching any of the given tags.
func HasAncestor(n *parser.Node, tags ...string) bool {
	for _, tag := range tags {
		if AncestorWithTag(n, tag) != nil {
			return true
		}
	}
	return false
}

// ChildElements returns only element node children (excludes text, comments, etc.).
func ChildElements(n *parser.Node) []*parser.Node {
	var children []*parser.Node
	for _, child := range n.Children {
		if child.Type == html.ElementNode {
			children = append(children, child)
		}
	}
	return children
}

// FirstChildElement returns the first element child, or nil if none.
func FirstChildElement(n *parser.Node) *parser.Node {
	for _, child := range n.Children {
		if child.Type == html.ElementNode {
			return child
		}
	}
	return nil
}

// CountChildrenWithTag counts direct children matching the given tag.
func CountChildrenWithTag(n *parser.Node, tag string) int {
	tag = strings.ToLower(tag)
	count := 0
	for _, child := range n.Children {
		if child.Type == html.ElementNode && strings.ToLower(child.Data) == tag {
			count++
		}
	}
	return count
}

// HasChildWithTag returns true if there's a direct child with the given tag.
func HasChildWithTag(n *parser.Node, tag string) bool {
	return CountChildrenWithTag(n, tag) > 0
}

// IsCustomElement returns true if the element is a valid custom element name.
// Custom elements must contain a hyphen and start with a lowercase letter.
func IsCustomElement(tagName string) bool {
	if tagName == "" {
		return false
	}
	// Must start with lowercase ASCII letter
	first := tagName[0]
	if first < 'a' || first > 'z' {
		return false
	}
	// Must contain a hyphen
	return strings.Contains(tagName, "-")
}
