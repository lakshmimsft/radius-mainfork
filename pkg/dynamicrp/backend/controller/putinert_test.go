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
	"testing"

	ctrl "github.com/radius-project/radius/pkg/armrpc/asyncoperation/controller"
	"github.com/radius-project/radius/test/testcontext"
	"github.com/stretchr/testify/require"
)

func Test_InertPutController_Run(t *testing.T) {
	setup := func() *InertPutController {
		opts := ctrl.Options{}
		controller, err := NewInertPutController(opts)
		require.NoError(t, err)
		return controller.(*InertPutController)
	}

	controller := setup()

	request := &ctrl.Request{}
	result, err := controller.Run(testcontext.New(t), request)
	require.NoError(t, err)
	require.Equal(t, ctrl.Result{}, result)
}
