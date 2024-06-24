// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/tokenizer/tokenizer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	jwt "github.com/dgrijalva/jwt-go/v4"
	gomock "github.com/golang/mock/gomock"
)

// MockITokenizer is a mock of ITokenizer interface.
type MockITokenizer struct {
	ctrl     *gomock.Controller
	recorder *MockITokenizerMockRecorder
}

// MockITokenizerMockRecorder is the mock recorder for MockITokenizer.
type MockITokenizerMockRecorder struct {
	mock *MockITokenizer
}

// NewMockITokenizer creates a new mock instance.
func NewMockITokenizer(ctrl *gomock.Controller) *MockITokenizer {
	mock := &MockITokenizer{ctrl: ctrl}
	mock.recorder = &MockITokenizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITokenizer) EXPECT() *MockITokenizerMockRecorder {
	return m.recorder
}

// IsBasicAuthorized mocks base method.
func (m *MockITokenizer) IsBasicAuthorized(token string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsBasicAuthorized", token)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsBasicAuthorized indicates an expected call of IsBasicAuthorized.
func (mr *MockITokenizerMockRecorder) IsBasicAuthorized(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBasicAuthorized", reflect.TypeOf((*MockITokenizer)(nil).IsBasicAuthorized), token)
}

// VerifyJWTToken mocks base method.
func (m *MockITokenizer) VerifyJWTToken(token string) (*jwt.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyJWTToken", token)
	ret0, _ := ret[0].(*jwt.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyJWTToken indicates an expected call of VerifyJWTToken.
func (mr *MockITokenizerMockRecorder) VerifyJWTToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyJWTToken", reflect.TypeOf((*MockITokenizer)(nil).VerifyJWTToken), token)
}
