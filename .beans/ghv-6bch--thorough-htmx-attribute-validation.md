---
# ghv-6bch
title: Thorough htmx attribute validation
status: in-progress
type: feature
priority: normal
created_at: 2026-01-17T23:31:34Z
updated_at: 2026-01-17T23:39:11Z
---

Add comprehensive htmx attribute validation beyond the current basic attribute recognition. Currently htmx validation only checks if attributes exist and are version-appropriate on `<input>` elements. This feature would add semantic validation of attribute values and patterns.

## Background

Current state (`rules/htmx.go`, `rules/input_attributes.go`):
- Validates htmx attributes only on `<input>` elements
- Checks version compatibility (v2 vs v4)
- Warns on deprecated/unknown attributes
- `attribute_misuse.go` does NOT skip hx-* attributes (causes false positives)

## Proposed Validation Rules

### Phase 1: Fix gaps and add value validation

- [x] Make `attribute_misuse` htmx-aware (skip hx-* when htmx enabled)
- [x] Validate hx-swap values: `innerHTML`, `outerHTML`, `beforebegin`, `afterbegin`, `beforeend`, `afterend`, `delete`, `none` + modifiers like `swap:1s`, `settle:500ms`, `scroll:top`
- [x] Validate hx-trigger syntax: event names, modifiers (`once`, `changed`, `delay:Xs`, `throttle:Xs`, `from:selector`, etc.)
- [x] Validate hx-target values: CSS selectors, special values (`this`, `closest`, `find`, `next`, `previous`)

### Phase 2: Pattern matching and semantic checks

- [x] Validate hx-on:* event names match known htmx events (e.g., `hx-on:htmx:after-request` vs typos)
- [x] Warn on hx-post/hx-get on submit buttons instead of parent form (common mistake)
- [x] Validate hx-vals/hx-headers JSON syntax
- [x] Check hx-include CSS selector validity
- [x] Validate timing values in modifiers (e.g., `delay:1s`, `throttle:500ms` - valid time formats) *(done in Phase 1)*

### Phase 3: htmx 4 specific

- [x] Validate hx-status:* patterns for valid HTTP status codes
- [ ] Check new htmx 4 event naming convention compliance
- [ ] Validate :inherited suffix usage

## Research Notes

- htmx 4.0 uses colon-separated event naming: `htmx:phase:action[:sub-action]`
- Common mistake: putting hx-post on button instead of form bypasses htmx validation
- hx-trigger modifiers have specific syntax: `delay:1s`, `throttle:500ms`, `from:closest form`
- Known htmx events: `htmx:load`, `htmx:configRequest`, `htmx:beforeRequest`, `htmx:afterRequest`, `htmx:beforeSwap`, `htmx:afterSwap`, etc.

## Implementation Approach

Could use regex for value validation. Consider a new dedicated rule like `htmx-attributes` that handles all htmx-specific validation, separate from `input_attributes` and `attribute_misuse`.

## References

- htmx hx-trigger: https://htmx.org/attributes/hx-trigger/
- htmx hx-swap: https://htmx.org/attributes/hx-swap/
- htmx events: https://htmx.org/events/
- htmx 4 changes: https://four.htmx.org/htmx-4/