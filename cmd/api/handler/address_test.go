// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// AddressTestSuite -
type AddressTestSuite struct {
	suite.Suite
	address *mock.MockIAddress
	txs     *mock.MockITx
	actions *mock.MockIAction
	rollups *mock.MockIRollup
	state   *mock.MockIState
	echo    *echo.Echo
	handler *AddressHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *AddressTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.address = mock.NewMockIAddress(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.actions = mock.NewMockIAction(s.ctrl)
	s.rollups = mock.NewMockIRollup(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewAddressHandler(s.address, s.txs, s.actions, s.rollups, s.state, testIndexerName)
}

// TearDownSuite -
func (s *AddressTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteAddress_Run(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}

func (s *AddressTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash")
	c.SetParamNames("hash")
	c.SetParamValues(testAddressHash)

	s.address.EXPECT().
		ByHash(gomock.Any(), testAddress.Hash).
		Return(testAddress, nil).
		Times(1)

	s.rollups.EXPECT().
		ByBridgeAddress(gomock.Any(), testAddress.Id).
		Return(testRollup, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var address responses.Address
	err := json.NewDecoder(rec.Body).Decode(&address)
	s.Require().NoError(err)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(1, address.ActionsCount)
	s.Require().EqualValues(1, address.SignedTxCount)
	s.Require().EqualValues(0, address.Height)
	s.Require().EqualValues(10, address.Nonce)
	s.Require().Equal(testAddressHash, address.Hash)
	s.Require().Equal(testRollup.String(), address.BridgedRollup)
}

func (s *AddressTestSuite) TestGetWithoutBridge() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash")
	c.SetParamNames("hash")
	c.SetParamValues(testAddressHash)

	s.address.EXPECT().
		ByHash(gomock.Any(), testAddress.Hash).
		Return(testAddress, nil).
		Times(1)

	s.rollups.EXPECT().
		ByBridgeAddress(gomock.Any(), testAddress.Id).
		Return(storage.Rollup{}, sql.ErrNoRows).
		Times(1)

	s.address.EXPECT().
		IsNoRows(gomock.Any()).
		Return(true).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var address responses.Address
	err := json.NewDecoder(rec.Body).Decode(&address)
	s.Require().NoError(err)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(1, address.ActionsCount)
	s.Require().EqualValues(1, address.SignedTxCount)
	s.Require().EqualValues(0, address.Height)
	s.Require().EqualValues(10, address.Nonce)
	s.Require().Equal(testAddressHash, address.Hash)
	s.Require().Equal("", address.BridgedRollup)
}

func (s *AddressTestSuite) TestGetInvalidAddress() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash")
	c.SetParamNames("hash")
	c.SetParamValues("invalid")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *AddressTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address")

	s.address.EXPECT().
		ListWithBalance(gomock.Any(), storage.AddressListFilter{
			Limit:  10,
			Offset: 0,
			Sort:   pgSort("asc"),
		}).
		Return([]storage.Address{
			testAddress,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var addresses []responses.Address
	err := json.NewDecoder(rec.Body).Decode(&addresses)
	s.Require().NoError(err)

	address := addresses[0]
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(1, address.ActionsCount)
	s.Require().EqualValues(1, address.SignedTxCount)
	s.Require().EqualValues(0, address.Height)
	s.Require().EqualValues(10, address.Nonce)
	s.Require().Equal(testAddressHash, address.Hash)
	s.Require().Equal("1000", address.Balance.Value)
	s.Require().Equal("nria", address.Balance.Currency)
}

func (s *AddressTestSuite) TestTransactions() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("status", "success")
	q.Set("action_types", "sequence")
	q.Set("height", "1000")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/txs")
	c.SetParamNames("hash")
	c.SetParamValues(testAddressHash)

	s.address.EXPECT().
		ByHash(gomock.Any(), testAddress.Hash).
		Return(testAddress, nil).
		Times(1)

	s.txs.EXPECT().
		ByAddress(gomock.Any(), uint64(1), gomock.Any()).
		Return([]storage.Tx{
			testTx,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Transactions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&txs)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().EqualValues(1, tx.Position)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, tx.Hash)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(10, tx.GasWanted)
	s.Require().EqualValues(8, tx.GasUsed)
	s.Require().EqualValues(10, tx.Nonce)
	s.Require().EqualValues([]string{types.ActionTypeSequence.String()}, tx.ActionTypes)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().Equal("codespace", tx.Codespace)
	s.Require().Equal(types.StatusSuccess, tx.Status)
}

func (s *AddressTestSuite) TestActions() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("sort", "desc")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/actions")
	c.SetParamNames("hash")
	c.SetParamValues(testAddressHash)

	s.address.EXPECT().
		ByHash(gomock.Any(), testAddress.Hash).
		Return(testAddress, nil).
		Times(1)

	s.actions.EXPECT().
		ByAddress(gomock.Any(), uint64(1), gomock.Any()).
		Return([]storage.AddressAction{
			{
				AddressId:  1,
				ActionId:   1,
				TxId:       1,
				ActionType: types.ActionTypeSequence,
				Height:     100,
				Time:       testTime,
				Address:    &testAddress,
				Action: &storage.Action{
					Data:     map[string]any{},
					Position: 1,
				},
				Tx: &testTx,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Actions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var actions []responses.Action
	err := json.NewDecoder(rec.Body).Decode(&actions)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(1, action.Id)
	s.Require().EqualValues(100, action.Height)
	s.Require().EqualValues(1, action.Position)
	s.Require().EqualValues(types.ActionTypeSequence, action.Type)
}

func (s *AddressTestSuite) TestCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/count")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(testState, nil).
		Times(1)

	s.Require().NoError(s.handler.Count(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count uint64
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(123, count)
}

func (s *AddressTestSuite) TestRollups() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("sort", "desc")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/address/:hash/actions")
	c.SetParamNames("hash")
	c.SetParamValues(testAddressHash)

	s.address.EXPECT().
		ByHash(gomock.Any(), testAddress.Hash).
		Return(testAddress, nil).
		Times(1)

	s.rollups.EXPECT().
		ListRollupsByAddress(gomock.Any(), uint64(1), 10, 0, sdk.SortOrderDesc).
		Return([]storage.RollupAddress{
			{
				AddressId: 1,
				RollupId:  1,
				Rollup:    &testRollup,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Rollups(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollups []responses.Rollup
	err := json.NewDecoder(rec.Body).Decode(&rollups)
	s.Require().NoError(err)
	s.Require().Len(rollups, 1)

	rollup := rollups[0]
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues(1, rollup.ActionsCount)
	s.Require().EqualValues(100, rollup.FirstHeight)
	s.Require().EqualValues(10, rollup.Size)
	s.Require().Equal(testRollupHash, rollup.AstriaId)
}
