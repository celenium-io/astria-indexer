// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// Code generated by MockGen. DO NOT EDIT.
// Source: validator.go
//
// Generated by this command:
//
//	mockgen -source=validator.go -destination=mock/validator.go -package=mock -typed
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

// MockIValidator is a mock of IValidator interface.
type MockIValidator struct {
	ctrl     *gomock.Controller
	recorder *MockIValidatorMockRecorder
}

// MockIValidatorMockRecorder is the mock recorder for MockIValidator.
type MockIValidatorMockRecorder struct {
	mock *MockIValidator
}

// NewMockIValidator creates a new mock instance.
func NewMockIValidator(ctrl *gomock.Controller) *MockIValidator {
	mock := &MockIValidator{ctrl: ctrl}
	mock.recorder = &MockIValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIValidator) EXPECT() *MockIValidatorMockRecorder {
	return m.recorder
}

// CursorList mocks base method.
func (m *MockIValidator) CursorList(ctx context.Context, id, limit uint64, order storage0.SortOrder, cmp storage0.Comparator) ([]*storage.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]*storage.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockIValidatorMockRecorder) CursorList(ctx, id, limit, order, cmp any) *MockIValidatorCursorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockIValidator)(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockIValidatorCursorListCall{Call: call}
}

// MockIValidatorCursorListCall wrap *gomock.Call
type MockIValidatorCursorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIValidatorCursorListCall) Return(arg0 []*storage.Validator, arg1 error) *MockIValidatorCursorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIValidatorCursorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Validator, error)) *MockIValidatorCursorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIValidatorCursorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder, storage0.Comparator) ([]*storage.Validator, error)) *MockIValidatorCursorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockIValidator) GetByID(ctx context.Context, id uint64) (*storage.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*storage.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIValidatorMockRecorder) GetByID(ctx, id any) *MockIValidatorGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIValidator)(nil).GetByID), ctx, id)
	return &MockIValidatorGetByIDCall{Call: call}
}

// MockIValidatorGetByIDCall wrap *gomock.Call
type MockIValidatorGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIValidatorGetByIDCall) Return(arg0 *storage.Validator, arg1 error) *MockIValidatorGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIValidatorGetByIDCall) Do(f func(context.Context, uint64) (*storage.Validator, error)) *MockIValidatorGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIValidatorGetByIDCall) DoAndReturn(f func(context.Context, uint64) (*storage.Validator, error)) *MockIValidatorGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockIValidator) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockIValidatorMockRecorder) IsNoRows(err any) *MockIValidatorIsNoRowsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockIValidator)(nil).IsNoRows), err)
	return &MockIValidatorIsNoRowsCall{Call: call}
}

// MockIValidatorIsNoRowsCall wrap *gomock.Call
type MockIValidatorIsNoRowsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIValidatorIsNoRowsCall) Return(arg0 bool) *MockIValidatorIsNoRowsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIValidatorIsNoRowsCall) Do(f func(error) bool) *MockIValidatorIsNoRowsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIValidatorIsNoRowsCall) DoAndReturn(f func(error) bool) *MockIValidatorIsNoRowsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockIValidator) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockIValidatorMockRecorder) LastID(ctx any) *MockIValidatorLastIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockIValidator)(nil).LastID), ctx)
	return &MockIValidatorLastIDCall{Call: call}
}

// MockIValidatorLastIDCall wrap *gomock.Call
type MockIValidatorLastIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIValidatorLastIDCall) Return(arg0 uint64, arg1 error) *MockIValidatorLastIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIValidatorLastIDCall) Do(f func(context.Context) (uint64, error)) *MockIValidatorLastIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIValidatorLastIDCall) DoAndReturn(f func(context.Context) (uint64, error)) *MockIValidatorLastIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockIValidator) List(ctx context.Context, limit, offset uint64, order storage0.SortOrder) ([]*storage.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]*storage.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIValidatorMockRecorder) List(ctx, limit, offset, order any) *MockIValidatorListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIValidator)(nil).List), ctx, limit, offset, order)
	return &MockIValidatorListCall{Call: call}
}

// MockIValidatorListCall wrap *gomock.Call
type MockIValidatorListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIValidatorListCall) Return(arg0 []*storage.Validator, arg1 error) *MockIValidatorListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIValidatorListCall) Do(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Validator, error)) *MockIValidatorListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIValidatorListCall) DoAndReturn(f func(context.Context, uint64, uint64, storage0.SortOrder) ([]*storage.Validator, error)) *MockIValidatorListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListByPower mocks base method.
func (m *MockIValidator) ListByPower(ctx context.Context, limit, offset int, order storage0.SortOrder) ([]storage.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByPower", ctx, limit, offset, order)
	ret0, _ := ret[0].([]storage.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByPower indicates an expected call of ListByPower.
func (mr *MockIValidatorMockRecorder) ListByPower(ctx, limit, offset, order any) *MockIValidatorListByPowerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByPower", reflect.TypeOf((*MockIValidator)(nil).ListByPower), ctx, limit, offset, order)
	return &MockIValidatorListByPowerCall{Call: call}
}

// MockIValidatorListByPowerCall wrap *gomock.Call
type MockIValidatorListByPowerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIValidatorListByPowerCall) Return(arg0 []storage.Validator, arg1 error) *MockIValidatorListByPowerCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIValidatorListByPowerCall) Do(f func(context.Context, int, int, storage0.SortOrder) ([]storage.Validator, error)) *MockIValidatorListByPowerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIValidatorListByPowerCall) DoAndReturn(f func(context.Context, int, int, storage0.SortOrder) ([]storage.Validator, error)) *MockIValidatorListByPowerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockIValidator) Save(ctx context.Context, m *storage.Validator) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIValidatorMockRecorder) Save(ctx, m any) *MockIValidatorSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIValidator)(nil).Save), ctx, m)
	return &MockIValidatorSaveCall{Call: call}
}

// MockIValidatorSaveCall wrap *gomock.Call
type MockIValidatorSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIValidatorSaveCall) Return(arg0 error) *MockIValidatorSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIValidatorSaveCall) Do(f func(context.Context, *storage.Validator) error) *MockIValidatorSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIValidatorSaveCall) DoAndReturn(f func(context.Context, *storage.Validator) error) *MockIValidatorSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockIValidator) Update(ctx context.Context, m *storage.Validator) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIValidatorMockRecorder) Update(ctx, m any) *MockIValidatorUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIValidator)(nil).Update), ctx, m)
	return &MockIValidatorUpdateCall{Call: call}
}

// MockIValidatorUpdateCall wrap *gomock.Call
type MockIValidatorUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockIValidatorUpdateCall) Return(arg0 error) *MockIValidatorUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockIValidatorUpdateCall) Do(f func(context.Context, *storage.Validator) error) *MockIValidatorUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockIValidatorUpdateCall) DoAndReturn(f func(context.Context, *storage.Validator) error) *MockIValidatorUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
