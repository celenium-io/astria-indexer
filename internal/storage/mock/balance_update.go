// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: balance_update.go
//
// Generated by this command:
//
//	mockgen -source=balance_update.go -destination=mock/balance_update.go -package=mock -typed
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

// MockIBalanceUpdate is a mock of IBalanceUpdate interface.
type MockIBalanceUpdate struct {
	ctrl     *gomock.Controller
	recorder *MockIBalanceUpdateMockRecorder
}

// MockIBalanceUpdateMockRecorder is the mock recorder for MockIBalanceUpdate.
type MockIBalanceUpdateMockRecorder struct {
	mock *MockIBalanceUpdate
}

// NewMockIBalanceUpdate creates a new mock instance.
func NewMockIBalanceUpdate(ctrl *gomock.Controller) *MockIBalanceUpdate {
	mock := &MockIBalanceUpdate{ctrl: ctrl}
	mock.recorder = &MockIBalanceUpdateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBalanceUpdate) EXPECT() *MockIBalanceUpdateMockRecorder {
	return m.recorder
}

// CursorList mocks base method.
func (m *MockIBalanceUpdate) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.BalanceUpdate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.BalanceUpdate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIBalanceUpdateMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIBalanceUpdateCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIBalanceUpdate)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIBalanceUpdateCursorListCall{Call: call}
}

// MockIBalanceUpdateCursorListCall wrap *gomock.Call
type MockIBalanceUpdateCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBalanceUpdateCursorListCall) Return(arg0 []*storage.BalanceUpdate, arg1 error) *MockIBalanceUpdateCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBalanceUpdateCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.BalanceUpdate, error)) *MockIBalanceUpdateCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBalanceUpdateCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.BalanceUpdate, error)) *MockIBalanceUpdateCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIBalanceUpdate) GetByID(ctx context.Context, id uint64) (*storage.BalanceUpdate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.BalanceUpdate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIBalanceUpdateMockRecorder) GetByID(ctx, id any) *MockIBalanceUpdateGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIBalanceUpdate)(nil).GetByID), ctx, id)
	return &MockIBalanceUpdateGetByIDCall{Call: call}
}

// MockIBalanceUpdateGetByIDCall wrap *gomock.Call
type MockIBalanceUpdateGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBalanceUpdateGetByIDCall) Return(arg0 *storage.BalanceUpdate, arg1 error) *MockIBalanceUpdateGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBalanceUpdateGetByIDCall) Do(f func(context.Context, uint64) (*storage.BalanceUpdate, error)) *MockIBalanceUpdateGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBalanceUpdateGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.BalanceUpdate, error)) *MockIBalanceUpdateGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIBalanceUpdate) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIBalanceUpdateMockRecorder) IsNoRows(err any) *MockIBalanceUpdateIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIBalanceUpdate)(nil).IsNoRows), err)
	return &MockIBalanceUpdateIsNoRowsCall{Call: call}
}

// MockIBalanceUpdateIsNoRowsCall wrap *gomock.Call
type MockIBalanceUpdateIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBalanceUpdateIsNoRowsCall) Return(arg0 bool) *MockIBalanceUpdateIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBalanceUpdateIsNoRowsCall) Do(f func(error) bool) *MockIBalanceUpdateIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBalanceUpdateIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIBalanceUpdateIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIBalanceUpdate) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIBalanceUpdateMockRecorder) LastID(ctx any) *MockIBalanceUpdateLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIBalanceUpdate)(nil).LastID), ctx)
	return &MockIBalanceUpdateLastIDCall{Call: call}
}

// MockIBalanceUpdateLastIDCall wrap *gomock.Call
type MockIBalanceUpdateLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBalanceUpdateLastIDCall) Return(arg0 uint64, arg1 error) *MockIBalanceUpdateLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBalanceUpdateLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIBalanceUpdateLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBalanceUpdateLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIBalanceUpdateLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIBalanceUpdate) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.BalanceUpdate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.BalanceUpdate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIBalanceUpdateMockRecorder) List(ctx, limit, offset, order any) *MockIBalanceUpdateListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIBalanceUpdate)(nil).List), ctx, limit, offset, order)
	return &MockIBalanceUpdateListCall{Call: call}
}

// MockIBalanceUpdateListCall wrap *gomock.Call
type MockIBalanceUpdateListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBalanceUpdateListCall) Return(arg0 []*storage.BalanceUpdate, arg1 error) *MockIBalanceUpdateListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBalanceUpdateListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.BalanceUpdate, error)) *MockIBalanceUpdateListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBalanceUpdateListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.BalanceUpdate, error)) *MockIBalanceUpdateListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIBalanceUpdate) Save(ctx context.Context, m *storage.BalanceUpdate) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIBalanceUpdateMockRecorder) Save(ctx, m any) *MockIBalanceUpdateSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIBalanceUpdate)(nil).Save), ctx, m)
	return &MockIBalanceUpdateSaveCall{Call: call}
}

// MockIBalanceUpdateSaveCall wrap *gomock.Call
type MockIBalanceUpdateSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBalanceUpdateSaveCall) Return(arg0 error) *MockIBalanceUpdateSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBalanceUpdateSaveCall) Do(f func(context.Context, *storage.BalanceUpdate) error) *MockIBalanceUpdateSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBalanceUpdateSaveCall) DoAndReturn(f func(context.Context, *storage.BalanceUpdate) error) *MockIBalanceUpdateSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIBalanceUpdate) Update(ctx context.Context, m *storage.BalanceUpdate) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIBalanceUpdateMockRecorder) Update(ctx, m any) *MockIBalanceUpdateUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIBalanceUpdate)(nil).Update), ctx, m)
	return &MockIBalanceUpdateUpdateCall{Call: call}
}

// MockIBalanceUpdateUpdateCall wrap *gomock.Call
type MockIBalanceUpdateUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBalanceUpdateUpdateCall) Return(arg0 error) *MockIBalanceUpdateUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBalanceUpdateUpdateCall) Do(f func(context.Context, *storage.BalanceUpdate) error) *MockIBalanceUpdateUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBalanceUpdateUpdateCall) DoAndReturn(f func(context.Context, *storage.BalanceUpdate) error) *MockIBalanceUpdateUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
