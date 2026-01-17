package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/STR-Consulting/go-html-validate/config"
)

func TestLoadIgnoreFile(t *testing.T) {
	dir := t.TempDir()
	content := `# This is a comment
node_modules/
vendor/

**/*.generated.html
*.min.html
`
	path := filepath.Join(dir, config.IgnoreFileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	patterns, err := config.LoadIgnoreFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"node_modules/",
		"vendor/",
		"**/*.generated.html",
		"*.min.html",
	}

	if len(patterns) != len(expected) {
		t.Fatalf("got %d patterns, want %d", len(patterns), len(expected))
	}

	for i, p := range patterns {
		if p != expected[i] {
			t.Errorf("pattern[%d] = %q, want %q", i, p, expected[i])
		}
	}
}

func TestMatchesIgnorePattern(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{
			name:     "simple glob match",
			path:     "test.min.html",
			patterns: []string{"*.min.html"},
			want:     true,
		},
		{
			name:     "simple glob no match",
			path:     "test.html",
			patterns: []string{"*.min.html"},
			want:     false,
		},
		{
			name:     "directory pattern match",
			path:     "node_modules/package/index.js",
			patterns: []string{"node_modules/"},
			want:     true,
		},
		{
			name:     "directory pattern match nested",
			path:     "src/node_modules/package/index.js",
			patterns: []string{"node_modules/"},
			want:     true,
		},
		{
			name:     "directory pattern no match",
			path:     "src/modules/index.js",
			patterns: []string{"node_modules/"},
			want:     false,
		},
		{
			name:     "doublestar pattern match",
			path:     "src/components/Button.generated.html",
			patterns: []string{"**/*.generated.html"},
			want:     true,
		},
		{
			name:     "doublestar pattern root match",
			path:     "test.generated.html",
			patterns: []string{"**/*.generated.html"},
			want:     true,
		},
		{
			name:     "doublestar pattern no match",
			path:     "src/components/Button.html",
			patterns: []string{"**/*.generated.html"},
			want:     false,
		},
		{
			name:     "multiple patterns first match",
			path:     "vendor/lib.js",
			patterns: []string{"node_modules/", "vendor/"},
			want:     true,
		},
		{
			name:     "multiple patterns second match",
			path:     "node_modules/lib.js",
			patterns: []string{"vendor/", "node_modules/"},
			want:     true,
		},
		{
			name:     "prefix doublestar",
			path:     "dist/js/bundle.js",
			patterns: []string{"dist/**"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.MatchesIgnorePattern(tt.path, tt.patterns)
			if got != tt.want {
				t.Errorf("MatchesIgnorePattern(%q, %v) = %v, want %v",
					tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}
