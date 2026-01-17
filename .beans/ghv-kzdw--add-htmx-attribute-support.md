---
# ghv-kzdw
title: Add htmx attribute support
status: completed
type: feature
priority: normal
created_at: 2026-01-17T23:06:23Z
updated_at: 2026-01-17T23:09:22Z
---

Add configuration option to allow htmx attributes (hx-*) to pass validation in the input-attributes rule. Since htmx 4 introduces new attributes not in htmx 2, add a version option to enforce correct attributes per version.

## Checklist
- [x] Add FrameworkConfig struct to config/config.go
- [x] Add Frameworks field to linter/config.go Config struct
- [x] Create rules/htmx.go with v2/v4 attribute definitions
- [x] Modify rules/input_attributes.go to use htmx validation
- [x] Add tests for htmx attribute scenarios
- [x] Run golangci-lint and go test

## Implementation

### Configuration

Add to `.htmlvalidate.json`:
```json
{
  "frameworks": {
    "htmx": true,
    "htmx-version": "2"
  }
}
```

### Behavior

- `htmx: false` (default): htmx attributes on `<input>` elements trigger a warning
- `htmx: true`: htmx attributes are validated against the specified version
- `htmx-version: "2"` (default): Warns on htmx 4-only attributes (`hx-optimistic`, `hx-preload`, etc.)
- `htmx-version: "4"`: Warns on deprecated v2 attributes (`hx-vars`, `hx-disinherit`, etc.)

### Files Changed

| File | Changes |
|------|---------|
| `config/config.go` | Added `FrameworkConfig` struct, updated `merge()` and `ToLinterConfig()` |
| `linter/config.go` | Added `FrameworkConfig` struct and `Frameworks` field to `Config` |
| `rules/htmx.go` (new) | htmx v2/v4 attribute sets and `ValidateHTMXAttribute()` helper |
| `rules/rule.go` | Added `HTMXConfigurable` interface |
| `rules/input_attributes.go` | Added htmx config and validation logic |
| `linter/linter.go` | Configures htmx-aware rules on initialization |
| `linter/linter_validation_test.go` | Added 9 test cases for htmx scenarios |
| `README.md` | Documented htmx configuration options |