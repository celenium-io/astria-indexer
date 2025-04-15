// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: price.go
//
// Generated by this command:
//
//	mockgen -source=price.go -destination=mock/price.go -package=mock -typed
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

// MockIPrice is a mock of IPrice interface.
type MockIPrice struct {
	ctrl     *gomock.Controller
	recorder *MockIPriceMockRecorder
}

// MockIPriceMockRecorder is the mock recorder for MockIPrice.
type MockIPriceMockRecorder struct {
	mock *MockIPrice
}

// NewMockIPrice creates a new mock instance.
func NewMockIPrice(ctrl *gomock.Controller) *MockIPrice {
	mock := &MockIPrice{ctrl: ctrl}
	mock.recorder = &MockIPriceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPrice) EXPECT() *MockIPriceMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockIPrice) All(ctx context.Context, limit, offset int) ([]storage.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", ctx, limit, offset)
	ret0, _ := ret[0].([]storage.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockIPriceMockRecorder) All(ctx, limit, offset any) *MockIPriceAllCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockIPrice)(nil).All), ctx, limit, offset)
	return &MockIPriceAllCall{Call: call}
}

// MockIPriceAllCall wrap *gomock.Call
type MockIPriceAllCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceAllCall) Return(arg0 []storage.Price, arg1 error) *MockIPriceAllCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceAllCall) Do(f func(context.Context, int, int) ([]storage.Price, error)) *MockIPriceAllCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceAllCall) DoAndReturn(f func(context.Context, int, int) ([]storage.Price, error)) *MockIPriceAllCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIPrice) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIPriceMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIPriceCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIPrice)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIPriceCursorListCall{Call: call}
}

// MockIPriceCursorListCall wrap *gomock.Call
type MockIPriceCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceCursorListCall) Return(arg0 []*storage.Price, arg1 error) *MockIPriceCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Price, error)) *MockIPriceCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Price, error)) *MockIPriceCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIPrice) GetByID(ctx context.Context, id uint64) (*storage.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIPriceMockRecorder) GetByID(ctx, id any) *MockIPriceGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIPrice)(nil).GetByID), ctx, id)
	return &MockIPriceGetByIDCall{Call: call}
}

// MockIPriceGetByIDCall wrap *gomock.Call
type MockIPriceGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceGetByIDCall) Return(arg0 *storage.Price, arg1 error) *MockIPriceGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceGetByIDCall) Do(f func(context.Context, uint64) (*storage.Price, error)) *MockIPriceGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Price, error)) *MockIPriceGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIPrice) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIPriceMockRecorder) IsNoRows(err any) *MockIPriceIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIPrice)(nil).IsNoRows), err)
	return &MockIPriceIsNoRowsCall{Call: call}
}

// MockIPriceIsNoRowsCall wrap *gomock.Call
type MockIPriceIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceIsNoRowsCall) Return(arg0 bool) *MockIPriceIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceIsNoRowsCall) Do(f func(error) bool) *MockIPriceIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIPriceIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Last mocks base method.
func (m *MockIPrice) Last(ctx context.Context, currencyPair string) (storage.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Last", ctx, currencyPair)
	ret0, _ := ret[0].(storage.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Last indicates an expected call of Last.
func (mr *MockIPriceMockRecorder) Last(ctx, currencyPair any) *MockIPriceLastCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Last", reflect.TypeOf((*MockIPrice)(nil).Last), ctx, currencyPair)
	return &MockIPriceLastCall{Call: call}
}

// MockIPriceLastCall wrap *gomock.Call
type MockIPriceLastCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceLastCall) Return(arg0 storage.Price, arg1 error) *MockIPriceLastCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceLastCall) Do(f func(context.Context, string) (storage.Price, error)) *MockIPriceLastCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceLastCall) DoAndReturn(f func(context.Context, string) (storage.Price, error)) *MockIPriceLastCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIPrice) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIPriceMockRecorder) LastID(ctx any) *MockIPriceLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIPrice)(nil).LastID), ctx)
	return &MockIPriceLastIDCall{Call: call}
}

// MockIPriceLastIDCall wrap *gomock.Call
type MockIPriceLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceLastIDCall) Return(arg0 uint64, arg1 error) *MockIPriceLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIPriceLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIPriceLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIPrice) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIPriceMockRecorder) List(ctx, limit, offset, order any) *MockIPriceListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIPrice)(nil).List), ctx, limit, offset, order)
	return &MockIPriceListCall{Call: call}
}

// MockIPriceListCall wrap *gomock.Call
type MockIPriceListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceListCall) Return(arg0 []*storage.Price, arg1 error) *MockIPriceListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Price, error)) *MockIPriceListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Price, error)) *MockIPriceListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIPrice) Save(ctx context.Context, m *storage.Price) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIPriceMockRecorder) Save(ctx, m any) *MockIPriceSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIPrice)(nil).Save), ctx, m)
	return &MockIPriceSaveCall{Call: call}
}

// MockIPriceSaveCall wrap *gomock.Call
type MockIPriceSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceSaveCall) Return(arg0 error) *MockIPriceSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceSaveCall) Do(f func(context.Context, *storage.Price) error) *MockIPriceSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceSaveCall) DoAndReturn(f func(context.Context, *storage.Price) error) *MockIPriceSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Series mocks base method.
func (m *MockIPrice) Series(ctx context.Context, currencyPair string, timeframe storage.Timeframe) ([]storage.Candle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Series", ctx, currencyPair, timeframe)
	ret0, _ := ret[0].([]storage.Candle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Series indicates an expected call of Series.
func (mr *MockIPriceMockRecorder) Series(ctx, currencyPair, timeframe any) *MockIPriceSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Series", reflect.TypeOf((*MockIPrice)(nil).Series), ctx, currencyPair, timeframe)
	return &MockIPriceSeriesCall{Call: call}
}

// MockIPriceSeriesCall wrap *gomock.Call
type MockIPriceSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceSeriesCall) Return(arg0 []storage.Candle, arg1 error) *MockIPriceSeriesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceSeriesCall) Do(f func(context.Context, string, storage.Timeframe) ([]storage.Candle, error)) *MockIPriceSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceSeriesCall) DoAndReturn(f func(context.Context, string, storage.Timeframe) ([]storage.Candle, error)) *MockIPriceSeriesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIPrice) Update(ctx context.Context, m *storage.Price) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIPriceMockRecorder) Update(ctx, m any) *MockIPriceUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIPrice)(nil).Update), ctx, m)
	return &MockIPriceUpdateCall{Call: call}
}

// MockIPriceUpdateCall wrap *gomock.Call
type MockIPriceUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIPriceUpdateCall) Return(arg0 error) *MockIPriceUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIPriceUpdateCall) Do(f func(context.Context, *storage.Price) error) *MockIPriceUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIPriceUpdateCall) DoAndReturn(f func(context.Context, *storage.Price) error) *MockIPriceUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
