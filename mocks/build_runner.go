// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildpack/pack (interfaces: BuildRunner)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockBuildRunner is a mock of BuildRunner interface
type MockBuildRunner struct {
	ctrl     *gomock.Controller
	recorder *MockBuildRunnerMockRecorder
}

// MockBuildRunnerMockRecorder is the mock recorder for MockBuildRunner
type MockBuildRunnerMockRecorder struct {
	mock *MockBuildRunner
}

// NewMockBuildRunner creates a new mock instance
func NewMockBuildRunner(ctrl *gomock.Controller) *MockBuildRunner {
	mock := &MockBuildRunner{ctrl: ctrl}
	mock.recorder = &MockBuildRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBuildRunner) EXPECT() *MockBuildRunnerMockRecorder {
	return m.recorder
}

// Run mocks base method
func (m *MockBuildRunner) Run(arg0 context.Context) error {
	ret := m.ctrl.Call(m, "Run", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run
func (mr *MockBuildRunnerMockRecorder) Run(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockBuildRunner)(nil).Run), arg0)
}
