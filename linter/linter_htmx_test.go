package linter_test

import (
	"strings"
	"testing"

	"github.com/toba/go-html-validate/linter"
	"github.com/toba/go-html-validate/rules"
)

// htmxTestCase defines a test case for htmx attribute validation.
type htmxTestCase struct {
	name       string
	html       string
	wantRule   string
	wantSubstr string
	severity   rules.Severity
}

// runHTMXTests runs a set of htmx test cases with the given linter.
func runHTMXTests(t *testing.T, l *linter.Linter, tests []htmxTestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}

			if tt.wantRule == "" {
				for _, r := range results {
					if r.Rule == rules.RuleHTMXAttributes {
						t.Errorf("expected no htmx-attributes results, got %v", results)
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
					if tt.wantSubstr != "" && !strings.Contains(r.Message, tt.wantSubstr) {
						t.Errorf("expected message containing %q, got %q", tt.wantSubstr, r.Message)
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

func TestLintContent_HTMXSwap(t *testing.T) {
	tests := []htmxTestCase{
		// Valid cases
		{
			name: "valid innerHTML",
			html: `<div hx-get="/api" hx-swap="innerHTML">content</div>`,
		},
		{
			name: "valid outerHTML",
			html: `<div hx-get="/api" hx-swap="outerHTML">content</div>`,
		},
		{
			name: "valid beforebegin",
			html: `<div hx-get="/api" hx-swap="beforebegin">content</div>`,
		},
		{
			name: "valid afterend",
			html: `<div hx-get="/api" hx-swap="afterend">content</div>`,
		},
		{
			name: "valid delete",
			html: `<div hx-get="/api" hx-swap="delete">content</div>`,
		},
		{
			name: "valid none",
			html: `<div hx-get="/api" hx-swap="none">content</div>`,
		},
		{
			name: "valid with swap modifier",
			html: `<div hx-get="/api" hx-swap="innerHTML swap:1s">content</div>`,
		},
		{
			name: "valid with settle modifier",
			html: `<div hx-get="/api" hx-swap="innerHTML settle:500ms">content</div>`,
		},
		{
			name: "valid with scroll modifier",
			html: `<div hx-get="/api" hx-swap="innerHTML scroll:top">content</div>`,
		},
		{
			name: "valid with show modifier",
			html: `<div hx-get="/api" hx-swap="innerHTML show:bottom">content</div>`,
		},
		{
			name: "valid with focus-scroll modifier",
			html: `<div hx-get="/api" hx-swap="innerHTML focus-scroll:true">content</div>`,
		},
		{
			name: "valid with multiple modifiers",
			html: `<div hx-get="/api" hx-swap="innerHTML swap:1s settle:500ms scroll:top">content</div>`,
		},
		{
			name: "empty swap (uses default)",
			html: `<div hx-get="/api" hx-swap="">content</div>`,
		},
		{
			name: "template expression in swap",
			html: `<div hx-get="/api" hx-swap="{{ .SwapMode }}">content</div>`,
		},
		// Invalid cases
		{
			name:       "invalid swap value",
			html:       `<div hx-get="/api" hx-swap="invalid">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "invalid hx-swap value",
			severity:   rules.Error,
		},
		{
			name:       "invalid modifier format",
			html:       `<div hx-get="/api" hx-swap="innerHTML notamodifier">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "missing colon",
			severity:   rules.Error,
		},
		{
			name:       "invalid swap time",
			html:       `<div hx-get="/api" hx-swap="innerHTML swap:nottime">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "time value",
			severity:   rules.Error,
		},
		{
			name:       "invalid focus-scroll value",
			html:       `<div hx-get="/api" hx-swap="innerHTML focus-scroll:maybe">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "focus-scroll",
			severity:   rules.Error,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXSwapV4Only(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		version  string
		wantRule string
	}{
		{
			name:    "textContent valid in v4",
			html:    `<div hx-get="/api" hx-swap="textContent">content</div>`,
			version: "4",
		},
		{
			name:    "upsert valid in v4",
			html:    `<div hx-get="/api" hx-swap="upsert">content</div>`,
			version: "4",
		},
		{
			name:     "textContent warns in v2",
			html:     `<div hx-get="/api" hx-swap="textContent">content</div>`,
			version:  "2",
			wantRule: rules.RuleHTMXAttributes,
		},
		{
			name:     "upsert warns in v2",
			html:     `<div hx-get="/api" hx-swap="upsert">content</div>`,
			version:  "2",
			wantRule: rules.RuleHTMXAttributes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := linter.DefaultConfig()
			cfg.Frameworks.HTMX = true
			cfg.Frameworks.HTMXVersion = tt.version
			l := linter.New(cfg)

			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}

			checkRule(t, results, rules.RuleHTMXAttributes, tt.wantRule)
		})
	}
}

func TestLintContent_HTMXTrigger(t *testing.T) {
	tests := []htmxTestCase{
		// Valid cases
		{
			name: "simple click",
			html: `<div hx-get="/api" hx-trigger="click">content</div>`,
		},
		{
			name: "click with once",
			html: `<div hx-get="/api" hx-trigger="click once">content</div>`,
		},
		{
			name: "click with changed",
			html: `<div hx-get="/api" hx-trigger="click changed">content</div>`,
		},
		{
			name: "click with delay",
			html: `<div hx-get="/api" hx-trigger="click delay:1s">content</div>`,
		},
		{
			name: "click with throttle",
			html: `<div hx-get="/api" hx-trigger="click throttle:500ms">content</div>`,
		},
		{
			name: "click with queue",
			html: `<div hx-get="/api" hx-trigger="click queue:last">content</div>`,
		},
		{
			name: "click with consume",
			html: `<div hx-get="/api" hx-trigger="click consume">content</div>`,
		},
		{
			name: "click with from",
			html: `<div hx-get="/api" hx-trigger="click from:body">content</div>`,
		},
		{
			name: "click with target",
			html: `<div hx-get="/api" hx-trigger="click target:#myid">content</div>`,
		},
		{
			name: "every polling",
			html: `<div hx-get="/api" hx-trigger="every 2s">content</div>`,
		},
		{
			name: "intersect",
			html: `<div hx-get="/api" hx-trigger="intersect">content</div>`,
		},
		{
			name: "intersect with threshold",
			html: `<div hx-get="/api" hx-trigger="intersect threshold:0.5">content</div>`,
		},
		{
			name: "revealed",
			html: `<div hx-get="/api" hx-trigger="revealed">content</div>`,
		},
		{
			name: "multiple triggers",
			html: `<div hx-get="/api" hx-trigger="click, keyup delay:500ms">content</div>`,
		},
		{
			name: "filter expression",
			html: `<div hx-get="/api" hx-trigger="click[ctrlKey]">content</div>`,
		},
		{
			name: "multiple modifiers",
			html: `<div hx-get="/api" hx-trigger="click once delay:1s">content</div>`,
		},
		{
			name: "template expression",
			html: `<div hx-get="/api" hx-trigger="{{ .Trigger }}">content</div>`,
		},
		// Invalid cases
		{
			name:       "invalid delay time",
			html:       `<div hx-get="/api" hx-trigger="click delay:nottime">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "time value",
			severity:   rules.Error,
		},
		{
			name:       "invalid throttle time",
			html:       `<div hx-get="/api" hx-trigger="click throttle:bad">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "time value",
			severity:   rules.Error,
		},
		{
			name:       "invalid queue mode",
			html:       `<div hx-get="/api" hx-trigger="click queue:invalid">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "queue mode",
			severity:   rules.Error,
		},
		{
			name:       "every without time",
			html:       `<div hx-get="/api" hx-trigger="every">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "time value",
			severity:   rules.Error,
		},
		{
			name:       "every with invalid time",
			html:       `<div hx-get="/api" hx-trigger="every bad">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "time value",
			severity:   rules.Error,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXTarget(t *testing.T) {
	tests := []struct {
		name       string
		html       string
		wantRule   string
		wantSubstr string
	}{
		// Valid cases
		{
			name: "this",
			html: `<div hx-get="/api" hx-target="this">content</div>`,
		},
		{
			name: "next",
			html: `<div hx-get="/api" hx-target="next">content</div>`,
		},
		{
			name: "previous",
			html: `<div hx-get="/api" hx-target="previous">content</div>`,
		},
		{
			name: "body",
			html: `<div hx-get="/api" hx-target="body">content</div>`,
		},
		{
			name: "closest selector",
			html: `<div hx-get="/api" hx-target="closest div">content</div>`,
		},
		{
			name: "find selector",
			html: `<div hx-get="/api" hx-target="find .result">content</div>`,
		},
		{
			name: "next with selector",
			html: `<div hx-get="/api" hx-target="next div">content</div>`,
		},
		{
			name: "previous with selector",
			html: `<div hx-get="/api" hx-target="previous .sibling">content</div>`,
		},
		{
			name: "id selector",
			html: `<div hx-get="/api" hx-target="#result">content</div>`,
		},
		{
			name: "class selector",
			html: `<div hx-get="/api" hx-target=".result">content</div>`,
		},
		{
			name: "element selector",
			html: `<div hx-get="/api" hx-target="div">content</div>`,
		},
		{
			name: "template expression",
			html: `<div hx-get="/api" hx-target="{{ .Target }}">content</div>`,
		},
		// Invalid cases
		{
			name:       "invalid keyword",
			html:       `<div hx-get="/api" hx-target="invalid selector">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "invalid hx-target keyword",
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}

			checkRule(t, results, rules.RuleHTMXAttributes, tt.wantRule)
		})
	}
}

func TestLintContent_HTMXDisabled(t *testing.T) {
	// When htmx is disabled, no htmx-attributes errors should be reported
	html := `<div hx-get="/api" hx-swap="invalid">content</div>`

	l := linter.New(nil) // nil config = htmx disabled

	results, err := l.LintContent("test.html", []byte(html))
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	for _, r := range results {
		if r.Rule == rules.RuleHTMXAttributes {
			t.Errorf("expected no htmx-attributes results when htmx disabled, got %v", results)
		}
	}
}

func TestLintContent_AttributeMisuseSkipsHTMX(t *testing.T) {
	// When htmx is enabled, attribute-misuse should not flag hx-* attributes
	html := `<div hx-get="/api" hx-swap="innerHTML">content</div>`

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	results, err := l.LintContent("test.html", []byte(html))
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	for _, r := range results {
		if r.Rule == rules.RuleAttributeMisuse {
			t.Errorf("expected no attribute-misuse results for hx-* when htmx enabled, got %v", results)
		}
	}
}

func TestLintContent_HTMXOnEvent(t *testing.T) {
	tests := []htmxTestCase{
		// Valid DOM events
		{
			name: "valid click event",
			html: `<button hx-on:click="alert('clicked')">Click</button>`,
		},
		{
			name: "valid submit event",
			html: `<form hx-on:submit="validate()">Form</form>`,
		},
		{
			name: "valid keydown event",
			html: `<input hx-on:keydown="handle()">`,
		},
		{
			name: "valid focus event",
			html: `<input hx-on:focus="highlight()">`,
		},
		// Valid htmx v2 events
		{
			name: "valid htmx:afterRequest",
			html: `<div hx-get="/api" hx-on:htmx:afterRequest="done()">content</div>`,
		},
		{
			name: "valid htmx:beforeSwap",
			html: `<div hx-get="/api" hx-on:htmx:beforeSwap="prep()">content</div>`,
		},
		{
			name: "valid htmx:configRequest",
			html: `<div hx-get="/api" hx-on:htmx:configRequest="config()">content</div>`,
		},
		{
			name: "valid htmx:load",
			html: `<div hx-get="/api" hx-on:htmx:load="loaded()">content</div>`,
		},
		// Valid deprecated hx-on- syntax
		{
			name: "valid hx-on- click",
			html: `<button hx-on-click="alert('clicked')">Click</button>`,
		},
		// Invalid cases
		{
			name:       "unknown htmx event",
			html:       `<div hx-get="/api" hx-on:htmx:unknownEvent="handle()">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "unknown htmx event",
			severity:   rules.Warning,
		},
		{
			name:       "unknown event",
			html:       `<div hx-on:customEvent="handle()">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "unknown event",
			severity:   rules.Warning,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXOnEventV4(t *testing.T) {
	tests := []htmxTestCase{
		// Valid htmx v4 events
		{
			name: "valid htmx:after:request",
			html: `<div hx-get="/api" hx-on:htmx:after:request="done()">content</div>`,
		},
		{
			name: "valid htmx:before:swap",
			html: `<div hx-get="/api" hx-on:htmx:before:swap="prep()">content</div>`,
		},
		{
			name: "valid htmx:error",
			html: `<div hx-get="/api" hx-on:htmx:error="handleError()">content</div>`,
		},
		{
			name: "valid htmx:finally:request",
			html: `<div hx-get="/api" hx-on:htmx:finally:request="cleanup()">content</div>`,
		},
		{
			name: "valid htmx:load (standalone)",
			html: `<div hx-get="/api" hx-on:htmx:load="loaded()">content</div>`,
		},
		// v4 shorthand syntax (hx-on::event)
		{
			name: "valid shorthand htmx:after:request",
			html: `<div hx-get="/api" hx-on::after:request="done()">content</div>`,
		},
		// Invalid cases
		{
			name:       "unknown v4 phase",
			html:       `<div hx-get="/api" hx-on:htmx:unknown:request="handle()">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "unknown htmx 4 event phase",
			severity:   rules.Warning,
		},
		{
			name:       "unknown v4 action",
			html:       `<div hx-get="/api" hx-on:htmx:after:unknownAction="handle()">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "unknown htmx 4 event action",
			severity:   rules.Warning,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "4"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXSubmitButtonInForm(t *testing.T) {
	tests := []htmxTestCase{
		// Valid cases - no warning needed
		{
			name: "hx-post on form element",
			html: `<form hx-post="/api"><button type="submit">Submit</button></form>`,
		},
		{
			name: "hx-post on button outside form",
			html: `<button hx-post="/api">Submit</button>`,
		},
		{
			name: "hx-post on button type=button inside form",
			html: `<form><button type="button" hx-post="/api">Click</button></form>`,
		},
		{
			name: "hx-post on div inside form",
			html: `<form><div hx-post="/api">Click</div></form>`,
		},
		{
			name: "submit button without htmx inside form",
			html: `<form action="/api"><button type="submit">Submit</button></form>`,
		},
		// Warning cases
		{
			name:       "hx-post on submit button inside form",
			html:       `<form><button type="submit" hx-post="/api">Submit</button></form>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "bypass form validation",
			severity:   rules.Warning,
		},
		{
			name:       "hx-post on default button inside form",
			html:       `<form><button hx-post="/api">Submit</button></form>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "bypass form validation",
			severity:   rules.Warning,
		},
		{
			name:       "hx-get on submit button inside form",
			html:       `<form><button type="submit" hx-get="/api">Submit</button></form>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "bypass form validation",
			severity:   rules.Warning,
		},
		{
			name:       "hx-post on input submit inside form",
			html:       `<form><input type="submit" hx-post="/api" value="Submit"></form>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "bypass form validation",
			severity:   rules.Warning,
		},
		{
			name:       "hx-delete on submit button inside nested form",
			html:       `<div><form><button type="submit" hx-delete="/api">Delete</button></form></div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "bypass form validation",
			severity:   rules.Warning,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXValsHeaders(t *testing.T) {
	tests := []htmxTestCase{
		// Valid hx-vals cases
		{
			name: "valid hx-vals JSON object",
			html: `<div hx-post="/api" hx-vals='{"key": "value"}'>content</div>`,
		},
		{
			name: "valid hx-vals with multiple keys",
			html: `<div hx-post="/api" hx-vals='{"name": "John", "age": 30}'>content</div>`,
		},
		{
			name: "valid hx-vals empty object",
			html: `<div hx-post="/api" hx-vals="{}">content</div>`,
		},
		{
			name: "valid hx-vals with js: prefix",
			html: `<div hx-post="/api" hx-vals="js:{name: getName()}">content</div>`,
		},
		{
			name: "valid hx-vals with javascript: prefix",
			html: `<div hx-post="/api" hx-vals="javascript:{name: getName()}">content</div>`,
		},
		{
			name: "valid hx-vals with template expression",
			html: `<div hx-post="/api" hx-vals='{"id": {{ .ID }}}'>content</div>`,
		},
		{
			name: "empty hx-vals",
			html: `<div hx-post="/api" hx-vals="">content</div>`,
		},
		// Valid hx-headers cases
		{
			name: "valid hx-headers JSON object",
			html: `<div hx-post="/api" hx-headers='{"X-Custom": "value"}'>content</div>`,
		},
		{
			name: "valid hx-headers with auth",
			html: `<div hx-post="/api" hx-headers='{"Authorization": "Bearer token"}'>content</div>`,
		},
		// Invalid hx-vals cases
		{
			name:       "invalid hx-vals missing quotes",
			html:       `<div hx-post="/api" hx-vals='{key: "value"}'>content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "invalid JSON",
			severity:   rules.Error,
		},
		{
			name:       "invalid hx-vals trailing comma",
			html:       `<div hx-post="/api" hx-vals='{"key": "value",}'>content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "invalid JSON",
			severity:   rules.Error,
		},
		{
			name:       "invalid hx-vals broken JSON",
			html:       `<div hx-post="/api" hx-vals='{"key":'>content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "invalid JSON",
			severity:   rules.Error,
		},
		// Invalid hx-headers cases
		{
			name:       "invalid hx-headers syntax",
			html:       `<div hx-post="/api" hx-headers='{X-Custom: value}'>content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "invalid JSON",
			severity:   rules.Error,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXInclude(t *testing.T) {
	tests := []htmxTestCase{
		// Valid hx-include cases
		{
			name: "valid id selector",
			html: `<div hx-post="/api" hx-include="#myform">content</div>`,
		},
		{
			name: "valid class selector",
			html: `<div hx-post="/api" hx-include=".inputs">content</div>`,
		},
		{
			name: "valid element selector",
			html: `<div hx-post="/api" hx-include="input">content</div>`,
		},
		{
			name: "valid attribute selector",
			html: `<div hx-post="/api" hx-include="[name]">content</div>`,
		},
		{
			name: "valid attribute value selector",
			html: `<div hx-post="/api" hx-include='[name="email"]'>content</div>`,
		},
		{
			name: "valid this keyword",
			html: `<div hx-post="/api" hx-include="this">content</div>`,
		},
		{
			name: "valid closest keyword",
			html: `<div hx-post="/api" hx-include="closest form">content</div>`,
		},
		{
			name: "valid find keyword",
			html: `<div hx-post="/api" hx-include="find input">content</div>`,
		},
		{
			name: "valid next keyword",
			html: `<div hx-post="/api" hx-include="next .sibling">content</div>`,
		},
		{
			name: "valid multiple selectors",
			html: `<div hx-post="/api" hx-include="#form1, #form2">content</div>`,
		},
		{
			name: "valid descendant combinator",
			html: `<div hx-post="/api" hx-include="form input">content</div>`,
		},
		{
			name: "valid child combinator",
			html: `<div hx-post="/api" hx-include="form > input">content</div>`,
		},
		{
			name: "valid universal selector",
			html: `<div hx-post="/api" hx-include="*">content</div>`,
		},
		{
			name: "valid pseudo-class",
			html: `<div hx-post="/api" hx-include=":not(.hidden)">content</div>`,
		},
		{
			name: "empty hx-include",
			html: `<div hx-post="/api" hx-include="">content</div>`,
		},
		{
			name: "valid template expression",
			html: `<div hx-post="/api" hx-include="{{ .Selector }}">content</div>`,
		},
		// Invalid hx-include cases
		{
			name:       "unclosed bracket",
			html:       `<div hx-post="/api" hx-include="[name">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "unclosed bracket",
			severity:   rules.Error,
		},
		{
			name:       "unbalanced brackets",
			html:       `<div hx-post="/api" hx-include="[name]]">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "unbalanced brackets",
			severity:   rules.Error,
		},
		{
			name:       "empty attribute selector",
			html:       `<div hx-post="/api" hx-include="input[]">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "empty attribute selector",
			severity:   rules.Error,
		},
		{
			name:       "invalid start character",
			html:       `<div hx-post="/api" hx-include="@invalid">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "invalid character",
			severity:   rules.Error,
		},
		{
			name:       "ends with combinator",
			html:       `<div hx-post="/api" hx-include="form >">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "ends with combinator",
			severity:   rules.Error,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXStatusV4(t *testing.T) {
	tests := []htmxTestCase{
		// Valid hx-status patterns (htmx 4)
		{
			name: "valid exact status 404",
			html: `<div hx-get="/api" hx-status:404="none">content</div>`,
		},
		{
			name: "valid exact status 200",
			html: `<div hx-get="/api" hx-status:200="innerHTML">content</div>`,
		},
		{
			name: "valid exact status 500",
			html: `<div hx-get="/api" hx-status:500="target:#error">content</div>`,
		},
		{
			name: "valid wildcard 2xx",
			html: `<div hx-get="/api" hx-status:2xx="innerHTML">content</div>`,
		},
		{
			name: "valid wildcard 4xx",
			html: `<div hx-get="/api" hx-status:4xx="none">content</div>`,
		},
		{
			name: "valid wildcard 5xx",
			html: `<div hx-get="/api" hx-status:5xx="target:#error">content</div>`,
		},
		// Invalid hx-status patterns
		{
			name:       "invalid status 99",
			html:       `<div hx-get="/api" hx-status:99="none">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "not a valid HTTP status",
			severity:   rules.Error,
		},
		{
			name:       "invalid status 600",
			html:       `<div hx-get="/api" hx-status:600="none">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "not a valid HTTP status",
			severity:   rules.Error,
		},
		{
			name:       "invalid wildcard 6xx",
			html:       `<div hx-get="/api" hx-status:6xx="none">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "not a valid HTTP status",
			severity:   rules.Error,
		},
		{
			name:       "invalid non-numeric",
			html:       `<div hx-get="/api" hx-status:abc="none">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "not a valid HTTP status",
			severity:   rules.Error,
		},
		{
			name: "valid partial wildcard 40x",
			html: `<div hx-get="/api" hx-status:40x="none">content</div>`,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "4"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXStatusV2Warning(t *testing.T) {
	// hx-status:* should warn when using htmx v2
	html := `<div hx-get="/api" hx-status:404="none">content</div>`

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	results, err := l.LintContent("test.html", []byte(html))
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	found := false
	for _, r := range results {
		if r.Rule == rules.RuleHTMXAttributes && strings.Contains(r.Message, "only available in htmx 4") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected warning about htmx 4 only, got %v", results)
	}
}

func TestLintContent_HTMXInheritedV4(t *testing.T) {
	tests := []htmxTestCase{
		// Valid :inherited suffix in htmx 4
		{
			name: "valid hx-boost:inherited",
			html: `<div hx-boost:inherited="true">content</div>`,
		},
		{
			name: "valid hx-swap:inherited",
			html: `<div hx-get="/api" hx-swap:inherited>content</div>`,
		},
		{
			name: "valid hx-target:inherited",
			html: `<div hx-get="/api" hx-target:inherited>content</div>`,
		},
		{
			name: "hx-swap:inherited still validates value",
			html: `<div hx-get="/api" hx-swap:inherited="innerHTML">content</div>`,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "4"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXInheritedV2Warning(t *testing.T) {
	// :inherited suffix should warn when using htmx v2
	html := `<div hx-boost:inherited="true">content</div>`

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	results, err := l.LintContent("test.html", []byte(html))
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	found := false
	for _, r := range results {
		if r.Rule == rules.RuleHTMXAttributes && strings.Contains(r.Message, ":inherited suffix is only available in htmx 4") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected warning about :inherited being htmx 4 only, got %v", results)
	}
}

func TestLintContent_HTMXSwapV4NewValues(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		version  string
		wantRule string
	}{
		{
			name:    "innerMorph valid in v4",
			html:    `<div hx-get="/api" hx-swap="innerMorph">content</div>`,
			version: "4",
		},
		{
			name:    "outerMorph valid in v4",
			html:    `<div hx-get="/api" hx-swap="outerMorph">content</div>`,
			version: "4",
		},
		{
			name:    "before alias valid in v4",
			html:    `<div hx-get="/api" hx-swap="before">content</div>`,
			version: "4",
		},
		{
			name:    "after alias valid in v4",
			html:    `<div hx-get="/api" hx-swap="after">content</div>`,
			version: "4",
		},
		{
			name:    "prepend alias valid in v4",
			html:    `<div hx-get="/api" hx-swap="prepend">content</div>`,
			version: "4",
		},
		{
			name:    "append alias valid in v4",
			html:    `<div hx-get="/api" hx-swap="append">content</div>`,
			version: "4",
		},
		{
			name:     "innerMorph warns in v2",
			html:     `<div hx-get="/api" hx-swap="innerMorph">content</div>`,
			version:  "2",
			wantRule: rules.RuleHTMXAttributes,
		},
		{
			name:     "before alias warns in v2",
			html:     `<div hx-get="/api" hx-swap="before">content</div>`,
			version:  "2",
			wantRule: rules.RuleHTMXAttributes,
		},
		{
			name:     "after alias warns in v2",
			html:     `<div hx-get="/api" hx-swap="after">content</div>`,
			version:  "2",
			wantRule: rules.RuleHTMXAttributes,
		},
		{
			name:     "prepend alias warns in v2",
			html:     `<div hx-get="/api" hx-swap="prepend">content</div>`,
			version:  "2",
			wantRule: rules.RuleHTMXAttributes,
		},
		{
			name:     "append alias warns in v2",
			html:     `<div hx-get="/api" hx-swap="append">content</div>`,
			version:  "2",
			wantRule: rules.RuleHTMXAttributes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := linter.DefaultConfig()
			cfg.Frameworks.HTMX = true
			cfg.Frameworks.HTMXVersion = tt.version
			l := linter.New(cfg)

			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}

			checkRule(t, results, rules.RuleHTMXAttributes, tt.wantRule)
		})
	}
}

func TestLintContent_HTMXSwapModifiersTargetStrip(t *testing.T) {
	tests := []htmxTestCase{
		{
			name: "valid target modifier",
			html: `<div hx-get="/api" hx-swap="innerHTML target:#result">content</div>`,
		},
		{
			name: "valid strip true",
			html: `<div hx-get="/api" hx-swap="innerHTML strip:true">content</div>`,
		},
		{
			name: "valid strip false",
			html: `<div hx-get="/api" hx-swap="innerHTML strip:false">content</div>`,
		},
		{
			name:       "invalid strip value",
			html:       `<div hx-get="/api" hx-swap="innerHTML strip:maybe">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "strip modifier should be 'true' or 'false'",
			severity:   rules.Error,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "4"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXRevealedTrigger(t *testing.T) {
	tests := []htmxTestCase{
		{
			name: "revealed with once",
			html: `<div hx-get="/api" hx-trigger="revealed once">content</div>`,
		},
		{
			name:       "revealed with invalid modifier",
			html:       `<div hx-get="/api" hx-trigger="revealed root:#foo">content</div>`,
			wantRule:   rules.RuleHTMXAttributes,
			wantSubstr: "unknown revealed modifier",
			severity:   rules.Warning,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXTargetNewKeywords(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "document target",
			html: `<div hx-get="/api" hx-target="document">content</div>`,
		},
		{
			name: "window target",
			html: `<div hx-get="/api" hx-target="window">content</div>`,
		},
		{
			name: "host target",
			html: `<div hx-get="/api" hx-target="host">content</div>`,
		},
		{
			name: "nextElementSibling target",
			html: `<div hx-get="/api" hx-target="nextElementSibling">content</div>`,
		},
		{
			name: "previousElementSibling target",
			html: `<div hx-get="/api" hx-target="previousElementSibling">content</div>`,
		},
		{
			name: "findAll keyword",
			html: `<div hx-get="/api" hx-target="findAll .item">content</div>`,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "4"
	l := linter.New(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}

			checkRule(t, results, rules.RuleHTMXAttributes, tt.wantRule)
		})
	}
}

func TestLintContent_HTMXv4NewEvents(t *testing.T) {
	tests := []htmxTestCase{
		{
			name: "valid htmx:config:request",
			html: `<div hx-get="/api" hx-on:htmx:config:request="cfg()">content</div>`,
		},
		{
			name: "valid htmx:before:init",
			html: `<div hx-get="/api" hx-on:htmx:before:init="setup()">content</div>`,
		},
		{
			name: "valid htmx:after:implicitInheritance",
			html: `<div hx-get="/api" hx-on:htmx:after:implicitinheritance="done()">content</div>`,
		},
		{
			name: "valid htmx:every standalone",
			html: `<div hx-get="/api" hx-on:htmx:every="poll()">content</div>`,
		},
		// Compound sub-actions
		{
			name: "valid compound htmx:before:sse:reconnect",
			html: `<div hx-get="/api" hx-on:htmx:before:sse:reconnect="handle()">content</div>`,
		},
		{
			name: "valid compound htmx:after:push:into:history",
			html: `<div hx-get="/api" hx-on:htmx:after:push:into:history="done()">content</div>`,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "4"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXStatusPartialWildcard(t *testing.T) {
	tests := []htmxTestCase{
		{
			name: "valid 40x pattern",
			html: `<div hx-get="/api" hx-status:40x="none">content</div>`,
		},
		{
			name: "valid 50x pattern",
			html: `<div hx-get="/api" hx-status:50x="none">content</div>`,
		},
		{
			name: "valid 20x pattern",
			html: `<div hx-get="/api" hx-status:20x="innerHTML">content</div>`,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "4"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXAppendSuffix(t *testing.T) {
	tests := []htmxTestCase{
		// Valid :append in v4
		{
			name: "valid hx-vals:append",
			html: `<div hx-post="/api" hx-vals:append='{"extra": "val"}'>content</div>`,
		},
		{
			name: "valid hx-headers:append",
			html: `<div hx-post="/api" hx-headers:append='{"X-Extra": "val"}'>content</div>`,
		},
		{
			name: "valid hx-swap:inherited:append",
			html: `<div hx-get="/api" hx-swap:inherited:append="innerHTML">content</div>`,
		},
	}

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "4"
	l := linter.New(cfg)

	runHTMXTests(t, l, tests)
}

func TestLintContent_HTMXAppendV2Warning(t *testing.T) {
	html := `<div hx-post="/api" hx-vals:append='{"extra": "val"}'>content</div>`

	cfg := linter.DefaultConfig()
	cfg.Frameworks.HTMX = true
	cfg.Frameworks.HTMXVersion = "2"
	l := linter.New(cfg)

	results, err := l.LintContent("test.html", []byte(html))
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	found := false
	for _, r := range results {
		if r.Rule == rules.RuleHTMXAttributes && strings.Contains(r.Message, ":append suffix is only available in htmx 4") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected warning about :append being htmx 4 only, got %v", results)
	}
}
