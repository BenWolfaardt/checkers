// Code generated by MockGen. DO NOT EDIT.
// Source: x/leaderboard/types/expected_keepers.go

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	types "github.com/BenWolfaardt/checkers/x/checkers/types"
	types0 "github.com/cosmos/cosmos-sdk/types"
	types1 "github.com/cosmos/cosmos-sdk/x/auth/types"
	gomock "github.com/golang/mock/gomock"
)

// MockAccountKeeper is a mock of AccountKeeper interface.
type MockAccountKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountKeeperMockRecorder
}

// MockAccountKeeperMockRecorder is the mock recorder for MockAccountKeeper.
type MockAccountKeeperMockRecorder struct {
	mock *MockAccountKeeper
}

// NewMockAccountKeeper creates a new mock instance.
func NewMockAccountKeeper(ctrl *gomock.Controller) *MockAccountKeeper {
	mock := &MockAccountKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountKeeper) EXPECT() *MockAccountKeeperMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockAccountKeeper) GetAccount(ctx types0.Context, addr types0.AccAddress) types1.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, addr)
	ret0, _ := ret[0].(types1.AccountI)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountKeeperMockRecorder) GetAccount(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), ctx, addr)
}

// MockBankKeeper is a mock of BankKeeper interface.
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper.
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance.
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// SpendableCoins mocks base method.
func (m *MockBankKeeper) SpendableCoins(ctx types0.Context, addr types0.AccAddress) types0.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendableCoins", ctx, addr)
	ret0, _ := ret[0].(types0.Coins)
	return ret0
}

// SpendableCoins indicates an expected call of SpendableCoins.
func (mr *MockBankKeeperMockRecorder) SpendableCoins(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendableCoins", reflect.TypeOf((*MockBankKeeper)(nil).SpendableCoins), ctx, addr)
}

// MockPlayerInfoKeeper is a mock of PlayerInfoKeeper interface.
type MockPlayerInfoKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockPlayerInfoKeeperMockRecorder
}

// MockPlayerInfoKeeperMockRecorder is the mock recorder for MockPlayerInfoKeeper.
type MockPlayerInfoKeeperMockRecorder struct {
	mock *MockPlayerInfoKeeper
}

// NewMockPlayerInfoKeeper creates a new mock instance.
func NewMockPlayerInfoKeeper(ctrl *gomock.Controller) *MockPlayerInfoKeeper {
	mock := &MockPlayerInfoKeeper{ctrl: ctrl}
	mock.recorder = &MockPlayerInfoKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlayerInfoKeeper) EXPECT() *MockPlayerInfoKeeperMockRecorder {
	return m.recorder
}

// PlayerInfoAll mocks base method.
func (m *MockPlayerInfoKeeper) PlayerInfoAll(c context.Context, req *types.QueryAllPlayerInfoRequest) (*types.QueryAllPlayerInfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlayerInfoAll", c, req)
	ret0, _ := ret[0].(*types.QueryAllPlayerInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlayerInfoAll indicates an expected call of PlayerInfoAll.
func (mr *MockPlayerInfoKeeperMockRecorder) PlayerInfoAll(c, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlayerInfoAll", reflect.TypeOf((*MockPlayerInfoKeeper)(nil).PlayerInfoAll), c, req)
}
