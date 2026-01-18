---
# ghv-n3ln
title: Implement template-syntax-valid linting rule
status: completed
type: feature
priority: normal
created_at: 2026-01-18T01:06:00Z
updated_at: 2026-01-18T01:09:00Z
---

Validate basic template syntax correctness.

## Checks
- Balanced `{{` and `}}`
- Balanced control structures (`if`/`end`, `range`/`end`, etc.)
- Valid trim marker syntax (`{{-` must have space after `-`)

## Checklist
- [x] Add RuleTemplateSyntaxValid constant to rules/rule.go
- [x] Create rules/template_syntax_valid.go implementing the rule
- [x] Register rule in NewRegistry()
- [x] Add tests in linter/linter_templates_test.go
- [x] Run go build, golangci-lint, and tests