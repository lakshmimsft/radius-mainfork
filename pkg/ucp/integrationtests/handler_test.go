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

package integrationtests

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/ucp"
	"github.com/radius-project/radius/pkg/ucp/testhost"
	"github.com/stretchr/testify/require"
)

func Test_Handler_MethodNotAllowed(t *testing.T) {
	ucp := testhost.Start(t, testhost.NoModules())

	response := ucp.MakeRequest(http.MethodDelete, "/planes?api-version=2023-10-01-preview", nil)
	require.Equal(t, "failed to parse route: undefined route path", response.Error.Error.Details[0].Message)
}

func Test_Handler_NotFound(t *testing.T) {
	ucp := testhost.Start(t, testhost.NoModules())

	response := ucp.MakeRequest(http.MethodGet, "/abc", nil)
	response.EqualsErrorCode(http.StatusNotFound, v1.CodeNotFound)
	require.Regexp(t, "The request 'GET /abc' is invalid.", response.Error.Error.Message)
}

func Test_Handler_NotFound_PathBase(t *testing.T) {
	ucp := testhost.Start(t, testhost.NoModules(), testhost.TestHostOptionFunc(func(options *ucp.Options) {
		options.Config.Server.PathBase = "/" + uuid.New().String()
	}))

	response := ucp.MakeRequest(http.MethodGet, "/abc", nil)
	response.EqualsErrorCode(http.StatusNotFound, v1.CodeNotFound)
	require.Regexp(t, "The request 'GET /.*/abc' is invalid.", response.Error.Error.Message)
}
