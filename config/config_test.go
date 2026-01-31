package config_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/toba/go-html-validate/config"
	"github.com/toba/go-html-validate/rules"
)

func TestLoadFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantRoot bool
		wantErr  bool
	}{
		{
			name:     "empty config",
			content:  `{}`,
			wantRoot: false,
		},
		{
			name:     "root true",
			content:  `{"root": true}`,
			wantRoot: true,
		},
		{
			name:    "invalid json",
			content: `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, config.ConfigFileName)
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatal(err)
			}

			cfg, err := config.LoadFile(path)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.Root != tt.wantRoot {
				t.Errorf("Root = %v, want %v", cfg.Root, tt.wantRoot)
			}
		})
	}
}

func TestRuleConfigUnmarshal(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		ruleName     string
		wantSeverity string
	}{
		{
			name:         "string severity error",
			content:      `{"rules": {"img-alt": "error"}}`,
			ruleName:     "img-alt",
			wantSeverity: "error",
		},
		{
			name:         "string severity warn",
			content:      `{"rules": {"img-alt": "warn"}}`,
			ruleName:     "img-alt",
			wantSeverity: "warn",
		},
		{
			name:         "string severity off",
			content:      `{"rules": {"img-alt": "off"}}`,
			ruleName:     "img-alt",
			wantSeverity: "off",
		},
		{
			name:         "number severity 2",
			content:      `{"rules": {"img-alt": 2}}`,
			ruleName:     "img-alt",
			wantSeverity: "error",
		},
		{
			name:         "number severity 1",
			content:      `{"rules": {"img-alt": 1}}`,
			ruleName:     "img-alt",
			wantSeverity: "warn",
		},
		{
			name:         "number severity 0",
			content:      `{"rules": {"img-alt": 0}}`,
			ruleName:     "img-alt",
			wantSeverity: "off",
		},
		{
			name:         "array format",
			content:      `{"rules": {"img-alt": ["error", {}]}}`,
			ruleName:     "img-alt",
			wantSeverity: "error",
		},
		{
			name:         "array with number",
			content:      `{"rules": {"img-alt": [1, {}]}}`,
			ruleName:     "img-alt",
			wantSeverity: "warn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, config.ConfigFileName)
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatal(err)
			}

			cfg, err := config.LoadFile(path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			rule, ok := cfg.Rules[tt.ruleName]
			if !ok {
				t.Fatalf("rule %q not found in config", tt.ruleName)
			}
			if rule.Severity != tt.wantSeverity {
				t.Errorf("severity = %q, want %q", rule.Severity, tt.wantSeverity)
			}
		})
	}
}

func TestExtends(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantSeverity string
	}{
		{
			name:         "string extends",
			content:      `{"extends": "html-validate:recommended"}`,
			wantSeverity: "",
		},
		{
			name:         "array extends",
			content:      `{"extends": ["html-validate:recommended"]}`,
			wantSeverity: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, config.ConfigFileName)
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatal(err)
			}

			cfg, err := config.LoadFile(path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(cfg.Extends) == 0 {
				t.Error("expected extends to be parsed")
			}
		})
	}
}

func TestToLinterConfig(t *testing.T) {
	fileCfg := &config.FileConfig{
		Rules: map[string]config.RuleConfig{
			"img-alt":         {Severity: "off"},
			"no-inline-style": {Severity: "warn"},
			"button-name":     {Severity: "error"},
		},
	}

	linterCfg := config.ToLinterConfig(fileCfg, "/test/.htmlvalidate.json")

	// Check disabled rules
	if !slices.Contains(linterCfg.DisabledRules, "img-alt") {
		t.Error("expected img-alt to be disabled")
	}

	// Check severity overrides
	if sev, ok := linterCfg.RuleSeverity["no-inline-style"]; !ok || sev != rules.Warning {
		t.Errorf("expected no-inline-style to have Warning severity, got %v", sev)
	}
	if sev, ok := linterCfg.RuleSeverity["button-name"]; !ok || sev != rules.Error {
		t.Errorf("expected button-name to have Error severity, got %v", sev)
	}

	// Check config path
	if linterCfg.ConfigPath != "/test/.htmlvalidate.json" {
		t.Errorf("ConfigPath = %q, want %q", linterCfg.ConfigPath, "/test/.htmlvalidate.json")
	}
}

func TestToLinterConfig_HTMXCustomEvents(t *testing.T) {
	fileCfg := &config.FileConfig{
		Frameworks: config.FrameworkConfig{
			HTMX:             true,
			HTMXVersion:      "2",
			HTMXCustomEvents: []string{"count", "notification", "status"},
		},
	}

	linterCfg := config.ToLinterConfig(fileCfg, "")

	if !linterCfg.Frameworks.HTMX {
		t.Error("expected HTMX to be enabled")
	}
	if len(linterCfg.Frameworks.HTMXCustomEvents) != 3 {
		t.Fatalf("expected 3 custom events, got %d", len(linterCfg.Frameworks.HTMXCustomEvents))
	}
	want := []string{"count", "notification", "status"}
	for i, ev := range want {
		if linterCfg.Frameworks.HTMXCustomEvents[i] != ev {
			t.Errorf("custom event %d = %q, want %q", i, linterCfg.Frameworks.HTMXCustomEvents[i], ev)
		}
	}
}

func TestLoadFile_HTMXCustomEvents(t *testing.T) {
	dir := t.TempDir()
	content := `{
		"frameworks": {
			"htmx": true,
			"htmx-version": "2",
			"htmx-custom-events": ["count", "notification"]
		}
	}`
	path := filepath.Join(dir, config.ConfigFileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.Frameworks.HTMX {
		t.Error("expected HTMX to be enabled")
	}
	if len(cfg.Frameworks.HTMXCustomEvents) != 2 {
		t.Fatalf("expected 2 custom events, got %d", len(cfg.Frameworks.HTMXCustomEvents))
	}
	if cfg.Frameworks.HTMXCustomEvents[0] != "count" {
		t.Errorf("custom event 0 = %q, want %q", cfg.Frameworks.HTMXCustomEvents[0], "count")
	}
	if cfg.Frameworks.HTMXCustomEvents[1] != "notification" {
		t.Errorf("custom event 1 = %q, want %q", cfg.Frameworks.HTMXCustomEvents[1], "notification")
	}
}

func TestResolveWithPreset(t *testing.T) {
	dir := t.TempDir()
	content := `{
		"extends": ["html-validate:a11y"],
		"rules": {
			"img-alt": "warn"
		}
	}`
	path := filepath.Join(dir, config.ConfigFileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, resolvedPath, err := config.Resolve(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolvedPath != path {
		t.Errorf("path = %q, want %q", resolvedPath, path)
	}

	// img-alt should be overridden to warn
	if rule, ok := cfg.Rules["img-alt"]; !ok || rule.Severity != "warn" {
		t.Errorf("expected img-alt severity to be warn, got %+v", cfg.Rules["img-alt"])
	}

	// prefer-tbody should be off from a11y preset
	if rule, ok := cfg.Rules["prefer-tbody"]; !ok || rule.Severity != "off" {
		t.Errorf("expected prefer-tbody severity to be off from preset, got %+v", cfg.Rules["prefer-tbody"])
	}
}
