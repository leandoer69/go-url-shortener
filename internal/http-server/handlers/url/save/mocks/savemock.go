// Code generated by MockGen. DO NOT EDIT.
// Source: save.go

// Package savemock is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockURLSaver is a mock of URLSaver interface.
type MockURLSaver struct {
	ctrl     *gomock.Controller
	recorder *MockURLSaverMockRecorder
}

// MockURLSaverMockRecorder is the mock recorder for MockURLSaver.
type MockURLSaverMockRecorder struct {
	mock *MockURLSaver
}

// NewMockURLSaver creates a new mock instance.
func NewMockURLSaver(ctrl *gomock.Controller) *MockURLSaver {
	mock := &MockURLSaver{ctrl: ctrl}
	mock.recorder = &MockURLSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLSaver) EXPECT() *MockURLSaverMockRecorder {
	return m.recorder
}

// SaveURL mocks base method.
func (m *MockURLSaver) SaveURL(urlToSave, alias string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveURL", urlToSave, alias)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveURL indicates an expected call of SaveURL.
func (mr *MockURLSaverMockRecorder) SaveURL(urlToSave, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveURL", reflect.TypeOf((*MockURLSaver)(nil).SaveURL), urlToSave, alias)
}
