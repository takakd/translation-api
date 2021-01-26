// Code generated by MockGen. DO NOT EDIT.
// Source: ../internal/app/driver/google/clientwrapper.go

// Package google is a generated GoMock package.
package google

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	gax "github.com/googleapis/gax-go/v2"
	translate "google.golang.org/genproto/googleapis/cloud/translate/v3"
	reflect "reflect"
)

// MockClientWrapper is a mock of ClientWrapper interface
type MockClientWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockClientWrapperMockRecorder
}

// MockClientWrapperMockRecorder is the mock recorder for MockClientWrapper
type MockClientWrapperMockRecorder struct {
	mock *MockClientWrapper
}

// NewMockClientWrapper creates a new mock instance
func NewMockClientWrapper(ctrl *gomock.Controller) *MockClientWrapper {
	mock := &MockClientWrapper{ctrl: ctrl}
	mock.recorder = &MockClientWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClientWrapper) EXPECT() *MockClientWrapperMockRecorder {
	return m.recorder
}

// TranslateText mocks base method
func (m *MockClientWrapper) TranslateText(ctx context.Context, req *translate.TranslateTextRequest, opts ...gax.CallOption) (*translate.TranslateTextResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, req}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TranslateText", varargs...)
	ret0, _ := ret[0].(*translate.TranslateTextResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TranslateText indicates an expected call of TranslateText
func (mr *MockClientWrapperMockRecorder) TranslateText(ctx, req interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, req}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TranslateText", reflect.TypeOf((*MockClientWrapper)(nil).TranslateText), varargs...)
}

// Close mocks base method
func (m *MockClientWrapper) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockClientWrapperMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClientWrapper)(nil).Close))
}
