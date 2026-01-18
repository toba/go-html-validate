---
# ghv-1r9s
title: Implement template-whitespace-trim linting rule
status: completed
type: feature
priority: normal
created_at: 2026-01-18T00:47:37Z
updated_at: 2026-01-18T01:04:45Z
---

Add linting rules for Go html/template syntax. Initial focus on whitespace trim markers to prevent unwanted empty lines.

## Implementation

Added a new `RawRule` interface that receives raw file content before template preprocessing, enabling linting of template syntax itself.

### Files Modified
- `rules/rule.go`: Added `RawRule` interface and `RuleTemplateWhitespaceTrim` constant
- `rules/template_whitespace_trim.go`: New file implementing the rule
- `linter/linter.go`: Modified `LintContent()` to call `CheckRaw()` for `RawRule` implementations
- `linter/linter_templates_test.go`: Added 14 test cases

## Checklist
- [x] Add RawRule interface to rules/rule.go
- [x] Add RuleTemplateWhitespaceTrim constant
- [x] Create rules/template_whitespace_trim.go implementing the rule
- [x] Modify linter/linter.go to call CheckRaw() for RawRule implementations
- [x] Add tests in linter/linter_templates_test.go
- [x] Run go build, golangci-lint, and tests