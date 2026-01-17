# go-html-validate

Fast HTML linter written in Go with special handling for Go templates. Validates HTML for accessibility, best practices, and common mistakes.

Implements rules from [html-validate.org](https://html-validate.org).

## Installation

```bash
go install github.com/STR-Consulting/go-html-validate@latest
```

## Usage

```bash
# Lint files or directories
go-html-validate web/
go-html-validate index.html about.html

# Errors only (no warnings)
go-html-validate -q web/

# JSON output
go-html-validate --format=json web/

# Disable specific rules
go-html-validate --disable=prefer-aria --disable=no-inline-style web/

# Ignore files by pattern
go-html-validate --ignore="*_test.html" web/

# List available rules
go-html-validate --list-rules
```

## Options

| Flag | Description |
|------|-------------|
| `-f, --format` | Output format: `text` (default), `json` |
| `-q, --quiet` | Only show errors, suppress warnings |
| `--no-color` | Disable colored output |
| `--ignore PATTERN` | Glob pattern to ignore (repeatable) |
| `--disable RULE` | Disable specific rule (repeatable) |
| `--list-rules` | List all available rules |
| `-h, --help` | Show help |
| `--config PATH` | Use specific config file |
| `--no-config` | Disable config file loading |
| `--print-config` | Print resolved configuration |

## Configuration

This tool uses the same configuration format as [html-validate](https://html-validate.org/usage/index.html).

### Config File

Create `.htmlvalidate.json` in your project root:

```json
{
  "extends": ["html-validate:recommended"],
  "rules": {
    "no-inline-style": "warn",
    "prefer-tbody": "off"
  }
}
```

The linter searches for `.htmlvalidate.json` in the target directory and parent directories.

### Rule Severity

- `"error"` or `2` - Error (fails CI)
- `"warn"` or `1` - Warning
- `"off"` or `0` - Disabled

### Built-in Presets

| Preset | Description |
|--------|-------------|
| `html-validate:recommended` | All rules enabled (default) |
| `html-validate:standard` | Core rules, fewer style preferences |
| `html-validate:a11y` | Accessibility-focused rules only |

### Ignore File

Create `.htmlvalidateignore` for gitignore-style patterns:

```
node_modules/
vendor/
**/*.generated.html
```

For full configuration options, see the [html-validate configuration documentation](https://html-validate.org/usage/index.html).

## Supported File Types

- `.html`
- `.htm`
- `.gohtml`
- `.tmpl`

## Rule Categories

### Accessibility (WCAG)
- `area-alt` - `<area>` elements must have alt text
- `aria-hidden-body` - `<body>` must not have aria-hidden
- `aria-label-misuse` - aria-label only on interactive elements
- `button-name` - Buttons must have accessible names
- `heading-content` - Headings must have text content
- `heading-level` - Heading levels must not be skipped
- `hidden-focusable` - Hidden elements must not be focusable
- `img-alt` - Images must have alt attributes
- `input-label` - Form inputs must have labels
- `link-name` - Links must have accessible names
- `meta-refresh` - Avoid meta refresh redirects
- `multiple-labeled-controls` - Labels must reference single controls
- `no-abstract-role` - No abstract ARIA roles
- `no-autoplay` - Avoid autoplay on media
- `no-redundant-role` - No redundant ARIA roles
- `prefer-native-element` - Prefer native HTML over ARIA
- `require-lang` - `<html>` must have lang attribute
- `svg-focusable` - SVGs must have focusable="false"
- `tabindex` - Avoid positive tabindex values
- `unique-landmark` - Landmark regions must be unique

### Validation
- `attribute-allowed-values` - Valid attribute values
- `attribute-misuse` - Attributes used correctly
- `doctype` - Document must have DOCTYPE
- `duplicate-id` - IDs must be unique
- `element-name` - Valid element names
- `element-permitted-content` - Valid child elements
- `element-permitted-occurrences` - Element count limits
- `element-permitted-order` - Correct element order
- `element-permitted-parent` - Valid parent elements
- `element-required-ancestor` - Required ancestor elements
- `element-required-attributes` - Required attributes present
- `element-required-content` - Required child content
- `no-dup-attr` - No duplicate attributes
- `no-dup-class` - No duplicate classes
- `valid-autocomplete` - Valid autocomplete values
- `valid-id` - Valid ID syntax
- `void-content` - Void elements have no content

### Deprecated
- `deprecated` - No deprecated elements
- `no-deprecated-attr` - No deprecated attributes
- `no-conditional-comment` - No IE conditional comments

### Best Practices
- `button-type` - Buttons should have explicit type
- `empty-title` - Title elements must not be empty
- `form-dup-name` - Unique form control names
- `form-submit` - Forms should have submit buttons
- `long-title` - Avoid overly long titles
- `map-dup-name` - Unique map names
- `map-id-name` - Map id and name should match
- `no-implicit-input-type` - Explicit input types
- `no-missing-references` - Valid ID references
- `no-multiple-main` - Single main element
- `no-redundant-for` - No redundant label for
- `no-utf8-bom` - No UTF-8 BOM
- `prefer-aria` - Use ARIA attributes
- `prefer-button` - Prefer button over input
- `prefer-semantic` - Use semantic elements
- `prefer-tbody` - Tables should have tbody
- `script-element` - Valid script elements
- `script-type` - Valid script types
- `tel-non-breaking` - Tel links with proper spacing

### Security
- `allowed-links` - Validate link protocols
- `no-inline-style` - Avoid inline styles
- `no-style-tag` - Avoid style tags
- `require-csp-nonce` - CSP nonce on scripts/styles
- `require-sri` - Subresource integrity

## License

MIT - See [LICENSE](LICENSE) for details.

Rule specifications derived from [html-validate](https://html-validate.org) (MIT License).
