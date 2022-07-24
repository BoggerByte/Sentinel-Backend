// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory (interfaces: Store)

// Package mockmemdb is a generated GoMock package.
package mockmemdb

import (
	context "context"
	reflect "reflect"
	time "time"

	memdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// DeleteOauth2Flow mocks base method.
func (m *MockStore) DeleteOauth2Flow(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOauth2Flow", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOauth2Flow indicates an expected call of DeleteOauth2Flow.
func (mr *MockStoreMockRecorder) DeleteOauth2Flow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOauth2Flow", reflect.TypeOf((*MockStore)(nil).DeleteOauth2Flow), arg0, arg1)
}

// GetOauth2Flow mocks base method.
func (m *MockStore) GetOauth2Flow(arg0 context.Context, arg1 string) (memdb.Oauth2Flow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOauth2Flow", arg0, arg1)
	ret0, _ := ret[0].(memdb.Oauth2Flow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOauth2Flow indicates an expected call of GetOauth2Flow.
func (mr *MockStoreMockRecorder) GetOauth2Flow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOauth2Flow", reflect.TypeOf((*MockStore)(nil).GetOauth2Flow), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockStore) GetSession(arg0 context.Context, arg1 uuid.UUID) (memdb.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(memdb.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockStoreMockRecorder) GetSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockStore)(nil).GetSession), arg0, arg1)
}

// SetOauth2Flow mocks base method.
func (m *MockStore) SetOauth2Flow(arg0 context.Context, arg1 string, arg2 memdb.Oauth2Flow, arg3 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOauth2Flow", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOauth2Flow indicates an expected call of SetOauth2Flow.
func (mr *MockStoreMockRecorder) SetOauth2Flow(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOauth2Flow", reflect.TypeOf((*MockStore)(nil).SetOauth2Flow), arg0, arg1, arg2, arg3)
}

// SetSession mocks base method.
func (m *MockStore) SetSession(arg0 context.Context, arg1 memdb.Session, arg2 time.Duration) (memdb.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSession", arg0, arg1, arg2)
	ret0, _ := ret[0].(memdb.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetSession indicates an expected call of SetSession.
func (mr *MockStoreMockRecorder) SetSession(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSession", reflect.TypeOf((*MockStore)(nil).SetSession), arg0, arg1, arg2)
}
