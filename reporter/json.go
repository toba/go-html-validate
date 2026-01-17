package reporter

import (
	"encoding/json"
	"io"
	"os"

	"github.com/STR-Consulting/go-html-validate/rules"
)

// JSON outputs results in JSON format for CI integration.
type JSON struct {
	Writer io.Writer
	Pretty bool
}

// NewJSON creates a JSON reporter writing to stdout.
func NewJSON() *JSON {
	return &JSON{
		Writer: os.Stdout,
		Pretty: false,
	}
}

// JSONResult is the JSON representation of a lint result.
type JSONResult struct {
	Rule     string `json:"rule"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
}

// JSONOutput is the top-level JSON structure.
type JSONOutput struct {
	Results []JSONResult `json:"results"`
	Summary Summary      `json:"summary"`
}

// Summary contains aggregate counts.
type Summary struct {
	Total    int `json:"total"`
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Info     int `json:"info"`
}

// Report outputs results as JSON.
func (j *JSON) Report(results []rules.Result) error {
	output := JSONOutput{
		Results: make([]JSONResult, 0, len(results)),
	}

	for _, r := range results {
		output.Results = append(output.Results, JSONResult{
			Rule:     r.Rule,
			Message:  r.Message,
			Filename: r.Filename,
			Line:     r.Line,
			Column:   r.Col,
			Severity: r.Severity.String(),
		})

		output.Summary.Total++
		switch r.Severity {
		case rules.Error:
			output.Summary.Errors++
		case rules.Warning:
			output.Summary.Warnings++
		case rules.Info:
			output.Summary.Info++
		}
	}

	encoder := json.NewEncoder(j.Writer)
	if j.Pretty {
		encoder.SetIndent("", "  ")
	}

	return encoder.Encode(output)
}
