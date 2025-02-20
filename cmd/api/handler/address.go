// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"
	"time"

	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type AddressHandler struct {
	constantCache *cache.ConstantsCache
	address       storage.IAddress
	txs           storage.ITx
	actions       storage.IAction
	rollups       storage.IRollup
	fees          storage.IFee
	bridge        storage.IBridge
	deposits      storage.IDeposit
	state         storage.IState
	indexerName   string
}

func NewAddressHandler(
	constantCache *cache.ConstantsCache,
	address storage.IAddress,
	txs storage.ITx,
	actions storage.IAction,
	rollups storage.IRollup,
	fees storage.IFee,
	bridge storage.IBridge,
	deposits storage.IDeposit,
	state storage.IState,
	indexerName string,
) *AddressHandler {
	return &AddressHandler{
		constantCache: constantCache,
		address:       address,
		txs:           txs,
		actions:       actions,
		rollups:       rollups,
		fees:          fees,
		bridge:        bridge,
		deposits:      deposits,
		state:         state,
		indexerName:   indexerName,
	}
}

var _ Handler = (*AddressHandler)(nil)

func (handler *AddressHandler) InitRoutes(srvr *echo.Group) {
	addressesGroup := srvr.Group("/address")
	{
		addressesGroup.GET("", handler.List)
		addressesGroup.GET("/count", handler.Count)
		addressGroup := addressesGroup.Group("/:hash")
		{
			addressGroup.GET("", handler.Get)
			addressGroup.GET("/txs", handler.Transactions)
			addressGroup.GET("/actions", handler.Actions)
			addressGroup.GET("/rollups", handler.Rollups)
			addressGroup.GET("/roles", handler.Roles)
			addressGroup.GET("/fees", handler.Fees)
			addressGroup.GET("/deposits", handler.Deposits)
		}
	}
}

type getAddressRequest struct {
	Hash string `param:"hash" validate:"required,address"`
}

// Get godoc
//
//	@Summary		Get address info
//	@Description	Get address info
//	@Tags			address
//	@ID				get-address
//	@Param			hash	path	string	true	"Hash"	minlength(48)	maxlength(48)
//	@Produce		json
//	@Success		200	{object}	responses.Address
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash} [get]
func (handler *AddressHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getAddressRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	address, err := handler.address.ByHash(c.Request().Context(), req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	sudoAddress, _ := handler.constantCache.Get(types.ModuleNameGeneric, "authority_sudo_address")
	ibcSudoAddress, _ := handler.constantCache.Get(types.ModuleNameGeneric, "ibc_sudo_address")

	bridge, err := handler.bridge.ByAddress(c.Request().Context(), address.Id)
	if err != nil {
		if !handler.bridge.IsNoRows(err) {
			return handleError(c, err, handler.address)
		}

		return c.JSON(http.StatusOK, responses.NewAddress(address, nil, sudoAddress, ibcSudoAddress))
	}

	return c.JSON(http.StatusOK, responses.NewAddress(address, &bridge, sudoAddress, ibcSudoAddress))
}

type listAddressRequest struct {
	Limit  uint64 `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset uint64 `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
	Asset  string `query:"asset"  validate:"omitempty"`
}

func (p *listAddressRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

// List godoc
//
//	@Summary		List address info
//	@Description	List address info
//	@Tags			address
//	@ID				list-address
//	@Param			limit		query	integer	false	"Count of requested entities"		mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"							mininum(1)
//	@Param			sort		query	string	false	"Sort order"						Enums(asc, desc)
//	@Param			asset		query	string	false	"Required balance asset"
//	@Produce		json
//	@Success		200	{array}		responses.Address
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address [get]
func (handler *AddressHandler) List(c echo.Context) error {
	req, err := bindAndValidate[listAddressRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs := storage.AddressListFilter{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
		Sort:   pgSort(req.Sort),
		Asset:  req.Asset,
	}

	address, err := handler.address.ListWithBalance(c.Request().Context(), fltrs)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	sudoAddress, _ := handler.constantCache.Get(types.ModuleNameGeneric, "authority_sudo_address")
	ibcSudoAddress, _ := handler.constantCache.Get(types.ModuleNameGeneric, "ibc_sudo_address")

	response := make([]responses.Address, len(address))
	for i := range address {
		response[i] = responses.NewAddress(address[i], nil, sudoAddress, ibcSudoAddress)
	}

	return returnArray(c, response)
}

type addressTxRequest struct {
	Hash        string      `param:"hash"         validate:"required,address"`
	Limit       uint64      `query:"limit"        validate:"omitempty,min=1,max=100"`
	Offset      uint64      `query:"offset"       validate:"omitempty,min=0"`
	Sort        string      `query:"sort"         validate:"omitempty,oneof=asc desc"`
	Height      uint64      `query:"height"       validate:"omitempty,min=1"`
	Status      StringArray `query:"status"       validate:"omitempty,dive,status"`
	ActionTypes StringArray `query:"action_types" validate:"omitempty,dive,action_type"`

	From int64 `example:"1692892095" query:"from" swaggertype:"integer" validate:"omitempty,min=1"`
	To   int64 `example:"1692892095" query:"to"   swaggertype:"integer" validate:"omitempty,min=1"`
}

func (p *addressTxRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

// Transactions godoc
//
//	@Summary		Get address transactions
//	@Description	Get address transactions
//	@Tags			address
//	@ID				address-transactions
//	@Param			hash			path	string					true	"Hash"							minlength(48)	maxlength(48)
//	@Param			limit			query	integer					false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset			query	integer					false	"Offset"						minimum(1)
//	@Param			sort			query	string					false	"Sort order"					Enums(asc, desc)
//	@Param			status			query	types.Status	    	false	"Comma-separated status list"
//	@Param			action_types	query	types.ActionType	    false	"Comma-separated action types list"
//	@Param			from			query	integer					false	"Time from in unix timestamp"	minimum(1)
//	@Param			to				query	integer					false	"Time to in unix timestamp"		minimum(1)
//	@Param			height			query	integer					false	"Block number"					minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Tx
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash}/txs [get]
func (handler *AddressHandler) Transactions(c echo.Context) error {
	req, err := bindAndValidate[addressTxRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	address, err := handler.address.ByHash(c.Request().Context(), req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	fltrs := storage.TxFilter{
		Limit:       int(req.Limit),
		Offset:      int(req.Offset),
		Sort:        pgSort(req.Sort),
		Status:      req.Status,
		Height:      req.Height,
		ActionTypes: types.NewActionTypeMask(),
	}
	if req.From > 0 {
		fltrs.TimeFrom = time.Unix(req.From, 0).UTC()
	}
	if req.To > 0 {
		fltrs.TimeTo = time.Unix(req.To, 0).UTC()
	}
	for i := range req.ActionTypes {
		fltrs.ActionTypes.SetType(types.ActionType(req.ActionTypes[i]))
	}

	txs, err := handler.txs.ByAddress(c.Request().Context(), address.Id, fltrs)
	if err != nil {
		return handleError(c, err, handler.address)
	}
	response := make([]responses.Tx, len(txs))
	for i := range txs {
		response[i] = responses.NewTx(txs[i])
	}
	return returnArray(c, response)
}

type getAddressMessages struct {
	Hash        string      `param:"hash"         validate:"required,address"`
	Limit       uint64      `query:"limit"        validate:"omitempty,min=1,max=100"`
	Offset      uint64      `query:"offset"       validate:"omitempty,min=0"`
	Sort        string      `query:"sort"         validate:"omitempty,oneof=asc desc"`
	ActionTypes StringArray `query:"action_types" validate:"omitempty,dive,action_type"`
}

func (p *getAddressMessages) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
	if p.ActionTypes == nil {
		p.ActionTypes = make(StringArray, 0)
	}
}

func (p *getAddressMessages) ToFilters() storage.AddressActionsFilter {
	fltrs := storage.AddressActionsFilter{
		Limit:       int(p.Limit),
		Offset:      int(p.Offset),
		Sort:        pgSort(p.Sort),
		ActionTypes: types.NewActionTypeMask(),
	}

	for i := range p.ActionTypes {
		fltrs.ActionTypes.SetType(types.ActionType(p.ActionTypes[i]))
	}

	return fltrs
}

// Actions godoc
//
//	@Summary		Get address actions
//	@Description	Get address actions
//	@Tags			address
//	@ID				address-actions
//	@Param			hash			path	string					true	"Hash"									minlength(48)	maxlength(48)
//	@Param			limit			query	integer					false	"Count of requested entities"			minimum(1)		maximum(100)
//	@Param			offset			query	integer					false	"Offset"								minimum(1)
//	@Param			sort			query	string					false	"Sort order"							Enums(asc, desc)
//	@Param			action_types	query	types.ActionType     	false	"Comma-separated action types list"
//	@Produce		json
//	@Success		200	{array}		responses.Action
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash}/actions [get]
func (handler *AddressHandler) Actions(c echo.Context) error {
	req, err := bindAndValidate[getAddressMessages](c)
	if err != nil {
		return badRequestError(c, err)
	}

	req.SetDefault()

	address, err := handler.address.ByHash(c.Request().Context(), req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	filters := req.ToFilters()
	actions, err := handler.actions.ByAddress(c.Request().Context(), address.Id, filters)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Action, len(actions))
	for i := range actions {
		response[i] = responses.NewAddressAction(actions[i])
	}

	return returnArray(c, response)
}

// Count godoc
//
//	@Summary		Get count of addresses in network
//	@Description	Get count of addresses in network
//	@Tags			address
//	@ID				get-address-count
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		500	{object}	Error
//	@Router			/v1/address/count [get]
func (handler *AddressHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.address)
	}
	return c.JSON(http.StatusOK, state.TotalAccounts)
}

type getAddressRollups struct {
	Hash   string `param:"hash"   validate:"required,address"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
}

func (p *getAddressRollups) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

// Rollups godoc
//
//	@Summary		Get rollups in which the address pushed something
//	@Description	Get rollups in which the address pushed something
//	@Tags			address
//	@ID				address-rollups
//	@Param			hash			path	string		true	"Hash"									minlength(48)	maxlength(48)
//	@Param			limit			query	integer		false	"Count of requested entities"			minimum(1)		maximum(100)
//	@Param			offset			query	integer		false	"Offset"								minimum(1)
//	@Param			sort			query	string		false	"Sort order"							Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.Rollup
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash}/rollups [get]
func (handler *AddressHandler) Rollups(c echo.Context) error {
	req, err := bindAndValidate[getAddressRollups](c)
	if err != nil {
		return badRequestError(c, err)
	}

	req.SetDefault()

	address, err := handler.address.ByHash(c.Request().Context(), req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	rollups, err := handler.rollups.ListRollupsByAddress(c.Request().Context(), address.Id, req.Limit, req.Offset, pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Rollup, len(rollups))
	for i := range rollups {
		response[i] = responses.NewRollup(rollups[i].Rollup)
	}

	return returnArray(c, response)
}

type getAddressRoles struct {
	Hash   string `param:"hash"   validate:"required,address"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
}

func (p *getAddressRoles) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// Roles godoc
//
//	@Summary		Get address roles in bridges
//	@Description	Get address roles in bridges
//	@Tags			address
//	@ID				address-roles
//	@Param			hash		path	string	true	"Hash"								minlength(48)	maxlength(48)
//	@Param			limit		query	integer	false	"Count of requested entities"		mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"							mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Bridge
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash}/roles [get]
func (handler *AddressHandler) Roles(c echo.Context) error {
	req, err := bindAndValidate[getAddressRoles](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	address, err := handler.address.ByHash(c.Request().Context(), req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	roles, err := handler.bridge.ByRoles(c.Request().Context(), address.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Bridge, len(roles))
	for i := range roles {
		response[i] = responses.NewBridge(roles[i])
	}
	return returnArray(c, response)
}

type getAddressFees struct {
	Hash   string `param:"hash"   validate:"required,address"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
}

func (p *getAddressFees) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = desc
	}
}

// Fees godoc
//
//	@Summary		Get address paid fees
//	@Description	Get address paid fees
//	@Tags			address
//	@ID				get-address-fees
//	@Param			hash		path	string	true	"Hash"								minlength(48)	maxlength(48)
//	@Param			limit		query	integer	false	"Count of requested entities"		mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"							mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.FullFee
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash}/fees [get]
func (handler *AddressHandler) Fees(c echo.Context) error {
	req, err := bindAndValidate[getAddressFees](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	address, err := handler.address.ByHash(c.Request().Context(), req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	fees, err := handler.fees.ByPayerId(c.Request().Context(), address.Id, req.Limit, req.Offset, pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.address)
	}
	response := make([]responses.FullFee, len(fees))
	for i := range fees {
		response[i] = responses.NewFullFee(fees[i])
	}
	return returnArray(c, response)
}

type getAddressDeposits struct {
	Hash   string `param:"hash"   validate:"required,address"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
}

func (p *getAddressDeposits) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = desc
	}
}

// Deposits godoc
//
//	@Summary		Get bridge deposits
//	@Description	Get bridge deposits
//	@Tags			address
//	@ID				get-address-deposits
//	@Param			hash		path	string	true	"Hash"								minlength(48)	maxlength(48)
//	@Param			limit		query	integer	false	"Count of requested entities"		mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"							mininum(1)
//	@Param			sort		query	string		false	"Sort order"					Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.Deposit
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash}/deposits [get]
func (handler *AddressHandler) Deposits(c echo.Context) error {
	req, err := bindAndValidate[getAddressDeposits](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	address, err := handler.address.ByHash(c.Request().Context(), req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	if !address.IsBridge {
		return handleError(c, errors.Errorf("address %s is not bridge", req.Hash), handler.address)
	}

	bridge, err := handler.bridge.ByAddress(c.Request().Context(), address.Id)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	deposits, err := handler.deposits.ByBridgeId(c.Request().Context(), bridge.Id, req.Limit, req.Offset, pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.address)
	}
	response := make([]responses.Deposit, len(deposits))
	for i := range deposits {
		response[i] = responses.NewDeposit(deposits[i])
	}
	return returnArray(c, response)
}
