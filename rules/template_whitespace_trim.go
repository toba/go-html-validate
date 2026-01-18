package rules

import (
	"bytes"
	"regexp"

	"github.com/STR-Consulting/go-html-validate/parser"
)

// TemplateWhitespaceTrim checks for Go template actions that may create unwanted
// whitespace in output. It suggests using trim markers ({{- and -}}) to prevent
// empty lines in rendered output.
type TemplateWhitespaceTrim struct{}

func (r *TemplateWhitespaceTrim) Name() string { return RuleTemplateWhitespaceTrim }

func (r *TemplateWhitespaceTrim) Description() string {
	return "suggest trim markers to prevent unwanted whitespace in template output"
}

// Check implements Rule but returns nil - this rule uses CheckRaw instead.
func (r *TemplateWhitespaceTrim) Check(_ *parser.Document) []Result {
	return nil
}

// templateActionPattern matches Go template actions and captures:
// - group 1: opening ({{ or {{-)
// - group 2: action content
// - group 3: closing (-}} or }})
var templateActionPattern = regexp.MustCompile(`(\{\{-?)\s*(.*?)\s*(-?\}\})`)

// controlFlowKeywords are template actions that don't produce output
// and commonly appear alone on lines.
var controlFlowKeywords = map[string]bool{
	"if":       true,
	"else":     true,
	"end":      true,
	"range":    true,
	"with":     true,
	"block":    true,
	"define":   true,
	"template": true,
}

// CheckRaw examines the raw template content for whitespace trim issues.
func (r *TemplateWhitespaceTrim) CheckRaw(filename string, content []byte) []Result {
	var results []Result

	lines := bytes.Split(content, []byte("\n"))

	for lineNum, line := range lines {
		// Find all template actions on this line
		matches := templateActionPattern.FindAllSubmatchIndex(line, -1)
		if len(matches) == 0 {
			continue
		}

		for _, match := range matches {
			// Check if this action is alone on the line (only whitespace around it)
			if !isActionAloneOnLine(line, match[0], match[1]) {
				continue
			}

			// Extract the action content
			actionContent := string(line[match[4]:match[5]])
			if !isControlFlowAction(actionContent) {
				continue
			}

			// Check if the closing already has a trim marker
			closing := string(line[match[6]:match[7]])
			if closing == "-}}" {
				continue
			}

			// Report a warning
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "control flow action alone on line should use trailing trim marker (-}}) to prevent blank lines",
				Filename: filename,
				Line:     lineNum + 1,
				Col:      match[0] + 1,
				Severity: Warning,
			})
		}
	}

	return results
}

// isActionAloneOnLine checks if the action at the given position is the only
// non-whitespace content on the line.
func isActionAloneOnLine(line []byte, start, end int) bool {
	// Check that everything before the action is whitespace
	for i := 0; i < start; i++ {
		if line[i] != ' ' && line[i] != '\t' {
			return false
		}
	}

	// Check that everything after the action is whitespace
	for i := end; i < len(line); i++ {
		if line[i] != ' ' && line[i] != '\t' && line[i] != '\r' {
			return false
		}
	}

	return true
}

// isControlFlowAction checks if the action content is a control flow keyword.
func isControlFlowAction(content string) bool {
	// Handle "else if" as a special case
	if len(content) >= 7 && content[:7] == "else if" {
		return true
	}

	// Extract the first word (the keyword)
	keyword := content
	for i, c := range content {
		if c == ' ' || c == '\t' {
			keyword = content[:i]
			break
		}
	}

	return controlFlowKeywords[keyword]
}
