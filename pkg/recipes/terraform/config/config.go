/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/radius-project/radius/pkg/recipes"
	"github.com/radius-project/radius/pkg/recipes/recipecontext"
	"github.com/radius-project/radius/pkg/recipes/terraform/config/backends"
	"github.com/radius-project/radius/pkg/recipes/terraform/config/providers"
	"github.com/radius-project/radius/pkg/ucp/ucplog"
)

const (
	// modeConfigFile is read/write mode only for the owner of the TF config file.
	modeConfigFile fs.FileMode = 0600
)

// New creates TerraformConfig with the given module name and its inputs (module source, version, parameters)
// Parameters are populated from environment recipe and resource recipe metadata.
func New(ctx context.Context, moduleName string, envRecipe *recipes.EnvironmentDefinition, resourceRecipe *recipes.ResourceMetadata) (*TerraformConfig, error) {
	// Resource parameter gets precedence over environment level parameter,
	// if same parameter is defined in both environment and resource recipe metadata.
	moduleData := newModuleConfig(envRecipe.TemplatePath, envRecipe.TemplateVersion, envRecipe.Parameters, resourceRecipe.Parameters)

	return &TerraformConfig{
		Terraform: nil,
		Provider:  nil,
		Module: map[string]TFModuleConfig{
			moduleName: moduleData,
		},
	}, nil
}

// getMainConfigFilePath returns the path of the Terraform main config file.
func getMainConfigFilePath(workingDir string) string {
	return fmt.Sprintf("%s/%s", workingDir, mainConfigFileName)
}

// Save writes the Terraform config to main.tf.json file in the working directory.
// This overwrites the existing file if it exists.
func (cfg *TerraformConfig) Save(ctx context.Context, workingDir string) error {
	logger := ucplog.FromContextOrDiscard(ctx)

	// Write the JSON data to a file in the working directory.
	// JSON configuration syntax for Terraform requires the file to be named with .tf.json suffix.
	// https://developer.hashicorp.com/terraform/language/syntax/json

	// Create a buffer to write the JSON to
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ") // Indent with 2 spaces to make the JSON file human-readable and consistent with codebase.

	// Encode the Terraform config to JSON. JSON encoding is being used to ensure that special characters
	// in the original text are preserved when writing to the file.
	// For example, when writing this text to file with JSON encoding (using enc.Encode(cfg)),
	//   the special characters in the following text will be preserved:
	//	"required_providers": {
	//			"aws": {
	//				"source": "hashicorp/aws",
	//				"version": ">= 3.0"
	//			},
	//		}
	//   However, if we were to write the text directly to the file without JSON encoding, t
	//	the special characters would be escaped and be written as follows:
	//	"required_providers": {
	//		"aws": {
	//				"source": "hashicorp/aws",
	//				"version": "\u003e= 2.0"
	//			},
	//		}
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	logger.Info(fmt.Sprintf("Writing Terraform JSON config to file: %s", getMainConfigFilePath(workingDir)))
	if err := os.WriteFile(getMainConfigFilePath(workingDir), buf.Bytes(), modeConfigFile); err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	return nil
}

// AddProviders adds provider configurations to Terraform configuration file based on input of environment recipe configuration, requiredProviders and ucpConfiguredProviders.
// It also updates module provider block if aliases exist and required_provider configuration to the file.
// Save() must be called to save the generated providers config. requiredProviders contains a list of provider names
// that are required for the module.
func (cfg *TerraformConfig) AddProviders(ctx context.Context, requiredProviders map[string]*RequiredProviderInfo, ucpConfiguredProviders map[string]providers.Provider, envConfig *recipes.Configuration, secrets map[string]recipes.SecretData) error {
	logger := ucplog.FromContextOrDiscard(ctx)
	providerConfigs, err := getProviderConfigs(ctx, requiredProviders, ucpConfiguredProviders, envConfig, secrets)
	if err != nil {
		return err
	}

	// Add generated provider configs for required providers to the existing terraform json config file
	if len(providerConfigs) > 0 {
		cfg.Provider = providerConfigs
	}

	// Update module configuration with aliased provider names, if they exist.
	logger.Info("Updating module config with providers aliases")
	if err := cfg.updateModuleWithProviderAliases(requiredProviders); err != nil {
		return err
	}

	// Set the required providers for the Terraform configuration.
	logger.Info("Update Terraform configuration with required providers")
	if cfg.Terraform == nil {
		cfg.Terraform = &TerraformDefinition{}
	}
	cfg.Terraform.RequiredProviders = requiredProviders

	return nil
}

// updateModuleWithProviderAliases updates the module provider configuration in the Terraform config
// by adding aliases to the provider configurations.
// https://developer.hashicorp.com/terraform/language/syntax/json#module-blocks
func (cfg *TerraformConfig) updateModuleWithProviderAliases(requiredProviders map[string]*RequiredProviderInfo) error {
	if cfg == nil {
		return fmt.Errorf("terraform configuration is not initialized")
	}
	moduleAliasConfig := map[string]string{}

	for providerName, providerConfigList := range cfg.Provider {
		// For each provider in the providerConfigs, if provider has a property "alias",
		// add entry to the module provider configuration.
		// Provider configurations (those with the alias argument set) are never inherited automatically by modules,
		// and so must always be passed explicitly using the providers map.
		// https://developer.hashicorp.com/terraform/language/modules/develop/providers#legacy-shared-modules-with-provider-configurations

		// Note: We're building configuration from user input, we're mapping the provider.alias names in
		// the required provider configuration (ConfigurationAliases) to the environment recipe provider configuration data.
		// This is being done to ensure that the provider configuration is passed to the module correctly.

		for _, providerConfig := range providerConfigList {
			if alias, ok := providerConfig["alias"]; ok {
				aliasProviderConfig := providerName + "." + fmt.Sprintf("%v", alias)

				// Check if the alias is in the required providers' configuration aliases. If there is a match, add the alias to the module provider configuration.
				if requiredProviders[providerName] != nil && len(requiredProviders[providerName].ConfigurationAliases) > 0 {
					for _, alias := range requiredProviders[providerName].ConfigurationAliases {
						if alias == aliasProviderConfig {
							moduleAliasConfig[alias] = alias
							break
						}
					}
				}
			}
		}
	}

	// Update the module provider configuration in the Terraform config.
	if len(moduleAliasConfig) > 0 {
		moduleConfig := cfg.Module
		for _, module := range moduleConfig {
			module["providers"] = moduleAliasConfig
		}
	}

	return nil
}

// AddRecipeContext adds RecipeContext to TerraformConfig module parameters if recipeCtx is not nil.
// Save() must be called after adding recipe context to the module config.
func (cfg *TerraformConfig) AddRecipeContext(ctx context.Context, moduleName string, recipeCtx *recipecontext.Context) error {
	mod, ok := cfg.Module[moduleName]
	if !ok {
		// must not happen because module key is set when the config is initialized in New().
		return fmt.Errorf("module %q not found in the initialized terraform config", moduleName)
	}

	if recipeCtx != nil {
		mod.SetParams(RecipeParams{recipecontext.RecipeContextParamKey: recipeCtx})
	}

	return nil
}

// newModuleConfig creates a new TFModuleConfig object with the given module source and version
// and also populates RecipeParams in TF module config. If same parameter key exists across params
// then the last map specified gets precedence.
func newModuleConfig(moduleSource string, moduleVersion string, params ...RecipeParams) TFModuleConfig {
	moduleConfig := TFModuleConfig{
		moduleSourceKey: moduleSource,
	}

	// Not all sources use versions, so only add the version if it's specified.
	// Registries require versions, but HTTP or filesystem sources do not.
	if moduleVersion != "" {
		moduleConfig[moduleVersionKey] = moduleVersion
	}

	// Populate recipe parameters
	for _, param := range params {
		moduleConfig.SetParams(param)
	}

	return moduleConfig
}

// getProviderConfigs generates the Terraform provider configurations. This is built from a combination of environment level recipe configuration for
// providers and the provider configurations registered with UCP. The environment level recipe configuration for providers takes precedence over UCP provider configurations.
// The function returns a map where the keys are provider names and the values are slices of maps.
// Each map in the slice represents a specific configuration for the corresponding provider.
// This structure allows for multiple configurations per provider.
func getProviderConfigs(ctx context.Context, requiredProviders map[string]*RequiredProviderInfo, ucpConfiguredProviders map[string]providers.Provider, envConfig *recipes.Configuration, secrets map[string]recipes.SecretData) (map[string][]map[string]any, error) {
	// Get recipe provider configurations from the environment configuration
	providerConfigs, err := providers.GetRecipeProviderConfigs(ctx, envConfig, secrets)
	if err != nil {
		return nil, err
	}

	// Build provider configurations for required providers excluding the ones already present in providerConfigs (environment level configuration).
	// Required providers that are not configured with UCP will be skipped.
	for provider := range requiredProviders {
		if _, ok := providerConfigs[provider]; ok {
			// Environment level recipe configuration for providers will take precedence over
			// UCP provider configuration (currently these include azurerm, aws, kubernetes providers)
			continue
		}

		builder, ok := ucpConfiguredProviders[provider]
		if !ok {
			// No-op: For any other provider under required_providers, Radius doesn't generate any custom configuration.
			continue
		}

		config, err := builder.BuildConfig(ctx, envConfig)
		if err != nil {
			return nil, err
		}

		if len(config) > 0 {
			providerConfigs[provider] = []map[string]any{config}
		}
	}

	return providerConfigs, nil
}

// AddTerraformBackend adds backend configurations to store Terraform state file for the deployment.
// Save() must be called to save the generated backend config.
// Currently, the supported backend for Terraform Recipes is Kubernetes secret. https://developer.hashicorp.com/terraform/language/settings/backends/kubernetes
func (cfg *TerraformConfig) AddTerraformBackend(resourceRecipe *recipes.ResourceMetadata, backend backends.Backend) (map[string]any, error) {
	backendConfig, err := backend.BuildBackend(resourceRecipe)
	if err != nil {
		return nil, err
	}

	if cfg.Terraform == nil {
		cfg.Terraform = &TerraformDefinition{}
	}
	cfg.Terraform.Backend = backendConfig

	return backendConfig, nil
}

// Add outputs to the config file referencing module outputs to populate expected Radius resource outputs.
// Outputs of modules are accessible through this format: module.<MODULE NAME>.<OUTPUT NAME>
// https://developer.hashicorp.com/terraform/language/modules/syntax#accessing-module-output-values
// This function only updates config in memory, Save() must be called to persist the updated config.
func (cfg *TerraformConfig) AddOutputs(localModuleName string) error {
	if localModuleName == "" {
		return errors.New("module name cannot be empty")
	}

	cfg.Output = map[string]any{
		recipes.ResultPropertyName: map[string]any{
			"value":     "${module." + localModuleName + "." + recipes.ResultPropertyName + "}",
			"sensitive": true, // since secret and non-secret values are combined in the result, mark the entire output sensitive
		},
	}

	return nil
}
