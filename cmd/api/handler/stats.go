// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/hex"
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type StatsHandler struct {
	repo    storage.IStats
	rollups storage.IRollup
}

func NewStatsHandler(repo storage.IStats, rollups storage.IRollup) StatsHandler {
	return StatsHandler{
		repo:    repo,
		rollups: rollups,
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
func (sh StatsHandler) Summary(c echo.Context) error {
	summary, err := sh.repo.Summary(c.Request().Context())
	if err != nil {
		return internalServerError(c, err)
	}
	return c.JSON(http.StatusOK, responses.NewNetworkSummary(summary))
}

type seriesRequest struct {
	Timeframe  string `example:"hour"       param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day month"`
	SeriesName string `example:"tps"        param:"name"      swaggertype:"string"  validate:"required,oneof=data_size tps bps rbps fee supply_change block_time tx_count bytes_in_block gas_price gas_efficiency gas_used gas_wanted"`
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
//	@Param			name		path	string	true	"Series name"					Enums(data_size, tps, bps, rbps, fee, supply_change, block_time, tx_count, bytes_in_block, gas_price, gas_efficiency, gas_used, gas_wanted)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.SeriesItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/series/{name}/{timeframe} [get]
func (sh StatsHandler) Series(c echo.Context) error {
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
		return internalServerError(c, err)
	}

	response := make([]responses.SeriesItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewSeriesItem(histogram[i])
	}
	return returnArray(c, response)
}

type rollupSeriesRequest struct {
	Hash       string `example:"00112233"   param:"hash"      swaggertype:"string"  validate:"required,rollup_id"`
	Timeframe  string `example:"hour"       param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day month"`
	SeriesName string `example:"size"       param:"name"      swaggertype:"string"  validate:"required,oneof=size avg_size min_size max_size actions_count"`
	From       int64  `example:"1692892095" query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To         int64  `example:"1692892095" query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// RollupSeries godoc
//
//	@Summary		Get histogram with precomputed rollup stats
//	@Description	Get histogram with precomputed rollup by series name and timeframe
//	@Tags			stats
//	@ID				stats-rollup-series
//	@Param			hash		path	string	true	"Hash"							minlength(48)	maxlength(48)
//	@Param			timeframe	path	string	true	"Timeframe"						Enums(hour, day, month)
//	@Param			name		path	string	true	"Series name"					Enums(size, avg_size, min_size, max_size, actions_count)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.RollupSeriesItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/rollup/series/{hash}/{name}/{timeframe} [get]
func (sh StatsHandler) RollupSeries(c echo.Context) error {
	req, err := bindAndValidate[rollupSeriesRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := hex.DecodeString(req.Hash)
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
		return internalServerError(c, err)
	}

	response := make([]responses.RollupSeriesItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewRollupSeriesItem(histogram[i])
	}
	return returnArray(c, response)
}
