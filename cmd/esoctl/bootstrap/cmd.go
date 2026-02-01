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

// Package bootstrap provides commands for bootstrapping new resources.
package bootstrap

import (
	"github.com/spf13/cobra"

	"github.com/external-secrets/external-secrets/cmd/esoctl/bootstrap/generator"
	"github.com/external-secrets/external-secrets/cmd/esoctl/bootstrap/provider"
)

// NewBootstrapCommand creates the bootstrap command with all subcommands.
func NewBootstrapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap new resources for external-secrets",
		Long:  `Bootstrap new resources like generators and providers for external-secrets operator.`,
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Usage()
		},
	}

	// Register subcommands
	cmd.AddCommand(generator.NewGeneratorCommand())
	cmd.AddCommand(provider.NewProviderCommand())

	return cmd
}
