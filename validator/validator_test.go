package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateFile_ValidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "valid.yaml")
	content := []byte("name: hello\nversion: 1\nitems:\n  - one\n  - two\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	r := ValidateFile(path)
	if !r.Valid {
		t.Errorf("expected valid, got error: %s (line %d)", r.Error, r.Line)
	}
}

func TestValidateFile_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.yaml")
	content := []byte("name: hello\n  - bad indentation\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	r := ValidateFile(path)
	if r.Valid {
		t.Error("expected invalid, got valid")
	}
}

func TestValidateFile_FileNotFound(t *testing.T) {
	r := ValidateFile("/nonexistent/path/file.yaml")
	if r.Valid {
		t.Error("expected invalid for nonexistent file")
	}
}

func TestValidateData_ValidYAML(t *testing.T) {
	data := []byte("key: value\ncount: 42\n")
	r := ValidateData(data)
	if !r.Valid {
		t.Errorf("expected valid, got: %s", r.Error)
	}
}

func TestValidateData_InvalidYAML(t *testing.T) {
	data := []byte("key: value\n\tbad tab\n")
	r := ValidateData(data)
	if r.Valid {
		t.Error("expected invalid for tab-indented YAML")
	}
}

func TestValidateData_EmptyYAML(t *testing.T) {
	r := ValidateData([]byte(""))
	if !r.Valid {
		t.Errorf("empty document should be valid, got: %s", r.Error)
	}
}
