package linter_test

import (
	"testing"

	"github.com/STR-Consulting/go-html-validate/rules"
)

// hasRule returns true if the results contain a finding from the given rule.
func hasRule(results []rules.Result, ruleName string) bool {
	for _, r := range results {
		if r.Rule == ruleName {
			return true
		}
	}
	return false
}

// checkRule asserts that a rule appears when wantRule matches ruleName,
// or does not appear when wantRule is empty.
func checkRule(t *testing.T, results []rules.Result, ruleName, wantRule string) {
	t.Helper()
	found := hasRule(results, ruleName)
	if wantRule != "" && !found {
		t.Errorf("expected %s rule, got %v", ruleName, results)
	}
	if wantRule == "" && found {
		t.Errorf("expected no %s rule, but got one", ruleName)
	}
}
