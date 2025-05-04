// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type ConstantHandler struct {
	constants              storage.IConstant
	defaultMiddlewareCache echo.MiddlewareFunc
}

func NewConstantHandler(
	constants storage.IConstant,
	defaultMiddlewareCache echo.MiddlewareFunc,
) *ConstantHandler {
	return &ConstantHandler{
		constants:              constants,
		defaultMiddlewareCache: defaultMiddlewareCache,
	}
}

var _ Handler = (*ConstantHandler)(nil)

func (handler *ConstantHandler) InitRoutes(srvr *echo.Group) {
	srvr.GET("/constants", handler.Get)
	srvr.GET("/enums", handler.Enums, handler.defaultMiddlewareCache)
}

// Get godoc
//
//	@Summary		Get network constants
//	@Description	Get network constants
//	@Tags			general
//	@ID				get-constants
//	@Produce		json
//	@Success		200	{object}	responses.Constants
//	@Success		204
//	@Failure		500	{object}	Error
//	@Router			/v1/constants [get]
func (handler *ConstantHandler) Get(c echo.Context) error {
	consts, err := handler.constants.All(c.Request().Context())
	if err != nil {
		return handleError(c, err, handler.constants)
	}
	return c.JSON(http.StatusOK, responses.NewConstants(consts))
}

// Enums godoc
//
//	@Summary		Get astria explorer enumerators
//	@Description	Get astria explorer enumerators
//	@Tags			general
//	@ID				get-enums
//	@Produce		json
//	@Success		200	{object}	responses.Enums
//	@Router			/v1/enums [get]
func (handler *ConstantHandler) Enums(c echo.Context) error {
	return c.JSON(http.StatusOK, responses.NewEnums())
}
