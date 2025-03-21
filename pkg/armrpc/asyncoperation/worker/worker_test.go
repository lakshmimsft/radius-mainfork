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

package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/components/database"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDefaultOptions(t *testing.T) {
	worker := New(Options{}, nil, nil, nil)
	require.Equal(t, defaultDeduplicationDuration, worker.options.DeduplicationDuration)
	require.Equal(t, defaultMaxOperationRetryCount, worker.options.MaxOperationRetryCount)
	require.Equal(t, defaultMessageExtendMargin, worker.options.MessageExtendMargin)
	require.Equal(t, defaultMinMessageLockDuration, worker.options.MinMessageLockDuration)
	require.Equal(t, defaultMaxOperationConcurrency, worker.options.MaxOperationConcurrency)
}

func TestUpdateResourceState(t *testing.T) {
	updateStates := []struct {
		tc          string
		in          map[string]any
		updateState v1.ProvisioningState
		outErr      error
		callSave    bool
	}{
		{
			tc: "not found provisioningState",
			in: map[string]any{
				"name":       "env0",
				"properties": map[string]any{},
			},
			updateState: v1.ProvisioningStateAccepted,
			outErr:      nil,
			callSave:    true,
		},
		{
			tc: "not update state",
			in: map[string]any{
				"name":              "env0",
				"provisioningState": "Accepted",
				"properties":        map[string]any{},
			},
			updateState: v1.ProvisioningStateAccepted,
			outErr:      nil,
			callSave:    false,
		},
		{
			tc: "update state",
			in: map[string]any{
				"name":              "env0",
				"provisioningState": "Updating",
				"properties":        map[string]any{},
			},
			updateState: v1.ProvisioningStateAccepted,
			outErr:      nil,
			callSave:    true,
		},
	}

	for _, tt := range updateStates {
		t.Run(tt.tc, func(t *testing.T) {
			mctrl := gomock.NewController(t)
			defer mctrl.Finish()

			databaseClient := database.NewMockClient(mctrl)
			ctx := context.Background()

			databaseClient.
				EXPECT().
				Get(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, id string, options ...database.GetOptions) (*database.Object, error) {
					return &database.Object{
						Data: tt.in,
					}, nil
				})

			if tt.callSave {
				databaseClient.
					EXPECT().
					Save(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, obj *database.Object, options ...database.SaveOptions) error {
						k := obj.Data.(map[string]any)
						require.Equal(t, k["provisioningState"].(string), string(tt.updateState))
						return nil
					})
			}

			err := updateResourceState(ctx, databaseClient, "fakeid", tt.updateState)
			require.ErrorIs(t, err, tt.outErr)
		})
	}

}

func TestGetMessageExtendDuration(t *testing.T) {
	tests := []struct {
		in  time.Time
		out time.Duration
	}{
		{
			in:  time.Now().Add(defaultMessageExtendMargin),
			out: defaultMinMessageLockDuration,
		}, {
			in:  time.Now().Add(-defaultMessageExtendMargin),
			out: defaultMinMessageLockDuration,
		}, {
			in:  time.Now().Add(time.Duration(180) * time.Second),
			out: time.Duration(180)*time.Second - defaultMessageExtendMargin,
		},
	}

	for _, tt := range tests {
		worker := New(Options{}, nil, nil, nil)
		d := worker.getMessageExtendDuration(tt.in)
		require.Equal(t, tt.out, d.Round(time.Second))
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		err            error
		expectedArmErr v1.ErrorDetails
	}{
		{
			err:            v1.NewClientErrInvalidRequest("client error"),
			expectedArmErr: v1.ErrorDetails{Code: v1.CodeInvalid, Message: "client error"},
		},
		{
			err:            errors.New("internal error"),
			expectedArmErr: v1.ErrorDetails{Code: v1.CodeInternal, Message: "internal error"},
		},
	}

	for _, tt := range tests {
		armErr := extractError(tt.err)
		require.Equal(t, tt.expectedArmErr, armErr)
	}
}
