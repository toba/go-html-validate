package linter_test

import (
	"testing"

	"github.com/STR-Consulting/go-html-validate/linter"
	"github.com/STR-Consulting/go-html-validate/rules"
)

func TestLintContent_GoTemplates(t *testing.T) {
	tests := []struct {
		name string
		html string
	}{
		{
			name: "template variable in alt",
			html: `<img src="test.jpg" alt="{{ .Title }}">`,
		},
		{
			name: "template conditional",
			html: `<button type="button">{{ if .Label }}{{ .Label }}{{ else }}Click{{ end }}</button>`,
		},
		{
			name: "template range",
			html: `{{ range .Items }}<li>{{ .Name }}</li>{{ end }}`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic or error on template syntax
			_, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
		})
	}
}

func TestLintContent_TemplateConditionalAttributes(t *testing.T) {
	// Tests that template conditionals with duplicate attributes are handled correctly.
	// The preprocessor should keep only the if-branch content.
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "conditional placeholder in if-else",
			html: `<input type="search"
				{{if .User.IsAdmin}}
				placeholder="Admin search..."
				{{else}}
				placeholder="User search..."
				{{end}}
			>`,
			// Should not flag no-dup-attr because preprocessor keeps only if-branch
		},
		{
			name: "conditional class in if-else",
			html: `<div
				{{if .Active}}
				class="active"
				{{else}}
				class="inactive"
				{{end}}
			>Content</div>`,
		},
		{
			name: "if without else keeps content",
			html: `<div {{if .ShowClass}}class="visible"{{end}}>Content</div>`,
		},
		{
			name:     "actual duplicate attribute (not template)",
			html:     `<input type="text" name="foo" name="bar">`,
			wantRule: rules.RuleNoDupAttr,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoDupAttr, tt.wantRule)
		})
	}
}

func TestLintContent_TemplateFragmentOrphanedElements(t *testing.T) {
	// Tests that template fragments (files starting with {{define) allow orphaned
	// elements that would get ancestors from parent templates.
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "li in template fragment (allowed)",
			html: `{{define "list-items"}}
<li>Item 1</li>
<li>Item 2</li>
{{end}}`,
			// Should not flag element-required-ancestor because it's a fragment
		},
		{
			name: "td in template fragment (allowed)",
			html: `{{define "table-cells"}}
<td>Cell 1</td>
<td>Cell 2</td>
{{end}}`,
		},
		{
			name:     "li without list parent (not a fragment)",
			html:     `<div><li>orphan</li></div>`,
			wantRule: rules.RuleElementRequiredAncestor,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleElementRequiredAncestor, tt.wantRule)
		})
	}
}
