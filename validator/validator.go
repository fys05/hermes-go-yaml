// Package validator provides YAML file format validation.
package validator

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Result holds the outcome of validating a single YAML file.
type Result struct {
	Path   string
	Valid  bool
	Error  string
	Line   int
	Column int
}

// ValidateFile parses a YAML file and returns a validation result.
// It checks both syntax and structural validity via yaml.v3 unmarshaling.
func ValidateFile(path string) Result {
	data, err := os.ReadFile(path)
	if err != nil {
		return Result{
			Path:  path,
			Valid: false,
			Error: fmt.Sprintf("cannot read file: %v", err),
		}
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		line, col := extractPosition(err)
		return Result{
			Path:   path,
			Valid:  false,
			Error:  err.Error(),
			Line:   line,
			Column: col,
		}
	}

	return Result{
		Path:  path,
		Valid: true,
	}
}

// ValidateData parses raw YAML bytes and returns a validation result.
func ValidateData(data []byte) Result {
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		line, col := extractPosition(err)
		return Result{
			Valid:  false,
			Error:  err.Error(),
			Line:   line,
			Column: col,
		}
	}
	return Result{Valid: true}
}

// extractPosition tries to extract line and column from a yaml parse error.
func extractPosition(err error) (int, int) {
	type positioner interface {
		Position() (int, int)
	}
	if p, ok := err.(positioner); ok {
		line, col := p.Position()
		return line, col
	}
	return 0, 0
}
