package linter_test

import (
	"testing"

	"github.com/STR-Consulting/go-html-validate/linter"
	"github.com/STR-Consulting/go-html-validate/rules"
)

func TestLintContent_ImgAlt(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
		wantErr  bool
	}{
		{
			name:     "img without alt",
			html:     `<img src="test.jpg">`,
			wantRule: "img-alt",
		},
		{
			name: "img with alt",
			html: `<img src="test.jpg" alt="A test image">`,
		},
		{
			name: "img with empty alt (decorative)",
			html: `<img src="test.jpg" alt="">`,
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
				if len(results) > 0 {
					t.Errorf("expected no results, got %d: %v", len(results), results)
				}
				return
			}

			found := false
			for _, r := range results {
				if r.Rule == tt.wantRule {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected rule %q in results, got %v", tt.wantRule, results)
			}
		})
	}
}

func TestLintContent_InputLabel(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "input without label",
			html:     `<input type="text" name="foo">`,
			wantRule: "input-label",
		},
		{
			name: "input with aria-label",
			html: `<input type="text" name="foo" aria-label="Foo field">`,
		},
		{
			name: "input with label for",
			html: `<label for="foo">Foo</label><input type="text" id="foo" name="foo">`,
		},
		{
			name: "input inside label",
			html: `<label>Foo <input type="text" name="foo"></label>`,
		},
		{
			name: "hidden input (no label needed)",
			html: `<input type="hidden" name="csrf">`,
		},
		{
			name: "submit button (no label needed)",
			html: `<input type="submit" value="Submit">`,
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
					if r.Rule == "input-label" {
						t.Errorf("expected no input-label results, got %v", results)
					}
				}
				return
			}

			found := false
			for _, r := range results {
				if r.Rule == tt.wantRule {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected rule %q in results, got %v", tt.wantRule, results)
			}
		})
	}
}

func TestLintContent_HeadingLevel(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "h1 to h3 skips level",
			html:     `<h1>Title</h1><h3>Subtitle</h3>`,
			wantRule: "heading-level",
		},
		{
			name: "h1 to h2 is valid",
			html: `<h1>Title</h1><h2>Subtitle</h2>`,
		},
		{
			name:     "h2 to h4 skips level",
			html:     `<h2>Title</h2><h4>Subtitle</h4>`,
			wantRule: "heading-level",
		},
		{
			name: "can decrease levels freely",
			html: `<h1>Title</h1><h2>Sub</h2><h1>Another Title</h1>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleHeadingLevel, tt.wantRule)
		})
	}
}

func TestLintContent_ButtonType(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "button without type",
			html:     `<button>Click</button>`,
			wantRule: "button-type",
		},
		{
			name: "button with type submit",
			html: `<button type="submit">Submit</button>`,
		},
		{
			name: "button with type button",
			html: `<button type="button">Click</button>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleButtonType, tt.wantRule)
		})
	}
}

func TestLintContent_TabindexNoPositive(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "positive tabindex",
			html:     `<button type="button" tabindex="5">Click</button>`,
			wantRule: "tabindex-no-positive",
		},
		{
			name: "tabindex zero is fine",
			html: `<div tabindex="0">Focusable</div>`,
		},
		{
			name: "tabindex negative is fine",
			html: `<div tabindex="-1">Programmatic focus only</div>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleTabindexNoPositive, tt.wantRule)
		})
	}
}

func TestLintContent_HiddenFocusable(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "button inside aria-hidden",
			html:     `<div aria-hidden="true"><button type="button">Hidden</button></div>`,
			wantRule: "hidden-focusable",
		},
		{
			name:     "link inside aria-hidden",
			html:     `<div aria-hidden="true"><a href="/test">Link</a></div>`,
			wantRule: "hidden-focusable",
		},
		{
			name: "aria-hidden without focusable content",
			html: `<div aria-hidden="true"><span>Decorative</span></div>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleHiddenFocusable, tt.wantRule)
		})
	}
}

func TestLintContent_EmptyTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "empty title",
			html:     `<html><head><title></title></head></html>`,
			wantRule: "empty-title",
		},
		{
			name:     "whitespace-only title",
			html:     `<html><head><title>   </title></head></html>`,
			wantRule: "empty-title",
		},
		{
			name: "title with content",
			html: `<html><head><title>My Page</title></head></html>`,
		},
		{
			name: "title with child elements containing text",
			html: `<html><head><title><span>lorem ipsum</span></title></head></html>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleEmptyTitle, tt.wantRule)
		})
	}
}

func TestLintContent_NoRedundantRole(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "button with role=button",
			html:     `<button type="button" role="button">Click</button>`,
			wantRule: "no-redundant-role",
		},
		{
			name:     "nav with role=navigation",
			html:     `<nav role="navigation">Menu</nav>`,
			wantRule: "no-redundant-role",
		},
		{
			name:     "li with role=listitem",
			html:     `<ul><li role="listitem">Item</li></ul>`,
			wantRule: "no-redundant-role",
		},
		{
			name:     "a with href and role=link",
			html:     `<a href="/test" role="link">Link</a>`,
			wantRule: "no-redundant-role",
		},
		{
			name: "button without role",
			html: `<button type="button">Click</button>`,
		},
		{
			name: "div with role (not redundant)",
			html: `<div role="button">Click</div>`,
		},
		{
			name: "li with role=presentation (not redundant)",
			html: `<ul><li role="presentation">Item</li></ul>`,
		},
		{
			name: "a without href and role=link (not redundant)",
			html: `<a role="link">Link</a>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoRedundantRole, tt.wantRule)
		})
	}
}

func TestLintContent_NoAbstractRole(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "abstract role widget",
			html:     `<div role="widget">Content</div>`,
			wantRule: "no-abstract-role",
		},
		{
			name:     "abstract role command",
			html:     `<div role="command">Content</div>`,
			wantRule: "no-abstract-role",
		},
		{
			name:     "abstract role input",
			html:     `<div role="input">Content</div>`,
			wantRule: "no-abstract-role",
		},
		{
			name:     "abstract role landmark",
			html:     `<div role="landmark">Content</div>`,
			wantRule: "no-abstract-role",
		},
		{
			name:     "abstract role structure",
			html:     `<div role="structure">Content</div>`,
			wantRule: "no-abstract-role",
		},
		{
			name:     "multiple roles with abstract",
			html:     `<div role="window none widget">Content</div>`,
			wantRule: "no-abstract-role",
		},
		{
			name: "concrete role button",
			html: `<div role="button">Content</div>`,
		},
		{
			name: "concrete role grid",
			html: `<div role="grid">Content</div>`,
		},
		{
			name: "non-role attribute with abstract value",
			html: `<div foo="command">Content</div>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoAbstractRole, tt.wantRule)
		})
	}
}

func TestLintContent_AriaLabelMisuse(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "aria-label on div",
			html:     `<div aria-label="Label">Content</div>`,
			wantRule: "aria-label-misuse",
		},
		{
			name:     "aria-label on span",
			html:     `<span aria-label="Label">Content</span>`,
			wantRule: "aria-label-misuse",
		},
		{
			name:     "aria-label on p",
			html:     `<p aria-label="foobar">Content</p>`,
			wantRule: "aria-label-misuse",
		},
		{
			name:     "aria-labelledby on span",
			html:     `<span aria-labelledby="foo">Content</span>`,
			wantRule: "aria-label-misuse",
		},
		{
			name: "aria-label on button (allowed)",
			html: `<button type="button" aria-label="Close">X</button>`,
		},
		{
			name: "aria-label on input (allowed)",
			html: `<input type="text" aria-label="Search">`,
		},
		{
			name: "aria-label on nav (allowed)",
			html: `<nav aria-label="Main navigation">Menu</nav>`,
		},
		{
			name: "aria-label on main (allowed)",
			html: `<main aria-label="Primary content">Content</main>`,
		},
		{
			name: "aria-label on div with role (allowed)",
			html: `<div role="button" aria-label="Close">X</div>`,
		},
		{
			name: "aria-label on div with tabindex (allowed)",
			html: `<div tabindex="0" aria-label="Focus me">Content</div>`,
		},
		{
			name: "aria-label on img (allowed)",
			html: `<img src="test.jpg" aria-label="Description">`,
		},
		{
			name: "aria-label on table (allowed)",
			html: `<table aria-label="Data table"><tr><td>Cell</td></tr></table>`,
		},
		{
			name: "empty aria-label (no flag)",
			html: `<p aria-label="">Content</p>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleAriaLabelMisuse, tt.wantRule)
		})
	}
}

func TestLintContent_UniqueLandmark(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "duplicate nav without labels",
			html:     `<nav>Menu 1</nav><nav>Menu 2</nav>`,
			wantRule: "unique-landmark",
		},
		{
			name:     "duplicate aside without labels",
			html:     `<aside>Sidebar 1</aside><aside>Sidebar 2</aside>`,
			wantRule: "unique-landmark",
		},
		{
			name:     "duplicate nav with same labels",
			html:     `<nav aria-label="Menu">Menu 1</nav><nav aria-label="Menu">Menu 2</nav>`,
			wantRule: "unique-landmark",
		},
		{
			name:     "duplicate nav with empty labels",
			html:     `<nav aria-label="">Menu 1</nav><nav aria-label="">Menu 2</nav>`,
			wantRule: "unique-landmark",
		},
		{
			name: "duplicate nav with unique labels",
			html: `<nav aria-label="Primary">Menu 1</nav><nav aria-label="Secondary">Menu 2</nav>`,
		},
		{
			name: "single nav (no flag)",
			html: `<nav>Menu</nav>`,
		},
		{
			name: "no landmarks (no flag)",
			html: `<div>Content</div>`,
		},
		{
			name: "duplicate nav with role=presentation (no flag)",
			html: `<nav role="presentation">Menu 1</nav><nav role="presentation">Menu 2</nav>`,
		},
		{
			name: "multiple unnamed forms (no flag - forms without names aren't landmarks)",
			html: `<form><input></form><form><input></form>`,
		},
		{
			name: "multiple unnamed sections (no flag - sections without names aren't landmarks)",
			html: `<section>Content 1</section><section>Content 2</section>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleUniqueLandmark, tt.wantRule)
		})
	}
}

func TestLintContent_PreferNativeElement(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "div with role=main",
			html:     `<div role="main">Content</div>`,
			wantRule: "prefer-native-element",
		},
		{
			name:     "div with role=navigation",
			html:     `<div role="navigation">Menu</div>`,
			wantRule: "prefer-native-element",
		},
		{
			name:     "div with role=banner",
			html:     `<div role="banner">Header</div>`,
			wantRule: "prefer-native-element",
		},
		{
			name:     "div with role=contentinfo",
			html:     `<div role="contentinfo">Footer</div>`,
			wantRule: "prefer-native-element",
		},
		{
			name:     "span with role=article",
			html:     `<span role="article">Content</span>`,
			wantRule: "prefer-native-element",
		},
		{
			name: "native main element",
			html: `<main>Content</main>`,
		},
		{
			name: "native nav element",
			html: `<nav>Menu</nav>`,
		},
		{
			name: "div with role=button (no native mapping in this rule)",
			html: `<div role="button">Click</div>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RulePreferNativeElement, tt.wantRule)
		})
	}
}

func TestLintContent_TelNonBreaking(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "tel link with regular spaces",
			html:     `<a href="tel:+1-555-123-4567">555 123 4567</a>`,
			wantRule: rules.RuleTelNonBreaking,
		},
		{
			name: "tel link with non-breaking spaces",
			html: `<a href="tel:+1-555-123-4567">555` + "\u00A0" + `123` + "\u00A0" + `4567</a>`,
		},
		{
			name: "tel link with no spaces",
			html: `<a href="tel:+1-555-123-4567">5551234567</a>`,
		},
		{
			name: "regular link (not tel)",
			html: `<a href="https://example.com">555 123 4567</a>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleTelNonBreaking, tt.wantRule)
		})
	}
}

func TestLintContent_WcagH36(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "input image without alt",
			html:     `<input type="image" src="submit.png">`,
			wantRule: rules.RuleWcagH36,
		},
		{
			name: "input image with alt",
			html: `<input type="image" src="submit.png" alt="Submit form">`,
		},
		{
			name:     "input image with empty alt",
			html:     `<input type="image" src="submit.png" alt="">`,
			wantRule: rules.RuleWcagH36,
		},
		{
			name: "regular input (not image)",
			html: `<input type="text" name="foo">`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleWcagH36, tt.wantRule)
		})
	}
}

func TestLintContent_WcagH63(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "th without scope",
			html:     `<table><tr><th>Header</th></tr></table>`,
			wantRule: rules.RuleWcagH63,
		},
		{
			name: "th with scope col",
			html: `<table><tr><th scope="col">Header</th></tr></table>`,
		},
		{
			name: "th with scope row",
			html: `<table><tr><th scope="row">Header</th></tr></table>`,
		},
		{
			name:     "th with invalid scope",
			html:     `<table><tr><th scope="invalid">Header</th></tr></table>`,
			wantRule: rules.RuleWcagH63,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleWcagH63, tt.wantRule)
		})
	}
}
