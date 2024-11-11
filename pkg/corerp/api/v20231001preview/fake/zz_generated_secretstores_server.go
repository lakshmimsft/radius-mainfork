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

// SecretStoresServer is a fake server for instances of the v20231001preview.SecretStoresClient type.
type SecretStoresServer struct{
	// BeginCreateOrUpdate is the fake for method SecretStoresClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, secretStoreName string, resource v20231001preview.SecretStoreResource, options *v20231001preview.SecretStoresClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[v20231001preview.SecretStoresClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method SecretStoresClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, secretStoreName string, options *v20231001preview.SecretStoresClientBeginDeleteOptions) (resp azfake.PollerResponder[v20231001preview.SecretStoresClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method SecretStoresClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, secretStoreName string, options *v20231001preview.SecretStoresClientGetOptions) (resp azfake.Responder[v20231001preview.SecretStoresClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByScopePager is the fake for method SecretStoresClient.NewListByScopePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByScopePager func(options *v20231001preview.SecretStoresClientListByScopeOptions) (resp azfake.PagerResponder[v20231001preview.SecretStoresClientListByScopeResponse])

	// ListSecrets is the fake for method SecretStoresClient.ListSecrets
	// HTTP status codes to indicate success: http.StatusOK
	ListSecrets func(ctx context.Context, secretStoreName string, body map[string]any, options *v20231001preview.SecretStoresClientListSecretsOptions) (resp azfake.Responder[v20231001preview.SecretStoresClientListSecretsResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method SecretStoresClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, secretStoreName string, properties v20231001preview.SecretStoreResourceUpdate, options *v20231001preview.SecretStoresClientBeginUpdateOptions) (resp azfake.PollerResponder[v20231001preview.SecretStoresClientUpdateResponse], errResp azfake.ErrorResponder)

}

// NewSecretStoresServerTransport creates a new instance of SecretStoresServerTransport with the provided implementation.
// The returned SecretStoresServerTransport instance is connected to an instance of v20231001preview.SecretStoresClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSecretStoresServerTransport(srv *SecretStoresServer) *SecretStoresServerTransport {
	return &SecretStoresServerTransport{
		srv: srv,
		beginCreateOrUpdate: newTracker[azfake.PollerResponder[v20231001preview.SecretStoresClientCreateOrUpdateResponse]](),
		beginDelete: newTracker[azfake.PollerResponder[v20231001preview.SecretStoresClientDeleteResponse]](),
		newListByScopePager: newTracker[azfake.PagerResponder[v20231001preview.SecretStoresClientListByScopeResponse]](),
		beginUpdate: newTracker[azfake.PollerResponder[v20231001preview.SecretStoresClientUpdateResponse]](),
	}
}

// SecretStoresServerTransport connects instances of v20231001preview.SecretStoresClient to instances of SecretStoresServer.
// Don't use this type directly, use NewSecretStoresServerTransport instead.
type SecretStoresServerTransport struct {
	srv *SecretStoresServer
	beginCreateOrUpdate *tracker[azfake.PollerResponder[v20231001preview.SecretStoresClientCreateOrUpdateResponse]]
	beginDelete *tracker[azfake.PollerResponder[v20231001preview.SecretStoresClientDeleteResponse]]
	newListByScopePager *tracker[azfake.PagerResponder[v20231001preview.SecretStoresClientListByScopeResponse]]
	beginUpdate *tracker[azfake.PollerResponder[v20231001preview.SecretStoresClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for SecretStoresServerTransport.
func (s *SecretStoresServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return s.dispatchToMethodFake(req, method)
}

func (s *SecretStoresServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		 if secretStoresServerTransportInterceptor != nil {
			 res.resp, res.err, intercepted = secretStoresServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "SecretStoresClient.BeginCreateOrUpdate":
				res.resp, res.err = s.dispatchBeginCreateOrUpdate(req)
			case "SecretStoresClient.BeginDelete":
				res.resp, res.err = s.dispatchBeginDelete(req)
			case "SecretStoresClient.Get":
				res.resp, res.err = s.dispatchGet(req)
			case "SecretStoresClient.NewListByScopePager":
				res.resp, res.err = s.dispatchNewListByScopePager(req)
			case "SecretStoresClient.ListSecrets":
				res.resp, res.err = s.dispatchListSecrets(req)
			case "SecretStoresClient.BeginUpdate":
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

func (s *SecretStoresServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := s.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/secretStores/(?P<secretStoreName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[v20231001preview.SecretStoreResource](req)
	if err != nil {
		return nil, err
	}
	secretStoreNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("secretStoreName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.BeginCreateOrUpdate(req.Context(), secretStoreNameParam, body, nil)
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

func (s *SecretStoresServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if s.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := s.beginDelete.get(req)
	if beginDelete == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/secretStores/(?P<secretStoreName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	secretStoreNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("secretStoreName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.BeginDelete(req.Context(), secretStoreNameParam, nil)
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

func (s *SecretStoresServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/secretStores/(?P<secretStoreName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	secretStoreNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("secretStoreName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), secretStoreNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SecretStoreResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SecretStoresServerTransport) dispatchNewListByScopePager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByScopePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByScopePager not implemented")}
	}
	newListByScopePager := s.newListByScopePager.get(req)
	if newListByScopePager == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/secretStores`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
resp := s.srv.NewListByScopePager(nil)
		newListByScopePager = &resp
		s.newListByScopePager.add(req, newListByScopePager)
		server.PagerResponderInjectNextLinks(newListByScopePager, req, func(page *v20231001preview.SecretStoresClientListByScopeResponse, createLink func() string) {
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

func (s *SecretStoresServerTransport) dispatchListSecrets(req *http.Request) (*http.Response, error) {
	if s.srv.ListSecrets == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListSecrets not implemented")}
	}
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/secretStores/(?P<secretStoreName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listSecrets`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[map[string]any](req)
	if err != nil {
		return nil, err
	}
	secretStoreNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("secretStoreName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.ListSecrets(req.Context(), secretStoreNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SecretStoreListSecretsResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SecretStoresServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := s.beginUpdate.get(req)
	if beginUpdate == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/secretStores/(?P<secretStoreName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[v20231001preview.SecretStoreResourceUpdate](req)
	if err != nil {
		return nil, err
	}
	secretStoreNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("secretStoreName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.BeginUpdate(req.Context(), secretStoreNameParam, body, nil)
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

// set this to conditionally intercept incoming requests to SecretStoresServerTransport
var secretStoresServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}