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

package deploy

import (
	"context"
	"fmt"
	"testing"

	"github.com/radius-project/radius/pkg/cli/bicep"
	"github.com/radius-project/radius/pkg/cli/clients"
	"github.com/radius-project/radius/pkg/cli/config"
	"github.com/radius-project/radius/pkg/cli/connections"
	"github.com/radius-project/radius/pkg/cli/deploy"
	"github.com/radius-project/radius/pkg/cli/framework"
	"github.com/radius-project/radius/pkg/cli/output"
	"github.com/radius-project/radius/pkg/cli/workspaces"
	"github.com/radius-project/radius/pkg/corerp/api/v20231001preview"
	"github.com/radius-project/radius/pkg/to"
	"github.com/radius-project/radius/test/radcli"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_CommandValidation(t *testing.T) {
	radcli.SharedCommandValidation(t, NewCommand)
}

func Test_Validate(t *testing.T) {
	configWithWorkspace := radcli.LoadConfigWithWorkspace(t)
	testcases := []radcli.ValidateInput{

		{
			Name:          "rad deploy - valid",
			Input:         []string{"app.bicep"},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
			ConfigureMocks: func(mocks radcli.ValidateMocks) {
				mocks.ApplicationManagementClient.EXPECT().
					GetEnvironment(gomock.Any(), "/planes/radius/local/resourceGroups/test-resource-group/providers/Applications.Core/environments/test-environment").
					Return(v20231001preview.EnvironmentResource{}, nil).
					Times(1)
			},
		},
		{
			Name:          "rad deploy - valid with parameters",
			Input:         []string{"app.bicep", "-p", "foo=bar", "--parameters", "a=b"},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
			ConfigureMocks: func(mocks radcli.ValidateMocks) {
				mocks.ApplicationManagementClient.EXPECT().
					GetEnvironment(gomock.Any(), radcli.TestEnvironmentID).
					Return(v20231001preview.EnvironmentResource{}, nil).
					Times(1)

			},
		},
		{
			Name:          "rad deploy - valid with environment",
			Input:         []string{"app.bicep", "-e", "prod"},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
			ConfigureMocks: func(mocks radcli.ValidateMocks) {
				mocks.ApplicationManagementClient.EXPECT().
					GetEnvironment(gomock.Any(), "prod").
					Return(v20231001preview.EnvironmentResource{
						Properties: &v20231001preview.EnvironmentProperties{
							Providers: &v20231001preview.Providers{
								Azure: &v20231001preview.ProvidersAzure{
									Scope: to.Ptr("/subscriptions/test-subId/resourceGroups/test-rg"),
								},
							},
						},
					}, nil).
					Times(1)

			},
		},
		{
			Name:          "rad deploy - env does not exist invalid",
			Input:         []string{"app.bicep", "-e", "prod"},
			ExpectedValid: false,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
			ConfigureMocks: func(mocks radcli.ValidateMocks) {
				mocks.ApplicationManagementClient.EXPECT().
					GetEnvironment(gomock.Any(), "prod").
					Return(v20231001preview.EnvironmentResource{}, radcli.Create404Error()).
					Times(1)

			},
		},
		{
			Name:          "rad deploy - valid with env ID",
			Input:         []string{"app.bicep", "-e", "/planes/radius/local/resourceGroups/test-resource-group/providers/applications.core/environments/prod"},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
			ConfigureMocks: func(mocks radcli.ValidateMocks) {
				mocks.ApplicationManagementClient.EXPECT().
					GetEnvironment(gomock.Any(), "/planes/radius/local/resourceGroups/test-resource-group/providers/applications.core/environments/prod").
					Return(v20231001preview.EnvironmentResource{
						ID: to.Ptr("/planes/radius/local/resourceGroups/test-resource-group/providers/applications.core/environments/prod"),
					}, nil).
					Times(1)
			},
			ValidateCallback: func(t *testing.T, obj framework.Runner) {
				runner := obj.(*Runner)
				scope := "/planes/radius/local/resourceGroups/test-resource-group"
				environmentID := scope + "/providers/applications.core/environments/prod"
				require.Equal(t, scope, runner.Workspace.Scope)
				require.Equal(t, environmentID, runner.Workspace.Environment)
			},
		},
		{
			Name:          "rad deploy - valid with app and env",
			Input:         []string{"app.bicep", "-e", "prod", "-a", "my-app"},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
			ConfigureMocks: func(mocks radcli.ValidateMocks) {
				mocks.ApplicationManagementClient.EXPECT().
					GetEnvironment(gomock.Any(), "prod").
					Return(v20231001preview.EnvironmentResource{
						ID: to.Ptr("/planes/radius/local/resourceGroups/test-resource-group/providers/applications.core/environments/prod"),
					}, nil).
					Times(1)
			},
			ValidateCallback: func(t *testing.T, obj framework.Runner) {
				runner := obj.(*Runner)
				scope := "/planes/radius/local/resourceGroups/test-resource-group"
				environmentID := scope + "/providers/applications.core/environments/prod"
				applicationID := scope + "/providers/applications.core/applications/my-app"
				require.Equal(t, scope, runner.Workspace.Scope)
				require.Equal(t, environmentID, runner.Workspace.Environment)
				require.Equal(t, clients.RadiusProvider{ApplicationID: applicationID, EnvironmentID: environmentID}, *runner.Providers.Radius)
			},
		},
		{
			Name:          "rad deploy - app set by directory config",
			Input:         []string{"app.bicep", "-e", "prod"},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
				DirectoryConfig: &config.DirectoryConfig{
					Workspace: config.DirectoryWorkspaceConfig{
						Application: "my-app",
					},
				},
			},
			ConfigureMocks: func(mocks radcli.ValidateMocks) {
				mocks.ApplicationManagementClient.EXPECT().
					GetEnvironment(gomock.Any(), "prod").
					Return(v20231001preview.EnvironmentResource{}, nil).
					Times(1)
			},
		},
		{
			Name:          "rad deploy - fallback workspace",
			Input:         []string{"app.bicep", "--group", "my-group", "--environment", "prod"},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         radcli.LoadEmptyConfig(t),
			},
			ConfigureMocks: func(mocks radcli.ValidateMocks) {
				mocks.ApplicationManagementClient.EXPECT().
					GetEnvironment(gomock.Any(), "prod").
					Return(v20231001preview.EnvironmentResource{}, nil).
					Times(1)
			},
		},
		{
			Name:          "rad deploy - fallback workspace requires resource group",
			Input:         []string{"app.bicep", "--environment", "prod"},
			ExpectedValid: false,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         radcli.LoadEmptyConfig(t),
			},
		},
		{
			Name:          "rad deploy - too many args",
			Input:         []string{"app.bicep", "anotherfile.json"},
			ExpectedValid: false,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         radcli.LoadEmptyConfig(t),
			},
		},
		{
			Name:          "rad deploy - missing env and app succeeds",
			Input:         []string{"app.bicep", "--group", "new-group"},
			ExpectedValid: true,
			ConfigHolder: framework.ConfigHolder{
				ConfigFilePath: "",
				Config:         configWithWorkspace,
			},
			ConfigureMocks: func(mocks radcli.ValidateMocks) {
				mocks.ApplicationManagementClient.EXPECT().
					GetEnvironment(gomock.Any(), gomock.Any()).
					Return(v20231001preview.EnvironmentResource{}, radcli.Create404Error()).
					Times(1)
			},
		},
	}

	radcli.SharedValidateValidation(t, NewCommand, testcases)
}

func Test_Run(t *testing.T) {
	t.Run("Environment-scoped deployment with az provider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bicep := bicep.NewMockInterface(ctrl)
		bicep.EXPECT().
			PrepareTemplate("app.bicep").
			Return(map[string]any{}, nil).
			Times(1)

		workspace := &workspaces.Workspace{
			Connection: map[string]any{
				"kind":    "kubernetes",
				"context": "kind-kind",
			},

			Name: "kind-kind",
		}
		provider :=
			&clients.Providers{
				Azure: &clients.AzureProvider{
					Scope: "test-scope",
				},
				Radius: &clients.RadiusProvider{
					EnvironmentID: fmt.Sprintf("/planes/radius/local/resourceGroups/%s/providers/applications.core/environments/%s", radcli.TestEnvironmentName, radcli.TestEnvironmentName),
				},
			}

		filePath := "app.bicep"
		progressText := fmt.Sprintf(
			"Deploying template '%v' into environment '%v' from workspace '%v'...\n\n"+
				"Deployment In Progress...", filePath, radcli.TestEnvironmentID, workspace.Name)

		options := deploy.Options{
			Workspace:      *workspace,
			Parameters:     map[string]map[string]any{},
			CompletionText: "Deployment Complete",
			ProgressText:   progressText,
			Template:       map[string]any{},
			Providers:      provider,
		}

		deployMock := deploy.NewMockInterface(ctrl)
		deployMock.EXPECT().
			DeployWithProgress(gomock.Any(), options).
			DoAndReturn(func(ctx context.Context, o deploy.Options) (clients.DeploymentResult, error) {
				// Capture options for verification
				options = o
				return clients.DeploymentResult{}, nil
			}).
			Times(1)

		outputSink := &output.MockOutput{}
		runner := &Runner{
			Bicep:               bicep,
			Deploy:              deployMock,
			Output:              outputSink,
			FilePath:            filePath,
			EnvironmentNameOrID: radcli.TestEnvironmentID,
			Parameters:          map[string]map[string]any{},
			Workspace:           workspace,
			Providers:           provider,
		}

		err := runner.Run(context.Background())
		require.NoError(t, err)

		// Deployment is scoped to env
		require.Equal(t, "", options.Providers.Radius.ApplicationID)
		require.Equal(t, runner.Providers.Radius.EnvironmentID, options.Providers.Radius.EnvironmentID)

		// All of the output in this command is being done by functions that we mock for testing, so this
		// is always empty.
		require.Empty(t, outputSink.Writes)
	})

	t.Run("Environment-scoped deployment with aws provider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bicep := bicep.NewMockInterface(ctrl)
		bicep.EXPECT().
			PrepareTemplate("app.bicep").
			Return(map[string]any{}, nil).
			Times(1)

		workspace := &workspaces.Workspace{
			Connection: map[string]any{
				"kind":    "kubernetes",
				"context": "kind-kind",
			},
			Name: "kind-kind",
		}
		ProviderConfig := clients.Providers{
			AWS: &clients.AWSProvider{
				Scope: "test-scope",
			},
			Radius: &clients.RadiusProvider{
				EnvironmentID: fmt.Sprintf("/planes/radius/local/resourceGroups/%s/providers/applications.core/environments/%s", radcli.TestEnvironmentName, radcli.TestEnvironmentName),
			},
		}

		filePath := "app.bicep"
		progressText := fmt.Sprintf(
			"Deploying template '%v' into environment '%v' from workspace '%v'...\n\n"+
				"Deployment In Progress...", filePath, radcli.TestEnvironmentID, workspace.Name)

		options := deploy.Options{
			Workspace:      *workspace,
			Parameters:     map[string]map[string]any{},
			CompletionText: "Deployment Complete",
			ProgressText:   progressText,
			Template:       map[string]any{},
			Providers:      &ProviderConfig,
		}

		deployMock := deploy.NewMockInterface(ctrl)
		deployMock.EXPECT().
			DeployWithProgress(gomock.Any(), options).
			DoAndReturn(func(ctx context.Context, o deploy.Options) (clients.DeploymentResult, error) {
				// Capture options for verification
				options = o
				return clients.DeploymentResult{}, nil
			}).
			Times(1)

		outputSink := &output.MockOutput{}
		runner := &Runner{
			Bicep:               bicep,
			Deploy:              deployMock,
			Output:              outputSink,
			Providers:           &ProviderConfig,
			FilePath:            filePath,
			EnvironmentNameOrID: radcli.TestEnvironmentID,
			Parameters:          map[string]map[string]any{},
			Workspace:           workspace,
		}

		err := runner.Run(context.Background())
		require.NoError(t, err)

		// Deployment is scoped to env
		require.Equal(t, "", options.Providers.Radius.ApplicationID)
		require.Equal(t, runner.Providers.Radius.EnvironmentID, options.Providers.Radius.EnvironmentID)

		// All of the output in this command is being done by functions that we mock for testing, so this
		// is always empty.
		require.Empty(t, outputSink.Writes)
	})

	t.Run("Application-scoped deployment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bicep := bicep.NewMockInterface(ctrl)
		bicep.EXPECT().
			PrepareTemplate("app.bicep").
			Return(map[string]any{}, nil).
			Times(1)

		options := deploy.Options{}

		appManagmentMock := clients.NewMockApplicationsManagementClient(ctrl)
		appManagmentMock.EXPECT().
			GetEnvironment(gomock.Any(), radcli.TestEnvironmentName).
			Return(v20231001preview.EnvironmentResource{}, nil).
			Times(1)
		appManagmentMock.EXPECT().
			CreateApplicationIfNotFound(gomock.Any(), "test-application", gomock.Any()).
			Return(nil).
			Times(1)

		deployMock := deploy.NewMockInterface(ctrl)
		deployMock.EXPECT().
			DeployWithProgress(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, o deploy.Options) (clients.DeploymentResult, error) {
				// Capture options for verification
				options = o
				return clients.DeploymentResult{}, nil
			}).
			Times(1)

		workspace := &workspaces.Workspace{
			Connection: map[string]any{
				"kind":    "kubernetes",
				"context": "kind-kind",
			},
			Name: "kind-kind",
		}
		outputSink := &output.MockOutput{}
		providers := clients.Providers{
			Radius: &clients.RadiusProvider{
				EnvironmentID: fmt.Sprintf("/planes/radius/local/resourceGroups/%s/providers/applications.core/environments/%s", radcli.TestEnvironmentName, radcli.TestEnvironmentName),
				ApplicationID: fmt.Sprintf("/planes/radius/local/resourceGroups/%s/providers/applications.core/environments/%s/applications/test-application", radcli.TestEnvironmentName, radcli.TestEnvironmentName),
			},
		}

		runner := &Runner{
			Bicep:               bicep,
			ConnectionFactory:   &connections.MockFactory{ApplicationsManagementClient: appManagmentMock},
			Deploy:              deployMock,
			Output:              outputSink,
			Providers:           &providers,
			FilePath:            "app.bicep",
			ApplicationName:     "test-application",
			EnvironmentNameOrID: radcli.TestEnvironmentName,
			Parameters:          map[string]map[string]any{},
			Workspace:           workspace,
		}

		err := runner.Run(context.Background())
		require.NoError(t, err)

		// Deployment is scoped to app and env
		require.Equal(t, runner.Providers.Radius.ApplicationID, options.Providers.Radius.ApplicationID)
		require.Equal(t, runner.Providers.Radius.EnvironmentID, options.Providers.Radius.EnvironmentID)

		// All of the output in this command is being done by functions that we mock for testing, so this
		// is always empty.
		require.Empty(t, outputSink.Writes)
	})

	t.Run("Deployment that doesn't need an app or env", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bicep := bicep.NewMockInterface(ctrl)
		bicep.EXPECT().
			PrepareTemplate("app.bicep").
			Return(map[string]any{}, nil).
			Times(1)

		appManagmentMock := clients.NewMockApplicationsManagementClient(ctrl)

		// GetEnvironment returns a 404 error
		appManagmentMock.EXPECT().
			GetEnvironment(gomock.Any(), "envdoesntexist").
			Return(v20231001preview.EnvironmentResource{}, radcli.Create404Error()).
			Times(1)

		deployMock := deploy.NewMockInterface(ctrl)
		deployMock.EXPECT().
			DeployWithProgress(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, o deploy.Options) (clients.DeploymentResult, error) {
				return clients.DeploymentResult{}, nil
			}).
			Times(1)

		workspace := &workspaces.Workspace{
			Connection: map[string]any{
				"kind":    "kubernetes",
				"context": "kind-kind",
			},
			Name: "kind-kind",
		}
		outputSink := &output.MockOutput{}

		providers := clients.Providers{
			Radius: &clients.RadiusProvider{
				EnvironmentID: fmt.Sprintf("/planes/radius/local/resourceGroups/%s/providers/applications.core/environments/%s", radcli.TestEnvironmentName, radcli.TestEnvironmentName),
				ApplicationID: fmt.Sprintf("/planes/radius/local/resourceGroups/%s/providers/applications.core/environments/%s/applications/test-application", radcli.TestEnvironmentName, radcli.TestEnvironmentName),
			},
		}

		runner := &Runner{
			Bicep:               bicep,
			ConnectionFactory:   &connections.MockFactory{ApplicationsManagementClient: appManagmentMock},
			Deploy:              deployMock,
			Output:              outputSink,
			Providers:           &providers,
			FilePath:            "app.bicep",
			ApplicationName:     "appdoesntexist",
			EnvironmentNameOrID: "envdoesntexist",
			Parameters:          map[string]map[string]any{},
			Workspace:           workspace,
		}

		err := runner.Run(context.Background())

		// Even though GetEnvironment returns a 404 error, the deployment should still succeed
		require.NoError(t, err)

		// All of the output in this command is being done by functions that we mock for testing, so this
		// is always empty.
		require.Empty(t, outputSink.Writes)
	})

	t.Run("Deployment with missing parameters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bicep := bicep.NewMockInterface(ctrl)
		bicep.EXPECT().
			PrepareTemplate("app.bicep").
			Return(map[string]any{
				"parameters": map[string]any{
					"application": map[string]any{},
					"environment": map[string]any{},
					"location":    map[string]any{},
					"size":        map[string]any{"defaultValue": "BIG!"},
				},
			}, nil).
			Times(1)

		workspace := &workspaces.Workspace{
			Connection: map[string]any{
				"kind":    "kubernetes",
				"context": "kind-kind",
			},
			Name: "kind-kind",
		}
		outputSink := &output.MockOutput{}

		providers := clients.Providers{
			Radius: &clients.RadiusProvider{
				EnvironmentID: fmt.Sprintf("/planes/radius/local/resourceGroups/%s/providers/applications.core/environments/%s", radcli.TestEnvironmentName, radcli.TestEnvironmentName),
			},
		}

		runner := &Runner{
			Bicep:               bicep,
			ConnectionFactory:   &connections.MockFactory{},
			Output:              outputSink,
			Providers:           &providers,
			EnvironmentNameOrID: radcli.TestEnvironmentName,
			FilePath:            "app.bicep",
			Parameters:          map[string]map[string]any{},
			Workspace:           workspace,
		}

		err := runner.Run(context.Background())

		// Even though GetEnvironment returns a 404 error, the deployment should still succeed
		require.Error(t, err)

		expected := `The template "app.bicep" could not be deployed because of the following errors:

  - The template requires an application. Use --application to specify the application name.
  - The template requires a parameter "location". Use --parameters location=<value> to specify the value.`
		require.Equal(t, expected, err.Error())

		// All of the output in this command is being done by functions that we mock for testing, so this
		// is always empty.
		require.Empty(t, outputSink.Writes)
	})
}

func Test_injectAutomaticParameters(t *testing.T) {
	template := map[string]any{
		"parameters": map[string]any{
			"environment": map[string]any{},
			"application": map[string]any{},
		},
	}

	runner := Runner{
		Parameters: map[string]map[string]any{
			"a": {
				"value": "YO",
			},
		},
		Providers: &clients.Providers{
			Radius: &clients.RadiusProvider{
				ApplicationID: "test-app",
				EnvironmentID: "test-env",
			},
		},
	}
	err := runner.injectAutomaticParameters(template)
	require.NoError(t, err)

	expected := map[string]map[string]any{
		"environment": {
			"value": "test-env",
		},
		"application": {
			"value": "test-app",
		},
		"a": {
			"value": "YO",
		},
	}

	require.Equal(t, expected, runner.Parameters)
}

func Test_reportMissingParameters(t *testing.T) {
	template := map[string]any{
		"parameters": map[string]any{
			"a":                         map[string]any{},
			"b":                         map[string]any{},
			"parameterWithDefaultValue": map[string]any{"defaultValue": "!"},
		},
	}

	t.Run("Missing parameters", func(t *testing.T) {
		runner := Runner{
			FilePath: "app.bicep",
			Parameters: map[string]map[string]any{
				"b": {
					"value": "YO",
				},
			},
		}
		err := runner.reportMissingParameters(template)

		expected := `The template "app.bicep" could not be deployed because of the following errors:

  - The template requires a parameter "a". Use --parameters a=<value> to specify the value.`
		require.Equal(t, expected, err.Error())
	})

	t.Run("All parameters provided", func(t *testing.T) {
		runner := Runner{
			FilePath: "app.bicep",
			Parameters: map[string]map[string]any{
				"a": {
					"value": "YO",
				},
				"b": {
					"value": "YO",
				},
			},
		}
		err := runner.reportMissingParameters(template)
		require.NoError(t, err)
	})

	t.Run("All parameters provided (ignoring case)", func(t *testing.T) {
		runner := Runner{
			FilePath: "app.bicep",
			Parameters: map[string]map[string]any{
				"A": {
					"value": "YO",
				},
				"B": {
					"value": "YO",
				},
				"parameterWithDEfaultValue": {
					"value": "YO",
				},
			},
		}
		err := runner.reportMissingParameters(template)
		require.NoError(t, err)
	})
}
