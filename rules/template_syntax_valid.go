package rules

import (
	"bytes"
	"regexp"

	"github.com/STR-Consulting/go-html-validate/parser"
)

// TemplateSyntaxValid checks for basic Go template syntax errors.
type TemplateSyntaxValid struct{}

func (r *TemplateSyntaxValid) Name() string { return RuleTemplateSyntaxValid }

func (r *TemplateSyntaxValid) Description() string {
	return "validate Go template syntax for common errors"
}

// Check implements Rule but returns nil - this rule uses CheckRaw instead.
func (r *TemplateSyntaxValid) Check(_ *parser.Document) []Result {
	return nil
}

// controlStructures that require matching end.
var controlStructures = map[string]bool{
	"if":     true,
	"range":  true,
	"with":   true,
	"block":  true,
	"define": true,
}

// CheckRaw examines the raw template content for syntax errors.
func (r *TemplateSyntaxValid) CheckRaw(filename string, content []byte) []Result {
	// Check for unbalanced braces
	braceResults := r.checkBalancedBraces(filename, content)

	// Check for unbalanced control structures
	controlResults := r.checkBalancedControlStructures(filename, content)

	// Check for invalid trim marker syntax
	trimResults := r.checkTrimMarkerSyntax(filename, content)

	// Combine all results
	results := make([]Result, 0, len(braceResults)+len(controlResults)+len(trimResults))
	results = append(results, braceResults...)
	results = append(results, controlResults...)
	results = append(results, trimResults...)

	return results
}

// checkBalancedBraces verifies that {{ and }} are balanced.
func (r *TemplateSyntaxValid) checkBalancedBraces(filename string, content []byte) []Result {
	var results []Result

	openCount := bytes.Count(content, []byte("{{"))
	closeCount := bytes.Count(content, []byte("}}"))

	if openCount > closeCount {
		// Find the first unmatched {{
		line, col := r.findUnmatchedOpen(content)
		results = append(results, Result{
			Rule:     r.Name(),
			Message:  "unmatched '{{' - missing closing '}}'",
			Filename: filename,
			Line:     line,
			Col:      col,
			Severity: Error,
		})
	} else if closeCount > openCount {
		// Find the first unmatched }}
		line, col := r.findUnmatchedClose(content)
		results = append(results, Result{
			Rule:     r.Name(),
			Message:  "unmatched '}}' - missing opening '{{'",
			Filename: filename,
			Line:     line,
			Col:      col,
			Severity: Error,
		})
	}

	return results
}

// findUnmatchedOpen finds the position of an unmatched {{.
func (r *TemplateSyntaxValid) findUnmatchedOpen(content []byte) (line, col int) {
	depth := 0
	lastOpenLine, lastOpenCol := 1, 1
	currentLine, currentCol := 1, 1

	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			currentLine++
			currentCol = 1
			continue
		}

		if i+1 < len(content) && content[i] == '{' && content[i+1] == '{' {
			if depth == 0 {
				lastOpenLine, lastOpenCol = currentLine, currentCol
			}
			depth++
			i++ // skip second {
			currentCol += 2
			continue
		}

		if i+1 < len(content) && content[i] == '}' && content[i+1] == '}' {
			depth--
			i++ // skip second }
			currentCol += 2
			continue
		}

		currentCol++
	}

	return lastOpenLine, lastOpenCol
}

// findUnmatchedClose finds the position of an unmatched }}.
func (r *TemplateSyntaxValid) findUnmatchedClose(content []byte) (line, col int) {
	depth := 0
	currentLine, currentCol := 1, 1

	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			currentLine++
			currentCol = 1
			continue
		}

		if i+1 < len(content) && content[i] == '{' && content[i+1] == '{' {
			depth++
			i++ // skip second {
			currentCol += 2
			continue
		}

		if i+1 < len(content) && content[i] == '}' && content[i+1] == '}' {
			depth--
			if depth < 0 {
				return currentLine, currentCol
			}
			i++ // skip second }
			currentCol += 2
			continue
		}

		currentCol++
	}

	return 1, 1
}

// checkBalancedControlStructures verifies that if/range/with/block have matching end.
func (r *TemplateSyntaxValid) checkBalancedControlStructures(filename string, content []byte) []Result {
	var results []Result

	// Stack to track open control structures
	type openStruct struct {
		keyword string
		line    int
		col     int
	}
	var stack []openStruct

	lines := bytes.Split(content, []byte("\n"))

	// Pattern to extract template actions
	actionRegex := regexp.MustCompile(`\{\{-?\s*(\w+)`)

	for lineNum, line := range lines {
		matches := actionRegex.FindAllSubmatchIndex(line, -1)
		for _, match := range matches {
			if match[2] < 0 || match[3] < 0 {
				continue
			}
			keyword := string(line[match[2]:match[3]])

			switch {
			case controlStructures[keyword]:
				stack = append(stack, openStruct{
					keyword: keyword,
					line:    lineNum + 1,
					col:     match[0] + 1,
				})
			case keyword == "end":
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				} else {
					results = append(results, Result{
						Rule:     r.Name(),
						Message:  "unexpected '{{ end }}' - no matching control structure",
						Filename: filename,
						Line:     lineNum + 1,
						Col:      match[0] + 1,
						Severity: Error,
					})
				}
			case keyword == "else":
				// else doesn't pop the stack, it's part of an if/with
				if len(stack) == 0 {
					results = append(results, Result{
						Rule:     r.Name(),
						Message:  "unexpected '{{ else }}' - no matching 'if' or 'with'",
						Filename: filename,
						Line:     lineNum + 1,
						Col:      match[0] + 1,
						Severity: Error,
					})
				}
			}
		}
	}

	// Report any unclosed control structures
	for _, open := range stack {
		results = append(results, Result{
			Rule:     r.Name(),
			Message:  "unclosed '{{ " + open.keyword + " }}' - missing '{{ end }}'",
			Filename: filename,
			Line:     open.line,
			Col:      open.col,
			Severity: Error,
		})
	}

	return results
}

// checkTrimMarkerSyntax verifies that trim markers have proper spacing.
func (r *TemplateSyntaxValid) checkTrimMarkerSyntax(filename string, content []byte) []Result {
	var results []Result

	lines := bytes.Split(content, []byte("\n"))

	// Pattern to find {{- without space after or -}} without space before
	// Valid: {{- foo }}, {{ foo -}}
	// Invalid: {{-foo }}, {{ foo-}}
	leadingTrimNoSpace := regexp.MustCompile(`\{\{-[^\s]`)
	trailingTrimNoSpace := regexp.MustCompile(`[^\s]-\}\}`)

	for lineNum, line := range lines {
		// Check leading trim marker
		if match := leadingTrimNoSpace.FindIndex(line); match != nil {
			// Make sure it's not {{--}} (double dash edge case)
			if match[0]+3 < len(line) && line[match[0]+3] != '-' {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  "trim marker '{{-' must be followed by whitespace",
					Filename: filename,
					Line:     lineNum + 1,
					Col:      match[0] + 1,
					Severity: Error,
				})
			}
		}

		// Check trailing trim marker
		if match := trailingTrimNoSpace.FindIndex(line); match != nil {
			// Get the character before the dash
			charBefore := line[match[0]]
			// Skip if it's part of a comment or string (heuristic: skip if dash follows a dash)
			if charBefore != '-' {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  "trim marker '-}}' must be preceded by whitespace",
					Filename: filename,
					Line:     lineNum + 1,
					Col:      match[0] + 1,
					Severity: Error,
				})
			}
		}
	}

	return results
}
