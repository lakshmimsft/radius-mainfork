// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/radius-project/radius/pkg/cli/kubernetes (interfaces: Interface)
//
// Generated by this command:
//
//	mockgen -typed -destination=./mock_kubernetes.go -package=kubernetes -self_package github.com/radius-project/radius/pkg/cli/kubernetes github.com/radius-project/radius/pkg/cli/kubernetes Interface
//

// Package kubernetes is a generated GoMock package.
package kubernetes

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	api "k8s.io/client-go/tools/clientcmd/api"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// DeleteNamespace mocks base method.
func (m *MockInterface) DeleteNamespace(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNamespace", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNamespace indicates an expected call of DeleteNamespace.
func (mr *MockInterfaceMockRecorder) DeleteNamespace(arg0 any) *MockInterfaceDeleteNamespaceCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNamespace", reflect.TypeOf((*MockInterface)(nil).DeleteNamespace), arg0)
	return &MockInterfaceDeleteNamespaceCall{Call: call}
}

// MockInterfaceDeleteNamespaceCall wrap *gomock.Call
type MockInterfaceDeleteNamespaceCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockInterfaceDeleteNamespaceCall) Return(arg0 error) *MockInterfaceDeleteNamespaceCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockInterfaceDeleteNamespaceCall) Do(f func(string) error) *MockInterfaceDeleteNamespaceCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockInterfaceDeleteNamespaceCall) DoAndReturn(f func(string) error) *MockInterfaceDeleteNamespaceCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetKubeContext mocks base method.
func (m *MockInterface) GetKubeContext() (*api.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKubeContext")
	ret0, _ := ret[0].(*api.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKubeContext indicates an expected call of GetKubeContext.
func (mr *MockInterfaceMockRecorder) GetKubeContext() *MockInterfaceGetKubeContextCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKubeContext", reflect.TypeOf((*MockInterface)(nil).GetKubeContext))
	return &MockInterfaceGetKubeContextCall{Call: call}
}

// MockInterfaceGetKubeContextCall wrap *gomock.Call
type MockInterfaceGetKubeContextCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockInterfaceGetKubeContextCall) Return(arg0 *api.Config, arg1 error) *MockInterfaceGetKubeContextCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockInterfaceGetKubeContextCall) Do(f func() (*api.Config, error)) *MockInterfaceGetKubeContextCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockInterfaceGetKubeContextCall) DoAndReturn(f func() (*api.Config, error)) *MockInterfaceGetKubeContextCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
