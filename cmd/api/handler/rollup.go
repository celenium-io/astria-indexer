// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type RollupHandler struct {
	rollups     storage.IRollup
	actions     storage.IAction
	bridge      storage.IBridge
	state       storage.IState
	indexerName string
}

func NewRollupHandler(
	rollups storage.IRollup,
	actions storage.IAction,
	bridge storage.IBridge,
	state storage.IState,
	indexerName string,
) *RollupHandler {
	return &RollupHandler{
		rollups:     rollups,
		actions:     actions,
		bridge:      bridge,
		state:       state,
		indexerName: indexerName,
	}
}

type getRollupRequest struct {
	Hash string `param:"hash" validate:"required,base64url"`
}

// Get godoc
//
//	@Summary		Get rollup info
//	@Description	Get rollup info
//	@Tags			rollup
//	@ID				get-rollup
//	@Param			hash	path	string	true	"Base64Url encoded rollup id"
//	@Produce		json
//	@Success		200	{object}	responses.Rollup
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/{hash} [get]
func (handler *RollupHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getRollupRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := base64.URLEncoding.DecodeString(req.Hash)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	rollup, err := handler.rollups.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	return c.JSON(http.StatusOK, responses.NewRollup(&rollup))
}

type listRollupsRequest struct {
	Limit     int    `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset    int    `query:"offset"  validate:"omitempty,min=0"`
	Sort      string `query:"sort"    validate:"omitempty,oneof=asc desc"`
	SortField string `query:"sort_by" validate:"omitempty,oneof=size id"`
}

func (p *listRollupsRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = desc
	}
}

// List godoc
//
//	@Summary		List rollups info
//	@Description	List rollups info
//	@Tags			rollup
//	@ID				list-rollups
//	@Param			limit		query	integer	false	"Count of requested entities"			mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"								mininum(1)
//	@Param			sort		query	string	false	"Sort order"							Enums(asc, desc)
//	@Param			sort_by		query	string	false	"Field using for sorting. Default: id"	Enums(id, size)
//	@Produce		json
//	@Success		200	{array}		responses.Rollup
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup [get]
func (handler *RollupHandler) List(c echo.Context) error {
	req, err := bindAndValidate[listRollupsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs := storage.RollupListFilter{
		Limit:     req.Limit,
		Offset:    req.Offset,
		SortOrder: pgSort(req.Sort),
		SortField: req.SortField,
	}
	rollups, err := handler.rollups.ListExt(c.Request().Context(), fltrs)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	response := make([]responses.Rollup, len(rollups))
	for i := range rollups {
		response[i] = responses.NewRollup(&rollups[i])
	}

	return returnArray(c, response)
}

type getRollupList struct {
	Hash   string `param:"hash"   validate:"required,base64url"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
}

func (p *getRollupList) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

// Actions godoc
//
//	@Summary		Get rollup actions
//	@Description	Get rollup actions
//	@Tags			rollup
//	@ID				rollup-actions
//	@Param			hash			path	string					true	"Base64Url encoded rollup id"
//	@Param			limit			query	integer					false	"Count of requested entities"			minimum(1)		maximum(100)
//	@Param			offset			query	integer					false	"Offset"								minimum(1)
//	@Param			sort			query	string					false	"Sort order"							Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.RollupAction
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/{hash}/actions [get]
func (handler *RollupHandler) Actions(c echo.Context) error {
	req, err := bindAndValidate[getRollupList](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := base64.URLEncoding.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	rollup, err := handler.rollups.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	actions, err := handler.actions.ByRollup(c.Request().Context(), rollup.Id, req.Limit, req.Offset, pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	response := make([]responses.RollupAction, len(actions))
	for i := range actions {
		response[i] = responses.NewRollupAction(actions[i])
	}

	return returnArray(c, response)
}

// Count godoc
//
//	@Summary		Get count of rollups in network
//	@Description	Get count of rollups in network
//	@Tags			rollup
//	@ID				get-rollup-count
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/count [get]
func (handler *RollupHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	return c.JSON(http.StatusOK, state.TotalRollups)
}

// Addresses godoc
//
//	@Summary		List addresses which pushed something in the rollup
//	@Description	List addresses which pushed something in the rollup
//	@Tags			rollup
//	@ID				get-rollup-addresses
//	@Param			limit		query	integer	false	"Count of requested entities"		mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"							mininum(1)
//	@Param			sort		query	string	false	"Sort order"						Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}	    responses.Address
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/{hash}/addresses [get]
func (handler *RollupHandler) Addresses(c echo.Context) error {
	req, err := bindAndValidate[getRollupList](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := base64.URLEncoding.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	rollup, err := handler.rollups.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	addresses, err := handler.rollups.Addresses(c.Request().Context(), rollup.Id, req.Limit, req.Offset, pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	response := make([]responses.Address, len(addresses))
	for i := range addresses {
		if addresses[i].Address != nil {
			response[i] = responses.NewAddress(*addresses[i].Address, nil)
		}
	}

	return returnArray(c, response)
}

// Bridges godoc
//
//	@Summary		Get rollup bridges
//	@Description	Get rollup bridges
//	@Tags			rollup
//	@ID				rollup-bridges
//	@Param			hash			path	string					true	"Base64Url encoded rollup id"
//	@Param			limit			query	integer					false	"Count of requested entities"			minimum(1)		maximum(100)
//	@Param			offset			query	integer					false	"Offset"								minimum(1)
//	@Param			sort			query	string					false	"Sort order"							Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.Bridge
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/{hash}/bridges [get]
func (handler *RollupHandler) Bridges(c echo.Context) error {
	req, err := bindAndValidate[getRollupList](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := base64.URLEncoding.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	rollup, err := handler.rollups.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	bridges, err := handler.bridge.ByRollup(c.Request().Context(), rollup.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	response := make([]responses.Bridge, len(bridges))
	for i := range bridges {
		response[i] = responses.NewBridge(bridges[i])
	}

	return returnArray(c, response)
}
