package linter_test

import (
	"testing"

	"github.com/STR-Consulting/go-html-validate/linter"
	"github.com/STR-Consulting/go-html-validate/rules"
)

func TestLintContent_FormSubmit(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "form without submit",
			html:     `<form><input type="text" aria-label="Name"></form>`,
			wantRule: "form-submit",
		},
		{
			name: "form with submit button",
			html: `<form><input type="text" aria-label="Name"><button type="submit">Submit</button></form>`,
		},
		{
			name: "form with input submit",
			html: `<form><input type="text" aria-label="Name"><input type="submit" value="Go"></form>`,
		},
		{
			name: "form with default button (implicit submit)",
			html: `<form><input type="text" aria-label="Name"><button>Submit</button></form>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleFormSubmit, tt.wantRule)
		})
	}
}

func TestLintContent_PreferButton(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "input type=button",
			html:     `<input type="button" value="Click">`,
			wantRule: "prefer-button",
		},
		{
			name:     "input type=submit",
			html:     `<input type="submit" value="Submit">`,
			wantRule: "prefer-button",
		},
		{
			name: "proper button element",
			html: `<button type="button">Click</button>`,
		},
		{
			name: "text input (no flag)",
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
			checkRule(t, results, rules.RulePreferButton, tt.wantRule)
		})
	}
}

func TestLintContent_NoMultipleMain(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "two visible main elements",
			html:     `<main>Content 1</main><main>Content 2</main>`,
			wantRule: "no-multiple-main",
		},
		{
			name: "single main element",
			html: `<main>Content</main>`,
		},
		{
			name: "two main with one hidden",
			html: `<main>Content 1</main><main hidden>Content 2</main>`,
		},
		{
			name: "no main elements",
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
			checkRule(t, results, rules.RuleNoMultipleMain, tt.wantRule)
		})
	}
}

func TestLintContent_MultipleLabeledControls(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "label wrapping multiple inputs",
			html:     `<label>Name <input type="text"><input type="text"></label>`,
			wantRule: "multiple-labeled-controls",
		},
		{
			name:     "label with for and different wrapped control",
			html:     `<label for="other">Name <input type="text"></label><input id="other" type="text">`,
			wantRule: "multiple-labeled-controls",
		},
		{
			name: "label with single input",
			html: `<label>Name <input type="text"></label>`,
		},
		{
			name: "label with for attribute only",
			html: `<label for="name">Name</label><input id="name" type="text">`,
		},
		{
			name: "label with for and same wrapped control (not flagged)",
			html: `<label for="name">Name <input id="name" type="text"></label>`,
		},
		{
			name: "empty label (not flagged)",
			html: `<label></label>`,
		},
		{
			name: "label with hidden input and visible control (not flagged)",
			html: `<label>Name <input type="hidden"><input type="checkbox"></label>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleMultipleLabeledControls, tt.wantRule)
		})
	}
}

func TestLintContent_VoidContent(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "void element without content",
			html: `<br>`,
		},
		{
			name: "img without content",
			html: `<img src="test.jpg" alt="test">`,
		},
		// Note: HTML parser normalizes invalid void element content,
		// so we can't easily test malformed cases like <br>text</br>
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleVoidContent, tt.wantRule)
		})
	}
}

func TestLintContent_TextContent(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "summary without text",
			html:     `<details><summary></summary>Content</details>`,
			wantRule: rules.RuleTextContent,
		},
		{
			name: "summary with text",
			html: `<details><summary>Click to expand</summary>Content</details>`,
		},
		{
			name: "summary with aria-label",
			html: `<details><summary aria-label="Expand section"></summary>Content</details>`,
		},
		{
			name: "details without summary (uses default)",
			html: `<details>Some content</details>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleTextContent, tt.wantRule)
		})
	}
}

func TestLintContent_WcagH67(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "decorative img with title",
			html:     `<img src="decorative.png" alt="" title="Should not have title">`,
			wantRule: rules.RuleWcagH67,
		},
		{
			name: "decorative img without title",
			html: `<img src="decorative.png" alt="">`,
		},
		{
			name: "non-decorative img with title",
			html: `<img src="photo.jpg" alt="A photo" title="More info">`,
		},
		{
			name:     "decorative img with role=img",
			html:     `<img src="decorative.png" alt="" role="img">`,
			wantRule: rules.RuleWcagH67,
		},
		{
			name: "decorative img with role=presentation",
			html: `<img src="decorative.png" alt="" role="presentation">`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleWcagH67, tt.wantRule)
		})
	}
}

func TestLintContent_WcagH71(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "fieldset without legend",
			html:     `<fieldset><input type="text"></fieldset>`,
			wantRule: rules.RuleWcagH71,
		},
		{
			name: "fieldset with legend",
			html: `<fieldset><legend>Contact Info</legend><input type="text"></fieldset>`,
		},
		{
			name: "fieldset with legend not first",
			html: `<fieldset><input type="text"><legend>Contact Info</legend></fieldset>`,
		},
		{
			name:     "empty fieldset",
			html:     `<fieldset></fieldset>`,
			wantRule: rules.RuleWcagH71,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleWcagH71, tt.wantRule)
		})
	}
}

func TestLintContent_ElementRequiredAncestor(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "li inside ul",
			html: `<ul><li>item</li></ul>`,
		},
		{
			name: "td inside table",
			html: `<table><tr><td>cell</td></tr></table>`,
		},
		{
			name:     "li without list parent",
			html:     `<div><li>orphan</li></div>`,
			wantRule: rules.RuleElementRequiredAncestor,
		},
		// Note: HTML parser auto-creates table structure for <td>,
		// so orphan td test not reliable
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

func TestLintContent_ElementPermittedContent(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "ul with li children",
			html: `<ul><li>item</li></ul>`,
		},
		{
			name:     "ul with div child",
			html:     `<ul><div>not allowed</div></ul>`,
			wantRule: rules.RuleElementPermittedContent,
		},
		{
			name: "dl with dt and dd children",
			html: `<dl><dt>term</dt><dd>definition</dd></dl>`,
		},
		{
			name:     "dl with p child",
			html:     `<dl><p>not allowed</p></dl>`,
			wantRule: rules.RuleElementPermittedContent,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleElementPermittedContent, tt.wantRule)
		})
	}
}

func TestLintContent_ElementPermittedOrder(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "summary first in details",
			html: `<details><summary>Title</summary><p>Content</p></details>`,
		},
		{
			name:     "summary not first",
			html:     `<details><p>Content</p><summary>Title</summary></details>`,
			wantRule: rules.RuleElementPermittedOrder,
		},
		{
			name: "legend first in fieldset",
			html: `<fieldset><legend>Title</legend><input></fieldset>`,
		},
		{
			name:     "legend not first",
			html:     `<fieldset><input><legend>Title</legend></fieldset>`,
			wantRule: rules.RuleElementPermittedOrder,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleElementPermittedOrder, tt.wantRule)
		})
	}
}

func TestLintContent_DetailsAndFieldsetContent(t *testing.T) {
	// Tests that details and fieldset allow flow content after summary/legend.
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "details with div after summary",
			html: `<details>
				<summary>Click to expand</summary>
				<div>Content goes here</div>
			</details>`,
		},
		{
			name: "details with ul after summary",
			html: `<details>
				<summary>Show list</summary>
				<ul><li>Item</li></ul>
			</details>`,
		},
		{
			name: "details with button after summary",
			html: `<details>
				<summary>Actions</summary>
				<button type="button">Action 1</button>
			</details>`,
		},
		{
			name: "fieldset with div after legend",
			html: `<fieldset>
				<legend>Contact Info</legend>
				<div class="form-group"><input type="text"></div>
			</fieldset>`,
		},
		{
			name: "fieldset with multiple form controls",
			html: `<fieldset>
				<legend>Options</legend>
				<div><input type="checkbox" id="opt1"><label for="opt1">Option 1</label></div>
				<div><input type="checkbox" id="opt2"><label for="opt2">Option 2</label></div>
			</fieldset>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			// Should not flag element-permitted-content for flow content
			checkRule(t, results, rules.RuleElementPermittedContent, tt.wantRule)
		})
	}
}
