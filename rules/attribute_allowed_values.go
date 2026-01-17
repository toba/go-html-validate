package rules

import (
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// AttributeAllowedValues checks that attributes have valid values.
type AttributeAllowedValues struct{}

// Name returns the rule identifier.
func (r *AttributeAllowedValues) Name() string { return RuleAttributeAllowedValues }

// Description returns what this rule checks.
func (r *AttributeAllowedValues) Description() string {
	return "attributes must have allowed values"
}

// Check examines the document for attributes with invalid values.
func (r *AttributeAllowedValues) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Check element-specific attributes
		switch tag {
		case "input":
			results = append(results, r.checkInputType(n, doc)...)
		case "button":
			results = append(results, r.checkButtonType(n, doc)...)
		case "form":
			results = append(results, r.checkFormAttrs(n, doc)...)
		case "a":
			results = append(results, r.checkAnchorRel(n, doc)...)
		case "link":
			results = append(results, r.checkLinkRel(n, doc)...)
		case "th":
			results = append(results, r.checkThScope(n, doc)...)
		case "img", "iframe":
			results = append(results, r.checkLoadingDecoding(n, doc)...)
		}

		// Check global attributes
		results = append(results, r.checkDirAttr(n, doc)...)
		results = append(results, r.checkCrossOrigin(n, doc)...)
		results = append(results, r.checkReferrerPolicy(n, doc)...)

		return true
	})

	return results
}

func (r *AttributeAllowedValues) checkInputType(n *parser.Node, doc *parser.Document) []Result {
	val := n.GetAttr("type")
	if val == "" {
		return nil
	}
	val = strings.ToLower(val)
	if !ValidInputTypes[val] {
		return []Result{{
			Rule:     RuleAttributeAllowedValues,
			Message:  "invalid input type: " + val,
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		}}
	}
	return nil
}

func (r *AttributeAllowedValues) checkButtonType(n *parser.Node, doc *parser.Document) []Result {
	val := n.GetAttr("type")
	if val == "" {
		return nil
	}
	val = strings.ToLower(val)
	if !ValidButtonTypes[val] {
		return []Result{{
			Rule:     RuleAttributeAllowedValues,
			Message:  "invalid button type: " + val,
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		}}
	}
	return nil
}

func (r *AttributeAllowedValues) checkFormAttrs(n *parser.Node, doc *parser.Document) []Result {
	var results []Result

	// Check method
	if method := n.GetAttr("method"); method != "" {
		if !ValidFormMethods[strings.ToLower(method)] {
			results = append(results, Result{
				Rule:     RuleAttributeAllowedValues,
				Message:  "invalid form method: " + method,
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		}
	}

	// Check enctype
	if enctype := n.GetAttr("enctype"); enctype != "" {
		if !ValidFormEnctypes[strings.ToLower(enctype)] {
			results = append(results, Result{
				Rule:     RuleAttributeAllowedValues,
				Message:  "invalid form enctype: " + enctype,
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		}
	}

	return results
}

func (r *AttributeAllowedValues) checkAnchorRel(n *parser.Node, doc *parser.Document) []Result {
	val := n.GetAttr("rel")
	if val == "" {
		return nil
	}

	var results []Result
	// rel can be space-separated list
	for rel := range strings.FieldsSeq(val) {
		rel = strings.ToLower(rel)
		if !ValidAnchorRels[rel] {
			results = append(results, Result{
				Rule:     RuleAttributeAllowedValues,
				Message:  "invalid anchor rel value: " + rel,
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}
	}
	return results
}

func (r *AttributeAllowedValues) checkLinkRel(n *parser.Node, doc *parser.Document) []Result {
	val := n.GetAttr("rel")
	if val == "" {
		return nil
	}

	var results []Result
	// rel can be space-separated list
	for rel := range strings.FieldsSeq(val) {
		rel = strings.ToLower(rel)
		if !ValidLinkRels[rel] {
			results = append(results, Result{
				Rule:     RuleAttributeAllowedValues,
				Message:  "invalid link rel value: " + rel,
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}
	}
	return results
}

func (r *AttributeAllowedValues) checkThScope(n *parser.Node, doc *parser.Document) []Result {
	val := n.GetAttr("scope")
	if val == "" {
		return nil
	}
	val = strings.ToLower(val)
	if !ValidScopeValues[val] {
		return []Result{{
			Rule:     RuleAttributeAllowedValues,
			Message:  "invalid th scope value: " + val,
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		}}
	}
	return nil
}

func (r *AttributeAllowedValues) checkLoadingDecoding(n *parser.Node, doc *parser.Document) []Result {
	var results []Result

	if loading := n.GetAttr("loading"); loading != "" {
		if !ValidLoadingValues[strings.ToLower(loading)] {
			results = append(results, Result{
				Rule:     RuleAttributeAllowedValues,
				Message:  "invalid loading value: " + loading,
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		}
	}

	if decoding := n.GetAttr("decoding"); decoding != "" {
		if !ValidDecodingValues[strings.ToLower(decoding)] {
			results = append(results, Result{
				Rule:     RuleAttributeAllowedValues,
				Message:  "invalid decoding value: " + decoding,
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		}
	}

	return results
}

func (r *AttributeAllowedValues) checkDirAttr(n *parser.Node, doc *parser.Document) []Result {
	val := n.GetAttr("dir")
	if val == "" {
		return nil
	}
	val = strings.ToLower(val)
	if !ValidDirValues[val] {
		return []Result{{
			Rule:     RuleAttributeAllowedValues,
			Message:  "invalid dir value: " + val,
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		}}
	}
	return nil
}

func (r *AttributeAllowedValues) checkCrossOrigin(n *parser.Node, doc *parser.Document) []Result {
	if !n.HasAttr("crossorigin") {
		return nil
	}
	val := strings.ToLower(n.GetAttr("crossorigin"))
	if !ValidCrossOriginValues[val] {
		return []Result{{
			Rule:     RuleAttributeAllowedValues,
			Message:  "invalid crossorigin value: " + val,
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		}}
	}
	return nil
}

func (r *AttributeAllowedValues) checkReferrerPolicy(n *parser.Node, doc *parser.Document) []Result {
	if !n.HasAttr("referrerpolicy") {
		return nil
	}
	val := strings.ToLower(n.GetAttr("referrerpolicy"))
	if !ValidReferrerPolicies[val] {
		return []Result{{
			Rule:     RuleAttributeAllowedValues,
			Message:  "invalid referrerpolicy value: " + val,
			Filename: doc.Filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		}}
	}
	return nil
}
