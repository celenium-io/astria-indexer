// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
)

type TxHandler struct {
	tx          storage.ITx
	actions     storage.IAction
	rollups     storage.IRollup
	fees        storage.IFee
	state       storage.IState
	indexerName string
}

func NewTxHandler(
	tx storage.ITx,
	actions storage.IAction,
	rollups storage.IRollup,
	fees storage.IFee,
	state storage.IState,
	indexerName string,
) *TxHandler {
	return &TxHandler{
		tx:          tx,
		actions:     actions,
		rollups:     rollups,
		fees:        fees,
		state:       state,
		indexerName: indexerName,
	}
}

var _ Handler = (*TxHandler)(nil)

func (handler *TxHandler) InitRoutes(srvr *echo.Group) {
	txGroup := srvr.Group("/tx")
	{
		txGroup.GET("", handler.List)
		txGroup.GET("/count", handler.Count)
		hashGroup := txGroup.Group("/:hash")
		{
			hashGroup.GET("", handler.Get)
			hashGroup.GET("/actions", handler.GetActions)
			hashGroup.GET("/fees", handler.GetFees)
			hashGroup.GET("/rollup_actions", handler.RollupActions)
			hashGroup.GET("/rollup_actions/count", handler.RollupActionsCount)
		}
	}
}

type getTxRequest struct {
	Hash string `param:"hash" validate:"required,hexadecimal,len=64"`

	Fee bool `query:"fee" validate:"omitempty"`
}

// Get godoc
//
//	@Summary		Get transaction by hash
//	@Description	Get transaction by hash
//	@Tags			transactions
//	@ID				get-transaction
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Param			fee 	query	boolean	false	"Flag which indicates need join full transaction fees"
//	@Produce		json
//	@Success		200	{object}	responses.Tx
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash} [get]
func (handler *TxHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getTxRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := responses.NewTx(tx)

	if req.Fee {
		fees, err := handler.fees.FullTxFee(c.Request().Context(), tx.Id)
		if err != nil {
			return handleError(c, err, handler.tx)
		}
		response.Fees = make([]responses.TxFee, len(fees))
		for i := range fees {
			response.Fees[i] = responses.NewTxFee(fees[i])
		}
	}

	return c.JSON(http.StatusOK, response)
}

type txListRequest struct {
	Limit       uint64      `query:"limit"        validate:"omitempty,min=1,max=100"`
	Offset      uint64      `query:"offset"       validate:"omitempty,min=0"`
	Sort        string      `query:"sort"         validate:"omitempty,oneof=asc desc"`
	Height      uint64      `query:"height"       validate:"omitempty,min=1"`
	Status      StringArray `query:"status"       validate:"omitempty,dive,status"`
	ActionTypes StringArray `query:"action_types" validate:"omitempty,dive,action_type"`
	WithActions bool        `query:"with_actions" validate:"omitempty"`

	From int64 `example:"1692892095" query:"from" swaggertype:"integer" validate:"omitempty,min=1"`
	To   int64 `example:"1692892095" query:"to"   swaggertype:"integer" validate:"omitempty,min=1"`
}

func (p *txListRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

// List godoc
//
//	@Summary		List transactions info
//	@Description	List transactions info
//	@Tags			transactions
//	@ID				list-transactions
//	@Param			limit				query	integer				false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset				query	integer				false	"Offset"						mininum(1)
//	@Param			sort				query	string				false	"Sort order"					Enums(asc, desc)
//	@Param			status				query	types.Status		false	"Comma-separated status list"
//	@Param			action_types		query	types.ActionType	false	"Comma-separated action types list"
//	@Param			from				query	integer				false	"Time from in unix timestamp"	mininum(1)
//	@Param			to					query	integer				false	"Time to in unix timestamp"		mininum(1)
//	@Param			height				query	integer				false	"Block number"					mininum(1)
//	@Param			messages			query	boolean				false	"If true join actions"			mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Tx
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx [get]
func (handler *TxHandler) List(c echo.Context) error {
	req, err := bindAndValidate[txListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs := storage.TxFilter{
		Limit:       int(req.Limit),
		Offset:      int(req.Offset),
		Sort:        pgSort(req.Sort),
		Status:      req.Status,
		Height:      req.Height,
		ActionTypes: types.NewActionTypeMask(),
		WithActions: req.WithActions,
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

	txs, err := handler.tx.Filter(c.Request().Context(), fltrs)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := make([]responses.Tx, len(txs))
	for i := range txs {
		response[i] = responses.NewTx(txs[i])
	}
	return returnArray(c, response)
}

type txRequestWithPagination struct {
	Hash   string `param:"hash"   validate:"required,hexadecimal,len=64"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
}

func (p *txRequestWithPagination) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// GetActions godoc
//
//	@Summary		Get transaction actions
//	@Description	Get transaction actions
//	@Tags			transactions
//	@ID				get-transaction-actions
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Param			limit	query	integer	false	"Count of requested entities"		mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"							mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Action
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/actions [get]
func (handler *TxHandler) GetActions(c echo.Context) error {
	req, err := bindAndValidate[txRequestWithPagination](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	events, err := handler.actions.ByTxId(c.Request().Context(), tx.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := make([]responses.Action, len(events))
	for i := range events {
		response[i] = responses.NewAction(events[i])
	}
	return returnArray(c, response)
}

// Count godoc
//
//	@Summary		Get count of transactions in network
//	@Description	Get count of transactions in network
//	@Tags			transactions
//	@ID				get-transactions-count
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/count [get]
func (handler *TxHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	return c.JSON(http.StatusOK, state.TotalTx)
}

// RollupActions godoc
//
//	@Summary		List transaction's rollup actions
//	@Description	List transaction's rollup actions
//	@Tags			transactions
//	@ID				list-transactions-rollup-actions
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Param			limit	query	integer	false	"Count of requested entities"		mininum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"							mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.RollupAction
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/rollup_actions [get]
func (handler *TxHandler) RollupActions(c echo.Context) error {
	req, err := bindAndValidate[txRequestWithPagination](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	actions, err := handler.rollups.ActionsByTxId(c.Request().Context(), tx.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := make([]responses.RollupAction, len(actions))
	for i := range actions {
		response[i] = responses.NewRollupAction(actions[i])
	}
	return returnArray(c, response)
}

// RollupActionsCount godoc
//
//	@Summary		Count of rollup actions
//	@Description	Count of rollup actions
//	@Tags			transactions
//	@ID				list-transactions-rollup-actions-count
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/rollup_actions/count [get]
func (handler *TxHandler) RollupActionsCount(c echo.Context) error {
	req, err := bindAndValidate[getTxRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	count, err := handler.rollups.CountActionsByTxId(c.Request().Context(), tx.Id)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	return c.JSON(http.StatusOK, count)
}

// GetFees godoc
//
//	@Summary		Get transaction fees
//	@Description	Get transaction fees
//	@Tags			transactions
//	@ID				get-transaction-fees
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Param			limit	query	integer	false	"Count of requested entities"		mininum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"							mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.FullFee
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/fees [get]
func (handler *TxHandler) GetFees(c echo.Context) error {
	req, err := bindAndValidate[txRequestWithPagination](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	fees, err := handler.fees.ByTxId(c.Request().Context(), tx.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := make([]responses.FullFee, len(fees))
	for i := range fees {
		response[i] = responses.NewFullFee(fees[i])
	}
	return returnArray(c, response)
}
