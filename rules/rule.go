package rules

import (
	"github.com/toba/go-html-validate/parser"
)

// Severity indicates how serious a lint violation is.
type Severity int

const (
	// Error indicates a definite accessibility problem.
	Error Severity = iota
	// Warning indicates a potential issue or best practice violation.
	Warning
	// Info indicates a suggestion for improvement.
	Info
)

func (s Severity) String() string {
	switch s {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Info:
		return "info"
	default:
		return "unknown"
	}
}

// Rule name constants to avoid magic strings.
const (
	RuleImgAlt                      = "img-alt"
	RuleAreaAlt                     = "area-alt"
	RuleInputLabel                  = "input-label"
	RuleButtonName                  = "button-name"
	RuleLinkName                    = "link-name"
	RuleHeadingContent              = "heading-content"
	RuleHeadingLevel                = "heading-level"
	RuleTextContent                 = "text-content"
	RuleEmptyTitle                  = "empty-title"
	RuleLongTitle                   = "long-title"
	RulePreferAria                  = "prefer-aria"
	RuleAriaHiddenBody              = "aria-hidden-body"
	RuleHiddenFocusable             = "hidden-focusable"
	RuleRedundantAriaLabel          = "no-redundant-aria-label"
	RuleNoRedundantRole             = "no-redundant-role"
	RuleNoAbstractRole              = "no-abstract-role"
	RuleAriaLabelMisuse             = "aria-label-misuse"
	RuleUniqueLandmark              = "unique-landmark"
	RuleFormSubmit                  = "form-submit"
	RuleButtonType                  = "button-type"
	RuleMultipleLabeledControls     = "multiple-labeled-controls"
	RuleFormDupName                 = "form-dup-name"
	RuleNoRedundantFor              = "no-redundant-for"
	RuleValidAutocomplete           = "valid-autocomplete"
	RuleNoImplicitInputType         = "no-implicit-input-type"
	RuleInputAttributes             = "input-attributes"
	RuleTabindexNoPositive          = "tabindex-no-positive"
	RuleSVGFocusable                = "svg-focusable"
	RuleNoAutoplay                  = "no-autoplay"
	RuleMetaRefresh                 = "meta-refresh"
	RuleWcagH36                     = "wcag/h36"
	RuleWcagH63                     = "wcag/h63"
	RuleWcagH67                     = "wcag/h67"
	RuleWcagH71                     = "wcag/h71"
	RulePreferSemantic              = "prefer-semantic"
	RuleDuplicateID                 = "duplicate-id"
	RulePreferButton                = "prefer-button"
	RuleNoInlineStyle               = "no-inline-style"
	RulePreferNativeElement         = "prefer-native-element"
	RulePreferTbody                 = "prefer-tbody"
	RuleNoDupAttr                   = "no-dup-attr"
	RuleNoDupClass                  = "no-dup-class"
	RuleMapIDName                   = "map-id-name"
	RuleMapDupName                  = "map-dup-name"
	RuleElementName                 = "element-name"
	RuleScriptType                  = "script-type"
	RuleScriptElement               = "script-element"
	RuleAttributeAllowedValues      = "attribute-allowed-values"
	RuleAttributeMisuse             = "attribute-misuse"
	RuleDeprecated                  = "deprecated"
	RuleNoDeprecatedAttr            = "no-deprecated-attr"
	RuleNoConditionalComment        = "no-conditional-comment"
	RuleVoidContent                 = "void-content"
	RuleElementRequiredAncestor     = "element-required-ancestor"
	RuleElementPermittedParent      = "element-permitted-parent"
	RuleElementRequiredAttributes   = "element-required-attributes"
	RuleElementPermittedContent     = "element-permitted-content"
	RuleElementPermittedOccurrences = "element-permitted-occurrences"
	RuleElementRequiredContent      = "element-required-content"
	RuleElementPermittedOrder       = "element-permitted-order"
	RuleNoMultipleMain              = "no-multiple-main"
	RuleValidID                     = "valid-id"
	RuleRequireLang                 = "require-lang"
	RuleNoMissingReferences         = "no-missing-references"
	RuleAllowedLinks                = "allowed-links"
	RuleNoUTF8BOM                   = "no-utf8-bom"
	RuleTelNonBreaking              = "tel-non-breaking"
	RuleRequireSRI                  = "require-sri"
	RuleRequireCSPNonce             = "require-csp-nonce"
	RuleNoStyleTag                  = "no-style-tag"
	RuleClassPattern                = "class-pattern"
	RuleIDPattern                   = "id-pattern"
	RuleNamePattern                 = "name-pattern"
	RuleValidFor                    = "valid-for"
	RuleUnrecognizedCharRef         = "unrecognized-char-ref"
	RuleHTMXAttributes              = "htmx-attributes"
	RuleTemplateWhitespaceTrim      = "template-whitespace-trim"
	RuleTemplateSyntaxValid         = "template-syntax-valid"
)

// Result represents a single lint finding.
type Result struct {
	Rule     string   // Rule name (e.g., "img-alt")
	Message  string   // Human-readable description
	Filename string   // Source file path
	Line     int      // 1-indexed line number
	Col      int      // 1-indexed column number
	Severity Severity // Error, Warning, or Info
}

// Rule defines the interface for accessibility rules.
type Rule interface {
	// Name returns the rule identifier (e.g., "img-alt")
	Name() string
	// Description returns a brief explanation of what the rule checks
	Description() string
	// Check examines a document and returns any violations found
	Check(doc *parser.Document) []Result
}

// HTMXConfigurable is implemented by rules that need htmx configuration.
type HTMXConfigurable interface {
	Configure(htmxEnabled bool, htmxVersion string)
}

// HTMXCustomEventsConfigurable is implemented by rules that accept custom event names.
type HTMXCustomEventsConfigurable interface {
	ConfigureCustomEvents(events []string)
}

// RawRule is implemented by rules that need access to the raw file content
// before template preprocessing. This allows linting template syntax itself.
type RawRule interface {
	Rule
	CheckRaw(filename string, content []byte) []Result
}

// Registry holds all available rules.
type Registry struct {
	rules []Rule
}

// NewRegistry creates a registry with all default rules.
func NewRegistry() *Registry {
	return &Registry{
		rules: []Rule{
			// Accessibility - content
			&ImgAlt{},
			&InputLabel{},
			&ButtonName{},
			&LinkName{},
			&HeadingContent{},
			&HeadingLevel{},
			&EmptyTitle{},
			// Accessibility - ARIA
			&PreferAria{},
			&AriaHiddenBody{},
			&HiddenFocusable{},
			&RedundantAriaLabel{},
			&NoRedundantRole{},
			&NoAbstractRole{},
			&AriaLabelMisuse{},
			&UniqueLandmark{},
			// Accessibility - forms
			&FormSubmit{},
			&ButtonType{},
			&MultipleLabeledControls{},
			// Accessibility - focus/navigation
			&TabindexNoPositive{},
			&SVGFocusable{},
			// Accessibility - media
			&NoAutoplay{},
			&MetaRefresh{},
			// Best practices
			&PreferSemantic{},
			&DuplicateID{},
			&PreferButton{},
			&NoInlineStyle{},
			// SEO
			&LongTitle{},
			// Security
			&RequireSRI{},
			// Document structure
			&NoMultipleMain{},
			&ValidID{},
			&RequireLang{},
			&PreferNativeElement{},
			// WCAG accessibility rules
			&AreaAlt{},
			&TextContent{},
			&WcagH36{},
			&WcagH63{},
			&WcagH67{},
			&WcagH71{},
			&TelNonBreaking{},
			// Syntax validation rules
			&NoDupAttr{},
			&NoDupClass{},
			&NoRedundantFor{},
			&FormDupName{},
			&MapDupName{},
			&MapIDName{},
			&ElementName{},
			&ScriptType{},
			&ValidAutocomplete{},
			&ValidFor{},
			&UnrecognizedCharRef{},
			// Deprecated rules
			&Deprecated{},
			&NoDeprecatedAttr{},
			&NoConditionalComment{},
			// Content model rules
			&VoidContent{},
			&ElementRequiredAncestor{},
			&ElementPermittedParent{},
			&ElementRequiredAttributes{},
			&ElementPermittedContent{},
			&ElementPermittedOccurrences{},
			&ElementRequiredContent{},
			&ElementPermittedOrder{},
			&AttributeAllowedValues{},
			&AttributeMisuse{},
			&InputAttributes{},
			&ScriptElement{},
			// Document structure rules
			&DoctypeHTML{},
			&MissingDoctype{},
			&NoUTF8BOM{},
			&NoMissingReferences{},
			&AllowedLinks{},
			// Security rules
			&RequireCSPNonce{},
			// Style rules
			&NoStyleTag{},
			&PreferTbody{},
			&NoImplicitInputType{},
			&ClassPattern{},
			&IDPattern{},
			&NamePattern{},
			// htmx rules
			&HTMXAttributes{},
			// Template rules
			&TemplateWhitespaceTrim{},
			&TemplateSyntaxValid{},
		},
	}
}

// All returns all registered rules.
func (r *Registry) All() []Rule {
	return r.rules
}

// ByName returns a rule by name, or nil if not found.
func (r *Registry) ByName(name string) Rule {
	for _, rule := range r.rules {
		if rule.Name() == name {
			return rule
		}
	}
	return nil
}
