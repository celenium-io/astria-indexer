// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// TxTestSuite -
type TxTestSuite struct {
	suite.Suite
	tx      *mock.MockITx
	actions *mock.MockIAction
	rollups *mock.MockIRollup
	state   *mock.MockIState
	echo    *echo.Echo
	handler *TxHandler
	ctrl    *gomock.Controller
}

func TestSuiteTx_Run(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}

func (s *TxTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.tx = mock.NewMockITx(s.ctrl)
	s.actions = mock.NewMockIAction(s.ctrl)
	s.rollups = mock.NewMockIRollup(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewTxHandler(s.tx, s.actions, s.rollups, s.state, testIndexerName)
}

func (s *TxTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func (s *TxTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/:hash")
	c.SetParamNames("hash")
	c.SetParamValues(testTxHash)

	s.tx.EXPECT().
		ByHash(gomock.Any(), testTx.Hash).
		Return(testTx, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var tx responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&tx)
	s.Require().NoError(err)

	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, tx.Hash)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(10, tx.GasWanted)
	s.Require().EqualValues(8, tx.GasUsed)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().EqualValues(10, tx.Nonce)
	s.Require().EqualValues(testAddress.Hash, tx.Signer)
	s.Require().Equal("codespace", tx.Codespace)
	s.Require().Equal(types.StatusSuccess, tx.Status)
}

func (s *TxTestSuite) TestGetInvalidTx() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/:hash")
	c.SetParamNames("hash")
	c.SetParamValues("invalid")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *TxTestSuite) TestList() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("action_types", "sequence,transfer")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.tx.EXPECT().
		Filter(gomock.Any(), storage.TxFilter{
			Limit:       2,
			Offset:      0,
			Sort:        pgSort(desc),
			Status:      []string{"success"},
			ActionTypes: types.NewActionTypeMask(types.ActionTypeSequence.String(), types.ActionTypeTransfer.String()),
		}).
		Return([]storage.Tx{
			testTx,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&txs)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, tx.Hash)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(10, tx.GasWanted)
	s.Require().EqualValues(8, tx.GasUsed)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().EqualValues(10, tx.Nonce)
	s.Require().EqualValues(testAddress.Hash, tx.Signer)
	s.Require().Equal("codespace", tx.Codespace)
	s.Require().Equal(types.StatusSuccess, tx.Status)
}

func (s *TxTestSuite) TestListValidationStatusError() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "invalid")
	q.Set("action_types", "sequence,transfer")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *TxTestSuite) TestListValidationActionTypeError() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("action_types", "invalid")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *TxTestSuite) TestListTime() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("action_types", "sequence")
	q.Set("from", "1692880000")
	q.Set("to", "1692890000")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.tx.EXPECT().
		Filter(gomock.Any(), storage.TxFilter{
			Limit:       2,
			Offset:      0,
			Sort:        pgSort(desc),
			Status:      []string{"success"},
			ActionTypes: types.NewActionTypeMask(types.ActionTypeSequence.String()),
			TimeFrom:    time.Unix(1692880000, 0).UTC(),
			TimeTo:      time.Unix(1692890000, 0).UTC(),
		}).
		Return([]storage.Tx{
			testTx,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&txs)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, tx.Hash)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(10, tx.GasWanted)
	s.Require().EqualValues(8, tx.GasUsed)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().EqualValues(10, tx.Nonce)
	s.Require().EqualValues(testAddress.Hash, tx.Signer)
	s.Require().Equal("codespace", tx.Codespace)
	s.Require().Equal(types.StatusSuccess, tx.Status)
}

func (s *TxTestSuite) TestListWithActions() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("height", "100")
	q.Set("action_types", "sequence")
	q.Set("with_actions", "true")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx")

	s.tx.EXPECT().
		Filter(gomock.Any(), storage.TxFilter{
			Limit:       2,
			Offset:      0,
			Sort:        pgSort(desc),
			Status:      []string{"success"},
			Height:      100,
			ActionTypes: types.NewActionTypeMask(types.ActionTypeSequence.String()),
			WithActions: true,
		}).
		Return([]storage.Tx{
			testTx,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&txs)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, tx.Hash)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(10, tx.GasWanted)
	s.Require().EqualValues(8, tx.GasUsed)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().EqualValues(10, tx.Nonce)
	s.Require().EqualValues(testAddress.Hash, tx.Signer)
	s.Require().Equal("codespace", tx.Codespace)
	s.Require().Equal(types.StatusSuccess, tx.Status)
}

func (s *TxTestSuite) TestGetActions() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/:hash/actions")
	c.SetParamNames("hash")
	c.SetParamValues(testTxHash)

	s.tx.EXPECT().
		ByHash(gomock.Any(), testTx.Hash).
		Return(testTx, nil).
		Times(1)

	s.actions.EXPECT().
		ByTxId(gomock.Any(), uint64(1), 2, 0).
		Return([]storage.Action{
			*testRollupAction.Action,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.GetActions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var actions []responses.Action
	err := json.NewDecoder(rec.Body).Decode(&actions)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)
	s.Require().EqualValues(1, actions[0].Id)
	s.Require().EqualValues(100, actions[0].Height)
	s.Require().EqualValues(1, actions[0].Position)
	s.Require().Equal(testTime, actions[0].Time)
	s.Require().EqualValues(string(types.ActionTypeSequence), actions[0].Type)
}

func (s *TxTestSuite) TestCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/count")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(testState, nil).
		Times(1)

	s.Require().NoError(s.handler.Count(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count uint64
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(1234, count)
}

func (s *TxTestSuite) TestRollupActions() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/:hash/rollup_actions")
	c.SetParamNames("hash")
	c.SetParamValues(testTxHash)

	s.tx.EXPECT().
		ByHash(gomock.Any(), testTx.Hash).
		Return(testTx, nil).
		Times(1)

	s.rollups.EXPECT().
		ActionsByTxId(gomock.Any(), testTx.Id, 10, 0).
		Return([]storage.RollupAction{
			{
				TxId:     testTx.Id,
				Time:     testTime,
				Height:   1000,
				ActionId: 1,
				RollupId: 1,
				Rollup:   &testRollup,
				Action: &storage.Action{
					Type: types.ActionTypeSequence,
				},
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.RollupActions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var actions []responses.RollupAction
	err := json.NewDecoder(rec.Body).Decode(&actions)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)
}

func (s *TxTestSuite) TestRollupActionsCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/tx/:hash/rollup_actions/count")
	c.SetParamNames("hash")
	c.SetParamValues(testTxHash)

	s.tx.EXPECT().
		ByHash(gomock.Any(), testTx.Hash).
		Return(testTx, nil).
		Times(1)

	s.rollups.EXPECT().
		CountActionsByTxId(gomock.Any(), testTx.Id).
		Return(1234, nil).
		Times(1)

	s.Require().NoError(s.handler.RollupActionsCount(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count int
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(1234, count)
}
