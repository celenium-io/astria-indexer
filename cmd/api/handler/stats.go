// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type StatsHandler struct {
	repo    storage.IStats
	rollups storage.IRollup
	cache   cache.ICache
}

func NewStatsHandler(
	repo storage.IStats,
	rollups storage.IRollup,
	cache cache.ICache,
) *StatsHandler {
	return &StatsHandler{
		repo:    repo,
		rollups: rollups,
		cache:   cache,
	}
}

var _ Handler = (*StatsHandler)(nil)

func (sh *StatsHandler) InitRoutes(srvr *echo.Group) {
	middlewareCache := cache.NewStatMiddlewareCache(sh.cache)

	stats := srvr.Group("/stats")
	{
		stats.GET("/summary", sh.Summary)
		stats.GET("/summary/:timeframe", sh.SummaryTimeframe)
		stats.GET("/summary/active_addresses_count", sh.ActiveAddressesCount)
		stats.GET("/series/:name/:timeframe", sh.Series, middlewareCache)

		rollup := stats.Group("/rollup")
		{
			rollup.GET("/series/:hash/:name/:timeframe", sh.RollupSeries, middlewareCache)
		}

		fee := stats.Group("/fee")
		{
			fee.GET("/summary", sh.FeeSummary, middlewareCache)
		}

		token := stats.Group("/token")
		{
			token.GET("/transfer_distribution", sh.TokenTransferDistribution, middlewareCache)
		}
	}
}

// Summary godoc
//
//	@Summary		Get network summary
//	@Description	Get network summary
//	@Tags			stats
//	@ID				stats-summary
//	@Produce		json
//	@Success		200	{array}		responses.NetworkSummary
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/summary [get]
func (sh *StatsHandler) Summary(c echo.Context) error {
	summary, err := sh.repo.Summary(c.Request().Context())
	if err != nil {
		return handleError(c, err, sh.rollups)
	}
	return c.JSON(http.StatusOK, responses.NewNetworkSummary(summary))
}

type summaryTimeframeRequest struct {
	Timeframe storage.Timeframe `example:"day" param:"timeframe" swaggertype:"string" validate:"required,oneof=day week month"`
}

// SummaryTimeframe godoc
//
//	@Summary		Get network summary for the last period
//	@Description	Get network summary for the last period
//	@Tags			stats
//	@ID				stats-summary-timeframe
//	@Param			timeframe	path	string	true "Timeframe" Enums(day, week, month)
//	@Produce		json
//	@Success		200	{array}		responses.NetworkSummaryWithChange
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/summary/{timeframe} [get]
func (sh *StatsHandler) SummaryTimeframe(c echo.Context) error {
	req, err := bindAndValidate[summaryTimeframeRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	summary, err := sh.repo.SummaryTimeframe(c.Request().Context(), req.Timeframe)
	if err != nil {
		return handleError(c, err, sh.rollups)
	}
	return c.JSON(http.StatusOK, responses.NewNetworkSummaryWithChange(summary))
}

type seriesRequest struct {
	Timeframe  string `example:"hour"       param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day month"`
	SeriesName string `example:"tps"        param:"name"      swaggertype:"string"  validate:"required,oneof=data_size tps bps rbps supply_change block_time tx_count bytes_in_block"`
	From       int64  `example:"1692892095" query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To         int64  `example:"1692892095" query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// Series godoc
//
//	@Summary		Get histogram with precomputed stats
//	@Description	Get histogram with precomputed stats by series name and timeframe
//	@Tags			stats
//	@ID				stats-series
//	@Param			timeframe	path	string	true	"Timeframe"						Enums(hour, day, month)
//	@Param			name		path	string	true	"Series name"					Enums(data_size, tps, bps, rbps, supply_change, block_time, tx_count, bytes_in_block)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.SeriesItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/series/{name}/{timeframe} [get]
func (sh *StatsHandler) Series(c echo.Context) error {
	req, err := bindAndValidate[seriesRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	histogram, err := sh.repo.Series(
		c.Request().Context(),
		storage.Timeframe(req.Timeframe),
		req.SeriesName,
		storage.NewSeriesRequest(req.From, req.To),
	)
	if err != nil {
		return handleError(c, err, sh.rollups)
	}

	response := make([]responses.SeriesItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewSeriesItem(histogram[i])
	}
	return returnArray(c, response)
}

type rollupSeriesRequest struct {
	Hash       string `example:"O0Ia+lPYYMf3iFfxBaWXCSdlhphc6d4ZoBXINov6Tjc=" param:"hash"      swaggertype:"string"  validate:"required,base64url"`
	Timeframe  string `example:"hour"                                         param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day month"`
	SeriesName string `example:"size"                                         param:"name"      swaggertype:"string"  validate:"required,oneof=size avg_size min_size max_size actions_count"`
	From       int64  `example:"1692892095"                                   query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To         int64  `example:"1692892095"                                   query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// RollupSeries godoc
//
//	@Summary		Get histogram with precomputed rollup stats
//	@Description	Get histogram with precomputed rollup by series name and timeframe
//	@Tags			stats
//	@ID				stats-rollup-series
//	@Param			hash		path	string	true	"Base64Url encoded rollup id"
//	@Param			timeframe	path	string	true	"Timeframe"						Enums(hour, day, month)
//	@Param			name		path	string	true	"Series name"					Enums(size, avg_size, min_size, max_size, actions_count)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.RollupSeriesItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/rollup/series/{hash}/{name}/{timeframe} [get]
func (sh *StatsHandler) RollupSeries(c echo.Context) error {
	req, err := bindAndValidate[rollupSeriesRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := base64.URLEncoding.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	rollup, err := sh.rollups.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, sh.rollups)
	}

	histogram, err := sh.repo.RollupSeries(
		c.Request().Context(),
		rollup.Id,
		storage.Timeframe(req.Timeframe),
		req.SeriesName,
		storage.NewSeriesRequest(req.From, req.To),
	)
	if err != nil {
		return handleError(c, err, sh.rollups)
	}

	response := make([]responses.RollupSeriesItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewRollupSeriesItem(histogram[i])
	}
	return returnArray(c, response)
}

// FeeSummary godoc
//
//	@Summary		Get fee summary
//	@Description	Get fee summary
//	@Tags			stats
//	@ID				stats-fee-summary
//	@Produce		json
//	@Success		200	{array}		responses.FeeSummary
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/fee/summary [get]
func (sh *StatsHandler) FeeSummary(c echo.Context) error {
	summary, err := sh.repo.FeeSummary(c.Request().Context())
	if err != nil {
		return handleError(c, err, sh.rollups)
	}
	response := make([]responses.FeeSummary, len(summary))
	for i := range summary {
		response[i] = responses.NewFeeSummary(summary[i])
	}
	return c.JSON(http.StatusOK, response)
}

type tokenTransferDistributionRequest struct {
	Limit uint64 `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (p *tokenTransferDistributionRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// TokenTransferDistribution godoc
//
//	@Summary		Token transfer distribution
//	@Description	Token transfer distribution
//	@Tags			stats
//	@ID				stats-token-transfer-distribution
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Produce		json
//	@Success		200	{array}		responses.TokenTransferDistributionItem
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/token/transfer_distribution [get]
func (sh *StatsHandler) TokenTransferDistribution(c echo.Context) error {
	req, err := bindAndValidate[tokenTransferDistributionRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	items, err := sh.repo.TokenTransferDistribution(c.Request().Context(), int(req.Limit))
	if err != nil {
		return handleError(c, err, sh.rollups)
	}
	response := make([]responses.TokenTransferDistributionItem, len(items))
	for i := range items {
		response[i] = responses.NewTokenTransferDistributionItem(items[i])
	}
	return c.JSON(http.StatusOK, response)
}

// ActiveAddressesCount godoc
//
//	@Summary		Active adddresses count
//	@Description	Active adddresses count
//	@Tags			stats
//	@ID				stats-active-addresses-count
//	@Produce		json
//	@Success		200	{integer}	int64
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/summary/active_addresses_count [get]
func (sh *StatsHandler) ActiveAddressesCount(c echo.Context) error {
	count, err := sh.repo.ActiveAddressesCount(c.Request().Context())
	if err != nil {
		return handleError(c, err, sh.rollups)
	}
	return c.JSON(http.StatusOK, count)
}
