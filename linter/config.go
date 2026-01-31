package linter

import (
	"slices"

	"github.com/toba/go-html-validate/rules"
)

// FrameworkConfig configures framework-specific attribute handling.
type FrameworkConfig struct {
	// HTMX enables htmx attribute validation.
	HTMX bool
	// HTMXVersion specifies which htmx version to validate against ("2" or "4").
	// Defaults to "2" when HTMX is enabled.
	HTMXVersion string
	// HTMXCustomEvents lists custom event names that should not trigger
	// "unknown event" warnings in hx-on:* validation.
	HTMXCustomEvents []string
}

// Config holds linter configuration options.
type Config struct {
	// EnabledRules lists rules to enable (empty means all)
	EnabledRules []string
	// DisabledRules lists rules to disable
	DisabledRules []string
	// RuleSeverity overrides severity for specific rules
	RuleSeverity map[string]rules.Severity
	// MinSeverity filters results to this severity or higher
	MinSeverity rules.Severity
	// IgnorePatterns are glob patterns for files to skip
	IgnorePatterns []string
	// ConfigPath is the path to the loaded config file (for debugging)
	ConfigPath string
	// Frameworks configures framework-specific attribute handling.
	Frameworks FrameworkConfig
}

// DefaultConfig returns a configuration with all rules enabled.
func DefaultConfig() *Config {
	return &Config{
		EnabledRules:   nil, // nil means all enabled
		DisabledRules:  nil,
		RuleSeverity:   make(map[string]rules.Severity),
		MinSeverity:    rules.Info, // Show everything by default
		IgnorePatterns: nil,
	}
}

// IsRuleEnabled checks if a rule should be run.
func (c *Config) IsRuleEnabled(name string) bool {
	// Check disabled list first
	if slices.Contains(c.DisabledRules, name) {
		return false
	}

	// If enabled list is specified, rule must be in it
	if len(c.EnabledRules) > 0 {
		return slices.Contains(c.EnabledRules, name)
	}

	// Default: all rules enabled
	return true
}

// ErrorsOnly configures the linter to only report errors.
func (c *Config) ErrorsOnly() *Config {
	c.MinSeverity = rules.Error
	return c
}

// WarningsAndErrors configures the linter to report warnings and errors.
func (c *Config) WarningsAndErrors() *Config {
	c.MinSeverity = rules.Warning
	return c
}
