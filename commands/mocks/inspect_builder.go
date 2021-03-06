// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildpack/pack/commands (interfaces: BuilderInspector)

// Package mocks is a generated GoMock package.
package mocks

import (
	pack "github.com/buildpack/pack"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockBuilderInspector is a mock of BuilderInspector interface
type MockBuilderInspector struct {
	ctrl     *gomock.Controller
	recorder *MockBuilderInspectorMockRecorder
}

// MockBuilderInspectorMockRecorder is the mock recorder for MockBuilderInspector
type MockBuilderInspectorMockRecorder struct {
	mock *MockBuilderInspector
}

// NewMockBuilderInspector creates a new mock instance
func NewMockBuilderInspector(ctrl *gomock.Controller) *MockBuilderInspector {
	mock := &MockBuilderInspector{ctrl: ctrl}
	mock.recorder = &MockBuilderInspectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBuilderInspector) EXPECT() *MockBuilderInspectorMockRecorder {
	return m.recorder
}

// InspectBuilder mocks base method
func (m *MockBuilderInspector) InspectBuilder(arg0 string, arg1 bool) (*pack.BuilderInfo, error) {
	ret := m.ctrl.Call(m, "InspectBuilder", arg0, arg1)
	ret0, _ := ret[0].(*pack.BuilderInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InspectBuilder indicates an expected call of InspectBuilder
func (mr *MockBuilderInspectorMockRecorder) InspectBuilder(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InspectBuilder", reflect.TypeOf((*MockBuilderInspector)(nil).InspectBuilder), arg0, arg1)
}
