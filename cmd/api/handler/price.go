// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type PriceHandler struct {
	prices storage.IPrice
	market storage.IMarket
	cache  cache.ICache
}

func NewPriceHandler(
	prices storage.IPrice,
	market storage.IMarket,
	cache cache.ICache,
) *PriceHandler {
	return &PriceHandler{
		prices: prices,
		market: market,
		cache:  cache,
	}
}

var _ Handler = (*PriceHandler)(nil)

func (handler *PriceHandler) InitRoutes(srvr *echo.Group) {
	middlewareCache := cache.NewStatMiddlewareCache(handler.cache)

	price := srvr.Group("/price")
	{
		price.GET("", handler.List)
		pair := price.Group("/:pair")
		{
			pair.GET("", handler.Last)
			// pair.GET("/history", handler.History)
			pair.GET("/:timeframe", handler.Series, middlewareCache)
		}
	}
}

type priceListRequest struct {
	Limit  int `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int `query:"offset" validate:"omitempty,min=0"`
}

func (p *priceListRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// List godoc
//
//	@Summary		Get all currency pairs
//	@Description	Get all currency pairs
//	@Tags			price
//	@ID				list-markets
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}	responses.Market
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/price [get]
func (handler *PriceHandler) List(c echo.Context) error {
	req, err := bindAndValidate[priceListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	markets, err := handler.market.List(c.Request().Context(), req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.prices)
	}
	response := make([]responses.Market, len(markets))
	for i := range markets {
		response[i] = responses.NewMarket(markets[i])
	}
	return returnArray(c, response)
}

type priceLastRequest struct {
	Pair string `param:"pair" validate:"required"`
}

// Last godoc
//
//	@Summary		Get the latest price and market info
//	@Description	Get the latest price and market info
//	@Tags			price
//	@ID				get-market
//	@Param			pair		path	string	true	"Currency pair"
//	@Produce		json
//	@Success		200	{object}	responses.Market
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/price/:pair [get]
func (handler *PriceHandler) Last(c echo.Context) error {
	req, err := bindAndValidate[priceLastRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	market, err := handler.market.Get(c.Request().Context(), req.Pair)
	if err != nil {
		return handleError(c, err, handler.prices)
	}
	return c.JSON(http.StatusOK, responses.NewMarket(market))
}

type priceSeriesRequest struct {
	Pair      string            `example:"BTC-USDT"   param:"pair"      swaggertype:"string"  validate:"required"`
	Timeframe storage.Timeframe `example:"day"        param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day"`
	From      int64             `example:"1692892095" query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To        int64             `example:"1692892095" query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// Series godoc
//
//	@Summary		Get price series
//	@Description	Get price series
//	@Tags			price
//	@ID				get-price-series
//	@Param			pair		path	string	true	"Currency pair"
//	@Param			timeframe	path	string	true    "Timeframe" Enums(hour, day)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}	responses.Candle
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/price/:pair/series/:timeframe [get]
func (handler *PriceHandler) Series(c echo.Context) error {
	req, err := bindAndValidate[priceSeriesRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	prices, err := handler.prices.Series(c.Request().Context(), req.Pair, req.Timeframe, storage.NewSeriesRequest(req.From, req.To))
	if err != nil {
		return handleError(c, err, handler.prices)
	}
	response := make([]responses.Candle, len(prices))
	for i := range prices {
		response[i] = responses.NewCandle(prices[i])
	}
	return returnArray(c, response)
}

type marketHistory struct {
	Pair   string `example:"BTC-USDT" param:"pair"   validate:"required"`
	Limit  int    `example:"10"       query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `example:"10"       query:"offset" validate:"omitempty,min=0"`
}

func (p *marketHistory) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// History godoc
//
//	@Summary		Get market settings history changes
//	@Description	Get market settings history changes
//	@Tags			price
//	@ID				get-market-history
//	@Param			pair	path	string	true	"Currency pair"
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}	responses.Market
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/price/:pair/history [get]
func (handler *PriceHandler) History(c echo.Context) error {
	req, err := bindAndValidate[marketHistory](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	markets, err := handler.market.History(c.Request().Context(), req.Pair, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.prices)
	}
	response := make([]responses.Market, len(markets))
	for i := range markets {
		response[i] = responses.NewMarket(markets[i])
	}
	return returnArray(c, response)
}
