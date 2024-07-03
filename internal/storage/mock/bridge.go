// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
func (mr *MockIBridgeMockRecorder) ByAddress(ctx, addressId any) *IBridgeByAddressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByAddress", reflect.TypeOf((*MockIBridge)(nil).ByAddress), ctx, addressId)
	return &IBridgeByAddressCall{Call: call}
}

// IBridgeByAddressCall wrap *gomock.Call
type IBridgeByAddressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeByAddressCall) Return(arg0 storage.Bridge, arg1 error) *IBridgeByAddressCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeByAddressCall) Do(f func(context.Context, uint64) (storage.Bridge, error)) *IBridgeByAddressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeByAddressCall) DoAndReturn(f func(context.Context, uint64) (storage.Bridge, error)) *IBridgeByAddressCall {
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
func (mr *MockIBridgeMockRecorder) ByRoles(ctx, addressId, limit, offset any) *IBridgeByRolesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByRoles", reflect.TypeOf((*MockIBridge)(nil).ByRoles), ctx, addressId, limit, offset)
	return &IBridgeByRolesCall{Call: call}
}

// IBridgeByRolesCall wrap *gomock.Call
type IBridgeByRolesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeByRolesCall) Return(arg0 []storage.Bridge, arg1 error) *IBridgeByRolesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeByRolesCall) Do(f func(context.Context, uint64, int, int) ([]storage.Bridge, error)) *IBridgeByRolesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeByRolesCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Bridge, error)) *IBridgeByRolesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByRollup mocks base method.
func (m *MockIBridge) ByRollup(ctx context.Context, rollupId uint64) (storage.Bridge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByRollup", ctx, rollupId)
	ret0, _ := ret[0].(storage.Bridge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByRollup indicates an expected call of ByRollup.
func (mr *MockIBridgeMockRecorder) ByRollup(ctx, rollupId any) *IBridgeByRollupCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByRollup", reflect.TypeOf((*MockIBridge)(nil).ByRollup), ctx, rollupId)
	return &IBridgeByRollupCall{Call: call}
}

// IBridgeByRollupCall wrap *gomock.Call
type IBridgeByRollupCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeByRollupCall) Return(arg0 storage.Bridge, arg1 error) *IBridgeByRollupCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeByRollupCall) Do(f func(context.Context, uint64) (storage.Bridge, error)) *IBridgeByRollupCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeByRollupCall) DoAndReturn(f func(context.Context, uint64) (storage.Bridge, error)) *IBridgeByRollupCall {
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
func (mr *MockIBridgeMockRecorder) CursorList(ctx, id, limit, order, cmp any) *IBridgeCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIBridge)(nil).CursorList), ctx, id, limit, order, cmp)
	return &IBridgeCursorListCall{Call: call}
}

// IBridgeCursorListCall wrap *gomock.Call
type IBridgeCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeCursorListCall) Return(arg0 []*storage.Bridge, arg1 error) *IBridgeCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Bridge, error)) *IBridgeCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Bridge, error)) *IBridgeCursorListCall {
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
func (mr *MockIBridgeMockRecorder) GetByID(ctx, id any) *IBridgeGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIBridge)(nil).GetByID), ctx, id)
	return &IBridgeGetByIDCall{Call: call}
}

// IBridgeGetByIDCall wrap *gomock.Call
type IBridgeGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeGetByIDCall) Return(arg0 *storage.Bridge, arg1 error) *IBridgeGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeGetByIDCall) Do(f func(context.Context, uint64) (*storage.Bridge, error)) *IBridgeGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Bridge, error)) *IBridgeGetByIDCall {
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
func (mr *MockIBridgeMockRecorder) IsNoRows(err any) *IBridgeIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIBridge)(nil).IsNoRows), err)
	return &IBridgeIsNoRowsCall{Call: call}
}

// IBridgeIsNoRowsCall wrap *gomock.Call
type IBridgeIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeIsNoRowsCall) Return(arg0 bool) *IBridgeIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeIsNoRowsCall) Do(f func(error) bool) *IBridgeIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeIsNoRowsCall) DoAndReturn(f func(error) bool) *IBridgeIsNoRowsCall {
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
func (mr *MockIBridgeMockRecorder) LastID(ctx any) *IBridgeLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIBridge)(nil).LastID), ctx)
	return &IBridgeLastIDCall{Call: call}
}

// IBridgeLastIDCall wrap *gomock.Call
type IBridgeLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeLastIDCall) Return(arg0 uint64, arg1 error) *IBridgeLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeLastIDCall) Do(f func(context.Context) (uint64, error)) *IBridgeLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *IBridgeLastIDCall {
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
func (mr *MockIBridgeMockRecorder) List(ctx, limit, offset, order any) *IBridgeListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIBridge)(nil).List), ctx, limit, offset, order)
	return &IBridgeListCall{Call: call}
}

// IBridgeListCall wrap *gomock.Call
type IBridgeListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeListCall) Return(arg0 []*storage.Bridge, arg1 error) *IBridgeListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Bridge, error)) *IBridgeListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Bridge, error)) *IBridgeListCall {
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
func (mr *MockIBridgeMockRecorder) Save(ctx, m any) *IBridgeSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIBridge)(nil).Save), ctx, m)
	return &IBridgeSaveCall{Call: call}
}

// IBridgeSaveCall wrap *gomock.Call
type IBridgeSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeSaveCall) Return(arg0 error) *IBridgeSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeSaveCall) Do(f func(context.Context, *storage.Bridge) error) *IBridgeSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeSaveCall) DoAndReturn(f func(context.Context, *storage.Bridge) error) *IBridgeSaveCall {
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
func (mr *MockIBridgeMockRecorder) Update(ctx, m any) *IBridgeUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIBridge)(nil).Update), ctx, m)
	return &IBridgeUpdateCall{Call: call}
}

// IBridgeUpdateCall wrap *gomock.Call
type IBridgeUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *IBridgeUpdateCall) Return(arg0 error) *IBridgeUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *IBridgeUpdateCall) Do(f func(context.Context, *storage.Bridge) error) *IBridgeUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *IBridgeUpdateCall) DoAndReturn(f func(context.Context, *storage.Bridge) error) *IBridgeUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
