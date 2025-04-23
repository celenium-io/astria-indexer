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
	"time"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// PriceTestSuite -
type PriceTestSuite struct {
	suite.Suite
	prices  *mock.MockIPrice
	markets *mock.MockIMarket
	echo    *echo.Echo
	handler *PriceHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *PriceTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.prices = mock.NewMockIPrice(s.ctrl)
	s.markets = mock.NewMockIMarket(s.ctrl)
	s.handler = NewPriceHandler(s.prices, s.markets)
}

// TearDownSuite -
func (s *PriceTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuitePrice_Run(t *testing.T) {
	suite.Run(t, new(PriceTestSuite))
}

func (s *PriceTestSuite) TestLast() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/price/:pair")
	c.SetParamNames("pair")
	c.SetParamValues("BTC-USDT")

	s.markets.EXPECT().
		Get(gomock.Any(), "BTC-USDT").
		Return(storage.Market{
			Price: &storage.Price{
				Price: decimal.RequireFromString("50.00"),
				Time:  time.Now(),
			},
			Pair:             "BTC-USDT",
			Decimals:         8,
			Base:             "BTC",
			Quote:            "USDT",
			Enabled:          true,
			MinProviderCount: 1,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Last(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var market responses.Market
	err := json.NewDecoder(rec.Body).Decode(&market)
	s.Require().NoError(err)
	s.Require().NotNil(market.Price)
	s.Require().Equal("BTC-USDT", market.Pair)
	s.Require().Equal("0.0000005", market.Price.Price)
	s.Require().NotEmpty(market.Price.Time)
}

func (s *PriceTestSuite) TestSeries() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/price/:pair/:timeframe")
	c.SetParamNames("pair", "timeframe")
	c.SetParamValues("BTC-USDT", "hour")

	s.markets.EXPECT().
		Decimals(gomock.Any(), "BTC-USDT").
		Return(8, nil).
		Times(1)

	s.prices.EXPECT().
		Series(gomock.Any(), "BTC-USDT", storage.TimeframeHour).
		Return([]storage.Candle{
			{
				CurrencyPair: "BTC-USDT",
				Open:         decimal.RequireFromString("50.00"),
				Time:         time.Now().Add(-time.Hour),
			},
			{
				CurrencyPair: "BTC-USDT",
				Open:         decimal.RequireFromString("51.00"),
				Time:         time.Now(),
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Series(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var prices []responses.Candle
	err := json.NewDecoder(rec.Body).Decode(&prices)
	s.Require().NoError(err)
	s.Require().Len(prices, 2)
	s.Require().Equal("0.0000005", prices[0].Open)
	s.Require().NotEmpty(prices[0].Time)
	s.Require().Equal("0.00000051", prices[1].Open)
	s.Require().NotEmpty(prices[1].Time)
}

func (s *PriceTestSuite) TestList() {
	q := make(url.Values)
	q.Set("limit", "1")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/price")

	s.markets.EXPECT().
		List(gomock.Any(), 1, 0).
		Return([]storage.Market{
			{
				Pair: "BTC-USDT",
				Price: &storage.Price{
					Price: decimal.RequireFromString("50.00"),
					Time:  time.Now(),
				},
				Decimals:         8,
				Base:             "BTC",
				Quote:            "USDT",
				Enabled:          true,
				MinProviderCount: 1,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var markets []responses.Market
	err := json.NewDecoder(rec.Body).Decode(&markets)
	s.Require().NoError(err)
	s.Require().Len(markets, 1)

	s.Require().Equal("BTC-USDT", markets[0].Pair)
	s.Require().NotNil(markets[0].Price)
	s.Require().Equal("BTC-USDT", markets[0].Pair)
	s.Require().Equal("0.0000005", markets[0].Price.Price)
	s.Require().NotEmpty(markets[0].Price.Time)
}
