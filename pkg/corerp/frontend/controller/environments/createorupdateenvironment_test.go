// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package environments

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	v20220315privatepreview "github.com/project-radius/radius/pkg/corerp/api/v20220315privatepreview"
	radiustesting "github.com/project-radius/radius/pkg/corerp/testing"
	"github.com/project-radius/radius/pkg/store"
	"github.com/stretchr/testify/require"
)

func TestCreateOrUpdateEnvironmentRun_20220315PrivatePreview(t *testing.T) {
	mctrl := gomock.NewController(t)
	defer mctrl.Finish()

	mStorageClient := store.NewMockStorageClient(mctrl)
	ctx := context.Background()

	createNewResourceCases := []struct {
		desc               string
		headerKey          string
		headerValue        string
		resourceETag       string
		expectedStatusCode int
		shouldFail         bool
	}{
		{"create-new-resource-no-if-match", "If-Match", "", "", 200, false},
		{"create-new-resource-*-if-match", "If-Match", "*", "", 412, true},
		{"create-new-resource-etag-if-match", "If-Match", "randome-etag", "", 412, true},
		{"create-new-resource-*-if-none-match", "If-None-Match", "*", "", 200, false},
	}

	for _, tt := range createNewResourceCases {
		t.Run(tt.desc, func(t *testing.T) {
			envInput, envDataModel, expectedOutput := getTestModels20220315privatepreview()
			w := httptest.NewRecorder()
			req, _ := radiustesting.GetARMTestHTTPRequest(ctx, http.MethodGet, testHeaderfile, envInput)
			req.Header.Set(tt.headerKey, tt.headerValue)
			ctx := radiustesting.ARMTestContextFromRequest(req)

			mStorageClient.
				EXPECT().
				Get(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, id string, _ ...store.GetOptions) (*store.Object, error) {
					return nil, &store.ErrNotFound{}
				})

			expectedOutput.SystemData.CreatedAt = expectedOutput.SystemData.LastModifiedAt
			expectedOutput.SystemData.CreatedBy = expectedOutput.SystemData.LastModifiedBy
			expectedOutput.SystemData.CreatedByType = expectedOutput.SystemData.LastModifiedByType

			if !tt.shouldFail {
				mStorageClient.
					EXPECT().
					Save(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, obj *store.Object, opts ...store.SaveOptions) (*store.Object, error) {
						return &store.Object{
							Metadata: store.Metadata{ID: obj.ID, ETag: "new-resource-etag"},
							Data:     envDataModel,
						}, nil
					})
			}

			ctl, err := NewCreateOrUpdateEnvironment(mStorageClient, nil)
			require.NoError(t, err)
			resp, err := ctl.Run(ctx, req)
			require.NoError(t, err)
			_ = resp.Apply(ctx, w, req)
			require.Equal(t, tt.expectedStatusCode, w.Result().StatusCode)

			if !tt.shouldFail {
				actualOutput := &v20220315privatepreview.EnvironmentResource{}
				_ = json.Unmarshal(w.Body.Bytes(), actualOutput)
				require.Equal(t, expectedOutput, actualOutput)

				require.Equal(t, "new-resource-etag", w.Header().Get("ETag"))
			}
		})
	}

	updateExistingResourceCases := []struct {
		desc               string
		headerKey          string
		headerValue        string
		resourceETag       string
		expectedStatusCode int
		shouldFail         bool
	}{
		{"update-resource-no-if-match", "If-Match", "", "resource-etag", 200, false},
		{"update-resource-*-if-match", "If-Match", "*", "resource-etag", 200, false},
		{"update-resource-matching-if-match", "If-Match", "matching-etag", "matching-etag", 200, false},
		{"update-resource-not-matching-if-match", "If-Match", "not-matching-etag", "another-etag", 412, true},
		{"update-resource-*-if-none-match", "If-None-Match", "*", "another-etag", 412, true},
	}

	for _, tt := range updateExistingResourceCases {
		t.Run(tt.desc, func(t *testing.T) {
			envInput, envDataModel, expectedOutput := getTestModels20220315privatepreview()
			w := httptest.NewRecorder()
			req, _ := radiustesting.GetARMTestHTTPRequest(ctx, http.MethodGet, testHeaderfile, envInput)
			req.Header.Set(tt.headerKey, tt.headerValue)
			ctx := radiustesting.ARMTestContextFromRequest(req)

			mStorageClient.
				EXPECT().
				Get(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, id string, _ ...store.GetOptions) (*store.Object, error) {
					return &store.Object{
						Metadata: store.Metadata{ID: id, ETag: tt.resourceETag},
						Data:     envDataModel,
					}, nil
				})

			if !tt.shouldFail {
				mStorageClient.
					EXPECT().
					Save(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, obj *store.Object, opts ...store.SaveOptions) (*store.Object, error) {
						return &store.Object{
							Metadata: store.Metadata{ID: obj.ID, ETag: "updated-resource-etag"},
							Data:     envDataModel,
						}, nil
					})
			}

			ctl, err := NewCreateOrUpdateEnvironment(mStorageClient, nil)
			require.NoError(t, err)
			resp, err := ctl.Run(ctx, req)
			_ = resp.Apply(ctx, w, req)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatusCode, w.Result().StatusCode)

			if !tt.shouldFail {
				actualOutput := &v20220315privatepreview.EnvironmentResource{}
				_ = json.Unmarshal(w.Body.Bytes(), actualOutput)
				require.Equal(t, expectedOutput, actualOutput)

				require.Equal(t, "updated-resource-etag", w.Header().Get("ETag"))
			}
		})
	}
}