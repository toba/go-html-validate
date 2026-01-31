package rules

import "strings"

// HTMXv2Attributes contains all valid htmx 2.x attributes.
var HTMXv2Attributes = map[string]bool{
	// Core request attributes
	"hx-get":    true,
	"hx-post":   true,
	"hx-put":    true,
	"hx-patch":  true,
	"hx-delete": true,

	// Content swapping & targeting
	"hx-swap":       true,
	"hx-swap-oob":   true,
	"hx-target":     true,
	"hx-select":     true,
	"hx-select-oob": true,
	"hx-trigger":    true,
	"hx-sync":       true,

	// Request configuration
	"hx-boost":        true,
	"hx-push-url":     true,
	"hx-replace-url":  true,
	"hx-vals":         true,
	"hx-vars":         true, // deprecated alias for hx-vals
	"hx-headers":      true,
	"hx-params":       true,
	"hx-include":      true,
	"hx-encoding":     true,
	"hx-request":      true,
	"hx-confirm":      true,
	"hx-prompt":       true,
	"hx-validate":     true,
	"hx-disable":      true,
	"hx-disabled-elt": true,

	// Inheritance control
	"hx-disinherit": true,
	"hx-inherit":    true,

	// Extensions
	"hx-ext": true,

	// History
	"hx-history":     true,
	"hx-history-elt": true,

	// UI feedback
	"hx-indicator": true,
	"hx-preserve":  true,
}

// HTMXv4OnlyAttributes contains attributes that are only valid in htmx 4.x.
var HTMXv4OnlyAttributes = map[string]bool{
	"hx-action":     true,
	"hx-config":     true,
	"hx-ignore":     true,
	"hx-method":     true,
	"hx-optimistic": true,
	"hx-preload":    true,
}

// HTMXv4DeprecatedAttributes contains v2 attributes deprecated/removed in v4.
var HTMXv4DeprecatedAttributes = map[string]bool{
	"hx-disabled-elt": true,
	"hx-disinherit":   true,
	"hx-history-elt":  true,
	"hx-request":      true,
	"hx-vars":         true,
}

// HTMXv4Attributes contains all valid htmx 4.x attributes (v2 minus deprecated, plus v4-only).
var HTMXv4Attributes = func() map[string]bool {
	attrs := make(map[string]bool)
	for attr := range HTMXv2Attributes {
		if !HTMXv4DeprecatedAttributes[attr] {
			attrs[attr] = true
		}
	}
	for attr := range HTMXv4OnlyAttributes {
		attrs[attr] = true
	}
	return attrs
}()

// IsHTMXAttribute checks if an attribute name is an htmx attribute pattern.
// Returns true if the attribute starts with "hx-" or matches "hx-on:*" pattern.
func IsHTMXAttribute(name string) bool {
	return strings.HasPrefix(name, "hx-")
}

// ValidateHTMXAttribute checks if an htmx attribute is valid for the given version.
// Returns (valid, deprecated, v4Only).
// - valid: attribute is recognized as an htmx attribute for this version
// - deprecated: attribute is deprecated in this version (v4 only)
// - v4Only: attribute is only available in v4 (when using v2)
func ValidateHTMXAttribute(name, version string) (valid, deprecated, v4Only bool) {
	if !IsHTMXAttribute(name) {
		return false, false, false
	}

	// Handle hx-on:* event attributes (valid in both versions)
	if strings.HasPrefix(name, "hx-on:") || strings.HasPrefix(name, "hx-on-") {
		return true, false, false
	}

	// Handle hx-status:* pattern (v4 only)
	if strings.HasPrefix(name, "hx-status:") || strings.HasPrefix(name, "hx-status-") {
		if version == "4" {
			return true, false, false
		}
		return false, false, true
	}

	// Handle :inherited and :append suffixes (v4 only, but the base attribute should be valid)
	baseName := name
	switch {
	case strings.HasSuffix(name, ":inherited:append"):
		baseName = strings.TrimSuffix(name, ":inherited:append")
		if version != "4" {
			if HTMXv2Attributes[baseName] {
				return false, false, true // :inherited:append is v4-only
			}
			return false, false, false
		}
	case strings.HasSuffix(name, ":inherited"):
		baseName = strings.TrimSuffix(name, ":inherited")
		if version != "4" {
			// Check if base attribute is valid for v2
			if HTMXv2Attributes[baseName] {
				return false, false, true // :inherited is v4-only
			}
			return false, false, false
		}
	case strings.HasSuffix(name, ":append"):
		baseName = strings.TrimSuffix(name, ":append")
		if version != "4" {
			if HTMXv2Attributes[baseName] {
				return false, false, true // :append is v4-only
			}
			return false, false, false
		}
	}

	if version == "4" {
		if HTMXv4DeprecatedAttributes[baseName] {
			return true, true, false // Valid but deprecated
		}
		if HTMXv4Attributes[baseName] {
			return true, false, false
		}
		// Unknown htmx attribute
		return false, false, false
	}

	// Version 2 (default)
	if HTMXv4OnlyAttributes[baseName] {
		return false, false, true // v4-only attribute
	}
	if HTMXv2Attributes[baseName] {
		return true, false, false
	}

	// Unknown htmx attribute
	return false, false, false
}
