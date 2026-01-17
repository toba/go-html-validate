package config

import "github.com/STR-Consulting/go-html-validate/rules"

// Presets contains built-in configuration presets.
var Presets = map[string]*FileConfig{
	"html-validate:recommended": recommendedPreset(),
	"html-validate:standard":    standardPreset(),
	"html-validate:a11y":        a11yPreset(),
}

// recommendedPreset returns the recommended preset with all rules at default severity.
func recommendedPreset() *FileConfig {
	return &FileConfig{
		Rules: map[string]RuleConfig{
			// All rules enabled at their defaults (empty means use rule's default)
		},
	}
}

// standardPreset returns the standard preset with core rules enabled.
// Disables some style-preference rules.
func standardPreset() *FileConfig {
	return &FileConfig{
		Rules: map[string]RuleConfig{
			rules.RulePreferTbody:         {Severity: "off"},
			rules.RuleNoInlineStyle:       {Severity: "off"},
			rules.RulePreferSemantic:      {Severity: "off"},
			rules.RuleClassPattern:        {Severity: "off"},
			rules.RuleIDPattern:           {Severity: "off"},
			rules.RuleNamePattern:         {Severity: "off"},
			rules.RuleNoStyleTag:          {Severity: "off"},
			rules.RulePreferNativeElement: {Severity: "off"},
		},
	}
}

// a11yPreset returns the accessibility-focused preset.
// Enables all accessibility rules, disables validation-only rules.
func a11yPreset() *FileConfig {
	return &FileConfig{
		Rules: map[string]RuleConfig{
			// Enable accessibility rules at error level
			rules.RuleImgAlt:             {Severity: "error"},
			rules.RuleAreaAlt:            {Severity: "error"},
			rules.RuleInputLabel:         {Severity: "error"},
			rules.RuleButtonName:         {Severity: "error"},
			rules.RuleLinkName:           {Severity: "error"},
			rules.RuleHeadingContent:     {Severity: "error"},
			rules.RuleHeadingLevel:       {Severity: "error"},
			rules.RuleTextContent:        {Severity: "error"},
			rules.RuleEmptyTitle:         {Severity: "error"},
			rules.RulePreferAria:         {Severity: "warn"},
			rules.RuleAriaHiddenBody:     {Severity: "error"},
			rules.RuleHiddenFocusable:    {Severity: "error"},
			rules.RuleAriaLabelMisuse:    {Severity: "error"},
			rules.RuleUniqueLandmark:     {Severity: "warn"},
			rules.RuleFormSubmit:         {Severity: "warn"},
			rules.RuleButtonType:         {Severity: "warn"},
			rules.RuleTabindexNoPositive: {Severity: "error"},
			rules.RuleSVGFocusable:       {Severity: "warn"},
			rules.RuleNoAutoplay:         {Severity: "error"},
			rules.RuleMetaRefresh:        {Severity: "error"},
			rules.RuleWcagH36:            {Severity: "error"},
			rules.RuleWcagH63:            {Severity: "error"},
			rules.RuleWcagH67:            {Severity: "error"},
			rules.RuleWcagH71:            {Severity: "error"},
			rules.RuleRequireLang:        {Severity: "error"},

			// Disable validation-only and style rules
			rules.RulePreferTbody:                 {Severity: "off"},
			rules.RuleNoInlineStyle:               {Severity: "off"},
			rules.RuleClassPattern:                {Severity: "off"},
			rules.RuleIDPattern:                   {Severity: "off"},
			rules.RuleNamePattern:                 {Severity: "off"},
			rules.RuleNoStyleTag:                  {Severity: "off"},
			rules.RuleDeprecated:                  {Severity: "off"},
			rules.RuleNoDeprecatedAttr:            {Severity: "off"},
			rules.RuleNoConditionalComment:        {Severity: "off"},
			rules.RuleElementName:                 {Severity: "off"},
			rules.RuleScriptType:                  {Severity: "off"},
			rules.RuleAttributeAllowedValues:      {Severity: "off"},
			rules.RuleVoidContent:                 {Severity: "off"},
			rules.RuleElementRequiredAncestor:     {Severity: "off"},
			rules.RuleElementPermittedParent:      {Severity: "off"},
			rules.RuleElementPermittedContent:     {Severity: "off"},
			rules.RuleElementPermittedOccurrences: {Severity: "off"},
			rules.RuleElementRequiredContent:      {Severity: "off"},
			rules.RuleElementPermittedOrder:       {Severity: "off"},
		},
	}
}
