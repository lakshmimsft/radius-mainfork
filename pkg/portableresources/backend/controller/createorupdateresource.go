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

package controller

import (
	"context"
	"errors"
	"fmt"

	ctrl "github.com/radius-project/radius/pkg/armrpc/asyncoperation/controller"
	"github.com/radius-project/radius/pkg/components/database"
	"github.com/radius-project/radius/pkg/portableresources/datamodel"
	"github.com/radius-project/radius/pkg/portableresources/processors"
	"github.com/radius-project/radius/pkg/recipes"
	"github.com/radius-project/radius/pkg/recipes/configloader"
	"github.com/radius-project/radius/pkg/recipes/engine"
	"github.com/radius-project/radius/pkg/recipes/util"
	rpv1 "github.com/radius-project/radius/pkg/rp/v1"
	"github.com/radius-project/radius/pkg/ucp/ucplog"
)

// CreateOrUpdateResource is the async operation controller to create or update portable resources.
type CreateOrUpdateResource[P interface {
	*T
	rpv1.RadiusResourceModel
}, T any] struct {
	ctrl.BaseController
	processor           processors.ResourceProcessor[P, T]
	engine              engine.Engine
	configurationLoader configloader.ConfigurationLoader
}

// NewCreateOrUpdateResource creates a new controller for creating or updating a resource with the given processor, engine,
// client, configurationLoader and options. The processor function will be called to process updates to the resource.
func NewCreateOrUpdateResource[P interface {
	*T
	rpv1.RadiusResourceModel
}, T any](opts ctrl.Options, processor processors.ResourceProcessor[P, T], eng engine.Engine, configurationLoader configloader.ConfigurationLoader) (ctrl.Controller, error) {
	return &CreateOrUpdateResource[P, T]{
		ctrl.NewBaseAsyncController(opts),
		processor,
		eng,
		configurationLoader,
	}, nil
}

// Run retrieves an existing resource, executes a recipe if needed, loads runtime configuration,
// processes the resource, cleans up any obsolete output resources, and saves the updated resource.
func (c *CreateOrUpdateResource[P, T]) Run(ctx context.Context, req *ctrl.Request) (ctrl.Result, error) {
	logger := ucplog.FromContextOrDiscard(ctx)
	obj, err := c.DatabaseClient().Get(ctx, req.ResourceID)
	if errors.Is(&database.ErrNotFound{ID: req.ResourceID}, err) {
		return ctrl.Result{}, err
	} else if err != nil {
		return ctrl.Result{}, err
	}

	data := P(new(T))
	if err = obj.As(data); err != nil {
		return ctrl.Result{}, err
	}

	// Clone existing output resources so we can diff them later.
	previousOutputResources := c.copyOutputResources(data)

	// Load configuration
	metadata := recipes.ResourceMetadata{EnvironmentID: data.ResourceMetadata().EnvironmentID(), ApplicationID: data.ResourceMetadata().ApplicationID(), ResourceID: data.GetBaseResource().ID}
	config, err := c.configurationLoader.LoadConfiguration(ctx, metadata)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Now we're ready to process recipes (if needed).
	recipeDataModel, supportsRecipes := any(data).(datamodel.RecipeDataModel)
	recipeOutput, err := c.executeRecipeIfNeeded(ctx, data, previousOutputResources, config.Simulated)
	if err != nil {
		if recipeError, ok := err.(*recipes.RecipeError); ok {
			logger.Error(err, fmt.Sprintf("failed to execute recipe. Encountered error while processing %s ", recipeError.ErrorDetails.Target))
			// Set the deployment status to the recipe error code.
			recipeDataModel.GetRecipe().DeploymentStatus = util.RecipeDeploymentStatus(recipeError.DeploymentStatus)
			update := &database.Object{
				Metadata: database.Metadata{
					ID: req.ResourceID,
				},
				Data: recipeDataModel.(rpv1.RadiusResourceModel),
			}
			// Save portable resource with updated deployment status to track errors during deletion.
			err = c.DatabaseClient().Save(ctx, update, database.WithETag(obj.ETag))
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.NewFailedResult(recipeError.ErrorDetails), nil
		}
		return ctrl.Result{}, err
	}

	if config.Simulated {
		logger.Info("The recipe was executed in simulation mode. No resources were deployed.")
	} else {
		// Now we're ready to process the resource. This will handle the updates to any user-visible state.
		err = c.processor.Process(ctx, data, processors.Options{RecipeOutput: recipeOutput, RuntimeConfiguration: config.Runtime})
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if supportsRecipes {
		recipeData := recipeDataModel.GetRecipe()
		if recipeData != nil {
			recipeData.DeploymentStatus = util.Success
			recipeDataModel.SetRecipe(recipeData)
		}
	}

	update := &database.Object{
		Metadata: database.Metadata{
			ID: req.ResourceID,
		},
		Data: recipeDataModel.(rpv1.RadiusResourceModel),
	}
	err = c.DatabaseClient().Save(ctx, update, database.WithETag(obj.ETag))
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, err
}

func (c *CreateOrUpdateResource[P, T]) copyOutputResources(data P) []string {
	previousOutputResources := []string{}
	for _, outputResource := range data.OutputResources() {
		previousOutputResources = append(previousOutputResources, outputResource.ID.String())
	}
	return previousOutputResources
}

func (c *CreateOrUpdateResource[P, T]) executeRecipeIfNeeded(ctx context.Context, data P, prevState []string, simulated bool) (*recipes.RecipeOutput, error) {
	// 'any' is required here to convert to an interface type, only then can we use a type assertion.
	recipeDataModel, supportsRecipes := any(data).(datamodel.RecipeDataModel)
	if !supportsRecipes {
		return nil, nil
	}

	input := recipeDataModel.GetRecipe()
	if input == nil {
		return nil, nil
	}
	request := recipes.ResourceMetadata{
		Name:          input.Name,
		Parameters:    input.Parameters,
		EnvironmentID: data.ResourceMetadata().EnvironmentID(),
		ApplicationID: data.ResourceMetadata().ApplicationID(),
		ResourceID:    data.GetBaseResource().ID,
	}

	return c.engine.Execute(ctx, engine.ExecuteOptions{
		BaseOptions: engine.BaseOptions{
			Recipe: request,
		},
		PreviousState: prevState,
		Simulated:     simulated,
	})
}
