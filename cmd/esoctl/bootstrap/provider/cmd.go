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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/external-secrets/external-secrets/cmd/esoctl/bootstrap/common"
)

var (
	providerName        string
	providerDescription string
	providerPackage     string
)

// NewProviderCommand creates the provider bootstrap command.
func NewProviderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Bootstrap a new secret provider",
		Long:  `Bootstrap a new secret provider with implementation, tests, and registration.`,
		RunE:  runProviderBootstrap,
	}

	cmd.Flags().StringVar(&providerName, "name", "", "Name of the provider (e.g., Dummy, MyProvider)")
	cmd.Flags().StringVar(&providerDescription, "description", "", "Description of the provider")
	cmd.Flags().StringVar(&providerPackage, "package", "", "Package name (default: lowercase of name)")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runProviderBootstrap(_ *cobra.Command, _ []string) error {
	// Validate provider name (PascalCase)
	if !common.ValidatePascalCase(providerName) {
		return fmt.Errorf("provider name must be PascalCase and start with an uppercase letter")
	}

	// Set default package name if not provided
	if providerPackage == "" {
		providerPackage = strings.ToLower(providerName)
	}

	// Set default description if not provided
	if providerDescription == "" {
		providerDescription = fmt.Sprintf("%s secret provider", providerName)
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

	// Create provider configuration
	cfg := Config{
		ProviderName: providerName,
		PackageName:  providerPackage,
		Description:  providerDescription,
	}

	// Bootstrap the provider
	if err := Bootstrap(rootDir, cfg); err != nil {
		return err
	}

	fmt.Printf("✓ Successfully bootstrapped provider: %s\n", providerName)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Review and customize the generated code in: providers/v1/%s/%s.go\n", providerPackage, providerPackage)
	fmt.Printf("2. Add your provider configuration to: apis/externalsecrets/v1/secretstore_types.go\n")
	fmt.Printf("3. Implement the TODO sections in the provider code\n")
	fmt.Printf("4. Add authentication fields to the provider spec\n")
	fmt.Printf("5. Create registration file: pkg/register/%s.go\n", providerPackage)
	fmt.Printf("6. Run: cd providers/v1/%s && go mod tidy\n", providerPackage)
	fmt.Printf("7. Run: make test (from repository root)\n")
	fmt.Printf("8. Add integration tests and documentation\n")

	return nil
}
