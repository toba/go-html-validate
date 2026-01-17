package linter_test

import (
	"testing"

	"github.com/STR-Consulting/go-html-validate/linter"
	"github.com/STR-Consulting/go-html-validate/rules"
)

func TestLintContent_ValidID(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
		severity rules.Severity
	}{
		{
			name:     "empty id",
			html:     `<div id="">Content</div>`,
			wantRule: "valid-id",
			severity: rules.Error,
		},
		{
			name:     "id with space",
			html:     `<div id="foo bar">Content</div>`,
			wantRule: "valid-id",
			severity: rules.Error,
		},
		{
			name:     "id starts with digit",
			html:     `<div id="123abc">Content</div>`,
			wantRule: "valid-id",
			severity: rules.Warning,
		},
		{
			name: "valid id with hyphen",
			html: `<div id="my-id">Content</div>`,
		},
		{
			name: "valid id with underscore",
			html: `<div id="my_id">Content</div>`,
		},
		{
			name: "no id attribute",
			html: `<div>Content</div>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			if tt.wantRule == "" {
				for _, r := range results {
					if r.Rule == "valid-id" {
						t.Errorf("expected no valid-id results, got %v", results)
					}
				}
				return
			}
			found := false
			for _, r := range results {
				if r.Rule == tt.wantRule {
					found = true
					if r.Severity != tt.severity {
						t.Errorf("expected severity %v, got %v", tt.severity, r.Severity)
					}
					break
				}
			}
			if !found {
				t.Errorf("expected rule %q in results, got %v", tt.wantRule, results)
			}
		})
	}
}

func TestLintContent_RequireLang(t *testing.T) {
	// Note: LintContent uses ParseFragment which doesn't preserve html element structure.
	// The require-lang rule is tested via LintFile in integration tests for full documents.
	// These tests verify the rule doesn't flag fragments.
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "fragment without html element (no flag)",
			html: `<div>Content</div>`,
		},
		{
			name: "fragment with main content (no flag)",
			html: `<main>Content</main>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleRequireLang, tt.wantRule)
		})
	}
}

func TestLintContent_ElementName(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "valid element",
			html: `<div>content</div>`,
		},
		{
			name: "valid custom element",
			html: `<my-component>content</my-component>`,
		},
		{
			name:     "unknown element",
			html:     `<foobar>content</foobar>`,
			wantRule: rules.RuleElementName,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleElementName, tt.wantRule)
		})
	}
}

func TestLintContent_AttributeAllowedValues(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "valid input type",
			html: `<input type="text">`,
		},
		{
			name:     "invalid input type",
			html:     `<input type="foobar">`,
			wantRule: rules.RuleAttributeAllowedValues,
		},
		{
			name: "valid button type",
			html: `<button type="submit">Click</button>`,
		},
		{
			name:     "invalid button type",
			html:     `<button type="invalid">Click</button>`,
			wantRule: rules.RuleAttributeAllowedValues,
		},
		{
			name: "valid form method",
			html: `<form method="post"></form>`,
		},
		{
			name:     "invalid form method",
			html:     `<form method="put"></form>`,
			wantRule: rules.RuleAttributeAllowedValues,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleAttributeAllowedValues, tt.wantRule)
		})
	}
}

func TestLintContent_NoMissingReferences(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "for references existing id",
			html: `<label for="name">Name</label><input id="name">`,
		},
		{
			name:     "for references non-existent id",
			html:     `<label for="missing">Name</label><input id="name">`,
			wantRule: rules.RuleNoMissingReferences,
		},
		{
			name: "aria-labelledby references existing id",
			html: `<span id="label">Label</span><input aria-labelledby="label">`,
		},
		{
			name:     "aria-labelledby references non-existent id",
			html:     `<input aria-labelledby="missing">`,
			wantRule: rules.RuleNoMissingReferences,
		},
		{
			name: "template expression in for (skip)",
			html: `<label for="TMPL">Name</label>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoMissingReferences, tt.wantRule)
		})
	}
}

func TestLintContent_FormDupName(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "unique names",
			html: `<form><input name="a"><input name="b"></form>`,
		},
		{
			name:     "duplicate names",
			html:     `<form><input name="a"><input name="a"></form>`,
			wantRule: rules.RuleFormDupName,
		},
		{
			name: "radio buttons can share names",
			html: `<form><input type="radio" name="choice"><input type="radio" name="choice"></form>`,
		},
		{
			name: "checkboxes can share names",
			html: `<form><input type="checkbox" name="opts"><input type="checkbox" name="opts"></form>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleFormDupName, tt.wantRule)
		})
	}
}

func TestLintContent_MapIDName(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "map with matching id and name",
			html: `<map id="nav" name="nav"></map>`,
		},
		{
			name:     "map with mismatched id and name",
			html:     `<map id="nav1" name="nav2"></map>`,
			wantRule: rules.RuleMapIDName,
		},
		{
			name: "map with name only",
			html: `<map name="nav"></map>`,
		},
		{
			name:     "map without name",
			html:     `<map id="nav"></map>`,
			wantRule: rules.RuleMapIDName,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleMapIDName, tt.wantRule)
		})
	}
}

func TestLintContent_NoDupClass(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "no class attr",
			html: `<p>text</p>`,
		},
		{
			name: "unique classes",
			html: `<p class="foo bar">text</p>`,
		},
		{
			name: "other attrs ok",
			html: `<p attr="foo bar foo">text</p>`,
		},
		{
			name:     "duplicate class",
			html:     `<p class="foo bar foo">text</p>`,
			wantRule: rules.RuleNoDupClass,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoDupClass, tt.wantRule)
		})
	}
}

func TestLintContent_AllowedLinks(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "valid http link",
			html: `<a href="https://example.com">Link</a>`,
		},
		{
			name: "valid relative link",
			html: `<a href="/page">Link</a>`,
		},
		{
			name: "valid anchor link",
			html: `<a href="#section">Link</a>`,
		},
		{
			name:     "javascript protocol",
			html:     `<a href="javascript:alert(1)">Link</a>`,
			wantRule: rules.RuleAllowedLinks,
		},
		{
			name:     "data protocol",
			html:     `<a href="data:text/html,<h1>Hi</h1>">Link</a>`,
			wantRule: rules.RuleAllowedLinks,
		},
		{
			name: "template expression (skip)",
			html: `<a href="TMPL">Link</a>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleAllowedLinks, tt.wantRule)
		})
	}
}

func TestLintContent_LongTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "title over 70 chars",
			html:     `<html><head><title>This is a very long title that exceeds the recommended seventy character limit for SEO</title></head></html>`,
			wantRule: "long-title",
		},
		{
			name: "title under 70 chars",
			html: `<html><head><title>Short Title</title></head></html>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleLongTitle, tt.wantRule)
		})
	}
}

func TestLintContent_NoInlineStyle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "element with inline style",
			html:     `<div style="color: red;">Red text</div>`,
			wantRule: "no-inline-style",
		},
		{
			name: "element without style",
			html: `<div class="red">Red text</div>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoInlineStyle, tt.wantRule)
		})
	}
}
