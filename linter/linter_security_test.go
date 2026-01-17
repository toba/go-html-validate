package linter_test

import (
	"testing"

	"github.com/STR-Consulting/go-html-validate/linter"
	"github.com/STR-Consulting/go-html-validate/rules"
)

func TestLintContent_RequireSRI(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "external script without integrity",
			html:     `<script src="https://cdn.example.com/lib.js"></script>`,
			wantRule: "require-sri",
		},
		{
			name:     "external stylesheet without integrity",
			html:     `<link rel="stylesheet" href="https://cdn.example.com/style.css">`,
			wantRule: "require-sri",
		},
		{
			name:     "protocol-relative script without integrity",
			html:     `<script src="//cdn.example.com/lib.js"></script>`,
			wantRule: "require-sri",
		},
		{
			name:     "preload stylesheet without integrity",
			html:     `<link rel="preload" as="style" href="https://cdn.example.com/style.css">`,
			wantRule: "require-sri",
		},
		{
			name:     "modulepreload without integrity",
			html:     `<link rel="modulepreload" href="https://cdn.example.com/module.js">`,
			wantRule: "require-sri",
		},
		{
			name: "external script with integrity",
			html: `<script src="https://cdn.example.com/lib.js" integrity="sha384-abc123"></script>`,
		},
		{
			name: "external stylesheet with integrity",
			html: `<link rel="stylesheet" href="https://cdn.example.com/style.css" integrity="sha384-abc123">`,
		},
		{
			name: "local script (no flag)",
			html: `<script src="/js/app.js"></script>`,
		},
		{
			name: "relative script (no flag)",
			html: `<script src="./app.js"></script>`,
		},
		{
			name: "inline script (no flag)",
			html: `<script>console.log('hi');</script>`,
		},
		{
			name: "link rel=icon (no flag)",
			html: `<link rel="icon" href="https://cdn.example.com/favicon.ico">`,
		},
		{
			name: "link without href (no flag)",
			html: `<link rel="stylesheet">`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleRequireSRI, tt.wantRule)
		})
	}
}

func TestLintContent_ScriptElement(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "script with src",
			html: `<script src="app.js"></script>`,
		},
		{
			name: "inline script",
			html: `<script>console.log('hi')</script>`,
		},
		{
			name:     "async without src",
			html:     `<script async>console.log('hi')</script>`,
			wantRule: rules.RuleScriptElement,
		},
		{
			name:     "defer without src",
			html:     `<script defer>console.log('hi')</script>`,
			wantRule: rules.RuleScriptElement,
		},
		{
			name: "module script with async",
			html: `<script type="module" async src="app.mjs"></script>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleScriptElement, tt.wantRule)
		})
	}
}

func TestLintContent_ScriptType(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "script without type",
			html: `<script>console.log('hi')</script>`,
		},
		{
			name: "script with text/javascript",
			html: `<script type="text/javascript">console.log('hi')</script>`,
		},
		{
			name: "script with module",
			html: `<script type="module">import x from './x.js'</script>`,
		},
		{
			name: "script with importmap",
			html: `<script type="importmap">{}</script>`,
		},
		{
			name:     "script with invalid type",
			html:     `<script type="text/typescript">const x: number = 1</script>`,
			wantRule: rules.RuleScriptType,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleScriptType, tt.wantRule)
		})
	}
}

func TestLintContent_ValidAutocomplete(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "valid autocomplete token",
			html: `<input autocomplete="name">`,
		},
		{
			name: "multiple valid tokens",
			html: `<input autocomplete="shipping street-address">`,
		},
		{
			name: "section prefix",
			html: `<input autocomplete="section-billing name">`,
		},
		{
			name:     "invalid token",
			html:     `<input autocomplete="invalid-token">`,
			wantRule: rules.RuleValidAutocomplete,
		},
		{
			name: "on/off values",
			html: `<input autocomplete="off">`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleValidAutocomplete, tt.wantRule)
		})
	}
}

func TestLintContent_AreaAlt(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "area without alt",
			html:     `<map name="test"><area href="/link" shape="rect" coords="0,0,100,100"></map>`,
			wantRule: rules.RuleAreaAlt,
		},
		{
			name: "area with alt",
			html: `<map name="test"><area href="/link" alt="Link description" shape="rect" coords="0,0,100,100"></map>`,
		},
		{
			name: "area without href (no alt needed)",
			html: `<map name="test"><area shape="rect" coords="0,0,100,100"></map>`,
		},
		{
			name:     "area with empty alt",
			html:     `<map name="test"><area href="/link" alt="" shape="rect" coords="0,0,100,100"></map>`,
			wantRule: rules.RuleAreaAlt, // Warning for empty alt
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleAreaAlt, tt.wantRule)
		})
	}
}
