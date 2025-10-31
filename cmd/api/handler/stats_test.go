// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// StatsTestSuite -
type StatsTestSuite struct {
	suite.Suite
	stats   *mock.MockIStats
	rollups *mock.MockIRollup
	echo    *echo.Echo
	handler *StatsHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *StatsTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.stats = mock.NewMockIStats(s.ctrl)
	s.rollups = mock.NewMockIRollup(s.ctrl)
	s.handler = NewStatsHandler(s.stats, s.rollups, nil)
}

// TearDownSuite -
func (s *StatsTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteStats_Run(t *testing.T) {
	suite.Run(t, new(StatsTestSuite))
}

func (s *StatsTestSuite) TestBlockStatsHistogram() {
	for _, name := range []string{
		storage.SeriesBPS,
		storage.SeriesDataSize,
		storage.SeriesBlockTime,
		storage.SeriesBytesInBlock,
		storage.SeriesSupplyChange,
		storage.SeriesTPS,
		storage.SeriesTxCount,
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeHour,
			storage.TimeframeDay,
			storage.TimeframeMonth,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/series/:name/:timeframe")
			c.SetParamNames("name", "timeframe")
			c.SetParamValues(name, string(tf))

			s.stats.EXPECT().
				Series(gomock.Any(), tf, name, gomock.Any()).
				Return([]storage.SeriesItem{
					{
						Time:  testTime,
						Value: "11234",
						Max:   "782634",
						Min:   "69.6665479793",
					},
				}, nil)

			s.Require().NoError(s.handler.Series(c))
			s.Require().Equal(http.StatusOK, rec.Code)

			var response []responses.SeriesItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("11234", item.Value)
			s.Require().Equal("782634", item.Max)
			s.Require().Equal("69.6665479793", item.Min)
		}
	}
}

func (s *StatsTestSuite) TestRollupStatsHistogram() {
	for _, name := range []string{
		storage.RollupSeriesActionsCount,
		storage.RollupSeriesAvgSize,
		storage.RollupSeriesMaxSize,
		storage.RollupSeriesMinSize,
		storage.RollupSeriesSize,
	} {

		for _, tf := range []storage.Timeframe{
			storage.TimeframeHour,
			storage.TimeframeDay,
			storage.TimeframeMonth,
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/v1/stats/rollup/series/:hash/:name/:timeframe")
			c.SetParamNames("hash", "name", "timeframe")
			c.SetParamValues(testRollupURLHash, name, string(tf))

			s.rollups.EXPECT().
				ByHash(gomock.Any(), testRollup.AstriaId).
				Return(testRollup, nil).
				Times(1)

			s.stats.EXPECT().
				RollupSeries(gomock.Any(), testRollup.Id, tf, name, gomock.Any()).
				Return([]storage.SeriesItem{
					{
						Time:  testTime,
						Value: "11234",
					},
				}, nil)

			s.Require().NoError(s.handler.RollupSeries(c), name)
			s.Require().Equal(http.StatusOK, rec.Code, name)

			var response []responses.RollupSeriesItem
			err := json.NewDecoder(rec.Body).Decode(&response)
			s.Require().NoError(err)
			s.Require().Len(response, 1)

			item := response[0]
			s.Require().Equal("11234", item.Value)
		}
	}
}

func (s *StatsTestSuite) TestSummary() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/summary")

	s.stats.EXPECT().
		Summary(gomock.Any()).
		Return(storage.NetworkSummary{
			BlockTime:    1000,
			TPS:          0.1,
			BPS:          0.2,
			RBPS:         0.15,
			DataSize:     10,
			TxCount:      10,
			Fee:          decimal.RequireFromString("100"),
			Supply:       decimal.RequireFromString("10000"),
			BytesInBlock: 1000,
		}, nil)

	s.Require().NoError(s.handler.Summary(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var summary responses.NetworkSummary
	err := json.NewDecoder(rec.Body).Decode(&summary)
	s.Require().NoError(err)
	s.Require().EqualValues(0.2, summary.BPS)
	s.Require().EqualValues(1000, summary.BlockTime)
	s.Require().EqualValues(1000, summary.BytesInBlock)
	s.Require().EqualValues(10, summary.DataSize)
	s.Require().EqualValues("100", summary.Fee)
	s.Require().EqualValues(0.15, summary.RBPS)
	s.Require().EqualValues(0.1, summary.TPS)
	s.Require().EqualValues("10000", summary.Supply)
	s.Require().EqualValues(10, summary.TxCount)
}

func (s *StatsTestSuite) TestSummaryTimeframe() {
	for _, tf := range []storage.Timeframe{
		storage.TimeframeDay, storage.TimeframeWeek, storage.TimeframeMonth,
	} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/v1/stats/summary/:timeframe")
		c.SetParamNames("timeframe")
		c.SetParamValues(string(tf))

		s.stats.EXPECT().
			SummaryTimeframe(gomock.Any(), tf).
			Return(storage.NetworkSummaryWithChange{
				BlockTime:       1000,
				BlockTimePct:    1,
				TPS:             0.1,
				TPSPct:          2,
				BPS:             0.2,
				BPSPct:          3,
				RBPS:            0.15,
				RBPSPct:         4,
				DataSize:        10,
				DataSizePct:     5,
				TxCount:         10,
				TxCountPct:      6,
				BytesInBlock:    1000,
				BytesInBlockPct: 7,
			}, nil)

		s.Require().NoError(s.handler.SummaryTimeframe(c))
		s.Require().Equal(http.StatusOK, rec.Code)

		var summary responses.NetworkSummaryWithChange
		err := json.NewDecoder(rec.Body).Decode(&summary)
		s.Require().NoError(err)
		s.Require().EqualValues(0.2, summary.BPS)
		s.Require().EqualValues(1000, summary.BlockTime)
		s.Require().EqualValues(1000, summary.BytesInBlock)
		s.Require().EqualValues(10, summary.DataSize)
		s.Require().EqualValues(0.15, summary.RBPS)
		s.Require().EqualValues(0.1, summary.TPS)
		s.Require().EqualValues(10, summary.TxCount)

		s.Require().EqualValues(1, summary.BlockTimePct)
		s.Require().EqualValues(2, summary.TPSPct)
		s.Require().EqualValues(3, summary.BPSPct)
		s.Require().EqualValues(4, summary.RBPSPct)
		s.Require().EqualValues(5, summary.DataSizePct)
		s.Require().EqualValues(6, summary.TxCountPct)
		s.Require().EqualValues(7, summary.BytesInBlockPct)
	}
}

func (s *StatsTestSuite) TestFeeSummary() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/fee/summary")

	s.stats.EXPECT().
		FeeSummary(gomock.Any()).
		Return([]storage.FeeSummary{
			{
				Asset:    currency.DefaultCurrency,
				Amount:   "1000",
				FeeCount: 100,
			},
		}, nil)

	s.Require().NoError(s.handler.FeeSummary(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var result []responses.FeeSummary
	err := json.NewDecoder(rec.Body).Decode(&result)
	s.Require().NoError(err)
	s.Require().Len(result, 1)

	summary := result[0]
	s.Require().EqualValues("1000", summary.Amount)
	s.Require().EqualValues(100, summary.FeeCount)
	s.Require().EqualValues(currency.DefaultCurrency, summary.Asset)
}

func (s *StatsTestSuite) TestTokenTransferDistribution() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/token/transfer_distribution")

	s.stats.EXPECT().
		TokenTransferDistribution(gomock.Any(), 10).
		Return([]storage.TokenTransferDistributionItem{
			{
				Asset:          currency.DefaultCurrency,
				Amount:         "1000",
				TransfersCount: 100,
			},
		}, nil)

	s.Require().NoError(s.handler.TokenTransferDistribution(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var result []responses.TokenTransferDistributionItem
	err := json.NewDecoder(rec.Body).Decode(&result)
	s.Require().NoError(err)
	s.Require().Len(result, 1)

	summary := result[0]
	s.Require().EqualValues("1000", summary.Amount)
	s.Require().EqualValues(100, summary.TransfersCount)
	s.Require().EqualValues(currency.DefaultCurrency, summary.Asset)
}

func (s *StatsTestSuite) TestActiveAddressesCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/v1/stats/summary/active_addresses_count")

	s.stats.EXPECT().
		ActiveAddressesCount(gomock.Any()).
		Return(100, nil)

	s.Require().NoError(s.handler.ActiveAddressesCount(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var result int64
	err := json.NewDecoder(rec.Body).Decode(&result)
	s.Require().NoError(err)
	s.Require().EqualValues(100, result)
}
