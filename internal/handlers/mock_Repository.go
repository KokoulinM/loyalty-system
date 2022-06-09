package handlers

import (
	"context"
	"reflect"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
	"github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	//mock.Mock
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CheckPassword mocks base method.
func (m *MockRepository) CheckPassword(ctx context.Context, user models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPassword", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPassword indicates an expected call of CheckPassword.
func (mr *MockRepositoryMockRecorder) CheckPassword(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPassword", reflect.TypeOf((*MockRepository)(nil).CheckPassword), ctx, user)
}

// CreateOrder mocks base method.
func (m *MockRepository) CreateOrder(ctx context.Context, order models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockRepositoryMockRecorder) CreateOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockRepository)(nil).CreateOrder), ctx, order)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), ctx, user)
}

// CreateWithdraw mocks base method.
func (m *MockRepository) CreateWithdraw(ctx context.Context, withdraw models.Withdraw, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWithdraw", ctx, withdraw, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateWithdraw indicates an expected call of CreateWithdraw.
func (mr *MockRepositoryMockRecorder) CreateWithdraw(ctx, withdraw, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWithdraw", reflect.TypeOf((*MockRepository)(nil).CreateWithdraw), ctx, withdraw, userID)
}

// GetBalance mocks base method.
func (m *MockRepository) GetBalance(ctx context.Context, userID string) (models.UserBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, userID)
	ret0, _ := ret[0].(models.UserBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockRepositoryMockRecorder) GetBalance(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockRepository)(nil).GetBalance), ctx, userID)
}

// GetOrders mocks base method.
func (m *MockRepository) GetOrders(ctx context.Context, userID string) ([]models.ResponseOrderWithAccrual, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, userID)
	ret0, _ := ret[0].([]models.ResponseOrderWithAccrual)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockRepositoryMockRecorder) GetOrders(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockRepository)(nil).GetOrders), ctx, userID)
}

// GetWithdrawals mocks base method.
func (m *MockRepository) GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawals", ctx, userID)
	ret0, _ := ret[0].([]models.WithdrawOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawals indicates an expected call of GetWithdrawals.
func (mr *MockRepositoryMockRecorder) GetWithdrawals(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawals", reflect.TypeOf((*MockRepository)(nil).GetWithdrawals), ctx, userID)
}
