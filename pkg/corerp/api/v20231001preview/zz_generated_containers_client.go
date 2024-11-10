// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20231001preview

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ContainersClient contains the methods for the Containers group.
// Don't use this type directly, use NewContainersClient() instead.
type ContainersClient struct {
	internal *arm.Client
	rootScope string
}

// NewContainersClient creates a new instance of ContainersClient with the specified values.
//   - rootScope - The scope in which the resource is present. UCP Scope is /planes/{planeType}/{planeName}/resourceGroup/{resourcegroupID}
//     and Azure resource scope is
//     /subscriptions/{subscriptionID}/resourceGroup/{resourcegroupID}
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewContainersClient(rootScope string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ContainersClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ContainersClient{
		rootScope: rootScope,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a ContainerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - containerName - Container name
//   - resource - Resource create parameters.
//   - options - ContainersClientBeginCreateOrUpdateOptions contains the optional parameters for the ContainersClient.BeginCreateOrUpdate
//     method.
func (client *ContainersClient) BeginCreateOrUpdate(ctx context.Context, containerName string, resource ContainerResource, options *ContainersClientBeginCreateOrUpdateOptions) (*runtime.Poller[ContainersClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, containerName, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ContainersClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ContainersClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create a ContainerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *ContainersClient) createOrUpdate(ctx context.Context, containerName string, resource ContainerResource, options *ContainersClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ContainersClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, containerName, resource, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ContainersClient) createOrUpdateCreateRequest(ctx context.Context, containerName string, resource ContainerResource, _ *ContainersClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers/{containerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if containerName == "" {
		return nil, errors.New("parameter containerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerName}", url.PathEscape(containerName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
	return nil, err
}
;	return req, nil
}

// BeginDelete - Delete a ContainerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - containerName - Container name
//   - options - ContainersClientBeginDeleteOptions contains the optional parameters for the ContainersClient.BeginDelete method.
func (client *ContainersClient) BeginDelete(ctx context.Context, containerName string, options *ContainersClientBeginDeleteOptions) (*runtime.Poller[ContainersClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, containerName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ContainersClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ContainersClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete a ContainerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *ContainersClient) deleteOperation(ctx context.Context, containerName string, options *ContainersClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ContainersClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, containerName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ContainersClient) deleteCreateRequest(ctx context.Context, containerName string, _ *ContainersClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers/{containerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if containerName == "" {
		return nil, errors.New("parameter containerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerName}", url.PathEscape(containerName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a ContainerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - containerName - Container name
//   - options - ContainersClientGetOptions contains the optional parameters for the ContainersClient.Get method.
func (client *ContainersClient) Get(ctx context.Context, containerName string, options *ContainersClientGetOptions) (ContainersClientGetResponse, error) {
	var err error
	const operationName = "ContainersClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, containerName, options)
	if err != nil {
		return ContainersClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ContainersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ContainersClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ContainersClient) getCreateRequest(ctx context.Context, containerName string, _ *ContainersClientGetOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers/{containerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if containerName == "" {
		return nil, errors.New("parameter containerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerName}", url.PathEscape(containerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ContainersClient) getHandleResponse(resp *http.Response) (ContainersClientGetResponse, error) {
	result := ContainersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContainerResource); err != nil {
		return ContainersClientGetResponse{}, err
	}
	return result, nil
}

// NewListByScopePager - List ContainerResource resources by Scope
//
// Generated from API version 2023-10-01-preview
//   - options - ContainersClientListByScopeOptions contains the optional parameters for the ContainersClient.NewListByScopePager
//     method.
func (client *ContainersClient) NewListByScopePager(options *ContainersClientListByScopeOptions) (*runtime.Pager[ContainersClientListByScopeResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ContainersClientListByScopeResponse]{
		More: func(page ContainersClientListByScopeResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ContainersClientListByScopeResponse) (ContainersClientListByScopeResponse, error) {
		ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ContainersClient.NewListByScopePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByScopeCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return ContainersClientListByScopeResponse{}, err
			}
			return client.listByScopeHandleResponse(resp)
			},
		Tracer: client.internal.Tracer(),
	})
}

// listByScopeCreateRequest creates the ListByScope request.
func (client *ContainersClient) listByScopeCreateRequest(ctx context.Context, _ *ContainersClientListByScopeOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByScopeHandleResponse handles the ListByScope response.
func (client *ContainersClient) listByScopeHandleResponse(resp *http.Response) (ContainersClientListByScopeResponse, error) {
	result := ContainersClientListByScopeResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContainerResourceListResult); err != nil {
		return ContainersClientListByScopeResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update a ContainerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - containerName - Container name
//   - properties - The resource properties to be updated.
//   - options - ContainersClientBeginUpdateOptions contains the optional parameters for the ContainersClient.BeginUpdate method.
func (client *ContainersClient) BeginUpdate(ctx context.Context, containerName string, properties ContainerResourceUpdate, options *ContainersClientBeginUpdateOptions) (*runtime.Poller[ContainersClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, containerName, properties, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ContainersClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ContainersClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - Update a ContainerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *ContainersClient) update(ctx context.Context, containerName string, properties ContainerResourceUpdate, options *ContainersClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ContainersClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, containerName, properties, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *ContainersClient) updateCreateRequest(ctx context.Context, containerName string, properties ContainerResourceUpdate, _ *ContainersClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/containers/{containerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if containerName == "" {
		return nil, errors.New("parameter containerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerName}", url.PathEscape(containerName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, properties); err != nil {
	return nil, err
}
;	return req, nil
}

