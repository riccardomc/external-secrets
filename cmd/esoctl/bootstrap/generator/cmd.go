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

package generator

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/external-secrets/external-secrets/cmd/esoctl/bootstrap/common"
)

var (
	generatorName        string
	generatorDescription string
	generatorPackage     string
)

// NewGeneratorCommand creates the generator bootstrap command.
func NewGeneratorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generator",
		Short: "Bootstrap a new generator",
		Long:  `Bootstrap a new generator with CRD definition and provider implementation.`,
		RunE:  runGeneratorBootstrap,
	}

	cmd.Flags().StringVar(&generatorName, "name", "", "Name of the generator (e.g., MyGenerator)")
	cmd.Flags().StringVar(&generatorDescription, "description", "", "Description of the generator")
	cmd.Flags().StringVar(&generatorPackage, "package", "", "Package name (default: lowercase of name)")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runGeneratorBootstrap(_ *cobra.Command, _ []string) error {
	// Validate generator name
	if !common.ValidatePascalCase(generatorName) {
		return fmt.Errorf("generator name must be PascalCase and start with an uppercase letter")
	}

	// Set default package name if not provided
	if generatorPackage == "" {
		generatorPackage = strings.ToLower(generatorName)
	}

	// Set default description if not provided
	if generatorDescription == "" {
		generatorDescription = fmt.Sprintf("%s generator", generatorName)
	}

	// Get root directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Try to find the root directory
	rootDir := common.FindRootDir(wd)
	if rootDir == "" {
		return fmt.Errorf("could not find repository root directory")
	}

	// Create generator configuration
	cfg := Config{
		GeneratorName: generatorName,
		PackageName:   generatorPackage,
		Description:   generatorDescription,
		GeneratorKind: "GeneratorKind" + generatorName,
	}

	// Bootstrap the generator
	if err := Bootstrap(rootDir, cfg); err != nil {
		return err
	}

	fmt.Printf("✓ Successfully bootstrapped generator: %s\n", generatorName)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Review and customize: apis/generators/v1alpha1/types_%s.go\n", generatorPackage)
	fmt.Printf("2. Implement the generator logic in: generators/v1/%s/%s.go\n", generatorPackage, generatorPackage)
	fmt.Printf("3. Run: go mod tidy\n")
	fmt.Printf("4. Run: make generate\n")
	fmt.Printf("5. Run: make manifests\n")
	fmt.Printf("6. Add tests for your generator\n")

	return nil
}
