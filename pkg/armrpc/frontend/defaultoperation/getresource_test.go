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

package defaultoperation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	ctrl "github.com/radius-project/radius/pkg/armrpc/frontend/controller"
	"github.com/radius-project/radius/pkg/armrpc/rpctest"
	"github.com/radius-project/radius/pkg/components/database"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testDataModel struct {
	Name string `json:"name"`
}

func (e testDataModel) ResourceTypeName() string {
	return "Applications.Test/resource"
}

func (e testDataModel) GetSystemData() *v1.SystemData {
	return nil
}

func (e testDataModel) GetBaseResource() *v1.BaseResource {
	return nil
}

func (e testDataModel) ProvisioningState() v1.ProvisioningState {
	return v1.ProvisioningStateAccepted
}

func (e testDataModel) SetProvisioningState(state v1.ProvisioningState) {
}

func (e testDataModel) UpdateMetadata(ctx *v1.ARMRequestContext, oldResource *v1.BaseResource) {
}

type testVersionedModel struct {
	Name string `json:"name"`
}

func (v *testVersionedModel) ConvertFrom(src v1.DataModelInterface) error {
	dm := src.(*testDataModel)
	v.Name = dm.Name
	return nil
}

func (v *testVersionedModel) ConvertTo() (v1.DataModelInterface, error) {
	return nil, nil
}

func resourceToVersioned(model *testDataModel, version string) (v1.VersionedModelInterface, error) {
	switch version {
	case testAPIVersion:
		versioned := &testVersionedModel{}
		if err := versioned.ConvertFrom(model); err != nil {
			return nil, err
		}
		return versioned, nil

	default:
		return nil, v1.ErrUnsupportedAPIVersion
	}
}

func TestGetResourceRun(t *testing.T) {
	mctrl := gomock.NewController(t)
	defer mctrl.Finish()

	databaseClient := database.NewMockClient(mctrl)
	ctx := context.Background()

	testResourceDataModel := &testDataModel{
		Name: "ResourceName",
	}
	expectedOutput := &testVersionedModel{
		Name: "ResourceName",
	}

	t.Run("get non-existing resource", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := rpctest.NewHTTPRequestFromJSON(ctx, http.MethodGet, resourceTestHeaderFile, nil)
		require.NoError(t, err)
		ctx := rpctest.NewARMRequestContext(req)

		databaseClient.
			EXPECT().
			Get(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, id string, _ ...database.GetOptions) (*database.Object, error) {
				return nil, &database.ErrNotFound{ID: id}
			})

		opts := ctrl.Options{
			DatabaseClient: databaseClient,
		}

		ctrlOpts := ctrl.ResourceOptions[testDataModel]{
			ResponseConverter: resourceToVersioned,
		}

		ctl, err := NewGetResource(opts, ctrlOpts)

		require.NoError(t, err)
		resp, err := ctl.Run(ctx, w, req)
		require.NoError(t, err)
		_ = resp.Apply(ctx, w, req)
		require.Equal(t, 404, w.Result().StatusCode)
	})

	t.Run("get existing resource", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := rpctest.NewHTTPRequestFromJSON(ctx, http.MethodGet, resourceTestHeaderFile, nil)
		require.NoError(t, err)
		ctx := rpctest.NewARMRequestContext(req)

		databaseClient.
			EXPECT().
			Get(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, id string, _ ...database.GetOptions) (*database.Object, error) {
				return &database.Object{
					Metadata: database.Metadata{ID: id},
					Data:     testResourceDataModel,
				}, nil
			})

		opts := ctrl.Options{
			DatabaseClient: databaseClient,
		}

		ctrlOpts := ctrl.ResourceOptions[testDataModel]{
			ResponseConverter: resourceToVersioned,
		}

		ctl, err := NewGetResource(opts, ctrlOpts)

		require.NoError(t, err)
		resp, err := ctl.Run(ctx, w, req)
		require.NoError(t, err)
		_ = resp.Apply(ctx, w, req)
		require.Equal(t, 200, w.Result().StatusCode)

		actualOutput := &testVersionedModel{}
		_ = json.Unmarshal(w.Body.Bytes(), actualOutput)

		require.Equal(t, expectedOutput, actualOutput)
	})
}
