package rules

import (
	"regexp"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// Default patterns - kebab-case for CSS classes, camelCase or kebab-case for IDs/names.
var (
	// Default: lowercase letters, numbers, hyphens, underscores
	defaultClassPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
	// Default: start with letter, allow letters, numbers, hyphens, underscores
	defaultIDPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	// Default: start with letter, allow letters, numbers, underscores, brackets for arrays
	defaultNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\[\]]*$`)
)

// ClassPattern checks that class names follow a naming convention.
type ClassPattern struct {
	Pattern *regexp.Regexp
}

// Name returns the rule identifier.
func (r *ClassPattern) Name() string { return RuleClassPattern }

// Description returns what this rule checks.
func (r *ClassPattern) Description() string {
	return "class names should follow naming convention"
}

// Check examines the document for class names not matching pattern.
func (r *ClassPattern) Check(doc *parser.Document) []Result {
	var results []Result

	pattern := r.Pattern
	if pattern == nil {
		pattern = defaultClassPattern
	}

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		classAttr := n.GetAttr("class")
		if classAttr == "" {
			return true
		}

		// Check each class name
		for class := range strings.FieldsSeq(classAttr) {
			// Skip template expressions
			if IsTemplateExpr(class) {
				continue
			}
			if !pattern.MatchString(class) {
				results = append(results, Result{
					Rule:     RuleClassPattern,
					Message:  "class \"" + class + "\" does not match naming convention",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Info,
				})
			}
		}

		return true
	})

	return results
}

// IDPattern checks that id attributes follow a naming convention.
type IDPattern struct {
	Pattern *regexp.Regexp
}

// Name returns the rule identifier.
func (r *IDPattern) Name() string { return RuleIDPattern }

// Description returns what this rule checks.
func (r *IDPattern) Description() string {
	return "id attributes should follow naming convention"
}

// Check examines the document for id values not matching pattern.
func (r *IDPattern) Check(doc *parser.Document) []Result {
	var results []Result

	pattern := r.Pattern
	if pattern == nil {
		pattern = defaultIDPattern
	}

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		id := n.GetAttr("id")
		if id == "" {
			return true
		}

		// Skip template expressions
		if IsTemplateExpr(id) {
			return true
		}

		if !pattern.MatchString(id) {
			results = append(results, Result{
				Rule:     RuleIDPattern,
				Message:  "id \"" + id + "\" does not match naming convention",
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

// NamePattern checks that name attributes follow a naming convention.
type NamePattern struct {
	Pattern *regexp.Regexp
}

// Name returns the rule identifier.
func (r *NamePattern) Name() string { return RuleNamePattern }

// Description returns what this rule checks.
func (r *NamePattern) Description() string {
	return "name attributes should follow naming convention"
}

// Check examines the document for name values not matching pattern.
func (r *NamePattern) Check(doc *parser.Document) []Result {
	var results []Result

	pattern := r.Pattern
	if pattern == nil {
		pattern = defaultNamePattern
	}

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		name := n.GetAttr("name")
		if name == "" {
			return true
		}

		// Skip template expressions
		if IsTemplateExpr(name) {
			return true
		}

		if !pattern.MatchString(name) {
			results = append(results, Result{
				Rule:     RuleNamePattern,
				Message:  "name \"" + name + "\" does not match naming convention",
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
