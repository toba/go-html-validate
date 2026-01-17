// Package reporter provides output formatting for lint results.
package reporter

import (
	"github.com/STR-Consulting/go-html-validate/rules"
)

// Reporter defines the interface for outputting lint results.
type Reporter interface {
	Report(results []rules.Result) error
}
