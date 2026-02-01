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

package provider

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/external-secrets/external-secrets/cmd/esoctl/bootstrap/common"
)

//go:embed templates/*.tmpl
var templates embed.FS

// Config holds the configuration for bootstrapping a provider.
type Config struct {
	ProviderName string // e.g., "Dummy"
	PackageName  string // e.g., "dummy"
	Description  string // e.g., "Dummy secret provider"
}

// Bootstrap creates a new provider with implementation and test files.
func Bootstrap(rootDir string, cfg Config) error {
	// Create provider implementation directory and files
	if err := createProviderImplementation(rootDir, cfg); err != nil {
		return fmt.Errorf("failed to create provider implementation: %w", err)
	}

	return nil
}

func createProviderImplementation(rootDir string, cfg Config) error {
	providerDir := filepath.Join(rootDir, "providers", "v1", cfg.PackageName)

	// Create provider directory
	if err := common.EnsureDirectory(providerDir); err != nil {
		return fmt.Errorf("failed to create provider directory: %w", err)
	}

	// Create main provider file
	providerFile := filepath.Join(providerDir, fmt.Sprintf("%s.go", cfg.PackageName))
	if common.FileExists(providerFile) {
		return fmt.Errorf("provider file already exists: %s", providerFile)
	}

	if err := common.CreateFromTemplate(templates, "templates/provider.go.tmpl", providerFile, cfg); err != nil {
		return fmt.Errorf("failed to create provider file: %w", err)
	}
	fmt.Printf("✓ Created provider implementation: %s\n", providerFile)

	// Create test file
	testFile := filepath.Join(providerDir, fmt.Sprintf("%s_test.go", cfg.PackageName))
	if err := common.CreateFromTemplate(templates, "templates/provider_test.go.tmpl", testFile, cfg); err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	fmt.Printf("✓ Created test file: %s\n", testFile)

	// Create go.mod
	goModFile := filepath.Join(providerDir, "go.mod")
	if err := common.CreateFromTemplate(templates, "templates/provider_go.mod.tmpl", goModFile, cfg); err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}
	fmt.Printf("✓ Created go.mod: %s\n", goModFile)

	// Create empty go.sum
	goSumFile := filepath.Join(providerDir, "go.sum")
	if err := os.WriteFile(goSumFile, []byte(""), 0o600); err != nil {
		return fmt.Errorf("failed to create go.sum: %w", err)
	}
	fmt.Printf("✓ Created go.sum: %s\n", goSumFile)

	return nil
}
