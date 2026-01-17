package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/STR-Consulting/go-html-validate/rules"
)

// Text outputs human-readable lint results.
type Text struct {
	Writer    io.Writer
	NoColor   bool
	ShowRules bool // Include rule name in output
}

// NewText creates a text reporter writing to stdout.
func NewText() *Text {
	return &Text{
		Writer:    os.Stdout,
		NoColor:   false,
		ShowRules: true,
	}
}

// Report outputs results in human-readable format.
func (t *Text) Report(results []rules.Result) error {
	if len(results) == 0 {
		return nil
	}

	// Group by file
	byFile := make(map[string][]rules.Result)
	for _, r := range results {
		byFile[r.Filename] = append(byFile[r.Filename], r)
	}

	// Sort files
	files := make([]string, 0, len(byFile))
	for f := range byFile {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, file := range files {
		fileResults := byFile[file]
		// Sort by line, then column
		sort.Slice(fileResults, func(i, j int) bool {
			if fileResults[i].Line != fileResults[j].Line {
				return fileResults[i].Line < fileResults[j].Line
			}
			return fileResults[i].Col < fileResults[j].Col
		})

		for _, r := range fileResults {
			line := t.formatResult(r)
			_, _ = fmt.Fprintln(t.Writer, line)
		}
	}

	// Summary
	errorCount := 0
	warningCount := 0
	for _, r := range results {
		switch r.Severity {
		case rules.Error:
			errorCount++
		case rules.Warning:
			warningCount++
		}
	}

	_, _ = fmt.Fprintln(t.Writer)
	if errorCount > 0 || warningCount > 0 {
		parts := []string{}
		if errorCount > 0 {
			parts = append(parts, fmt.Sprintf("%d error(s)", errorCount))
		}
		if warningCount > 0 {
			parts = append(parts, fmt.Sprintf("%d warning(s)", warningCount))
		}
		_, _ = fmt.Fprintf(t.Writer, "Found %s\n", strings.Join(parts, ", "))
	}

	return nil
}

func (t *Text) formatResult(r rules.Result) string {
	// Format: file:line:col: severity: message [rule]
	severity := r.Severity.String()
	if !t.NoColor {
		severity = t.colorize(severity, r.Severity)
	}

	if t.ShowRules {
		return fmt.Sprintf("%s:%d:%d: %s: %s [%s]",
			r.Filename, r.Line, r.Col, severity, r.Message, r.Rule)
	}
	return fmt.Sprintf("%s:%d:%d: %s: %s",
		r.Filename, r.Line, r.Col, severity, r.Message)
}

func (t *Text) colorize(text string, severity rules.Severity) string {
	var code string
	switch severity {
	case rules.Error:
		code = "\033[31m" // Red
	case rules.Warning:
		code = "\033[33m" // Yellow
	case rules.Info:
		code = "\033[36m" // Cyan
	default:
		return text
	}
	return code + text + "\033[0m"
}
