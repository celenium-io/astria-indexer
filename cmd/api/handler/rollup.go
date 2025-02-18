// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/astria-indexer/internal/test_suite"
	"github.com/labstack/echo/v4"
)

type RollupHandler struct {
	constantCache *cache.ConstantsCache
	rollups       storage.IRollup
	actions       storage.IAction
	bridge        storage.IBridge
	deposits      storage.IDeposit
	app           storage.IApp
	state         storage.IState
	indexerName   string
}

func NewRollupHandler(
	constantCache *cache.ConstantsCache,
	rollups storage.IRollup,
	actions storage.IAction,
	bridge storage.IBridge,
	deposits storage.IDeposit,
	app storage.IApp,
	state storage.IState,
	indexerName string,
) *RollupHandler {
	return &RollupHandler{
		constantCache: constantCache,
		rollups:       rollups,
		actions:       actions,
		bridge:        bridge,
		deposits:      deposits,
		app:           app,
		state:         state,
		indexerName:   indexerName,
	}
}

var _ Handler = (*RollupHandler)(nil)

func (handler *RollupHandler) InitRoutes(srvr *echo.Group) {
	rollupsGroup := srvr.Group("/rollup")
	{
		rollupsGroup.GET("", handler.List)
		rollupsGroup.GET("/count", handler.Count)

		rollupGroup := rollupsGroup.Group("/:hash")
		{
			rollupGroup.GET("", handler.Get)
			rollupGroup.GET("/actions", handler.Actions)
			rollupGroup.GET("/all_actions", handler.AllActions)
			rollupGroup.GET("/addresses", handler.Addresses)
			rollupGroup.GET("/bridges", handler.Bridges)
			rollupGroup.GET("/deposits", handler.Deposits)
		}
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

	response := responses.NewRollup(&rollup)

	if app, err := handler.app.ByRollupId(c.Request().Context(), rollup.Id); err != nil {
		if !handler.app.IsNoRows(err) {
			return handleError(c, err, handler.rollups)
		}
	} else {
		appResp := responses.NewAppWithStats(app)
		response.App = &appResp
	}

	return c.JSON(http.StatusOK, response)
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
			sudoAddress, _ := handler.constantCache.Get(types.ModuleNameGeneric, "authority_sudo_address")
			ibcSudoAddress, _ := handler.constantCache.Get(types.ModuleNameGeneric, "ibc_sudo_address")
			response[i] = responses.NewAddress(*addresses[i].Address, nil, sudoAddress, ibcSudoAddress)
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

type allRollupActionsRequest struct {
	Hash          string      `param:"hash"           validate:"required,base64url"`
	Limit         int         `query:"limit"          validate:"omitempty,min=1,max=100"`
	Offset        int         `query:"offset"         validate:"omitempty,min=0"`
	Sort          string      `query:"sort"           validate:"omitempty,oneof=asc desc"`
	RollupActions *bool       `query:"rollup_actions" validate:"omitempty"`
	BridgeActions *bool       `query:"bridge_actions" validate:"omitempty"`
	ActionTypes   StringArray `query:"action_types"   validate:"omitempty,dive,action_type"`

	From int64 `example:"1692892095" query:"from" swaggertype:"integer" validate:"omitempty,min=1"`
	To   int64 `example:"1692892095" query:"to"   swaggertype:"integer" validate:"omitempty,min=1"`
}

func (p *allRollupActionsRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
	if p.BridgeActions == nil {
		p.BridgeActions = testsuite.Ptr(true)
	}
	if p.RollupActions == nil {
		p.RollupActions = testsuite.Ptr(true)
	}
}

func (p *allRollupActionsRequest) toDbRequest() storage.RollupAndBridgeActionsFilter {
	fltrs := storage.RollupAndBridgeActionsFilter{
		Limit:         p.Limit,
		Offset:        p.Offset,
		Sort:          pgSort(p.Sort),
		RollupActions: *p.RollupActions,
		BridgeActions: *p.BridgeActions,
		ActionTypes:   types.NewActionTypeMask(),
	}

	if p.From > 0 {
		fltrs.From = time.Unix(p.From, 0).UTC()
	}
	if p.To > 0 {
		fltrs.To = time.Unix(p.To, 0).UTC()
	}
	for i := range p.ActionTypes {
		fltrs.ActionTypes.SetType(types.ActionType(p.ActionTypes[i]))
	}

	return fltrs
}

// AllActions godoc
//
//	@Summary		Get rollup actions with actions of all connected bridges
//	@Description	Get rollup actions with actions of all connected bridges
//	@Tags			rollup
//	@ID				rollup-all-actions
//	@Param			hash			path	string				true	"Base64Url encoded rollup id"
//	@Param			limit			query	integer				false	"Count of requested entities"					minimum(1)		maximum(100)
//	@Param			offset			query	integer				false	"Offset"										minimum(1)
//	@Param			sort			query	string				false	"Sort order"									Enums(asc, desc)
//	@Param			rollup_actions	query	boolean				false	"If true join rollup actions. Default: true"
//	@Param			bridge_actions	query	boolean				false	"If true join brigde actions. Default: true"
//	@Param			action_types	query	types.ActionType	false	"Comma-separated action types list"
//	@Param			from			query	integer				false	"Time from in unix timestamp"					mininum(1)
//	@Param			to				query	integer				false	"Time to in unix timestamp"						mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Action
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/{hash}/all_actions [get]
func (handler *RollupHandler) AllActions(c echo.Context) error {
	req, err := bindAndValidate[allRollupActionsRequest](c)
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

	actions, err := handler.actions.ByRollupAndBridge(c.Request().Context(), rollup.Id, req.toDbRequest())
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	response := make([]responses.Action, len(actions))
	for i := range actions {
		response[i] = responses.NewActionWithTx(actions[i])
	}

	return returnArray(c, response)
}

type getRollupDeposits struct {
	Hash   string `param:"hash"   validate:"required,base64url"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
}

func (p *getRollupDeposits) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = desc
	}
}

// Deposits godoc
//
//	@Summary		Get rollup deposits
//	@Description	Get rollup deposits
//	@Tags			rollup
//	@ID				get-rollup-deposits
//	@Param			hash		path	string	true	"Base64Url encoded rollup id"
//	@Param			limit		query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"						mininum(1)
//	@Param			sort		query	string	false	"Sort order"					Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.Deposit
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/{hash}/deposits [get]
func (handler *RollupHandler) Deposits(c echo.Context) error {
	req, err := bindAndValidate[getRollupDeposits](c)
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

	deposits, err := handler.deposits.ByRollupId(c.Request().Context(), rollup.Id, req.Limit, req.Offset, pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	response := make([]responses.Deposit, len(deposits))
	for i := range deposits {
		response[i] = responses.NewDeposit(deposits[i])
	}
	return returnArray(c, response)
}
