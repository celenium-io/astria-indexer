package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	echo    *echo.Echo
	handler PriceHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *PriceTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.prices = mock.NewMockIPrice(s.ctrl)
	s.handler = NewPriceHandler(s.prices)
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

	s.prices.EXPECT().
		Last(gomock.Any(), "BTC-USDT").
		Return(storage.Price{
			CurrencyPair: "BTC-USDT",
			Price:        decimal.RequireFromString("50.00"),
			Time:         time.Now(),
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Last(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var price responses.Price
	err := json.NewDecoder(rec.Body).Decode(&price)
	s.Require().NoError(err)
	s.Require().Equal("BTC-USDT", price.CurrencyPair)
	s.Require().Equal("50", price.Price)
	s.Require().NotEmpty(price.Time)
}

func (s *PriceTestSuite) TestSeries() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/price/:pair/:timeframe")
	c.SetParamNames("pair", "timeframe")
	c.SetParamValues("BTC-USDT", "hour")

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
	s.Require().Equal("50", prices[0].Open)
	s.Require().NotEmpty(prices[0].Time)
	s.Require().Equal("51", prices[1].Open)
	s.Require().NotEmpty(prices[1].Time)
}
