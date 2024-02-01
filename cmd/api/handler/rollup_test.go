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

	"github.com/aopoltorzhicky/astria/cmd/api/handler/responses"
	"github.com/aopoltorzhicky/astria/internal/storage"
	"github.com/aopoltorzhicky/astria/internal/storage/mock"
	"github.com/aopoltorzhicky/astria/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// RollupTestSuite -
type RollupTestSuite struct {
	suite.Suite
	rollups *mock.MockIRollup
	actions *mock.MockIAction
	state   *mock.MockIState
	echo    *echo.Echo
	handler *RollupHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *RollupTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.rollups = mock.NewMockIRollup(s.ctrl)
	s.actions = mock.NewMockIAction(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewRollupHandler(s.rollups, s.actions, s.state, testIndexerName)
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
	c.SetParamValues(testRollupHash)

	s.rollups.EXPECT().
		ByHash(gomock.Any(), testRollup.AstriaId).
		Return(testRollup, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollup responses.Rollup
	err := json.NewDecoder(rec.Body).Decode(&rollup)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues(1, rollup.ActionsCount)
	s.Require().EqualValues(100, rollup.FirstHeight)
	s.Require().EqualValues(10, rollup.Size)
	s.Require().Equal(testRollupHash, rollup.AstriaId)
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
		List(gomock.Any(), uint64(10), uint64(0), pgSort("asc")).
		Return([]*storage.Rollup{
			&testRollup,
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
	s.Require().Equal(testRollupHash, rollup.AstriaId)
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
	c.SetParamValues(testRollupHash)

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
					Type:     types.ActionTypeSequence,
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
	s.Require().EqualValues(types.ActionTypeSequence, action.Type)
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
	c.SetParamValues(testRollupHash)

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
