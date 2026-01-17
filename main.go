// htmlint is an HTML accessibility linter for Go templates.
//
// Usage:
//
//	htmlint [options] <files or directories>
//
// Options:
//
//	-f, --format     Output format: text, json (default: text)
//	-q, --quiet      Only show errors, not warnings
//	--no-color       Disable colored output
//	--ignore         Glob patterns to ignore (can be repeated)
//	--disable        Disable specific rules (can be repeated)
//	--config         Path to config file
//	--no-config      Disable config file loading
//	--print-config   Print resolved configuration and exit
//	-h, --help       Show help
//
// Examples:
//
//	htmlint web/
//	htmlint -q web/**/*.html
//	htmlint --format=json web/ > lint-results.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/STR-Consulting/go-html-validate/config"
	"github.com/STR-Consulting/go-html-validate/linter"
	"github.com/STR-Consulting/go-html-validate/reporter"
	"github.com/STR-Consulting/go-html-validate/rules"
)

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	os.Exit(run())
}

func run() int {
	var (
		format       string
		quiet        bool
		noColor      bool
		ignoreFlags  stringSlice
		disableFlags stringSlice
		showHelp     bool
		showVersion  bool
		listRules    bool
		configPath   string
		noConfig     bool
		printConfig  bool
	)

	flag.StringVar(&format, "format", "text", "Output format: text, json")
	flag.StringVar(&format, "f", "text", "Output format (shorthand)")
	flag.BoolVar(&quiet, "quiet", false, "Only show errors")
	flag.BoolVar(&quiet, "q", false, "Only show errors (shorthand)")
	flag.BoolVar(&noColor, "no-color", false, "Disable colored output")
	flag.Var(&ignoreFlags, "ignore", "Glob pattern to ignore")
	flag.Var(&disableFlags, "disable", "Rule to disable")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showHelp, "h", false, "Show help (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version (shorthand)")
	flag.BoolVar(&listRules, "list-rules", false, "List available rules")
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.BoolVar(&noConfig, "no-config", false, "Disable config file loading")
	flag.BoolVar(&printConfig, "print-config", false, "Print resolved configuration")

	flag.Usage = usage
	flag.Parse()

	if showHelp {
		usage()
		return 0
	}

	if showVersion {
		fmt.Println(getVersion())
		return 0
	}

	if listRules {
		printRules()
		return 0
	}

	args := flag.Args()

	// Determine search directory for config
	searchDir := "."
	if len(args) > 0 {
		if info, err := os.Stat(args[0]); err == nil && info.IsDir() {
			searchDir = args[0]
		} else if err == nil {
			searchDir = filepath.Dir(args[0])
		}
	}

	// Load config file
	var fileCfg *config.FileConfig
	var loadedConfigPath string
	if !noConfig {
		var err error
		if configPath != "" {
			fileCfg, err = config.LoadFile(configPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return 1
			}
			loadedConfigPath = configPath
			// Resolve extends for explicit config
			fileCfg, err = resolveExtendsFromPath(fileCfg, configPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return 1
			}
		} else {
			fileCfg, loadedConfigPath, err = config.Resolve(searchDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
				return 1
			}
		}
	}

	// Load ignore patterns
	ignorePatterns, err := config.LoadIgnorePatterns(searchDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: error loading ignore file: %v\n", err)
	}

	// Build linter config from file config
	cfg := config.ToLinterConfig(fileCfg, loadedConfigPath)

	// Add ignore patterns from ignore file
	cfg.IgnorePatterns = append(cfg.IgnorePatterns, ignorePatterns...)

	// CLI flags override config file
	cfg.DisabledRules = append(cfg.DisabledRules, disableFlags...)
	cfg.IgnorePatterns = append(cfg.IgnorePatterns, ignoreFlags...)

	if quiet {
		cfg.ErrorsOnly()
	}

	// Print config and exit if requested
	if printConfig {
		printResolvedConfig(cfg, loadedConfigPath)
		return 0
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: no files or directories specified")
		fmt.Fprintln(os.Stderr, "usage: htmlint [options] <files or directories>")
		return 1
	}

	// Create linter
	l := linter.New(cfg)

	// Set reporter
	var rep linter.Reporter
	switch format {
	case "json":
		rep = reporter.NewJSON()
	default:
		textRep := reporter.NewText()
		textRep.NoColor = noColor
		rep = textRep
	}
	l.SetReporter(rep)

	// Run linting
	errorCount, err := l.Run(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if errorCount > 0 {
		return 1
	}
	return 0
}

// resolveExtendsFromPath resolves extends for a config loaded from an explicit path.
func resolveExtendsFromPath(cfg *config.FileConfig, path string) (*config.FileConfig, error) {
	if len(cfg.Extends) == 0 {
		return cfg, nil
	}

	// Use config package's Resolve by writing to temp and loading
	// This is a simplified approach - we merge presets directly
	result := &config.FileConfig{
		Root:  cfg.Root,
		Rules: make(map[string]config.RuleConfig),
	}

	// Apply extends
	for _, ext := range cfg.Extends {
		if preset, ok := config.Presets[ext]; ok {
			for name, rule := range preset.Rules {
				result.Rules[name] = rule
			}
		} else {
			// Try as file path relative to config
			extPath := ext
			if !filepath.IsAbs(extPath) {
				extPath = filepath.Join(filepath.Dir(path), ext)
			}
			extCfg, err := config.LoadFile(extPath)
			if err != nil {
				return nil, fmt.Errorf("loading extended config %q: %w", ext, err)
			}
			for name, rule := range extCfg.Rules {
				result.Rules[name] = rule
			}
		}
	}

	// Apply current config on top
	for name, rule := range cfg.Rules {
		result.Rules[name] = rule
	}

	return result, nil
}

func printResolvedConfig(cfg *linter.Config, configPath string) {
	output := struct {
		ConfigFile     string            `json:"configFile,omitempty"`
		DisabledRules  []string          `json:"disabledRules,omitempty"`
		RuleSeverities map[string]string `json:"ruleSeverities,omitempty"`
		IgnorePatterns []string          `json:"ignorePatterns,omitempty"`
	}{
		ConfigFile:     configPath,
		DisabledRules:  cfg.DisabledRules,
		IgnorePatterns: cfg.IgnorePatterns,
	}

	if len(cfg.RuleSeverity) > 0 {
		output.RuleSeverities = make(map[string]string)
		for name, sev := range cfg.RuleSeverity {
			output.RuleSeverities[name] = sev.String()
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(output)
}

func usage() {
	fmt.Fprintf(os.Stderr, `htmlint - HTML accessibility linter for Go templates

Usage:
  htmlint [options] <files or directories>

Options:
  -f, --format      Output format: text, json (default: text)
  -q, --quiet       Only show errors, not warnings
  --no-color        Disable colored output
  --ignore PATTERN  Glob pattern to ignore (can be repeated)
  --disable RULE    Disable specific rule (can be repeated)
  --config PATH     Path to config file (.htmlvalidate.json)
  --no-config       Disable config file loading
  --print-config    Print resolved configuration and exit
  --list-rules      List available rules
  -v, --version     Show version
  -h, --help        Show this help

Config files:
  htmlint looks for .htmlvalidate.json in the target directory and parent
  directories. Use .htmlvalidateignore for gitignore-style file patterns.

Examples:
  htmlint web/
  htmlint -q web/**/*.html
  htmlint --format=json web/ > lint-results.json
  htmlint --disable=prefer-aria web/
`)
}

func printRules() {
	registry := rules.NewRegistry()
	fmt.Println("Available rules:")
	fmt.Println()
	for _, rule := range registry.All() {
		fmt.Printf("  %-30s %s\n", rule.Name(), rule.Description())
	}
}

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "dev"
}
