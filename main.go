// hermes-go-yaml — a YAML file format validator.
//
// Usage:
//
//	yaml-validator [file1.yaml file2.yaml ...]
//
// If no files are given, reads YAML from stdin.
// Exits 0 when all files are valid, 1 otherwise.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fys05/hermes-go-yaml/validator"
)

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) > 1 {
		return validateFiles(os.Args[1:])
	}
	return validateStdin()
}

func validateFiles(paths []string) int {
	allValid := true
	for _, path := range paths {
		r := validator.ValidateFile(path)
		if r.Valid {
			fmt.Printf("✔ %s — valid YAML\n", r.Path)
		} else {
			allValid = false
			if r.Line > 0 {
				fmt.Printf("✘ %s:%d:%d — %s\n", r.Path, r.Line, r.Column, r.Error)
			} else {
				fmt.Printf("✘ %s — %s\n", r.Path, r.Error)
			}
		}
	}
	if allValid {
		return 0
	}
	return 1
}

func validateStdin() int {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		return 1
	}
	r := validator.ValidateData(data)
	if r.Valid {
		fmt.Println("✔ stdin — valid YAML")
		return 0
	}
	if r.Line > 0 {
		fmt.Printf("✘ stdin:%d:%d — %s\n", r.Line, r.Column, r.Error)
	} else {
		fmt.Printf("✘ stdin — %s\n", r.Error)
	}
	return 1
}
