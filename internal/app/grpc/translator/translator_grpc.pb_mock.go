// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/takakd/prj/retranslation/src/api/scripts/../internal/app/grpc/translator/translator_grpc.pb.go

// Package translator is a generated GoMock package.
package translator

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockTranslatorClient is a mock of TranslatorClient interface
type MockTranslatorClient struct {
	ctrl     *gomock.Controller
	recorder *MockTranslatorClientMockRecorder
}

// MockTranslatorClientMockRecorder is the mock recorder for MockTranslatorClient
type MockTranslatorClientMockRecorder struct {
	mock *MockTranslatorClient
}

// NewMockTranslatorClient creates a new mock instance
func NewMockTranslatorClient(ctrl *gomock.Controller) *MockTranslatorClient {
	mock := &MockTranslatorClient{ctrl: ctrl}
	mock.recorder = &MockTranslatorClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTranslatorClient) EXPECT() *MockTranslatorClientMockRecorder {
	return m.recorder
}

// Translate mocks base method
func (m *MockTranslatorClient) Translate(ctx context.Context, in *TranslateRequest, opts ...grpc.CallOption) (*TranslateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Translate", varargs...)
	ret0, _ := ret[0].(*TranslateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Translate indicates an expected call of Translate
func (mr *MockTranslatorClientMockRecorder) Translate(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Translate", reflect.TypeOf((*MockTranslatorClient)(nil).Translate), varargs...)
}

// MockTranslatorServer is a mock of TranslatorServer interface
type MockTranslatorServer struct {
	ctrl     *gomock.Controller
	recorder *MockTranslatorServerMockRecorder
}

// MockTranslatorServerMockRecorder is the mock recorder for MockTranslatorServer
type MockTranslatorServerMockRecorder struct {
	mock *MockTranslatorServer
}

// NewMockTranslatorServer creates a new mock instance
func NewMockTranslatorServer(ctrl *gomock.Controller) *MockTranslatorServer {
	mock := &MockTranslatorServer{ctrl: ctrl}
	mock.recorder = &MockTranslatorServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTranslatorServer) EXPECT() *MockTranslatorServerMockRecorder {
	return m.recorder
}

// Translate mocks base method
func (m *MockTranslatorServer) Translate(arg0 context.Context, arg1 *TranslateRequest) (*TranslateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Translate", arg0, arg1)
	ret0, _ := ret[0].(*TranslateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Translate indicates an expected call of Translate
func (mr *MockTranslatorServerMockRecorder) Translate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Translate", reflect.TypeOf((*MockTranslatorServer)(nil).Translate), arg0, arg1)
}

// mustEmbedUnimplementedTranslatorServer mocks base method
func (m *MockTranslatorServer) mustEmbedUnimplementedTranslatorServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedTranslatorServer")
}

// mustEmbedUnimplementedTranslatorServer indicates an expected call of mustEmbedUnimplementedTranslatorServer
func (mr *MockTranslatorServerMockRecorder) mustEmbedUnimplementedTranslatorServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedTranslatorServer", reflect.TypeOf((*MockTranslatorServer)(nil).mustEmbedUnimplementedTranslatorServer))
}

// MockUnsafeTranslatorServer is a mock of UnsafeTranslatorServer interface
type MockUnsafeTranslatorServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeTranslatorServerMockRecorder
}

// MockUnsafeTranslatorServerMockRecorder is the mock recorder for MockUnsafeTranslatorServer
type MockUnsafeTranslatorServerMockRecorder struct {
	mock *MockUnsafeTranslatorServer
}

// NewMockUnsafeTranslatorServer creates a new mock instance
func NewMockUnsafeTranslatorServer(ctrl *gomock.Controller) *MockUnsafeTranslatorServer {
	mock := &MockUnsafeTranslatorServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeTranslatorServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUnsafeTranslatorServer) EXPECT() *MockUnsafeTranslatorServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedTranslatorServer mocks base method
func (m *MockUnsafeTranslatorServer) mustEmbedUnimplementedTranslatorServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedTranslatorServer")
}

// mustEmbedUnimplementedTranslatorServer indicates an expected call of mustEmbedUnimplementedTranslatorServer
func (mr *MockUnsafeTranslatorServerMockRecorder) mustEmbedUnimplementedTranslatorServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedTranslatorServer", reflect.TypeOf((*MockUnsafeTranslatorServer)(nil).mustEmbedUnimplementedTranslatorServer))
}
