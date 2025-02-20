// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

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

	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// RollupTestSuite -
type RollupTestSuite struct {
	suite.Suite
	rollups  *mock.MockIRollup
	actions  *mock.MockIAction
	bridge   *mock.MockIBridge
	deposits *mock.MockIDeposit
	state    *mock.MockIState
	app      *mock.MockIApp
	echo     *echo.Echo
	handler  *RollupHandler
	ctrl     *gomock.Controller
}

// SetupSuite -
func (s *RollupTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.rollups = mock.NewMockIRollup(s.ctrl)
	s.actions = mock.NewMockIAction(s.ctrl)
	s.bridge = mock.NewMockIBridge(s.ctrl)
	s.deposits = mock.NewMockIDeposit(s.ctrl)
	s.app = mock.NewMockIApp(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	cc := cache.NewConstantsCache(nil)
	s.handler = NewRollupHandler(cc, s.rollups, s.actions, s.bridge, s.deposits, s.app, s.state, testIndexerName)
}

// TearDownSuite -
func (s *RollupTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteRollup_Run(t *testing.T) {
	suite.Run(t, new(RollupTestSuite))
}

func (s *RollupTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:hash")
	c.SetParamNames("hash")
	c.SetParamValues(testRollupURLHash)

	s.rollups.EXPECT().
		ByHash(gomock.Any(), testRollup.AstriaId).
		Return(testRollup, nil).
		Times(1)

	s.app.EXPECT().
		ByRollupId(gomock.Any(), testRollup.Id).
		Return(storage.AppWithStats{}, sql.ErrNoRows).
		Times(1)

	s.app.EXPECT().
		IsNoRows(sql.ErrNoRows).
		Return(true).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollup responses.Rollup
	err := json.NewDecoder(rec.Body).Decode(&rollup)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues(1, rollup.ActionsCount)
	s.Require().EqualValues(1, rollup.BridgeCount)
	s.Require().EqualValues(100, rollup.FirstHeight)
	s.Require().EqualValues(10, rollup.Size)
	s.Require().Equal(testRollup.AstriaId, rollup.AstriaId)
	s.Require().Nil(rollup.App)
}

func (s *RollupTestSuite) TestGetWithApp() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:hash")
	c.SetParamNames("hash")
	c.SetParamValues(testRollupURLHash)

	s.rollups.EXPECT().
		ByHash(gomock.Any(), testRollup.AstriaId).
		Return(testRollup, nil).
		Times(1)

	s.app.EXPECT().
		ByRollupId(gomock.Any(), testRollup.Id).
		Return(testAppWithStats, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollup responses.Rollup
	err := json.NewDecoder(rec.Body).Decode(&rollup)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues(1, rollup.ActionsCount)
	s.Require().EqualValues(1, rollup.BridgeCount)
	s.Require().EqualValues(100, rollup.FirstHeight)
	s.Require().EqualValues(10, rollup.Size)
	s.Require().Equal(testRollup.AstriaId, rollup.AstriaId)
	s.Require().NotNil(rollup.App)
}

func (s *RollupTestSuite) TestGetInvalidAddress() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:hash")
	c.SetParamNames("hash")
	c.SetParamValues("invalid")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "validation")
}

func (s *RollupTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup")

	s.rollups.EXPECT().
		ListExt(gomock.Any(), storage.RollupListFilter{
			Limit:     10,
			Offset:    0,
			SortField: "",
			SortOrder: sdk.SortOrderDesc,
		}).
		Return([]storage.Rollup{
			testRollup,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollups []responses.Rollup
	err := json.NewDecoder(rec.Body).Decode(&rollups)
	s.Require().NoError(err)

	rollup := rollups[0]
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues(1, rollup.ActionsCount)
	s.Require().EqualValues(100, rollup.FirstHeight)
	s.Require().EqualValues(10, rollup.Size)
	s.Require().Equal(testRollup.AstriaId, rollup.AstriaId)
}

func (s *RollupTestSuite) TestActions() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("sort", "desc")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:hash/actions")
	c.SetParamNames("hash")
	c.SetParamValues(testRollupURLHash)

	s.rollups.EXPECT().
		ByHash(gomock.Any(), testRollup.AstriaId).
		Return(testRollup, nil).
		Times(1)

	s.actions.EXPECT().
		ByRollup(gomock.Any(), uint64(1), 10, 0, sdk.SortOrderDesc).
		Return([]storage.RollupAction{
			{
				RollupId: 1,
				ActionId: 1,
				TxId:     1,
				Height:   100,
				Time:     testTime,
				Rollup:   &testRollup,
				Action: &storage.Action{
					Data:     map[string]any{},
					Position: 1,
					Type:     types.ActionTypeRollupDataSubmission,
					Id:       1,
					Height:   100,
				},
				Tx: &testTx,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Actions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var actions []responses.RollupAction
	err := json.NewDecoder(rec.Body).Decode(&actions)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(1, action.Id)
	s.Require().EqualValues(100, action.Height)
	s.Require().EqualValues(1, action.Position)
	s.Require().EqualValues(types.ActionTypeRollupDataSubmission, action.Type)
}

func (s *RollupTestSuite) TestCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/count")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(testState, nil).
		Times(1)

	s.Require().NoError(s.handler.Count(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count uint64
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(30, count)
}

func (s *RollupTestSuite) TestAddresses() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("sort", "desc")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:hash/addresses")
	c.SetParamNames("hash")
	c.SetParamValues(testRollupURLHash)

	s.rollups.EXPECT().
		ByHash(gomock.Any(), testRollup.AstriaId).
		Return(testRollup, nil).
		Times(1)

	s.rollups.EXPECT().
		Addresses(gomock.Any(), uint64(1), 10, 0, sdk.SortOrderDesc).
		Return([]storage.RollupAddress{
			{
				RollupId:  1,
				AddressId: 1,
				Height:    100,
				Address:   &testAddress,
				Rollup:    &testRollup,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Addresses(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var addresses []responses.Address
	err := json.NewDecoder(rec.Body).Decode(&addresses)
	s.Require().NoError(err)
	s.Require().Len(addresses, 1)

	address := addresses[0]
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(1, address.ActionsCount)
	s.Require().EqualValues(1, address.SignedTxCount)
	s.Require().EqualValues(0, address.Height)
	s.Require().EqualValues(10, address.Nonce)
	s.Require().Equal(testAddressHash, address.Hash)
}

func (s *RollupTestSuite) TestBridges() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:hash/bridges")
	c.SetParamNames("hash")
	c.SetParamValues(testRollupURLHash)

	s.rollups.EXPECT().
		ByHash(gomock.Any(), testRollup.AstriaId).
		Return(testRollup, nil).
		Times(1)

	s.bridge.EXPECT().
		ByRollup(gomock.Any(), uint64(1), 10, 0).
		Return([]storage.Bridge{
			{
				Asset:    currency.DefaultCurrency,
				FeeAsset: currency.DefaultCurrency,
				Address:  &testAddress,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Bridges(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var bridges []responses.Bridge
	err := json.NewDecoder(rec.Body).Decode(&bridges)
	s.Require().NoError(err)
	s.Require().Len(bridges, 1)

	bridge := bridges[0]
	s.Require().Equal(testAddressHash, bridge.Address)
	s.Require().Equal("nria", bridge.Asset)
	s.Require().Equal("nria", bridge.FeeAsset)
}

func (s *RollupTestSuite) TestAllActions() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")
	q.Set("sort", "desc")
	q.Set("rollup_actions", "false")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:hash/all_actions")
	c.SetParamNames("hash")
	c.SetParamValues(testRollupURLHash)

	s.rollups.EXPECT().
		ByHash(gomock.Any(), testRollup.AstriaId).
		Return(testRollup, nil).
		Times(1)

	s.actions.EXPECT().
		ByRollupAndBridge(gomock.Any(), uint64(1), storage.RollupAndBridgeActionsFilter{
			Limit:         10,
			Offset:        0,
			Sort:          sdk.SortOrderDesc,
			RollupActions: false,
			BridgeActions: true,
		}).
		Return([]storage.ActionWithTx{
			{
				Action: storage.Action{
					Data:     map[string]any{},
					Position: 1,
					Type:     types.ActionTypeRollupDataSubmission,
					Id:       1,
					Height:   100,
				},
				Tx: &testTx,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.AllActions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var actions []responses.RollupAction
	err := json.NewDecoder(rec.Body).Decode(&actions)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(1, action.Id)
	s.Require().EqualValues(100, action.Height)
	s.Require().EqualValues(1, action.Position)
	s.Require().EqualValues(types.ActionTypeRollupDataSubmission, action.Type)
}

func (s *RollupTestSuite) TestDeposits() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/rollup/:hash/deposits")
	c.SetParamNames("hash")
	c.SetParamValues(testRollupURLHash)

	s.rollups.EXPECT().
		ByHash(gomock.Any(), testRollup.AstriaId).
		Return(testRollup, nil).
		Times(1)

	s.deposits.EXPECT().
		ByRollupId(gomock.Any(), uint64(1), 10, 0, sdk.SortOrderDesc).
		Return([]storage.Deposit{
			{
				TxId:     testTx.Id,
				Time:     testTime,
				Height:   1000,
				ActionId: 1,
				Amount:   decimal.RequireFromString("1000"),
				Asset:    currency.DefaultCurrency,
				BridgeId: 1,
				Tx:       &testTx,
				Bridge: &storage.Bridge{
					Address: &testAddress,
				},
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Deposits(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var deposits []responses.Deposit
	err := json.NewDecoder(rec.Body).Decode(&deposits)
	s.Require().NoError(err)
	s.Require().Len(deposits, 1)
}
