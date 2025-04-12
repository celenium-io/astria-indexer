// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: fee.go
//
// Generated by this command:
//
//	mockgen -source=fee.go -destination=mock/fee.go -package=mock -typed
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

// MockIFee is a mock of IFee interface.
type MockIFee struct {
	ctrl     *gomock.Controller
	recorder *MockIFeeMockRecorder
}

// MockIFeeMockRecorder is the mock recorder for MockIFee.
type MockIFeeMockRecorder struct {
	mock *MockIFee
}

// NewMockIFee creates a new mock instance.
func NewMockIFee(ctrl *gomock.Controller) *MockIFee {
	mock := &MockIFee{ctrl: ctrl}
	mock.recorder = &MockIFeeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIFee) EXPECT() *MockIFeeMockRecorder {
	return m.recorder
}

// ByPayerId mocks base method.
func (m *MockIFee) ByPayerId(ctx context.Context, id uint64, limit, offset int, sort storage0.SortOrder) ([]storage.Fee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByPayerId", ctx, id, limit, offset, sort)
	ret0, _ := ret[0].([]storage.Fee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByPayerId indicates an expected call of ByPayerId.
func (mr *MockIFeeMockRecorder) ByPayerId(ctx, id, limit, offset, sort any) *MockIFeeByPayerIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByPayerId", reflect.TypeOf((*MockIFee)(nil).ByPayerId), ctx, id, limit, offset, sort)
	return &MockIFeeByPayerIdCall{Call: call}
}

// MockIFeeByPayerIdCall wrap *gomock.Call
type MockIFeeByPayerIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeByPayerIdCall) Return(arg0 []storage.Fee, arg1 error) *MockIFeeByPayerIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeByPayerIdCall) Do(f func(context.Context, uint64, int, int, storage0.SortOrder) ([]storage.Fee, error)) *MockIFeeByPayerIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeByPayerIdCall) DoAndReturn(f func(context.Context, uint64, int, int, storage0.SortOrder) ([]storage.Fee, error)) *MockIFeeByPayerIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ByTxId mocks base method.
func (m *MockIFee) ByTxId(ctx context.Context, id uint64, limit, offset int) ([]storage.Fee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByTxId", ctx, id, limit, offset)
	ret0, _ := ret[0].([]storage.Fee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByTxId indicates an expected call of ByTxId.
func (mr *MockIFeeMockRecorder) ByTxId(ctx, id, limit, offset any) *MockIFeeByTxIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByTxId", reflect.TypeOf((*MockIFee)(nil).ByTxId), ctx, id, limit, offset)
	return &MockIFeeByTxIdCall{Call: call}
}

// MockIFeeByTxIdCall wrap *gomock.Call
type MockIFeeByTxIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeByTxIdCall) Return(arg0 []storage.Fee, arg1 error) *MockIFeeByTxIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeByTxIdCall) Do(f func(context.Context, uint64, int, int) ([]storage.Fee, error)) *MockIFeeByTxIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeByTxIdCall) DoAndReturn(f func(context.Context, uint64, int, int) ([]storage.Fee, error)) *MockIFeeByTxIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIFee) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Fee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Fee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIFeeMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIFeeCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIFee)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIFeeCursorListCall{Call: call}
}

// MockIFeeCursorListCall wrap *gomock.Call
type MockIFeeCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeCursorListCall) Return(arg0 []*storage.Fee, arg1 error) *MockIFeeCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Fee, error)) *MockIFeeCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Fee, error)) *MockIFeeCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FullTxFee mocks base method.
func (m *MockIFee) FullTxFee(ctx context.Context, id uint64) ([]storage.Fee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullTxFee", ctx, id)
	ret0, _ := ret[0].([]storage.Fee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FullTxFee indicates an expected call of FullTxFee.
func (mr *MockIFeeMockRecorder) FullTxFee(ctx, id any) *MockIFeeFullTxFeeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullTxFee", reflect.TypeOf((*MockIFee)(nil).FullTxFee), ctx, id)
	return &MockIFeeFullTxFeeCall{Call: call}
}

// MockIFeeFullTxFeeCall wrap *gomock.Call
type MockIFeeFullTxFeeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeFullTxFeeCall) Return(arg0 []storage.Fee, arg1 error) *MockIFeeFullTxFeeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeFullTxFeeCall) Do(f func(context.Context, uint64) ([]storage.Fee, error)) *MockIFeeFullTxFeeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeFullTxFeeCall) DoAndReturn(f func(context.Context, uint64) ([]storage.Fee, error)) *MockIFeeFullTxFeeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIFee) GetByID(ctx context.Context, id uint64) (*storage.Fee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Fee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIFeeMockRecorder) GetByID(ctx, id any) *MockIFeeGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIFee)(nil).GetByID), ctx, id)
	return &MockIFeeGetByIDCall{Call: call}
}

// MockIFeeGetByIDCall wrap *gomock.Call
type MockIFeeGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeGetByIDCall) Return(arg0 *storage.Fee, arg1 error) *MockIFeeGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeGetByIDCall) Do(f func(context.Context, uint64) (*storage.Fee, error)) *MockIFeeGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Fee, error)) *MockIFeeGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIFee) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIFeeMockRecorder) IsNoRows(err any) *MockIFeeIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIFee)(nil).IsNoRows), err)
	return &MockIFeeIsNoRowsCall{Call: call}
}

// MockIFeeIsNoRowsCall wrap *gomock.Call
type MockIFeeIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeIsNoRowsCall) Return(arg0 bool) *MockIFeeIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeIsNoRowsCall) Do(f func(error) bool) *MockIFeeIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIFeeIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIFee) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIFeeMockRecorder) LastID(ctx any) *MockIFeeLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIFee)(nil).LastID), ctx)
	return &MockIFeeLastIDCall{Call: call}
}

// MockIFeeLastIDCall wrap *gomock.Call
type MockIFeeLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeLastIDCall) Return(arg0 uint64, arg1 error) *MockIFeeLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIFeeLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIFeeLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIFee) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Fee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Fee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIFeeMockRecorder) List(ctx, limit, offset, order any) *MockIFeeListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIFee)(nil).List), ctx, limit, offset, order)
	return &MockIFeeListCall{Call: call}
}

// MockIFeeListCall wrap *gomock.Call
type MockIFeeListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeListCall) Return(arg0 []*storage.Fee, arg1 error) *MockIFeeListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Fee, error)) *MockIFeeListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Fee, error)) *MockIFeeListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIFee) Save(ctx context.Context, m *storage.Fee) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIFeeMockRecorder) Save(ctx, m any) *MockIFeeSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIFee)(nil).Save), ctx, m)
	return &MockIFeeSaveCall{Call: call}
}

// MockIFeeSaveCall wrap *gomock.Call
type MockIFeeSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeSaveCall) Return(arg0 error) *MockIFeeSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeSaveCall) Do(f func(context.Context, *storage.Fee) error) *MockIFeeSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeSaveCall) DoAndReturn(f func(context.Context, *storage.Fee) error) *MockIFeeSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIFee) Update(ctx context.Context, m *storage.Fee) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIFeeMockRecorder) Update(ctx, m any) *MockIFeeUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIFee)(nil).Update), ctx, m)
	return &MockIFeeUpdateCall{Call: call}
}

// MockIFeeUpdateCall wrap *gomock.Call
type MockIFeeUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIFeeUpdateCall) Return(arg0 error) *MockIFeeUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIFeeUpdateCall) Do(f func(context.Context, *storage.Fee) error) *MockIFeeUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIFeeUpdateCall) DoAndReturn(f func(context.Context, *storage.Fee) error) *MockIFeeUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
