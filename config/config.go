// Package config handles .htmlvalidate.json configuration file loading.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/STR-Consulting/go-html-validate/linter"
	"github.com/STR-Consulting/go-html-validate/rules"
)

// ConfigFileName is the name of the configuration file.
const ConfigFileName = ".htmlvalidate.json"

// FrameworkConfig configures framework-specific attribute handling.
type FrameworkConfig struct {
	// HTMX enables htmx attribute validation.
	HTMX bool `json:"htmx"`
	// HTMXVersion specifies which htmx version to validate against ("2" or "4").
	// Defaults to "2" when HTMX is enabled.
	HTMXVersion string `json:"htmx-version"`
}

// FileConfig represents the JSON structure of .htmlvalidate.json.
type FileConfig struct {
	// Schema is the JSON schema URL (ignored, but allowed for IDE support).
	Schema string `json:"$schema"`
	// Root stops parent directory searching when true.
	Root bool `json:"root"`
	// Extends lists presets or config files to extend.
	Extends StringOrStrings `json:"extends"`
	// Rules configures individual rule severity.
	Rules map[string]RuleConfig `json:"rules"`
	// Frameworks configures framework-specific attribute handling.
	Frameworks FrameworkConfig `json:"frameworks"`
}

// StringOrStrings handles JSON that can be either a string or array of strings.
type StringOrStrings []string

func (s *StringOrStrings) UnmarshalJSON(data []byte) error {
	// Try as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = []string{str}
		return nil
	}
	// Try as array
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*s = arr
	return nil
}

// RuleConfig holds configuration for a single rule.
// Supports both simple ("error") and array (["error", {}]) formats.
type RuleConfig struct {
	Severity string
	Options  map[string]any
}

func (r *RuleConfig) UnmarshalJSON(data []byte) error {
	// Try as string: "error", "warn", "off"
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		r.Severity = str
		return nil
	}

	// Try as number: 0, 1, 2
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		switch num {
		case 0:
			r.Severity = "off"
		case 1:
			r.Severity = "warn"
		case 2:
			r.Severity = "error"
		default:
			return fmt.Errorf("invalid severity number: %d (must be 0, 1, or 2)", num)
		}
		return nil
	}

	// Try as array: ["error", {...}]
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("rule config must be string, number, or array")
	}
	if len(arr) == 0 {
		return errors.New("rule config array cannot be empty")
	}

	// First element is severity
	if err := json.Unmarshal(arr[0], &str); err == nil {
		r.Severity = str
	} else if err := json.Unmarshal(arr[0], &num); err == nil {
		switch num {
		case 0:
			r.Severity = "off"
		case 1:
			r.Severity = "warn"
		case 2:
			r.Severity = "error"
		default:
			return fmt.Errorf("invalid severity number: %d", num)
		}
	} else {
		return errors.New("first element of rule config must be severity")
	}

	// Second element is options (optional)
	if len(arr) > 1 {
		r.Options = make(map[string]any)
		if err := json.Unmarshal(arr[1], &r.Options); err != nil {
			return fmt.Errorf("invalid rule options: %w", err)
		}
	}

	return nil
}

// Load searches for and loads .htmlvalidate.json from dir upward.
// Returns nil config if no config file is found.
func Load(dir string) (*FileConfig, string, error) {
	path, err := FindConfigFile(dir)
	if err != nil {
		return nil, "", err
	}
	if path == "" {
		return nil, "", nil
	}
	cfg, err := LoadFile(path)
	return cfg, path, err
}

// LoadFile loads a specific configuration file.
func LoadFile(path string) (*FileConfig, error) {
	data, err := os.ReadFile(path) //nolint:gosec // user-specified config path
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg FileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	return &cfg, nil
}

// FindConfigFile searches for .htmlvalidate.json from dir upward.
// Returns empty string if no config file is found.
func FindConfigFile(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(absDir, ConfigFileName)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		parent := filepath.Dir(absDir)
		if parent == absDir {
			// Reached root
			return "", nil
		}
		absDir = parent
	}
}

// Resolve loads a config and resolves all extends.
func Resolve(dir string) (*FileConfig, string, error) {
	cfg, path, err := Load(dir)
	if err != nil {
		return nil, "", err
	}
	if cfg == nil {
		return nil, "", nil
	}

	resolved, err := resolveExtends(cfg, filepath.Dir(path))
	if err != nil {
		return nil, path, err
	}

	return resolved, path, nil
}

// resolveExtends merges extended configs into the base config.
func resolveExtends(cfg *FileConfig, baseDir string) (*FileConfig, error) {
	if len(cfg.Extends) == 0 {
		return cfg, nil
	}

	// Start with empty config, apply extends in order, then overlay current config
	result := &FileConfig{
		Rules: make(map[string]RuleConfig),
	}

	for _, ext := range cfg.Extends {
		var extCfg *FileConfig
		var err error

		if preset, ok := Presets[ext]; ok {
			extCfg = preset
		} else {
			// Try as file path
			extPath := ext
			if !filepath.IsAbs(extPath) {
				extPath = filepath.Join(baseDir, ext)
			}
			extCfg, err = LoadFile(extPath)
			if err != nil {
				return nil, fmt.Errorf("loading extended config %q: %w", ext, err)
			}
			// Recursively resolve extends
			extCfg, err = resolveExtends(extCfg, filepath.Dir(extPath))
			if err != nil {
				return nil, err
			}
		}

		// Merge extended config into result
		result = merge(result, extCfg)
	}

	// Apply current config on top
	result = merge(result, cfg)
	result.Root = cfg.Root
	result.Extends = nil // Already resolved

	return result, nil
}

// merge combines two configs, with overlay taking precedence.
func merge(base, overlay *FileConfig) *FileConfig {
	result := &FileConfig{
		Root:    overlay.Root || base.Root,
		Extends: overlay.Extends,
		Rules:   make(map[string]RuleConfig),
	}

	// Copy base rules
	for name, cfg := range base.Rules {
		result.Rules[name] = cfg
	}

	// Overlay rules take precedence
	for name, cfg := range overlay.Rules {
		result.Rules[name] = cfg
	}

	// Merge frameworks config (overlay takes precedence)
	result.Frameworks = base.Frameworks
	if overlay.Frameworks.HTMX {
		result.Frameworks.HTMX = true
	}
	if overlay.Frameworks.HTMXVersion != "" {
		result.Frameworks.HTMXVersion = overlay.Frameworks.HTMXVersion
	}

	return result
}

// ToLinterConfig converts a FileConfig to linter.Config.
func ToLinterConfig(fc *FileConfig, configPath string) *linter.Config {
	cfg := linter.DefaultConfig()
	cfg.ConfigPath = configPath

	if fc == nil {
		return cfg
	}

	for name, ruleCfg := range fc.Rules {
		switch ruleCfg.Severity {
		case "off", "0":
			cfg.DisabledRules = append(cfg.DisabledRules, name)
		case "error", "2":
			cfg.RuleSeverity[name] = rules.Error
		case "warn", "warning", "1":
			cfg.RuleSeverity[name] = rules.Warning
		}
	}

	// Copy frameworks config
	cfg.Frameworks = linter.FrameworkConfig{
		HTMX:        fc.Frameworks.HTMX,
		HTMXVersion: fc.Frameworks.HTMXVersion,
	}

	return cfg
}

// ParseSeverity converts a severity string to rules.Severity.
func ParseSeverity(s string) (rules.Severity, error) {
	switch s {
	case "error", "2":
		return rules.Error, nil
	case "warn", "warning", "1":
		return rules.Warning, nil
	case "info":
		return rules.Info, nil
	case "off", "0":
		return rules.Info, nil // off handled separately
	default:
		// Try parsing as number
		if n, err := strconv.Atoi(s); err == nil {
			switch n {
			case 0:
				return rules.Info, nil
			case 1:
				return rules.Warning, nil
			case 2:
				return rules.Error, nil
			}
		}
		return rules.Error, fmt.Errorf("invalid severity: %q", s)
	}
}
