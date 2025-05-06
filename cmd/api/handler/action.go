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

type ActionHandler struct {
	actions storage.IAction
	cache   cache.ICache
}

func NewActionHandler(
	actions storage.IAction,
	cache cache.ICache,
) *ActionHandler {
	return &ActionHandler{
		actions: actions,
		cache:   cache,
	}
}

var _ Handler = (*ActionHandler)(nil)

func (handler *ActionHandler) InitRoutes(srvr *echo.Group) {
	middlewareCache := cache.NewStatMiddlewareCache(handler.cache)

	srvr.GET("/action/:id", handler.Get, middlewareCache)
}

type getActionRequest struct {
	Id uint64 `param:"id" validate:"required"`
}

// Get godoc
//
//	@Summary		Get action by internal id
//	@Description	Get action by internal id
//	@Tags			actions
//	@ID				get-action
//	@Produce		json
//	@Success		200	{object}	responses.Action
//	@Success		204
//	@Failure		500	{object}	Error
//	@Router			/v1/action/{id} [get]
func (handler *ActionHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getActionRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	action, err := handler.actions.ById(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.actions)
	}
	return c.JSON(http.StatusOK, responses.NewActionWithTx(action))
}
