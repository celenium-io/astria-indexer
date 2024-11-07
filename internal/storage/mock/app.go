// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: app.go
//
// Generated by this command:
//
//	mockgen -source=app.go -destination=mock/app.go -package=mock -typed
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

// MockIApp is a mock of IApp interface.
type MockIApp struct {
	ctrl     *gomock.Controller
	recorder *MockIAppMockRecorder
}

// MockIAppMockRecorder is the mock recorder for MockIApp.
type MockIAppMockRecorder struct {
	mock *MockIApp
}

// NewMockIApp creates a new mock instance.
func NewMockIApp(ctrl *gomock.Controller) *MockIApp {
	mock := &MockIApp{ctrl: ctrl}
	mock.recorder = &MockIAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIApp) EXPECT() *MockIAppMockRecorder {
	return m.recorder
}

// Actions mocks base method.
func (m *MockIApp) Actions(ctx context.Context, slug string, limit, offset int, sort storage0.SortOrder) ([]storage.RollupAction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Actions", ctx, slug, limit, offset, sort)
	ret0, _ := ret[0].([]storage.RollupAction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Actions indicates an expected call of Actions.
func (mr *MockIAppMockRecorder) Actions(ctx, slug, limit, offset, sort any) *MockIAppActionsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Actions", reflect.TypeOf((*MockIApp)(nil).Actions), ctx, slug, limit, offset, sort)
	return &MockIAppActionsCall{Call: call}
}

// MockIAppActionsCall wrap *gomock.Call
type MockIAppActionsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppActionsCall) Return(arg0 []storage.RollupAction, arg1 error) *MockIAppActionsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppActionsCall) Do(f func(context.Context, string, int, int, storage0.SortOrder) ([]storage.RollupAction, error)) *MockIAppActionsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppActionsCall) DoAndReturn(f func(context.Context, string, int, int, storage0.SortOrder) ([]storage.RollupAction, error)) *MockIAppActionsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BySlug mocks base method.
func (m *MockIApp) BySlug(ctx context.Context, slug string) (storage.AppWithStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BySlug", ctx, slug)
	ret0, _ := ret[0].(storage.AppWithStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BySlug indicates an expected call of BySlug.
func (mr *MockIAppMockRecorder) BySlug(ctx, slug any) *MockIAppBySlugCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BySlug", reflect.TypeOf((*MockIApp)(nil).BySlug), ctx, slug)
	return &MockIAppBySlugCall{Call: call}
}

// MockIAppBySlugCall wrap *gomock.Call
type MockIAppBySlugCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppBySlugCall) Return(arg0 storage.AppWithStats, arg1 error) *MockIAppBySlugCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppBySlugCall) Do(f func(context.Context, string) (storage.AppWithStats, error)) *MockIAppBySlugCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppBySlugCall) DoAndReturn(f func(context.Context, string) (storage.AppWithStats, error)) *MockIAppBySlugCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CursorList mocks base method.
func (m *MockIApp) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIAppMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIAppCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIApp)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIAppCursorListCall{Call: call}
}

// MockIAppCursorListCall wrap *gomock.Call
type MockIAppCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppCursorListCall) Return(arg0 []*storage.App, arg1 error) *MockIAppCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.App, error)) *MockIAppCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.App, error)) *MockIAppCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIApp) GetByID(ctx context.Context, id uint64) (*storage.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIAppMockRecorder) GetByID(ctx, id any) *MockIAppGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIApp)(nil).GetByID), ctx, id)
	return &MockIAppGetByIDCall{Call: call}
}

// MockIAppGetByIDCall wrap *gomock.Call
type MockIAppGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppGetByIDCall) Return(arg0 *storage.App, arg1 error) *MockIAppGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppGetByIDCall) Do(f func(context.Context, uint64) (*storage.App, error)) *MockIAppGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.App, error)) *MockIAppGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIApp) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIAppMockRecorder) IsNoRows(err any) *MockIAppIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIApp)(nil).IsNoRows), err)
	return &MockIAppIsNoRowsCall{Call: call}
}

// MockIAppIsNoRowsCall wrap *gomock.Call
type MockIAppIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppIsNoRowsCall) Return(arg0 bool) *MockIAppIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppIsNoRowsCall) Do(f func(error) bool) *MockIAppIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIAppIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIApp) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIAppMockRecorder) LastID(ctx any) *MockIAppLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIApp)(nil).LastID), ctx)
	return &MockIAppLastIDCall{Call: call}
}

// MockIAppLastIDCall wrap *gomock.Call
type MockIAppLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppLastIDCall) Return(arg0 uint64, arg1 error) *MockIAppLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIAppLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIAppLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Leaderboard mocks base method.
func (m *MockIApp) Leaderboard(ctx context.Context, fltrs storage.LeaderboardFilters) ([]storage.AppWithStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Leaderboard", ctx, fltrs)
	ret0, _ := ret[0].([]storage.AppWithStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Leaderboard indicates an expected call of Leaderboard.
func (mr *MockIAppMockRecorder) Leaderboard(ctx, fltrs any) *MockIAppLeaderboardCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Leaderboard", reflect.TypeOf((*MockIApp)(nil).Leaderboard), ctx, fltrs)
	return &MockIAppLeaderboardCall{Call: call}
}

// MockIAppLeaderboardCall wrap *gomock.Call
type MockIAppLeaderboardCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppLeaderboardCall) Return(arg0 []storage.AppWithStats, arg1 error) *MockIAppLeaderboardCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppLeaderboardCall) Do(f func(context.Context, storage.LeaderboardFilters) ([]storage.AppWithStats, error)) *MockIAppLeaderboardCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppLeaderboardCall) DoAndReturn(f func(context.Context, storage.LeaderboardFilters) ([]storage.AppWithStats, error)) *MockIAppLeaderboardCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIApp) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIAppMockRecorder) List(ctx, limit, offset, order any) *MockIAppListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIApp)(nil).List), ctx, limit, offset, order)
	return &MockIAppListCall{Call: call}
}

// MockIAppListCall wrap *gomock.Call
type MockIAppListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppListCall) Return(arg0 []*storage.App, arg1 error) *MockIAppListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.App, error)) *MockIAppListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.App, error)) *MockIAppListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIApp) Save(ctx context.Context, m *storage.App) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIAppMockRecorder) Save(ctx, m any) *MockIAppSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIApp)(nil).Save), ctx, m)
	return &MockIAppSaveCall{Call: call}
}

// MockIAppSaveCall wrap *gomock.Call
type MockIAppSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppSaveCall) Return(arg0 error) *MockIAppSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppSaveCall) Do(f func(context.Context, *storage.App) error) *MockIAppSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppSaveCall) DoAndReturn(f func(context.Context, *storage.App) error) *MockIAppSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Series mocks base method.
func (m *MockIApp) Series(ctx context.Context, slug string, timeframe storage.Timeframe, column string, req storage.SeriesRequest) ([]storage.SeriesItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Series", ctx, slug, timeframe, column, req)
	ret0, _ := ret[0].([]storage.SeriesItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Series indicates an expected call of Series.
func (mr *MockIAppMockRecorder) Series(ctx, slug, timeframe, column, req any) *MockIAppSeriesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Series", reflect.TypeOf((*MockIApp)(nil).Series), ctx, slug, timeframe, column, req)
	return &MockIAppSeriesCall{Call: call}
}

// MockIAppSeriesCall wrap *gomock.Call
type MockIAppSeriesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppSeriesCall) Return(items []storage.SeriesItem, err error) *MockIAppSeriesCall {
	c.Call = c.Call.Return(items, err)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppSeriesCall) Do(f func(context.Context, string, storage.Timeframe, string, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIAppSeriesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppSeriesCall) DoAndReturn(f func(context.Context, string, storage.Timeframe, string, storage.SeriesRequest) ([]storage.SeriesItem, error)) *MockIAppSeriesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIApp) Update(ctx context.Context, m *storage.App) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIAppMockRecorder) Update(ctx, m any) *MockIAppUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIApp)(nil).Update), ctx, m)
	return &MockIAppUpdateCall{Call: call}
}

// MockIAppUpdateCall wrap *gomock.Call
type MockIAppUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIAppUpdateCall) Return(arg0 error) *MockIAppUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIAppUpdateCall) Do(f func(context.Context, *storage.App) error) *MockIAppUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIAppUpdateCall) DoAndReturn(f func(context.Context, *storage.App) error) *MockIAppUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
