// Code generated by MockGen. DO NOT EDIT.
// Source: api.go
//
// Generated by this command:
//
//	mockgen -source=api.go -destination=mock/api.go -package=mock -typed
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	types "github.com/celenium-io/astria-indexer/pkg/node/types"
	types0 "github.com/celenium-io/astria-indexer/pkg/types"
	gomock "go.uber.org/mock/gomock"
)

// MockApi is a mock of Api interface.
type MockApi struct {
	ctrl     *gomock.Controller
	recorder *MockApiMockRecorder
}

// MockApiMockRecorder is the mock recorder for MockApi.
type MockApiMockRecorder struct {
	mock *MockApi
}

// NewMockApi creates a new mock instance.
func NewMockApi(ctrl *gomock.Controller) *MockApi {
	mock := &MockApi{ctrl: ctrl}
	mock.recorder = &MockApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApi) EXPECT() *MockApiMockRecorder {
	return m.recorder
}

// Block mocks base method.
func (m *MockApi) Block(ctx context.Context, level types0.Level) (types0.ResultBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Block", ctx, level)
	ret0, _ := ret[0].(types0.ResultBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Block indicates an expected call of Block.
func (mr *MockApiMockRecorder) Block(ctx, level any) *ApiBlockCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Block", reflect.TypeOf((*MockApi)(nil).Block), ctx, level)
	return &ApiBlockCall{Call: call}
}

// ApiBlockCall wrap *gomock.Call
type ApiBlockCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiBlockCall) Return(arg0 types0.ResultBlock, arg1 error) *ApiBlockCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiBlockCall) Do(f func(context.Context, types0.Level) (types0.ResultBlock, error)) *ApiBlockCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiBlockCall) DoAndReturn(f func(context.Context, types0.Level) (types0.ResultBlock, error)) *ApiBlockCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BlockData mocks base method.
func (m *MockApi) BlockData(ctx context.Context, level types0.Level) (types0.BlockData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockData", ctx, level)
	ret0, _ := ret[0].(types0.BlockData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockData indicates an expected call of BlockData.
func (mr *MockApiMockRecorder) BlockData(ctx, level any) *ApiBlockDataCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockData", reflect.TypeOf((*MockApi)(nil).BlockData), ctx, level)
	return &ApiBlockDataCall{Call: call}
}

// ApiBlockDataCall wrap *gomock.Call
type ApiBlockDataCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiBlockDataCall) Return(arg0 types0.BlockData, arg1 error) *ApiBlockDataCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiBlockDataCall) Do(f func(context.Context, types0.Level) (types0.BlockData, error)) *ApiBlockDataCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiBlockDataCall) DoAndReturn(f func(context.Context, types0.Level) (types0.BlockData, error)) *ApiBlockDataCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BlockDataGet mocks base method.
func (m *MockApi) BlockDataGet(ctx context.Context, level types0.Level) (types0.BlockData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockDataGet", ctx, level)
	ret0, _ := ret[0].(types0.BlockData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockDataGet indicates an expected call of BlockDataGet.
func (mr *MockApiMockRecorder) BlockDataGet(ctx, level any) *ApiBlockDataGetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockDataGet", reflect.TypeOf((*MockApi)(nil).BlockDataGet), ctx, level)
	return &ApiBlockDataGetCall{Call: call}
}

// ApiBlockDataGetCall wrap *gomock.Call
type ApiBlockDataGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiBlockDataGetCall) Return(arg0 types0.BlockData, arg1 error) *ApiBlockDataGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiBlockDataGetCall) Do(f func(context.Context, types0.Level) (types0.BlockData, error)) *ApiBlockDataGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiBlockDataGetCall) DoAndReturn(f func(context.Context, types0.Level) (types0.BlockData, error)) *ApiBlockDataGetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BlockResults mocks base method.
func (m *MockApi) BlockResults(ctx context.Context, level types0.Level) (types0.ResultBlockResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockResults", ctx, level)
	ret0, _ := ret[0].(types0.ResultBlockResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockResults indicates an expected call of BlockResults.
func (mr *MockApiMockRecorder) BlockResults(ctx, level any) *ApiBlockResultsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockResults", reflect.TypeOf((*MockApi)(nil).BlockResults), ctx, level)
	return &ApiBlockResultsCall{Call: call}
}

// ApiBlockResultsCall wrap *gomock.Call
type ApiBlockResultsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiBlockResultsCall) Return(arg0 types0.ResultBlockResults, arg1 error) *ApiBlockResultsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiBlockResultsCall) Do(f func(context.Context, types0.Level) (types0.ResultBlockResults, error)) *ApiBlockResultsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiBlockResultsCall) DoAndReturn(f func(context.Context, types0.Level) (types0.ResultBlockResults, error)) *ApiBlockResultsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Genesis mocks base method.
func (m *MockApi) Genesis(ctx context.Context) (types.Genesis, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Genesis", ctx)
	ret0, _ := ret[0].(types.Genesis)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Genesis indicates an expected call of Genesis.
func (mr *MockApiMockRecorder) Genesis(ctx any) *ApiGenesisCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Genesis", reflect.TypeOf((*MockApi)(nil).Genesis), ctx)
	return &ApiGenesisCall{Call: call}
}

// ApiGenesisCall wrap *gomock.Call
type ApiGenesisCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiGenesisCall) Return(arg0 types.Genesis, arg1 error) *ApiGenesisCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiGenesisCall) Do(f func(context.Context) (types.Genesis, error)) *ApiGenesisCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiGenesisCall) DoAndReturn(f func(context.Context) (types.Genesis, error)) *ApiGenesisCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Head mocks base method.
func (m *MockApi) Head(ctx context.Context) (types0.ResultBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Head", ctx)
	ret0, _ := ret[0].(types0.ResultBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Head indicates an expected call of Head.
func (mr *MockApiMockRecorder) Head(ctx any) *ApiHeadCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Head", reflect.TypeOf((*MockApi)(nil).Head), ctx)
	return &ApiHeadCall{Call: call}
}

// ApiHeadCall wrap *gomock.Call
type ApiHeadCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiHeadCall) Return(arg0 types0.ResultBlock, arg1 error) *ApiHeadCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiHeadCall) Do(f func(context.Context) (types0.ResultBlock, error)) *ApiHeadCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiHeadCall) DoAndReturn(f func(context.Context) (types0.ResultBlock, error)) *ApiHeadCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Status mocks base method.
func (m *MockApi) Status(ctx context.Context) (types.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", ctx)
	ret0, _ := ret[0].(types.Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status.
func (mr *MockApiMockRecorder) Status(ctx any) *ApiStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockApi)(nil).Status), ctx)
	return &ApiStatusCall{Call: call}
}

// ApiStatusCall wrap *gomock.Call
type ApiStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiStatusCall) Return(arg0 types.Status, arg1 error) *ApiStatusCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiStatusCall) Do(f func(context.Context) (types.Status, error)) *ApiStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiStatusCall) DoAndReturn(f func(context.Context) (types.Status, error)) *ApiStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
