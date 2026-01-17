package rules

import (
	"slices"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// AttributeMisuse checks that attributes are used on correct elements.
type AttributeMisuse struct {
	htmxEnabled bool
}

// Configure implements HTMXConfigurable.
func (r *AttributeMisuse) Configure(htmxEnabled bool, _ string) {
	r.htmxEnabled = htmxEnabled
}

// Name returns the rule identifier.
func (r *AttributeMisuse) Name() string { return RuleAttributeMisuse }

// Description returns what this rule checks.
func (r *AttributeMisuse) Description() string {
	return "attributes must be used on appropriate elements"
}

// attributeElementMap defines which attributes are valid on which elements.
// Empty slice means the attribute is global.
var attributeElementMap = map[string][]string{
	// Form-related attributes
	"action":         {"form"},
	"method":         {"form"},
	"enctype":        {"form"},
	"novalidate":     {"form"},
	"accept":         {"input"},
	"autocomplete":   {"form", "input", "select", "textarea"},
	"autofocus":      {"button", "input", "select", "textarea"},
	"cols":           {"textarea"},
	"rows":           {"textarea"},
	"disabled":       {"button", "fieldset", "input", "optgroup", "option", "select", "textarea"},
	"for":            {"label", "output"},
	"form":           {"button", "fieldset", "input", "label", "meter", "object", "output", "progress", "select", "textarea"},
	"formaction":     {"button", "input"},
	"formenctype":    {"button", "input"},
	"formmethod":     {"button", "input"},
	"formnovalidate": {"button", "input"},
	"formtarget":     {"button", "input"},
	"maxlength":      {"input", "textarea"},
	"minlength":      {"input", "textarea"},
	"multiple":       {"input", "select"},
	"pattern":        {"input"},
	"placeholder":    {"input", "textarea"},
	"readonly":       {"input", "textarea"},
	"required":       {"input", "select", "textarea"},
	"size":           {"input", "select"},
	"step":           {"input"},
	"wrap":           {"textarea"},

	// Media attributes
	"autoplay": {"audio", "video"},
	"controls": {"audio", "video"},
	"loop":     {"audio", "video"},
	"muted":    {"audio", "video"},
	"poster":   {"video"},
	"preload":  {"audio", "video"},
	"src":      {"audio", "embed", "iframe", "img", "input", "script", "source", "track", "video"},
	"srcdoc":   {"iframe"},
	"srcset":   {"img", "source"},
	"sizes":    {"img", "link", "source"},
	"width":    {"canvas", "embed", "iframe", "img", "input", "object", "video"},
	"height":   {"canvas", "embed", "iframe", "img", "input", "object", "video"},

	// Table attributes
	"colspan": {"td", "th"},
	"rowspan": {"td", "th"},
	"headers": {"td", "th"},
	"scope":   {"th"},
	"span":    {"col", "colgroup"},

	// Link/anchor attributes
	"download": {"a", "area"},
	"href":     {"a", "area", "base", "link"},
	"hreflang": {"a", "link"},
	"ping":     {"a", "area"},
	"rel":      {"a", "area", "form", "link"},
	"target":   {"a", "area", "base", "form"},
	"type":     {"a", "button", "embed", "input", "link", "object", "ol", "script", "source", "style"},

	// Other element-specific
	"alt":             {"area", "img", "input"},
	"checked":         {"input"},
	"cite":            {"blockquote", "del", "ins", "q"},
	"datetime":        {"del", "ins", "time"},
	"default":         {"track"},
	"defer":           {"script"},
	"dirname":         {"input", "textarea"},
	"high":            {"meter"},
	"kind":            {"track"},
	"label":           {"optgroup", "option", "track"},
	"list":            {"input"},
	"low":             {"meter"},
	"max":             {"input", "meter", "progress"},
	"min":             {"input", "meter"},
	"name":            {"button", "fieldset", "form", "iframe", "input", "map", "meta", "object", "output", "param", "select", "slot", "textarea"},
	"optimum":         {"meter"},
	"selected":        {"option"},
	"srclang":         {"track"},
	"start":           {"ol"},
	"value":           {"button", "data", "input", "li", "meter", "option", "param", "progress"},
	"async":           {"script"},
	"charset":         {"meta", "script"},
	"content":         {"meta"},
	"http-equiv":      {"meta"},
	"integrity":       {"link", "script"},
	"crossorigin":     {"audio", "img", "link", "script", "video"},
	"referrerpolicy":  {"a", "area", "iframe", "img", "link", "script"},
	"loading":         {"iframe", "img"},
	"decoding":        {"img"},
	"ismap":           {"img"},
	"usemap":          {"img", "object"},
	"sandbox":         {"iframe"},
	"allow":           {"iframe"},
	"allowfullscreen": {"iframe"},
	"open":            {"details", "dialog"},
	"reversed":        {"ol"},
	"coords":          {"area"},
	"shape":           {"area"},
	"media":           {"link", "source", "style"},
}

// Check examines the document for misused attributes.
func (r *AttributeMisuse) Check(doc *parser.Document) []Result {
	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		tag := strings.ToLower(n.Data)

		// Skip custom elements
		if IsCustomElement(tag) {
			return true
		}

		// Check each attribute
		for _, attr := range n.Attr {
			attrName := strings.ToLower(attr.Key)

			// Skip data-* and aria-* attributes (always valid)
			if strings.HasPrefix(attrName, "data-") || strings.HasPrefix(attrName, "aria-") {
				continue
			}

			// Skip hx-* attributes when htmx is enabled (handled by htmx-attributes rule)
			if r.htmxEnabled && strings.HasPrefix(attrName, "hx-") {
				continue
			}

			// Skip global attributes
			if isGlobalAttribute(attrName) {
				continue
			}

			// Skip event handlers
			if strings.HasPrefix(attrName, "on") {
				continue
			}

			// Check if attribute is element-specific
			validElements, hasConstraint := attributeElementMap[attrName]
			if !hasConstraint {
				continue // Unknown attribute, let other rules handle
			}

			// Check if element is in valid list
			if !slices.Contains(validElements, tag) {
				results = append(results, Result{
					Rule:     RuleAttributeMisuse,
					Message:  "attribute '" + attrName + "' is not valid on <" + tag + ">",
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

// isGlobalAttribute returns true for attributes valid on all elements.
func isGlobalAttribute(attr string) bool {
	globalAttrs := map[string]bool{
		"accesskey":          true,
		"autocapitalize":     true,
		"autofocus":          false, // Has specific elements
		"class":              true,
		"contenteditable":    true,
		"dir":                true,
		"draggable":          true,
		"enterkeyhint":       true,
		"hidden":             true,
		"id":                 true,
		"inert":              true,
		"inputmode":          true,
		"is":                 true,
		"itemid":             true,
		"itemprop":           true,
		"itemref":            true,
		"itemscope":          true,
		"itemtype":           true,
		"lang":               true,
		"nonce":              true,
		"part":               true,
		"popover":            true,
		"role":               true,
		"slot":               true,
		"spellcheck":         true,
		"style":              true,
		"tabindex":           true,
		"title":              true,
		"translate":          true,
		"writingsuggestions": true,
	}
	return globalAttrs[attr]
}
