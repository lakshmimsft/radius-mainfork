// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"context"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/radius-project/radius/pkg/corerp/api/v20231001preview"
	"net/http"
	"net/url"
	"regexp"
)

// VolumesServer is a fake server for instances of the v20231001preview.VolumesClient type.
type VolumesServer struct{
	// BeginCreateOrUpdate is the fake for method VolumesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, volumeName string, resource v20231001preview.VolumeResource, options *v20231001preview.VolumesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[v20231001preview.VolumesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method VolumesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, volumeName string, options *v20231001preview.VolumesClientBeginDeleteOptions) (resp azfake.PollerResponder[v20231001preview.VolumesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method VolumesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, volumeName string, options *v20231001preview.VolumesClientGetOptions) (resp azfake.Responder[v20231001preview.VolumesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByScopePager is the fake for method VolumesClient.NewListByScopePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByScopePager func(options *v20231001preview.VolumesClientListByScopeOptions) (resp azfake.PagerResponder[v20231001preview.VolumesClientListByScopeResponse])

	// BeginUpdate is the fake for method VolumesClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, volumeName string, properties v20231001preview.VolumeResourceUpdate, options *v20231001preview.VolumesClientBeginUpdateOptions) (resp azfake.PollerResponder[v20231001preview.VolumesClientUpdateResponse], errResp azfake.ErrorResponder)

}

// NewVolumesServerTransport creates a new instance of VolumesServerTransport with the provided implementation.
// The returned VolumesServerTransport instance is connected to an instance of v20231001preview.VolumesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewVolumesServerTransport(srv *VolumesServer) *VolumesServerTransport {
	return &VolumesServerTransport{
		srv: srv,
		beginCreateOrUpdate: newTracker[azfake.PollerResponder[v20231001preview.VolumesClientCreateOrUpdateResponse]](),
		beginDelete: newTracker[azfake.PollerResponder[v20231001preview.VolumesClientDeleteResponse]](),
		newListByScopePager: newTracker[azfake.PagerResponder[v20231001preview.VolumesClientListByScopeResponse]](),
		beginUpdate: newTracker[azfake.PollerResponder[v20231001preview.VolumesClientUpdateResponse]](),
	}
}

// VolumesServerTransport connects instances of v20231001preview.VolumesClient to instances of VolumesServer.
// Don't use this type directly, use NewVolumesServerTransport instead.
type VolumesServerTransport struct {
	srv *VolumesServer
	beginCreateOrUpdate *tracker[azfake.PollerResponder[v20231001preview.VolumesClientCreateOrUpdateResponse]]
	beginDelete *tracker[azfake.PollerResponder[v20231001preview.VolumesClientDeleteResponse]]
	newListByScopePager *tracker[azfake.PagerResponder[v20231001preview.VolumesClientListByScopeResponse]]
	beginUpdate *tracker[azfake.PollerResponder[v20231001preview.VolumesClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for VolumesServerTransport.
func (v *VolumesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return v.dispatchToMethodFake(req, method)
}

func (v *VolumesServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		 if volumesServerTransportInterceptor != nil {
			 res.resp, res.err, intercepted = volumesServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "VolumesClient.BeginCreateOrUpdate":
				res.resp, res.err = v.dispatchBeginCreateOrUpdate(req)
			case "VolumesClient.BeginDelete":
				res.resp, res.err = v.dispatchBeginDelete(req)
			case "VolumesClient.Get":
				res.resp, res.err = v.dispatchGet(req)
			case "VolumesClient.NewListByScopePager":
				res.resp, res.err = v.dispatchNewListByScopePager(req)
			case "VolumesClient.BeginUpdate":
				res.resp, res.err = v.dispatchBeginUpdate(req)
				default:
		res.err = fmt.Errorf("unhandled API %s", method)
			}

		}
		select {
		case resultChan <- res:
		case <-req.Context().Done():
		}
	}()

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

func (v *VolumesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if v.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := v.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/volumes/(?P<volumeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[v20231001preview.VolumeResource](req)
	if err != nil {
		return nil, err
	}
	volumeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("volumeName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.BeginCreateOrUpdate(req.Context(), volumeNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
		beginCreateOrUpdate = &respr
		v.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		v.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		v.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (v *VolumesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if v.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := v.beginDelete.get(req)
	if beginDelete == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/volumes/(?P<volumeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	volumeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("volumeName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.BeginDelete(req.Context(), volumeNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
		beginDelete = &respr
		v.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		v.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		v.beginDelete.remove(req)
	}

	return resp, nil
}

func (v *VolumesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if v.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/volumes/(?P<volumeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	volumeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("volumeName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.Get(req.Context(), volumeNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).VolumeResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *VolumesServerTransport) dispatchNewListByScopePager(req *http.Request) (*http.Response, error) {
	if v.srv.NewListByScopePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByScopePager not implemented")}
	}
	newListByScopePager := v.newListByScopePager.get(req)
	if newListByScopePager == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/volumes`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
resp := v.srv.NewListByScopePager(nil)
		newListByScopePager = &resp
		v.newListByScopePager.add(req, newListByScopePager)
		server.PagerResponderInjectNextLinks(newListByScopePager, req, func(page *v20231001preview.VolumesClientListByScopeResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByScopePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		v.newListByScopePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByScopePager) {
		v.newListByScopePager.remove(req)
	}
	return resp, nil
}

func (v *VolumesServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if v.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := v.beginUpdate.get(req)
	if beginUpdate == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/volumes/(?P<volumeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[v20231001preview.VolumeResourceUpdate](req)
	if err != nil {
		return nil, err
	}
	volumeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("volumeName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.BeginUpdate(req.Context(), volumeNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
		beginUpdate = &respr
		v.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		v.beginUpdate.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to VolumesServerTransport
var volumesServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}