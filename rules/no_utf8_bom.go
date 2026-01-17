package rules

import (
	"github.com/STR-Consulting/go-html-validate/parser"
)

// NoUTF8BOM checks that files don't start with a UTF-8 BOM.
type NoUTF8BOM struct{}

// Name returns the rule identifier.
func (r *NoUTF8BOM) Name() string { return RuleNoUTF8BOM }

// Description returns what this rule checks.
func (r *NoUTF8BOM) Description() string {
	return "files should not have UTF-8 BOM"
}

// Check examines the document for UTF-8 BOM.
// Note: This check needs raw file content, which the parser may strip.
// The linter should check for BOM before parsing if needed.
func (r *NoUTF8BOM) Check(doc *parser.Document) []Result {
	// BOM checking should be done at the file read level, not after parsing.
	// The HTML parser typically handles/strips the BOM.
	// This rule serves as a placeholder; actual BOM detection should be
	// implemented in the linter's file reading logic.
	return nil
}

// HasUTF8BOM checks if content starts with UTF-8 BOM bytes.
// This helper can be used by the linter before parsing.
func HasUTF8BOM(content []byte) bool {
	// UTF-8 BOM is EF BB BF
	return len(content) >= 3 &&
		content[0] == 0xEF &&
		content[1] == 0xBB &&
		content[2] == 0xBF
}
