// Code generated by MockGen. DO NOT EDIT.
// Source: base.go

// Package runner is a generated GoMock package.
package runner

import (
	context "context"
	reflect "reflect"

	service "github.com/achu-1612/glcm/service"
	gomock "github.com/golang/mock/gomock"
)

// MockBase is a mock of Base interface.
type MockBase struct {
	ctrl     *gomock.Controller
	recorder *MockBaseMockRecorder
}

// MockBaseMockRecorder is the mock recorder for MockBase.
type MockBaseMockRecorder struct {
	mock *MockBase
}

// NewMockBase creates a new mock instance.
func NewMockBase(ctrl *gomock.Controller) *MockBase {
	mock := &MockBase{ctrl: ctrl}
	mock.recorder = &MockBaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBase) EXPECT() *MockBaseMockRecorder {
	return m.recorder
}

// BootUp mocks base method.
func (m *MockBase) BootUp(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BootUp", arg0)
}

// BootUp indicates an expected call of BootUp.
func (mr *MockBaseMockRecorder) BootUp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BootUp", reflect.TypeOf((*MockBase)(nil).BootUp), arg0)
}

// IsRunning mocks base method.
func (m *MockBase) IsRunning() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRunning indicates an expected call of IsRunning.
func (mr *MockBaseMockRecorder) IsRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRunning", reflect.TypeOf((*MockBase)(nil).IsRunning))
}

// RegisterService mocks base method.
func (m *MockBase) RegisterService(arg0 service.Service, arg1 ...service.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RegisterService", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterService indicates an expected call of RegisterService.
func (mr *MockBaseMockRecorder) RegisterService(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterService", reflect.TypeOf((*MockBase)(nil).RegisterService), varargs...)
}

// RestartAllServices mocks base method.
func (m *MockBase) RestartAllServices() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RestartAllServices")
}

// RestartAllServices indicates an expected call of RestartAllServices.
func (mr *MockBaseMockRecorder) RestartAllServices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestartAllServices", reflect.TypeOf((*MockBase)(nil).RestartAllServices))
}

// RestartService mocks base method.
func (m *MockBase) RestartService(arg0 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RestartService", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestartService indicates an expected call of RestartService.
func (mr *MockBaseMockRecorder) RestartService(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestartService", reflect.TypeOf((*MockBase)(nil).RestartService), arg0...)
}

// Shutdown mocks base method.
func (m *MockBase) Shutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown")
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockBaseMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockBase)(nil).Shutdown))
}

// StopAllServices mocks base method.
func (m *MockBase) StopAllServices() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopAllServices")
}

// StopAllServices indicates an expected call of StopAllServices.
func (mr *MockBaseMockRecorder) StopAllServices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopAllServices", reflect.TypeOf((*MockBase)(nil).StopAllServices))
}

// StopService mocks base method.
func (m *MockBase) StopService(arg0 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StopService", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopService indicates an expected call of StopService.
func (mr *MockBaseMockRecorder) StopService(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopService", reflect.TypeOf((*MockBase)(nil).StopService), arg0...)
}
