// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: bridge.go
//
// Generated by this command:
//
//	mockgen -source=bridge.go -destination=mock/bridge.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/celenium-io/astria-indexer/internal/storage"
	storage0 "github.com/dipdup-net/indexer-sdk/pkg/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIBridge is a mock of IBridge interface.
type MockIBridge struct {
	ctrl     *gomock.Controller
	recorder *MockIBridgeMockRecorder
}

// MockIBridgeMockRecorder is the mock recorder for MockIBridge.
type MockIBridgeMockRecorder struct {
	mock *MockIBridge
}

// NewMockIBridge creates a new mock instance.
func NewMockIBridge(ctrl *gomock.Controller) *MockIBridge {
	mock := &MockIBridge{ctrl: ctrl}
	mock.recorder = &MockIBridgeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBridge) EXPECT() *MockIBridgeMockRecorder {
	return m.recorder
}

// ByAddress mocks base method.
func (m *MockIBridge) ByAddress(ctx context.Context, addressId uint64) (storage.Bridge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByAddress", ctx, addressId)
	ret0, _ := ret[0].(storage.Bridge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByAddress indicates an expected call of ByAddress.
func (mr *MockIBridgeMockRecorder) ByAddress(ctx, addressId any) *MockIBridgeByAddressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByAddress", reflect.TypeOf((*MockIBridge)(nil).ByAddress), ctx, addressId)
	return &MockIBridgeByAddressCall{Call: call}
}

// MockIBridgeByAddressCall wrap *gomock.Call
type MockIBridgeByAddressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeByAddressCall) Return(arg0 storage.Bridge, arg1 error) *MockIBridgeByAddressCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeByAddressCall) Do(f func(context.Context, uint64) (storage.Bridge, error)) *MockIBridgeByAddressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeByAddressCall) DoAndReturn(f func(context.Context, uint64) (storage.Bridge, error)) *MockIBridgeByAddressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ById mocks base method.
func (m *MockIBridge) ById(ctx context.Context, id uint64) (storage.Bridge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ById", ctx, id)
	ret0, _ := ret[0].(storage.Bridge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ById indicates an expected call of ById.
func (mr *MockIBridgeMockRecorder) ById(ctx, id any) *MockIBridgeByIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ById", reflect.TypeOf((*MockIBridge)(nil).ById), ctx, id)
	return &MockIBridgeByIdCall{Call: call}
}

// MockIBridgeByIdCall wrap *gomock.Call
type MockIBridgeByIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeByIdCall) Return(arg0 storage.Bridge, arg1 error) *MockIBridgeByIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeByIdCall) Do(f func(context.Context, uint64) (storage.Bridge, error)) *MockIBridgeByIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeByIdCall) DoAndReturn(f func(context.Context, uint64) (storage.Bridge, error)) *MockIBridgeByIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByRoles mocks base method.
func (m *MockIBridge) ByRoles(ctx context.Context, addressId uint64, limit, offset int) ([]storage.Bridge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByRoles", ctx, addressId, limit, offset)
	ret0, _ := ret[0].([]storage.Bridge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByRoles indicates an expected call of ByRoles.
func (mr *MockIBridgeMockRecorder) ByRoles(ctx, addressId, limit, offset any) *MockIBridgeByRolesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByRoles", reflect.TypeOf((*MockIBridge)(nil).ByRoles), ctx, addressId, limit, offset)
	return &MockIBridgeByRolesCall{Call: call}
}

// MockIBridgeByRolesCall wrap *gomock.Call
type MockIBridgeByRolesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeByRolesCall) Return(arg0 []storage.Bridge, arg1 error) *MockIBridgeByRolesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeByRolesCall) Do(f func(context.Context, uint64, int, int) ([]storage.Bridge, error)) *MockIBridgeByRolesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeByRolesCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Bridge, error)) *MockIBridgeByRolesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByRollup mocks base method.
func (m *MockIBridge) ByRollup(ctx context.Context, rollupId uint64, limit, offset int) ([]storage.Bridge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByRollup", ctx, rollupId, limit, offset)
	ret0, _ := ret[0].([]storage.Bridge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByRollup indicates an expected call of ByRollup.
func (mr *MockIBridgeMockRecorder) ByRollup(ctx, rollupId, limit, offset any) *MockIBridgeByRollupCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByRollup", reflect.TypeOf((*MockIBridge)(nil).ByRollup), ctx, rollupId, limit, offset)
	return &MockIBridgeByRollupCall{Call: call}
}

// MockIBridgeByRollupCall wrap *gomock.Call
type MockIBridgeByRollupCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeByRollupCall) Return(arg0 []storage.Bridge, arg1 error) *MockIBridgeByRollupCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeByRollupCall) Do(f func(context.Context, uint64, int, int) ([]storage.Bridge, error)) *MockIBridgeByRollupCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeByRollupCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Bridge, error)) *MockIBridgeByRollupCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIBridge) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Bridge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Bridge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIBridgeMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIBridgeCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIBridge)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIBridgeCursorListCall{Call: call}
}

// MockIBridgeCursorListCall wrap *gomock.Call
type MockIBridgeCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeCursorListCall) Return(arg0 []*storage.Bridge, arg1 error) *MockIBridgeCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Bridge, error)) *MockIBridgeCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Bridge, error)) *MockIBridgeCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIBridge) GetByID(ctx context.Context, id uint64) (*storage.Bridge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Bridge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIBridgeMockRecorder) GetByID(ctx, id any) *MockIBridgeGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIBridge)(nil).GetByID), ctx, id)
	return &MockIBridgeGetByIDCall{Call: call}
}

// MockIBridgeGetByIDCall wrap *gomock.Call
type MockIBridgeGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeGetByIDCall) Return(arg0 *storage.Bridge, arg1 error) *MockIBridgeGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeGetByIDCall) Do(f func(context.Context, uint64) (*storage.Bridge, error)) *MockIBridgeGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Bridge, error)) *MockIBridgeGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIBridge) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIBridgeMockRecorder) IsNoRows(err any) *MockIBridgeIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIBridge)(nil).IsNoRows), err)
	return &MockIBridgeIsNoRowsCall{Call: call}
}

// MockIBridgeIsNoRowsCall wrap *gomock.Call
type MockIBridgeIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeIsNoRowsCall) Return(arg0 bool) *MockIBridgeIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeIsNoRowsCall) Do(f func(error) bool) *MockIBridgeIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIBridgeIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIBridge) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIBridgeMockRecorder) LastID(ctx any) *MockIBridgeLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIBridge)(nil).LastID), ctx)
	return &MockIBridgeLastIDCall{Call: call}
}

// MockIBridgeLastIDCall wrap *gomock.Call
type MockIBridgeLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeLastIDCall) Return(arg0 uint64, arg1 error) *MockIBridgeLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIBridgeLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIBridgeLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIBridge) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Bridge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Bridge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIBridgeMockRecorder) List(ctx, limit, offset, order any) *MockIBridgeListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIBridge)(nil).List), ctx, limit, offset, order)
	return &MockIBridgeListCall{Call: call}
}

// MockIBridgeListCall wrap *gomock.Call
type MockIBridgeListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeListCall) Return(arg0 []*storage.Bridge, arg1 error) *MockIBridgeListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Bridge, error)) *MockIBridgeListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Bridge, error)) *MockIBridgeListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListWithAddress mocks base method.
func (m *MockIBridge) ListWithAddress(ctx context.Context, limit, offset int) ([]storage.Bridge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWithAddress", ctx, limit, offset)
	ret0, _ := ret[0].([]storage.Bridge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWithAddress indicates an expected call of ListWithAddress.
func (mr *MockIBridgeMockRecorder) ListWithAddress(ctx, limit, offset any) *MockIBridgeListWithAddressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWithAddress", reflect.TypeOf((*MockIBridge)(nil).ListWithAddress), ctx, limit, offset)
	return &MockIBridgeListWithAddressCall{Call: call}
}

// MockIBridgeListWithAddressCall wrap *gomock.Call
type MockIBridgeListWithAddressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeListWithAddressCall) Return(arg0 []storage.Bridge, arg1 error) *MockIBridgeListWithAddressCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeListWithAddressCall) Do(f func(context.Context, int, int) ([]storage.Bridge, error)) *MockIBridgeListWithAddressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeListWithAddressCall) DoAndReturn(f func(context.Context, int, int) ([]storage.Bridge, error)) *MockIBridgeListWithAddressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIBridge) Save(ctx context.Context, m *storage.Bridge) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIBridgeMockRecorder) Save(ctx, m any) *MockIBridgeSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIBridge)(nil).Save), ctx, m)
	return &MockIBridgeSaveCall{Call: call}
}

// MockIBridgeSaveCall wrap *gomock.Call
type MockIBridgeSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeSaveCall) Return(arg0 error) *MockIBridgeSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeSaveCall) Do(f func(context.Context, *storage.Bridge) error) *MockIBridgeSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeSaveCall) DoAndReturn(f func(context.Context, *storage.Bridge) error) *MockIBridgeSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIBridge) Update(ctx context.Context, m *storage.Bridge) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIBridgeMockRecorder) Update(ctx, m any) *MockIBridgeUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIBridge)(nil).Update), ctx, m)
	return &MockIBridgeUpdateCall{Call: call}
}

// MockIBridgeUpdateCall wrap *gomock.Call
type MockIBridgeUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBridgeUpdateCall) Return(arg0 error) *MockIBridgeUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBridgeUpdateCall) Do(f func(context.Context, *storage.Bridge) error) *MockIBridgeUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBridgeUpdateCall) DoAndReturn(f func(context.Context, *storage.Bridge) error) *MockIBridgeUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
