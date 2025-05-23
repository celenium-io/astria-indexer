// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	celestialMock "github.com/celenium-io/celestial-module/pkg/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// SearchTestSuite -
type SearchTestSuite struct {
	suite.Suite
	search     *mock.MockISearch
	blocks     *mock.MockIBlock
	txs        *mock.MockITx
	address    *mock.MockIAddress
	rollups    *mock.MockIRollup
	validators *mock.MockIValidator
	bridges    *mock.MockIBridge
	celestials *celestialMock.MockICelestial
	app        *mock.MockIApp
	echo       *echo.Echo
	handler    *SearchHandler
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *SearchTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.search = mock.NewMockISearch(s.ctrl)
	s.address = mock.NewMockIAddress(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.blocks = mock.NewMockIBlock(s.ctrl)
	s.rollups = mock.NewMockIRollup(s.ctrl)
	s.validators = mock.NewMockIValidator(s.ctrl)
	s.bridges = mock.NewMockIBridge(s.ctrl)
	s.celestials = celestialMock.NewMockICelestial(s.ctrl)
	s.app = mock.NewMockIApp(s.ctrl)
	cc := cache.NewConstantsCache(nil)
	s.handler = NewSearchHandler(cc, s.search, s.address, s.blocks, s.txs, s.rollups, s.bridges, s.validators, s.celestials, s.app)
}

// TearDownSuite -
func (s *SearchTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteSearch_Run(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}

func (s *SearchTestSuite) TestSearchAddress() {
	q := make(url.Values)
	q.Add("query", testAddressHash)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), testAddressHash).
		Return([]storage.SearchResult{
			{
				Id:    testAddress.Id,
				Type:  "address",
				Value: testAddressHash,
			},
		}, nil).
		Times(1)

	s.address.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&testAddress, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []responses.SearchResult
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().Equal("address", result.Type)
	s.Require().Equal(testAddressHash, result.Value)
	s.Require().NotNil(result.Body)
}

func (s *SearchTestSuite) TestSearchBlock() {
	q := make(url.Values)
	q.Add("query", testBlockHash)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), testBlockHash).
		Return([]storage.SearchResult{
			{
				Type:  "block",
				Value: testBlockHash,
				Id:    1,
			},
		}, nil).
		Times(1)

	s.blocks.EXPECT().
		GetByID(gomock.Any(), testBlock.Id).
		Return(&testBlock, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []responses.SearchResult
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().Equal("block", result.Type)
	s.Require().Equal(testBlockHash, result.Value)
	s.Require().NotNil(result.Body)
}

func (s *SearchTestSuite) TestSearchTx() {
	q := make(url.Values)
	q.Add("query", testTxHash)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), testTxHash).
		Return([]storage.SearchResult{
			{
				Type:  "tx",
				Value: testTxHash,
				Id:    1,
			},
		}, nil).
		Times(1)

	s.txs.EXPECT().
		GetByID(gomock.Any(), testTx.Id).
		Return(&testTx, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []responses.SearchResult
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().Equal("tx", result.Type)
	s.Require().Equal(testTxHash, result.Value)
	s.Require().NotNil(result.Body)
}

func (s *SearchTestSuite) TestSearchRollup() {
	q := make(url.Values)
	q.Add("query", testRollupHash)

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), testRollupHash).
		Return([]storage.SearchResult{
			{
				Type:  "rollup",
				Value: testRollupHash,
				Id:    1,
			},
		}, nil).
		Times(1)

	s.rollups.EXPECT().
		GetByID(gomock.Any(), testRollup.Id).
		Return(&testRollup, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []responses.SearchResult
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().Equal("rollup", result.Type)
	s.Require().Equal(testRollupHash, result.Value)
	s.Require().NotNil(result.Body)
}

func (s *SearchTestSuite) TestSearchValidator() {
	q := make(url.Values)
	q.Add("query", "nam")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), "nam").
		Return([]storage.SearchResult{
			{
				Type:  "validator",
				Value: "name",
				Id:    1,
			},
		}, nil).
		Times(1)

	s.validators.EXPECT().
		GetByID(gomock.Any(), testValidator.Id).
		Return(&testValidator, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []responses.SearchResult
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().Equal("validator", result.Type)
	s.Require().Equal("name", result.Value)
	s.Require().NotNil(result.Body)
}

func (s *SearchTestSuite) TestSearchBridge() {
	q := make(url.Values)
	q.Add("query", "nam")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), "nam").
		Return([]storage.SearchResult{
			{
				Type:  "bridge",
				Value: "name",
				Id:    1,
			},
		}, nil).
		Times(1)

	s.bridges.EXPECT().
		ById(gomock.Any(), uint64(1)).
		Return(storage.Bridge{
			Id:         1,
			Asset:      "name",
			Address:    &testAddress,
			AddressId:  testAddress.Id,
			Rollup:     &testRollup,
			RollupId:   testRollup.Id,
			FeeAsset:   "fee_asset",
			InitHeight: 1000,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []responses.SearchResult
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().Equal("bridge", result.Type)
	s.Require().Equal("name", result.Value)
	s.Require().NotNil(result.Body)
}

func (s *SearchTestSuite) TestSearchApp() {
	q := make(url.Values)
	q.Add("query", "app")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), "app").
		Return([]storage.SearchResult{
			{
				Type:  "app",
				Value: testApplication.Name,
				Id:    1,
			},
		}, nil).
		Times(1)

	s.app.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&testApplication, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []responses.SearchResult
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().Equal("app", result.Type)
	s.Require().Equal("test app", result.Value)
	s.Require().NotNil(result.Body)
}

func (s *SearchTestSuite) TestSearchCelestial() {
	q := make(url.Values)
	q.Add("query", "celestial")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/search")

	s.search.EXPECT().
		Search(gomock.Any(), "celestial").
		Return([]storage.SearchResult{
			{
				Type:  "celestial",
				Value: "celestial 1",
				Id:    1,
			},
		}, nil).
		Times(1)

	s.address.EXPECT().
		GetByID(gomock.Any(), uint64(1)).
		Return(&testAddress, nil).
		Times(1)

	s.celestials.EXPECT().
		ById(gomock.Any(), "celestial 1").
		Return(*testAddress.Celestials, nil).
		Times(1)

	s.Require().NoError(s.handler.Search(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var results []responses.SearchResult
	err := json.NewDecoder(rec.Body).Decode(&results)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().Equal("address", result.Type)
	s.Require().Equal("celestial 1", result.Value)
	s.Require().NotNil(result.Body)
}
