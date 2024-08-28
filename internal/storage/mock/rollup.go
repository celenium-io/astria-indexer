// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: rollup.go
//
// Generated by this command:
//
//	mockgen -source=rollup.go -destination=mock/rollup.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/celenium-io/astria-indexer/internal/storage"
	types "github.com/celenium-io/astria-indexer/pkg/types"
	storage0 "github.com/dipdup-net/indexer-sdk/pkg/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockIRollup is a mock of IRollup interface.
type MockIRollup struct {
	ctrl     *gomock.Controller
	recorder *MockIRollupMockRecorder
}

// MockIRollupMockRecorder is the mock recorder for MockIRollup.
type MockIRollupMockRecorder struct {
	mock *MockIRollup
}

// NewMockIRollup creates a new mock instance.
func NewMockIRollup(ctrl *gomock.Controller) *MockIRollup {
	mock := &MockIRollup{ctrl: ctrl}
	mock.recorder = &MockIRollupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRollup) EXPECT() *MockIRollupMockRecorder {
	return m.recorder
}

// ActionsByHeight mocks base method.
func (m *MockIRollup) ActionsByHeight(ctx context.Context, height types.Level, limit, offset int) ([]storage.RollupAction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActionsByHeight", ctx, height, limit, offset)
	ret0, _ := ret[0].([]storage.RollupAction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ActionsByHeight indicates an expected call of ActionsByHeight.
func (mr *MockIRollupMockRecorder) ActionsByHeight(ctx, height, limit, offset any) *MockIRollupActionsByHeightCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActionsByHeight", reflect.TypeOf((*MockIRollup)(nil).ActionsByHeight), ctx, height, limit, offset)
	return &MockIRollupActionsByHeightCall{Call: call}
}

// MockIRollupActionsByHeightCall wrap *gomock.Call
type MockIRollupActionsByHeightCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupActionsByHeightCall) Return(arg0 []storage.RollupAction, arg1 error) *MockIRollupActionsByHeightCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupActionsByHeightCall) Do(f func(context.Context, types.Level, int, int) ([]storage.RollupAction, error)) *MockIRollupActionsByHeightCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupActionsByHeightCall) DoAndReturn(f func(context.Context, types.Level, int, int) ([]storage.RollupAction, error)) *MockIRollupActionsByHeightCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ActionsByTxId mocks base method.
func (m *MockIRollup) ActionsByTxId(ctx context.Context, txId uint64, limit, offset int) ([]storage.RollupAction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActionsByTxId", ctx, txId, limit, offset)
	ret0, _ := ret[0].([]storage.RollupAction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ActionsByTxId indicates an expected call of ActionsByTxId.
func (mr *MockIRollupMockRecorder) ActionsByTxId(ctx, txId, limit, offset any) *MockIRollupActionsByTxIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActionsByTxId", reflect.TypeOf((*MockIRollup)(nil).ActionsByTxId), ctx, txId, limit, offset)
	return &MockIRollupActionsByTxIdCall{Call: call}
}

// MockIRollupActionsByTxIdCall wrap *gomock.Call
type MockIRollupActionsByTxIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupActionsByTxIdCall) Return(arg0 []storage.RollupAction, arg1 error) *MockIRollupActionsByTxIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupActionsByTxIdCall) Do(f func(context.Context, uint64, int, int) ([]storage.RollupAction, error)) *MockIRollupActionsByTxIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupActionsByTxIdCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.RollupAction, error)) *MockIRollupActionsByTxIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Addresses mocks base method.
func (m *MockIRollup) Addresses(ctx context.Context, rollupId uint64, limit, offset int, sort storage0.SortOrder) ([]storage.RollupAddress, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Addresses", ctx, rollupId, limit, offset, sort)
	ret0, _ := ret[0].([]storage.RollupAddress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Addresses indicates an expected call of Addresses.
func (mr *MockIRollupMockRecorder) Addresses(ctx, rollupId, limit, offset, sort any) *MockIRollupAddressesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Addresses", reflect.TypeOf((*MockIRollup)(nil).Addresses), ctx, rollupId, limit, offset, sort)
	return &MockIRollupAddressesCall{Call: call}
}

// MockIRollupAddressesCall wrap *gomock.Call
type MockIRollupAddressesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupAddressesCall) Return(arg0 []storage.RollupAddress, arg1 error) *MockIRollupAddressesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupAddressesCall) Do(f func(context.Context, uint64, int, int, storage0.SortOrder) ([]storage.RollupAddress, error)) *MockIRollupAddressesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupAddressesCall) DoAndReturn(f func(context.Context, uint64, int, int, storage0.SortOrder) ([]storage.RollupAddress, error)) *MockIRollupAddressesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByHash mocks base method.
func (m *MockIRollup) ByHash(ctx context.Context, hash []byte) (storage.Rollup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByHash", ctx, hash)
	ret0, _ := ret[0].(storage.Rollup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByHash indicates an expected call of ByHash.
func (mr *MockIRollupMockRecorder) ByHash(ctx, hash any) *MockIRollupByHashCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByHash", reflect.TypeOf((*MockIRollup)(nil).ByHash), ctx, hash)
	return &MockIRollupByHashCall{Call: call}
}

// MockIRollupByHashCall wrap *gomock.Call
type MockIRollupByHashCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupByHashCall) Return(arg0 storage.Rollup, arg1 error) *MockIRollupByHashCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupByHashCall) Do(f func(context.Context, []byte) (storage.Rollup, error)) *MockIRollupByHashCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupByHashCall) DoAndReturn(f func(context.Context, []byte) (storage.Rollup, error)) *MockIRollupByHashCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CountActionsByHeight mocks base method.
func (m *MockIRollup) CountActionsByHeight(ctx context.Context, height types.Level) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountActionsByHeight", ctx, height)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountActionsByHeight indicates an expected call of CountActionsByHeight.
func (mr *MockIRollupMockRecorder) CountActionsByHeight(ctx, height any) *MockIRollupCountActionsByHeightCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountActionsByHeight", reflect.TypeOf((*MockIRollup)(nil).CountActionsByHeight), ctx, height)
	return &MockIRollupCountActionsByHeightCall{Call: call}
}

// MockIRollupCountActionsByHeightCall wrap *gomock.Call
type MockIRollupCountActionsByHeightCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupCountActionsByHeightCall) Return(arg0 int64, arg1 error) *MockIRollupCountActionsByHeightCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupCountActionsByHeightCall) Do(f func(context.Context, types.Level) (int64, error)) *MockIRollupCountActionsByHeightCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupCountActionsByHeightCall) DoAndReturn(f func(context.Context, types.Level) (int64, error)) *MockIRollupCountActionsByHeightCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CountActionsByTxId mocks base method.
func (m *MockIRollup) CountActionsByTxId(ctx context.Context, txId uint64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountActionsByTxId", ctx, txId)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountActionsByTxId indicates an expected call of CountActionsByTxId.
func (mr *MockIRollupMockRecorder) CountActionsByTxId(ctx, txId any) *MockIRollupCountActionsByTxIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountActionsByTxId", reflect.TypeOf((*MockIRollup)(nil).CountActionsByTxId), ctx, txId)
	return &MockIRollupCountActionsByTxIdCall{Call: call}
}

// MockIRollupCountActionsByTxIdCall wrap *gomock.Call
type MockIRollupCountActionsByTxIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupCountActionsByTxIdCall) Return(arg0 int64, arg1 error) *MockIRollupCountActionsByTxIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupCountActionsByTxIdCall) Do(f func(context.Context, uint64) (int64, error)) *MockIRollupCountActionsByTxIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupCountActionsByTxIdCall) DoAndReturn(f func(context.Context, uint64) (int64, error)) *MockIRollupCountActionsByTxIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIRollup) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Rollup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Rollup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIRollupMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIRollupCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIRollup)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIRollupCursorListCall{Call: call}
}

// MockIRollupCursorListCall wrap *gomock.Call
type MockIRollupCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupCursorListCall) Return(arg0 []*storage.Rollup, arg1 error) *MockIRollupCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Rollup, error)) *MockIRollupCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Rollup, error)) *MockIRollupCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIRollup) GetByID(ctx context.Context, id uint64) (*storage.Rollup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Rollup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIRollupMockRecorder) GetByID(ctx, id any) *MockIRollupGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIRollup)(nil).GetByID), ctx, id)
	return &MockIRollupGetByIDCall{Call: call}
}

// MockIRollupGetByIDCall wrap *gomock.Call
type MockIRollupGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupGetByIDCall) Return(arg0 *storage.Rollup, arg1 error) *MockIRollupGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupGetByIDCall) Do(f func(context.Context, uint64) (*storage.Rollup, error)) *MockIRollupGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Rollup, error)) *MockIRollupGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIRollup) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIRollupMockRecorder) IsNoRows(err any) *MockIRollupIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIRollup)(nil).IsNoRows), err)
	return &MockIRollupIsNoRowsCall{Call: call}
}

// MockIRollupIsNoRowsCall wrap *gomock.Call
type MockIRollupIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupIsNoRowsCall) Return(arg0 bool) *MockIRollupIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupIsNoRowsCall) Do(f func(error) bool) *MockIRollupIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIRollupIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIRollup) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIRollupMockRecorder) LastID(ctx any) *MockIRollupLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIRollup)(nil).LastID), ctx)
	return &MockIRollupLastIDCall{Call: call}
}

// MockIRollupLastIDCall wrap *gomock.Call
type MockIRollupLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupLastIDCall) Return(arg0 uint64, arg1 error) *MockIRollupLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIRollupLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIRollupLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIRollup) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Rollup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Rollup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIRollupMockRecorder) List(ctx, limit, offset, order any) *MockIRollupListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIRollup)(nil).List), ctx, limit, offset, order)
	return &MockIRollupListCall{Call: call}
}

// MockIRollupListCall wrap *gomock.Call
type MockIRollupListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupListCall) Return(arg0 []*storage.Rollup, arg1 error) *MockIRollupListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Rollup, error)) *MockIRollupListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Rollup, error)) *MockIRollupListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListExt mocks base method.
func (m *MockIRollup) ListExt(ctx context.Context, fltrs storage.RollupListFilter) ([]storage.Rollup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListExt", ctx, fltrs)
	ret0, _ := ret[0].([]storage.Rollup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListExt indicates an expected call of ListExt.
func (mr *MockIRollupMockRecorder) ListExt(ctx, fltrs any) *MockIRollupListExtCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListExt", reflect.TypeOf((*MockIRollup)(nil).ListExt), ctx, fltrs)
	return &MockIRollupListExtCall{Call: call}
}

// MockIRollupListExtCall wrap *gomock.Call
type MockIRollupListExtCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupListExtCall) Return(arg0 []storage.Rollup, arg1 error) *MockIRollupListExtCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupListExtCall) Do(f func(context.Context, storage.RollupListFilter) ([]storage.Rollup, error)) *MockIRollupListExtCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupListExtCall) DoAndReturn(f func(context.Context, storage.RollupListFilter) ([]storage.Rollup, error)) *MockIRollupListExtCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListRollupsByAddress mocks base method.
func (m *MockIRollup) ListRollupsByAddress(ctx context.Context, addressId uint64, limit, offset int, sort storage0.SortOrder) ([]storage.RollupAddress, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRollupsByAddress", ctx, addressId, limit, offset, sort)
	ret0, _ := ret[0].([]storage.RollupAddress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRollupsByAddress indicates an expected call of ListRollupsByAddress.
func (mr *MockIRollupMockRecorder) ListRollupsByAddress(ctx, addressId, limit, offset, sort any) *MockIRollupListRollupsByAddressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRollupsByAddress", reflect.TypeOf((*MockIRollup)(nil).ListRollupsByAddress), ctx, addressId, limit, offset, sort)
	return &MockIRollupListRollupsByAddressCall{Call: call}
}

// MockIRollupListRollupsByAddressCall wrap *gomock.Call
type MockIRollupListRollupsByAddressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupListRollupsByAddressCall) Return(arg0 []storage.RollupAddress, arg1 error) *MockIRollupListRollupsByAddressCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupListRollupsByAddressCall) Do(f func(context.Context, uint64, int, int, storage0.SortOrder) ([]storage.RollupAddress, error)) *MockIRollupListRollupsByAddressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupListRollupsByAddressCall) DoAndReturn(f func(context.Context, uint64, int, int, storage0.SortOrder) ([]storage.RollupAddress, error)) *MockIRollupListRollupsByAddressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIRollup) Save(ctx context.Context, m *storage.Rollup) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIRollupMockRecorder) Save(ctx, m any) *MockIRollupSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIRollup)(nil).Save), ctx, m)
	return &MockIRollupSaveCall{Call: call}
}

// MockIRollupSaveCall wrap *gomock.Call
type MockIRollupSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupSaveCall) Return(arg0 error) *MockIRollupSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupSaveCall) Do(f func(context.Context, *storage.Rollup) error) *MockIRollupSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupSaveCall) DoAndReturn(f func(context.Context, *storage.Rollup) error) *MockIRollupSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIRollup) Update(ctx context.Context, m *storage.Rollup) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIRollupMockRecorder) Update(ctx, m any) *MockIRollupUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIRollup)(nil).Update), ctx, m)
	return &MockIRollupUpdateCall{Call: call}
}

// MockIRollupUpdateCall wrap *gomock.Call
type MockIRollupUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIRollupUpdateCall) Return(arg0 error) *MockIRollupUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIRollupUpdateCall) Do(f func(context.Context, *storage.Rollup) error) *MockIRollupUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIRollupUpdateCall) DoAndReturn(f func(context.Context, *storage.Rollup) error) *MockIRollupUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
