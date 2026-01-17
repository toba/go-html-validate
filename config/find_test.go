package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/STR-Consulting/go-html-validate/config"
)

func TestFindConfigFile(t *testing.T) {
	root, sub, deep := setupTestDirs(t)
	configPath := filepath.Join(sub, config.ConfigFileName)
	if err := os.WriteFile(configPath, []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}

	testFindFile(t, config.FindConfigFile, root, sub, deep, configPath)
}

func TestFindIgnoreFile(t *testing.T) {
	root, sub, deep := setupTestDirs(t)
	ignorePath := filepath.Join(sub, config.IgnoreFileName)
	if err := os.WriteFile(ignorePath, []byte("*.tmp"), 0o600); err != nil {
		t.Fatal(err)
	}

	testFindFile(t, config.FindIgnoreFile, root, sub, deep, ignorePath)
}

// testFindFile is a generic test helper for file finding functions.
func testFindFile(t *testing.T, findFn func(string) (string, error), root, sub, deep, expectedPath string) {
	t.Helper()

	tests := []struct {
		name      string
		searchDir string
		wantPath  string
	}{
		{"find in same directory", sub, expectedPath},
		{"find in parent directory", deep, expectedPath},
		{"not found in root", root, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findFn(tt.searchDir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantPath {
				t.Errorf("got %q, want %q", got, tt.wantPath)
			}
		})
	}
}

// setupTestDirs creates a test directory structure: root/sub/deep
func setupTestDirs(t *testing.T) (root, sub, deep string) {
	t.Helper()
	root = t.TempDir()
	sub = filepath.Join(root, "sub")
	deep = filepath.Join(sub, "deep")
	if err := os.MkdirAll(deep, 0o750); err != nil {
		t.Fatal(err)
	}
	return root, sub, deep
}
