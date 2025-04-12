package handler

import (
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type PriceHandler struct {
	prices storage.IPrice
}

func NewPriceHandler(prices storage.IPrice) PriceHandler {
	return PriceHandler{
		prices: prices,
	}
}

var _ Handler = (*PriceHandler)(nil)

func (handler *PriceHandler) InitRoutes(srvr *echo.Group) {
	price := srvr.Group("/price/:pair")
	{
		price.GET("", handler.Last)
		price.GET("/:timeframe", handler.Series)
	}
}

type priceLastRequest struct {
	Pair string `param:"pair" validate:"required"`
}

// Last godoc
//
//	@Summary		Get the latest price
//	@Description	Get the latest price
//	@Tags			price
//	@ID				get-price
//	@Param			pair		path	string	true	"Currency pair"
//	@Produce		json
//	@Success		200	{object}	responses.Price
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/price/:pair [get]
func (handler *PriceHandler) Last(c echo.Context) error {
	req, err := bindAndValidate[priceLastRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	price, err := handler.prices.Last(c.Request().Context(), req.Pair)
	if err != nil {
		return handleError(c, err, handler.prices)
	}
	return c.JSON(http.StatusOK, responses.NewPrice(price))
}

type priceSeriesRequest struct {
	Pair      string            `example:"BTC-USDT" param:"pair"      swaggertype:"string" validate:"required"`
	Timeframe storage.Timeframe `example:"day"      param:"timeframe" swaggertype:"string" validate:"required,oneof=hour day"`
}

// Series godoc
//
//	@Summary		Get price series
//	@Description	Get price series
//	@Tags			price
//	@ID				get-price-series
//	@Param			pair		path	string	true	"Currency pair"
//	@Param			timeframe	path	string	true    "Timeframe" Enums(hour, day)
//	@Produce		json
//	@Success		200	{array}	responses.Candle
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/price/:pair/:timeframe [get]
func (handler *PriceHandler) Series(c echo.Context) error {
	req, err := bindAndValidate[priceSeriesRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	prices, err := handler.prices.Series(c.Request().Context(), req.Pair, req.Timeframe)
	if err != nil {
		return handleError(c, err, handler.prices)
	}
	response := make([]responses.Candle, len(prices))
	for i := range prices {
		response[i] = responses.NewCandle(prices[i])
	}
	return returnArray(c, response)
}
