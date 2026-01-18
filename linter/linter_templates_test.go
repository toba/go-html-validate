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

func TestLintContent_TemplateSyntaxValid(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
		wantMsg  string
	}{
		// Balanced braces tests
		{
			name:     "unmatched opening brace",
			html:     `<div>{{ if .Show </div>`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name:     "unmatched closing brace",
			html:     `<div>.Show }}</div>`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name: "balanced braces - no error",
			html: `<div>{{ .Title }}</div>`,
		},

		// Control structure tests
		{
			name: "unclosed if",
			html: `{{ if .Show }}
content`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name: "unclosed range",
			html: `{{ range .Items }}
<li>{{ .Name }}</li>`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name: "unclosed with",
			html: `{{ with .User }}
<span>{{ .Name }}</span>`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name:     "unexpected end",
			html:     `<div>{{ end }}</div>`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name:     "unexpected else",
			html:     `<div>{{ else }}</div>`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name: "balanced if/end - no error",
			html: `{{ if .Show }}content{{ end }}`,
		},
		{
			name: "balanced if/else/end - no error",
			html: `{{ if .Show }}yes{{ else }}no{{ end }}`,
		},
		{
			name: "nested control structures - balanced",
			html: `{{ if .Show }}{{ range .Items }}{{ .Name }}{{ end }}{{ end }}`,
		},
		{
			name: "unclosed block",
			html: `{{ block "content" . }}
default`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name: "unclosed define",
			html: `{{ define "partial" }}
content`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name: "balanced define/end - no error",
			html: `{{ define "partial" }}content{{ end }}`,
		},

		// Trim marker syntax tests
		{
			name:     "trim marker no space after",
			html:     `<div>{{-if .Show }}</div>`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name:     "trim marker no space before",
			html:     `<div>{{ if .Show-}}</div>`,
			wantRule: rules.RuleTemplateSyntaxValid,
		},
		{
			name: "trim markers with proper spacing - no error",
			html: `<div>{{- if .Show -}}content{{- end -}}</div>`,
		},
		{
			name: "standard trim usage - no error",
			html: `{{ if .Show -}}
content
{{- end }}`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleTemplateSyntaxValid, tt.wantRule)
		})
	}
}

func TestLintContent_TemplateWhitespaceTrim(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
		wantMsg  string
	}{
		{
			name: "control flow alone on line - missing trailing trim",
			html: `<div>
{{ if .Show }}
content
{{ end }}
</div>`,
			wantRule: rules.RuleTemplateWhitespaceTrim,
		},
		{
			name: "has trailing trim marker - no warning",
			html: `<div>
{{ if .Show -}}
content
{{ end -}}
</div>`,
		},
		{
			name: "has both trim markers - no warning",
			html: `<div>
{{- if .Show -}}
content
{{- end -}}
</div>`,
		},
		{
			name: "action with content on same line - no warning",
			html: `<div>{{ .Title }}</div>`,
		},
		{
			name: "inline conditional - no warning",
			html: `<div class="{{ if .Active }}active{{ end }}">content</div>`,
		},
		{
			name: "output expression alone on line - no warning",
			html: `<div>
{{ .Title }}
</div>`,
		},
		{
			name: "range without trim",
			html: `{{ range .Items }}
<li>{{ .Name }}</li>
{{ end }}`,
			wantRule: rules.RuleTemplateWhitespaceTrim,
		},
		{
			name: "with without trim",
			html: `{{ with .User }}
<span>{{ .Name }}</span>
{{ end }}`,
			wantRule: rules.RuleTemplateWhitespaceTrim,
		},
		{
			name: "else without trim",
			html: `{{ if .Show }}
show
{{ else }}
hide
{{ end }}`,
			wantRule: rules.RuleTemplateWhitespaceTrim,
		},
		{
			name: "else if without trim",
			html: `{{ if .A }}
a
{{ else if .B }}
b
{{ end }}`,
			wantRule: rules.RuleTemplateWhitespaceTrim,
		},
		{
			name: "define without trim",
			html: `{{ define "foo" }}
content
{{ end }}`,
			wantRule: rules.RuleTemplateWhitespaceTrim,
		},
		{
			name: "template call without trim",
			html: `{{ template "header" . }}
<main>content</main>`,
			wantRule: rules.RuleTemplateWhitespaceTrim,
		},
		{
			name: "block without trim",
			html: `{{ block "content" . }}
default content
{{ end }}`,
			wantRule: rules.RuleTemplateWhitespaceTrim,
		},
		{
			name: "leading trim only - still warns for trailing",
			html: `<div>
{{- if .Show }}
content
{{- end }}
</div>`,
			wantRule: rules.RuleTemplateWhitespaceTrim,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleTemplateWhitespaceTrim, tt.wantRule)
		})
	}
}
