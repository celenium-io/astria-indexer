// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: block_signature.go
//
// Generated by this command:
//
//	mockgen -source=block_signature.go -destination=mock/block_signature.go -package=mock -typed
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

// MockIBlockSignature is a mock of IBlockSignature interface.
type MockIBlockSignature struct {
	ctrl     *gomock.Controller
	recorder *MockIBlockSignatureMockRecorder
}

// MockIBlockSignatureMockRecorder is the mock recorder for MockIBlockSignature.
type MockIBlockSignatureMockRecorder struct {
	mock *MockIBlockSignature
}

// NewMockIBlockSignature creates a new mock instance.
func NewMockIBlockSignature(ctrl *gomock.Controller) *MockIBlockSignature {
	mock := &MockIBlockSignature{ctrl: ctrl}
	mock.recorder = &MockIBlockSignatureMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBlockSignature) EXPECT() *MockIBlockSignatureMockRecorder {
	return m.recorder
}

// CursorList mocks base method.
func (m *MockIBlockSignature) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.BlockSignature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.BlockSignature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIBlockSignatureMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIBlockSignatureCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIBlockSignature)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIBlockSignatureCursorListCall{Call: call}
}

// MockIBlockSignatureCursorListCall wrap *gomock.Call
type MockIBlockSignatureCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockSignatureCursorListCall) Return(arg0 []*storage.BlockSignature, arg1 error) *MockIBlockSignatureCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockSignatureCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.BlockSignature, error)) *MockIBlockSignatureCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockSignatureCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.BlockSignature, error)) *MockIBlockSignatureCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIBlockSignature) GetByID(ctx context.Context, id uint64) (*storage.BlockSignature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.BlockSignature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIBlockSignatureMockRecorder) GetByID(ctx, id any) *MockIBlockSignatureGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIBlockSignature)(nil).GetByID), ctx, id)
	return &MockIBlockSignatureGetByIDCall{Call: call}
}

// MockIBlockSignatureGetByIDCall wrap *gomock.Call
type MockIBlockSignatureGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockSignatureGetByIDCall) Return(arg0 *storage.BlockSignature, arg1 error) *MockIBlockSignatureGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockSignatureGetByIDCall) Do(f func(context.Context, uint64) (*storage.BlockSignature, error)) *MockIBlockSignatureGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockSignatureGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.BlockSignature, error)) *MockIBlockSignatureGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIBlockSignature) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIBlockSignatureMockRecorder) IsNoRows(err any) *MockIBlockSignatureIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIBlockSignature)(nil).IsNoRows), err)
	return &MockIBlockSignatureIsNoRowsCall{Call: call}
}

// MockIBlockSignatureIsNoRowsCall wrap *gomock.Call
type MockIBlockSignatureIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockSignatureIsNoRowsCall) Return(arg0 bool) *MockIBlockSignatureIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockSignatureIsNoRowsCall) Do(f func(error) bool) *MockIBlockSignatureIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockSignatureIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIBlockSignatureIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIBlockSignature) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIBlockSignatureMockRecorder) LastID(ctx any) *MockIBlockSignatureLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIBlockSignature)(nil).LastID), ctx)
	return &MockIBlockSignatureLastIDCall{Call: call}
}

// MockIBlockSignatureLastIDCall wrap *gomock.Call
type MockIBlockSignatureLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockSignatureLastIDCall) Return(arg0 uint64, arg1 error) *MockIBlockSignatureLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockSignatureLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIBlockSignatureLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockSignatureLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIBlockSignatureLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LevelsByValidator mocks base method.
func (m *MockIBlockSignature) LevelsByValidator(ctx context.Context, validatorId uint64, startHeight types.Level) ([]types.Level, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LevelsByValidator", ctx, validatorId, startHeight)
	ret0, _ := ret[0].([]types.Level)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LevelsByValidator indicates an expected call of LevelsByValidator.
func (mr *MockIBlockSignatureMockRecorder) LevelsByValidator(ctx, validatorId, startHeight any) *MockIBlockSignatureLevelsByValidatorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LevelsByValidator", reflect.TypeOf((*MockIBlockSignature)(nil).LevelsByValidator), ctx, validatorId, startHeight)
	return &MockIBlockSignatureLevelsByValidatorCall{Call: call}
}

// MockIBlockSignatureLevelsByValidatorCall wrap *gomock.Call
type MockIBlockSignatureLevelsByValidatorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockSignatureLevelsByValidatorCall) Return(arg0 []types.Level, arg1 error) *MockIBlockSignatureLevelsByValidatorCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockSignatureLevelsByValidatorCall) Do(f func(context.Context, uint64, types.Level) ([]types.Level, error)) *MockIBlockSignatureLevelsByValidatorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockSignatureLevelsByValidatorCall) DoAndReturn(f func(context.Context, uint64, types.Level) ([]types.Level, error)) *MockIBlockSignatureLevelsByValidatorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIBlockSignature) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.BlockSignature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.BlockSignature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIBlockSignatureMockRecorder) List(ctx, limit, offset, order any) *MockIBlockSignatureListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIBlockSignature)(nil).List), ctx, limit, offset, order)
	return &MockIBlockSignatureListCall{Call: call}
}

// MockIBlockSignatureListCall wrap *gomock.Call
type MockIBlockSignatureListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockSignatureListCall) Return(arg0 []*storage.BlockSignature, arg1 error) *MockIBlockSignatureListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockSignatureListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.BlockSignature, error)) *MockIBlockSignatureListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockSignatureListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.BlockSignature, error)) *MockIBlockSignatureListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIBlockSignature) Save(ctx context.Context, m *storage.BlockSignature) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIBlockSignatureMockRecorder) Save(ctx, m any) *MockIBlockSignatureSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIBlockSignature)(nil).Save), ctx, m)
	return &MockIBlockSignatureSaveCall{Call: call}
}

// MockIBlockSignatureSaveCall wrap *gomock.Call
type MockIBlockSignatureSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockSignatureSaveCall) Return(arg0 error) *MockIBlockSignatureSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockSignatureSaveCall) Do(f func(context.Context, *storage.BlockSignature) error) *MockIBlockSignatureSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockSignatureSaveCall) DoAndReturn(f func(context.Context, *storage.BlockSignature) error) *MockIBlockSignatureSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIBlockSignature) Update(ctx context.Context, m *storage.BlockSignature) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIBlockSignatureMockRecorder) Update(ctx, m any) *MockIBlockSignatureUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIBlockSignature)(nil).Update), ctx, m)
	return &MockIBlockSignatureUpdateCall{Call: call}
}

// MockIBlockSignatureUpdateCall wrap *gomock.Call
type MockIBlockSignatureUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIBlockSignatureUpdateCall) Return(arg0 error) *MockIBlockSignatureUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIBlockSignatureUpdateCall) Do(f func(context.Context, *storage.BlockSignature) error) *MockIBlockSignatureUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIBlockSignatureUpdateCall) DoAndReturn(f func(context.Context, *storage.BlockSignature) error) *MockIBlockSignatureUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
