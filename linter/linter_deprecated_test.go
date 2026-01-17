package linter_test

import (
	"testing"

	"github.com/STR-Consulting/go-html-validate/linter"
	"github.com/STR-Consulting/go-html-validate/rules"
)

func TestLintContent_PreferAria(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
		severity rules.Severity
	}{
		{
			name:     "data-label should be aria-label",
			html:     `<button data-label="Close">X</button>`,
			wantRule: "prefer-aria",
			severity: rules.Warning,
		},
		{
			name:     "data-sort should be aria-sort",
			html:     `<table><tr><th data-sort="asc">Name</th></tr></table>`,
			wantRule: "prefer-aria",
			severity: rules.Warning,
		},
		{
			name: "data-custom is fine",
			html: `<div data-custom="value">Content</div>`,
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
					if r.Rule == "prefer-aria" {
						t.Errorf("expected no prefer-aria results, got %v", results)
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

func TestLintContent_Deprecated(t *testing.T) {
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
			name:     "deprecated marquee",
			html:     `<marquee>scrolling</marquee>`,
			wantRule: rules.RuleDeprecated,
		},
		{
			name:     "deprecated center",
			html:     `<center>centered</center>`,
			wantRule: rules.RuleDeprecated,
		},
		{
			name:     "deprecated font",
			html:     `<font color="red">text</font>`,
			wantRule: rules.RuleDeprecated,
		},
		{
			name:     "deprecated blink",
			html:     `<blink>flashing</blink>`,
			wantRule: rules.RuleDeprecated,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleDeprecated, tt.wantRule)
		})
	}
}

func TestLintContent_NoDeprecatedAttr(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "valid attributes",
			html: `<div class="foo">content</div>`,
		},
		{
			name:     "deprecated bgcolor on table",
			html:     `<table bgcolor="white"><tr><td>cell</td></tr></table>`,
			wantRule: rules.RuleNoDeprecatedAttr,
		},
		{
			name:     "deprecated align on div",
			html:     `<div align="center">content</div>`,
			wantRule: rules.RuleNoDeprecatedAttr,
		},
		{
			name: "width on img is not deprecated",
			html: `<img src="test.jpg" alt="test" width="100">`,
		},
		{
			name:     "border on table is deprecated",
			html:     `<table border="1"><tr><td>cell</td></tr></table>`,
			wantRule: rules.RuleNoDeprecatedAttr,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoDeprecatedAttr, tt.wantRule)
		})
	}
}

func TestLintContent_NoConditionalComment(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "regular comment",
			html: `<!-- This is a regular comment -->`,
		},
		{
			name:     "IE conditional comment",
			html:     `<!--[if IE]><p>IE only</p><![endif]-->`,
			wantRule: rules.RuleNoConditionalComment,
		},
		{
			name:     "IE version conditional",
			html:     `<!--[if lt IE 9]><script src="html5shiv.js"></script><![endif]-->`,
			wantRule: rules.RuleNoConditionalComment,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoConditionalComment, tt.wantRule)
		})
	}
}

func TestLintContent_NoRedundantFor(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "label with for pointing elsewhere",
			html: `<label for="other">Name</label><input id="name">`,
		},
		{
			name:     "label with for wrapping matching input",
			html:     `<label for="name"><input id="name"></label>`,
			wantRule: rules.RuleNoRedundantFor,
		},
		{
			name: "label without for wrapping input",
			html: `<label><input id="name"></label>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoRedundantFor, tt.wantRule)
		})
	}
}
