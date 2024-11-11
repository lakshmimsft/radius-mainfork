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

// GatewaysServer is a fake server for instances of the v20231001preview.GatewaysClient type.
type GatewaysServer struct{
	// BeginCreate is the fake for method GatewaysClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, gatewayName string, resource v20231001preview.GatewayResource, options *v20231001preview.GatewaysClientBeginCreateOptions) (resp azfake.PollerResponder[v20231001preview.GatewaysClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginCreateOrUpdate is the fake for method GatewaysClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreateOrUpdate func(ctx context.Context, gatewayName string, properties v20231001preview.GatewayResourceUpdate, options *v20231001preview.GatewaysClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[v20231001preview.GatewaysClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method GatewaysClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, gatewayName string, options *v20231001preview.GatewaysClientBeginDeleteOptions) (resp azfake.PollerResponder[v20231001preview.GatewaysClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method GatewaysClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, gatewayName string, options *v20231001preview.GatewaysClientGetOptions) (resp azfake.Responder[v20231001preview.GatewaysClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByScopePager is the fake for method GatewaysClient.NewListByScopePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByScopePager func(options *v20231001preview.GatewaysClientListByScopeOptions) (resp azfake.PagerResponder[v20231001preview.GatewaysClientListByScopeResponse])

}

// NewGatewaysServerTransport creates a new instance of GatewaysServerTransport with the provided implementation.
// The returned GatewaysServerTransport instance is connected to an instance of v20231001preview.GatewaysClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewGatewaysServerTransport(srv *GatewaysServer) *GatewaysServerTransport {
	return &GatewaysServerTransport{
		srv: srv,
		beginCreate: newTracker[azfake.PollerResponder[v20231001preview.GatewaysClientCreateResponse]](),
		beginCreateOrUpdate: newTracker[azfake.PollerResponder[v20231001preview.GatewaysClientCreateOrUpdateResponse]](),
		beginDelete: newTracker[azfake.PollerResponder[v20231001preview.GatewaysClientDeleteResponse]](),
		newListByScopePager: newTracker[azfake.PagerResponder[v20231001preview.GatewaysClientListByScopeResponse]](),
	}
}

// GatewaysServerTransport connects instances of v20231001preview.GatewaysClient to instances of GatewaysServer.
// Don't use this type directly, use NewGatewaysServerTransport instead.
type GatewaysServerTransport struct {
	srv *GatewaysServer
	beginCreate *tracker[azfake.PollerResponder[v20231001preview.GatewaysClientCreateResponse]]
	beginCreateOrUpdate *tracker[azfake.PollerResponder[v20231001preview.GatewaysClientCreateOrUpdateResponse]]
	beginDelete *tracker[azfake.PollerResponder[v20231001preview.GatewaysClientDeleteResponse]]
	newListByScopePager *tracker[azfake.PagerResponder[v20231001preview.GatewaysClientListByScopeResponse]]
}

// Do implements the policy.Transporter interface for GatewaysServerTransport.
func (g *GatewaysServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return g.dispatchToMethodFake(req, method)
}

func (g *GatewaysServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		 if gatewaysServerTransportInterceptor != nil {
			 res.resp, res.err, intercepted = gatewaysServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "GatewaysClient.BeginCreate":
				res.resp, res.err = g.dispatchBeginCreate(req)
			case "GatewaysClient.BeginCreateOrUpdate":
				res.resp, res.err = g.dispatchBeginCreateOrUpdate(req)
			case "GatewaysClient.BeginDelete":
				res.resp, res.err = g.dispatchBeginDelete(req)
			case "GatewaysClient.Get":
				res.resp, res.err = g.dispatchGet(req)
			case "GatewaysClient.NewListByScopePager":
				res.resp, res.err = g.dispatchNewListByScopePager(req)
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

func (g *GatewaysServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if g.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := g.beginCreate.get(req)
	if beginCreate == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/gateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[v20231001preview.GatewayResource](req)
	if err != nil {
		return nil, err
	}
	gatewayNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("gatewayName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.BeginCreate(req.Context(), gatewayNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
		beginCreate = &respr
		g.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		g.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		g.beginCreate.remove(req)
	}

	return resp, nil
}

func (g *GatewaysServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if g.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := g.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/gateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[v20231001preview.GatewayResourceUpdate](req)
	if err != nil {
		return nil, err
	}
	gatewayNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("gatewayName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.BeginCreateOrUpdate(req.Context(), gatewayNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
		beginCreateOrUpdate = &respr
		g.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		g.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		g.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (g *GatewaysServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if g.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := g.beginDelete.get(req)
	if beginDelete == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/gateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	gatewayNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("gatewayName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.BeginDelete(req.Context(), gatewayNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
		beginDelete = &respr
		g.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		g.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		g.beginDelete.remove(req)
	}

	return resp, nil
}

func (g *GatewaysServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if g.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/gateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	gatewayNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("gatewayName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.Get(req.Context(), gatewayNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).GatewayResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GatewaysServerTransport) dispatchNewListByScopePager(req *http.Request) (*http.Response, error) {
	if g.srv.NewListByScopePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByScopePager not implemented")}
	}
	newListByScopePager := g.newListByScopePager.get(req)
	if newListByScopePager == nil {
	const regexStr = `/(?P<rootScope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Applications\.Core/gateways`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
resp := g.srv.NewListByScopePager(nil)
		newListByScopePager = &resp
		g.newListByScopePager.add(req, newListByScopePager)
		server.PagerResponderInjectNextLinks(newListByScopePager, req, func(page *v20231001preview.GatewaysClientListByScopeResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByScopePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		g.newListByScopePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByScopePager) {
		g.newListByScopePager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to GatewaysServerTransport
var gatewaysServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}