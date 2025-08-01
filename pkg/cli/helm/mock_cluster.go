// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/radius-project/radius/pkg/cli/helm (interfaces: Interface)
//
// Generated by this command:
//
//	mockgen -typed -destination=./mock_cluster.go -package=helm -self_package github.com/radius-project/radius/pkg/cli/helm github.com/radius-project/radius/pkg/cli/helm Interface
//

// Package helm is a generated GoMock package.
package helm

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
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

// CheckRadiusInstall mocks base method.
func (m *MockInterface) CheckRadiusInstall(arg0 string) (InstallState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRadiusInstall", arg0)
	ret0, _ := ret[0].(InstallState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRadiusInstall indicates an expected call of CheckRadiusInstall.
func (mr *MockInterfaceMockRecorder) CheckRadiusInstall(arg0 any) *MockInterfaceCheckRadiusInstallCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRadiusInstall", reflect.TypeOf((*MockInterface)(nil).CheckRadiusInstall), arg0)
	return &MockInterfaceCheckRadiusInstallCall{Call: call}
}

// MockInterfaceCheckRadiusInstallCall wrap *gomock.Call
type MockInterfaceCheckRadiusInstallCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockInterfaceCheckRadiusInstallCall) Return(arg0 InstallState, arg1 error) *MockInterfaceCheckRadiusInstallCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockInterfaceCheckRadiusInstallCall) Do(f func(string) (InstallState, error)) *MockInterfaceCheckRadiusInstallCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockInterfaceCheckRadiusInstallCall) DoAndReturn(f func(string) (InstallState, error)) *MockInterfaceCheckRadiusInstallCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetLatestRadiusVersion mocks base method.
func (m *MockInterface) GetLatestRadiusVersion(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestRadiusVersion", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestRadiusVersion indicates an expected call of GetLatestRadiusVersion.
func (mr *MockInterfaceMockRecorder) GetLatestRadiusVersion(arg0 any) *MockInterfaceGetLatestRadiusVersionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestRadiusVersion", reflect.TypeOf((*MockInterface)(nil).GetLatestRadiusVersion), arg0)
	return &MockInterfaceGetLatestRadiusVersionCall{Call: call}
}

// MockInterfaceGetLatestRadiusVersionCall wrap *gomock.Call
type MockInterfaceGetLatestRadiusVersionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockInterfaceGetLatestRadiusVersionCall) Return(arg0 string, arg1 error) *MockInterfaceGetLatestRadiusVersionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockInterfaceGetLatestRadiusVersionCall) Do(f func(context.Context) (string, error)) *MockInterfaceGetLatestRadiusVersionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockInterfaceGetLatestRadiusVersionCall) DoAndReturn(f func(context.Context) (string, error)) *MockInterfaceGetLatestRadiusVersionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InstallRadius mocks base method.
func (m *MockInterface) InstallRadius(arg0 context.Context, arg1 ClusterOptions, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstallRadius", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// InstallRadius indicates an expected call of InstallRadius.
func (mr *MockInterfaceMockRecorder) InstallRadius(arg0, arg1, arg2 any) *MockInterfaceInstallRadiusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstallRadius", reflect.TypeOf((*MockInterface)(nil).InstallRadius), arg0, arg1, arg2)
	return &MockInterfaceInstallRadiusCall{Call: call}
}

// MockInterfaceInstallRadiusCall wrap *gomock.Call
type MockInterfaceInstallRadiusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockInterfaceInstallRadiusCall) Return(arg0 error) *MockInterfaceInstallRadiusCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockInterfaceInstallRadiusCall) Do(f func(context.Context, ClusterOptions, string) error) *MockInterfaceInstallRadiusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockInterfaceInstallRadiusCall) DoAndReturn(f func(context.Context, ClusterOptions, string) error) *MockInterfaceInstallRadiusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UninstallRadius mocks base method.
func (m *MockInterface) UninstallRadius(arg0 context.Context, arg1 ClusterOptions, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UninstallRadius", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UninstallRadius indicates an expected call of UninstallRadius.
func (mr *MockInterfaceMockRecorder) UninstallRadius(arg0, arg1, arg2 any) *MockInterfaceUninstallRadiusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UninstallRadius", reflect.TypeOf((*MockInterface)(nil).UninstallRadius), arg0, arg1, arg2)
	return &MockInterfaceUninstallRadiusCall{Call: call}
}

// MockInterfaceUninstallRadiusCall wrap *gomock.Call
type MockInterfaceUninstallRadiusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockInterfaceUninstallRadiusCall) Return(arg0 error) *MockInterfaceUninstallRadiusCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockInterfaceUninstallRadiusCall) Do(f func(context.Context, ClusterOptions, string) error) *MockInterfaceUninstallRadiusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockInterfaceUninstallRadiusCall) DoAndReturn(f func(context.Context, ClusterOptions, string) error) *MockInterfaceUninstallRadiusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpgradeRadius mocks base method.
func (m *MockInterface) UpgradeRadius(arg0 context.Context, arg1 ClusterOptions, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpgradeRadius", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpgradeRadius indicates an expected call of UpgradeRadius.
func (mr *MockInterfaceMockRecorder) UpgradeRadius(arg0, arg1, arg2 any) *MockInterfaceUpgradeRadiusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeRadius", reflect.TypeOf((*MockInterface)(nil).UpgradeRadius), arg0, arg1, arg2)
	return &MockInterfaceUpgradeRadiusCall{Call: call}
}

// MockInterfaceUpgradeRadiusCall wrap *gomock.Call
type MockInterfaceUpgradeRadiusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockInterfaceUpgradeRadiusCall) Return(arg0 error) *MockInterfaceUpgradeRadiusCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockInterfaceUpgradeRadiusCall) Do(f func(context.Context, ClusterOptions, string) error) *MockInterfaceUpgradeRadiusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockInterfaceUpgradeRadiusCall) DoAndReturn(f func(context.Context, ClusterOptions, string) error) *MockInterfaceUpgradeRadiusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
