package handlers

import (
	context "context"
	reflect "reflect"

	models "github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockJobStore is a mock of JobStore interface.
type MockJobStore struct {
	ctrl     *gomock.Controller
	recorder *MockJobStoreMockRecorder
}

// MockJobStoreMockRecorder is the mock recorder for MockJobStore.
type MockJobStoreMockRecorder struct {
	mock *MockJobStore
}

// NewMockJobStore creates a new mock instance.
func NewMockJobStore(ctrl *gomock.Controller) *MockJobStore {
	mock := &MockJobStore{ctrl: ctrl}
	mock.recorder = &MockJobStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJobStore) EXPECT() *MockJobStoreMockRecorder {
	return m.recorder
}

// AddJob mocks base method.
func (m *MockJobStore) AddJob(ctx context.Context, job models.JobStoreRow) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddJob", ctx, job)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddJob indicates an expected call of AddJob.
func (mr *MockJobStoreMockRecorder) AddJob(ctx, job interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddJob", reflect.TypeOf((*MockJobStore)(nil).AddJob), ctx, job)
}
