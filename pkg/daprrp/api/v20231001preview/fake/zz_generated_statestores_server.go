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
	"github.com/radius-project/radius/pkg/daprrp/api/v20231001preview"
	"net/http"
	"net/url"
	"regexp"
)

// StateStoresServer is a fake server for instances of the v20231001preview.StateStoresClient type.
type StateStoresServer struct{
	// BeginCreateOrUpdate is the fake for method StateStoresClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, stateStoreName string, resource v20231001preview.DaprStateStoreResource, options *v20231001preview.StateStoresClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[v20231001preview.StateStoresClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method StateStoresClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, stateStoreName string, options *v20231001preview.StateStoresClientBeginDeleteOptions) (resp azfake.PollerResponder[v20231001preview.StateStoresClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method StateStoresClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, stateStoreName string, options *v20231001preview.StateStoresClientGetOptions) (resp azfake.Responder[v20231001preview.StateStoresClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByScopePager is the fake for method StateStoresClient.NewListByScopePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByScopePager func(options *v20231001preview.StateStoresClientListByScopeOptions) (resp azfake.PagerResponder[v20231001preview.StateStoresClientListByScopeResponse])

	// BeginUpdate is the fake for method StateStoresClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, stateStoreName string, properties v20231001preview.DaprStateStoreResourceUpdate, options *v20231001preview.StateStoresClientBeginUpdateOptions) (resp azfake.PollerResponder[v20231001preview.StateStoresClientUpdateResponse], errResp azfake.ErrorResponder)

}

// NewStateStoresServerTransport creates a new instance of StateStoresServerTransport with the provided implementation.
// The returned StateStoresServerTransport instance is connected to an instance of v20231001preview.StateStoresClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewStateStoresServerTransport(srv *StateStoresServer) *StateStoresServerTransport {
	return &StateStoresServerTransport{
		srv: srv,
		beginCreateOrUpdate: newTracker[azfake.PollerResponder[v20231001preview.StateStoresClientCreateOrUpdateResponse]](),
		beginDelete: newTracker[azfake.PollerResponder[v20231001preview.StateStoresClientDeleteResponse]](),
		newListByScopePager: newTracker[azfake.PagerResponder[v20231001preview.StateStoresClientListByScopeResponse]](),
		beginUpdate: newTracker[azfake.PollerResponder[v20231001preview.StateStoresClientUpdateResponse]](),
	}
}

// StateStoresServerTransport connects instances of v20231001preview.StateStoresClient to instances of StateStoresServer.
// Don't use this type directly, use NewStateStoresServerTransport instead.
type StateStoresServerTransport struct {
	srv *StateStoresServer
	beginCreateOrUpdate *tracker[azfake.PollerResponder[v20231001preview.StateStoresClientCreateOrUpdateResponse]]
	beginDelete *tracker[azfake.PollerResponder[v20231001preview.StateStoresClientDeleteResponse]]
	newListByScopePager *tracker[azfake.PagerResponder[v20231001preview.StateStoresClientListByScopeResponse]]
	beginUpdate *tracker[azfake.PollerResponder[v20231001preview.StateStoresClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for StateStoresServerTransport.
func (s *StateStoresServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return s.dispatchToMethodFake(req, method)
}

func (s *StateStoresServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		 if stateStoresServerTransportInterceptor != nil {
			 res.resp, res.err, intercepted = stateStoresServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "StateStoresClient.BeginCreateOrUpdate":
				res.resp, res.err = s.dispatchBeginCreateOrUpdate(req)
			case "StateStoresClient.BeginDelete":
				res.resp, res.err = s.dispatchBeginDelete(req)
			case "StateStoresClient.Get":
				res.resp, res.err = s.dispatchGet(req)
			case "StateStoresClient.NewListByScopePager":
				res.resp, res.err = s.dispatchNewListByScopePager(req)
			case "StateStoresClient.BeginUpdate":
				res.resp, res.err = s.dispatchBeginUpdate(req)
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

func (s *StateStoresServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := s.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Dapr/stateStores/(?P<stateStoreName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[v20231001preview.DaprStateStoreResource](req)
	if err != nil {
		return nil, err
	}
	stateStoreNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("stateStoreName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.BeginCreateOrUpdate(req.Context(), stateStoreNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
		beginCreateOrUpdate = &respr
		s.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		s.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		s.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (s *StateStoresServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if s.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := s.beginDelete.get(req)
	if beginDelete == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Dapr/stateStores/(?P<stateStoreName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	stateStoreNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("stateStoreName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.BeginDelete(req.Context(), stateStoreNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
		beginDelete = &respr
		s.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		s.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		s.beginDelete.remove(req)
	}

	return resp, nil
}

func (s *StateStoresServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Dapr/stateStores/(?P<stateStoreName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	stateStoreNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("stateStoreName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), stateStoreNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DaprStateStoreResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *StateStoresServerTransport) dispatchNewListByScopePager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByScopePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByScopePager not implemented")}
	}
	newListByScopePager := s.newListByScopePager.get(req)
	if newListByScopePager == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Dapr/stateStores`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
resp := s.srv.NewListByScopePager(nil)
		newListByScopePager = &resp
		s.newListByScopePager.add(req, newListByScopePager)
		server.PagerResponderInjectNextLinks(newListByScopePager, req, func(page *v20231001preview.StateStoresClientListByScopeResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByScopePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByScopePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByScopePager) {
		s.newListByScopePager.remove(req)
	}
	return resp, nil
}

func (s *StateStoresServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := s.beginUpdate.get(req)
	if beginUpdate == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Dapr/stateStores/(?P<stateStoreName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[v20231001preview.DaprStateStoreResourceUpdate](req)
	if err != nil {
		return nil, err
	}
	stateStoreNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("stateStoreName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.BeginUpdate(req.Context(), stateStoreNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
		beginUpdate = &respr
		s.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		s.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		s.beginUpdate.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to StateStoresServerTransport
var stateStoresServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}