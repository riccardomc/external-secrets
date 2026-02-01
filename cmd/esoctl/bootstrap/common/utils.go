/*
Copyright © 2025 ESO Maintainer Team

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package common provides shared utilities for bootstrap operations.
package common

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// FindRootDir finds the root directory of the external-secrets repository.
func FindRootDir(startDir string) string {
	dir := startDir
	for {
		// Check if go.mod exists and contains external-secrets
		goModPath := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(filepath.Clean(goModPath)); err == nil {
			if strings.Contains(string(data), "module github.com/external-secrets/external-secrets") {
				return dir
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// CreateFromTemplate creates a file from an embedded template.
func CreateFromTemplate(templates embed.FS, tmplPath, outputFile string, data interface{}) error {
	tmplContent, err := templates.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	tmpl := template.Must(template.New(filepath.Base(tmplPath)).Parse(string(tmplContent)))
	f, err := os.Create(filepath.Clean(outputFile))
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return tmpl.Execute(f, data)
}

// ValidatePascalCase validates that a string is in PascalCase format.
func ValidatePascalCase(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Must start with uppercase letter
	if s[0] < 'A' || s[0] > 'Z' {
		return false
	}
	// Only alphanumeric characters allowed
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}

// EnsureDirectory creates a directory if it doesn't exist.
func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0o750)
}

// FileExists checks if a file exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
