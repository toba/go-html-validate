// Package linter provides the core HTML linting orchestration.
package linter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"github.com/STR-Consulting/go-html-validate/rules"
)

// Linter coordinates HTML template accessibility checking.
type Linter struct {
	rules    []rules.Rule
	config   *Config
	reporter Reporter
}

// Reporter defines the interface for outputting lint results.
type Reporter interface {
	Report(results []rules.Result) error
}

// New creates a new Linter with the given configuration.
func New(cfg *Config) *Linter {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	registry := rules.NewRegistry()
	enabledRules := make([]rules.Rule, 0)

	for _, rule := range registry.All() {
		if cfg.IsRuleEnabled(rule.Name()) {
			// Configure htmx-aware rules
			if htmxRule, ok := rule.(rules.HTMXConfigurable); ok {
				htmxRule.Configure(cfg.Frameworks.HTMX, cfg.Frameworks.HTMXVersion)
			}
			enabledRules = append(enabledRules, rule)
		}
	}

	return &Linter{
		rules:  enabledRules,
		config: cfg,
	}
}

// SetReporter sets the output reporter.
func (l *Linter) SetReporter(r Reporter) {
	l.reporter = r
}

// LintFile checks a single file and returns any violations.
func (l *Linter) LintFile(path string) ([]rules.Result, error) {
	content, err := os.ReadFile(path) //nolint:gosec // user-specified file path is intentional
	if err != nil {
		return nil, err
	}

	return l.LintContent(path, content)
}

// LintContent checks HTML content and returns any violations.
func (l *Linter) LintContent(filename string, content []byte) ([]rules.Result, error) {
	doc, err := parser.ParseFragment(filename, content)
	if err != nil {
		return nil, err
	}

	var allResults []rules.Result
	for _, rule := range l.rules {
		results := rule.Check(doc)
		for _, r := range results {
			// Apply severity override from config
			if severity, ok := l.config.RuleSeverity[r.Rule]; ok {
				r.Severity = severity
			}
			// Filter by minimum severity
			if r.Severity <= l.config.MinSeverity {
				allResults = append(allResults, r)
			}
		}
	}

	return allResults, nil
}

// LintFiles checks multiple files and returns all violations.
func (l *Linter) LintFiles(paths []string) ([]rules.Result, error) {
	var allResults []rules.Result

	for _, path := range paths {
		// Skip ignored patterns
		if l.shouldIgnore(path) {
			continue
		}

		results, err := l.LintFile(path)
		if err != nil {
			// Report error but continue with other files
			allResults = append(allResults, rules.Result{
				Rule:     "parse-error",
				Message:  err.Error(),
				Filename: path,
				Line:     1,
				Col:      1,
				Severity: rules.Error,
			})
			continue
		}
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

// LintDir recursively checks all HTML files in a directory.
func (l *Linter) LintDir(dir string) ([]rules.Result, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if isHTMLFile(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return l.LintFiles(files)
}

// Run executes linting and reports results.
func (l *Linter) Run(paths []string) (int, error) {
	var allResults []rules.Result

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return 0, err
		}

		var results []rules.Result
		if info.IsDir() {
			results, err = l.LintDir(path)
		} else {
			results, err = l.LintFiles([]string{path})
		}
		if err != nil {
			return 0, err
		}
		allResults = append(allResults, results...)
	}

	if l.reporter != nil {
		if err := l.reporter.Report(allResults); err != nil {
			return 0, err
		}
	}

	// Count errors (not warnings)
	errorCount := 0
	for _, r := range allResults {
		if r.Severity == rules.Error {
			errorCount++
		}
	}

	return errorCount, nil
}

func (l *Linter) shouldIgnore(path string) bool {
	for _, pattern := range l.config.IgnorePatterns {
		if matchIgnorePattern(path, pattern) {
			return true
		}
	}
	return false
}

// matchIgnorePattern checks if a path matches a gitignore-style pattern.
func matchIgnorePattern(path, pattern string) bool {
	// Handle directory patterns (ending with /)
	if strings.HasSuffix(pattern, "/") {
		dir := strings.TrimSuffix(pattern, "/")
		if strings.Contains(path, "/"+dir+"/") ||
			strings.HasPrefix(path, dir+"/") ||
			path == dir {
			return true
		}
		return false
	}

	// Handle ** recursive patterns
	if strings.Contains(pattern, "**") {
		return matchDoublestar(path, pattern)
	}

	// Simple glob matching against basename
	if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
		return true
	}

	// Try matching full path
	if matched, _ := filepath.Match(pattern, path); matched {
		return true
	}

	return false
}

// matchDoublestar handles ** patterns.
func matchDoublestar(path, pattern string) bool {
	parts := strings.Split(pattern, "**")
	if len(parts) != 2 {
		return false
	}

	prefix := strings.TrimSuffix(parts[0], "/")
	suffix := strings.TrimPrefix(parts[1], "/")

	// **/*.html - match any path ending with suffix
	if prefix == "" {
		if matched, _ := filepath.Match(suffix, filepath.Base(path)); matched {
			return true
		}
		if suffix != "" && strings.HasSuffix(path, suffix) {
			return true
		}
		return false
	}

	// prefix/**/suffix - match prefix at start, suffix at end
	if suffix != "" {
		hasPrefix := strings.HasPrefix(path, prefix+"/") || strings.HasPrefix(path, prefix)
		hasSuffix := strings.HasSuffix(path, suffix) ||
			func() bool { m, _ := filepath.Match(suffix, filepath.Base(path)); return m }()
		return hasPrefix && hasSuffix
	}

	// prefix/** - match anything under prefix
	return strings.HasPrefix(path, prefix+"/") || path == prefix
}

func isHTMLFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".html" || ext == ".htm" || ext == ".gohtml" || ext == ".tmpl"
}
