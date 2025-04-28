// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/astria-indexer/pkg/types"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type BlockHandler struct {
	block       storage.IBlock
	blockStats  storage.IBlockStats
	txs         storage.ITx
	actions     storage.IAction
	rollups     storage.IRollup
	price       storage.IPrice
	state       storage.IState
	indexerName string
}

func NewBlockHandler(
	block storage.IBlock,
	blockStats storage.IBlockStats,
	txs storage.ITx,
	actions storage.IAction,
	rollups storage.IRollup,
	price storage.IPrice,
	state storage.IState,
	indexerName string,
) *BlockHandler {
	return &BlockHandler{
		block:       block,
		blockStats:  blockStats,
		txs:         txs,
		actions:     actions,
		rollups:     rollups,
		price:       price,
		state:       state,
		indexerName: indexerName,
	}
}

var _ Handler = (*BlockHandler)(nil)

func (handler *BlockHandler) InitRoutes(srvr *echo.Group) {
	blockGroup := srvr.Group("/block")
	{
		blockGroup.GET("", handler.List)
		blockGroup.GET("/count", handler.Count)
		heightGroup := blockGroup.Group("/:height")
		{
			heightGroup.GET("", handler.Get)
			heightGroup.GET("/actions", handler.GetActions)
			heightGroup.GET("/txs", handler.GetTransactions)
			heightGroup.GET("/stats", handler.GetStats)
			heightGroup.GET("/rollup_actions", handler.GetRollupActions)
			heightGroup.GET("/rollup_actions/count", handler.GetRollupsActionsCount)
			heightGroup.GET("/prices", handler.GetPrices)
		}
	}
}

type getBlockByHeightRequest struct {
	Height types.Level `param:"height" validate:"min=0"`
}

type getBlockRequest struct {
	Height types.Level `param:"height" validate:"min=0"`

	Stats bool `query:"stats" validate:"omitempty"`
}

// Get godoc
//
//	@Summary		Get block info
//	@Description	Get block info
//	@Tags			block
//	@ID				get-block
//	@Param			height	path	integer	true	"Block height"	minimum(1)
//	@Param			stats	query	boolean	false	"Need join stats for block"
//	@Produce		json
//	@Success		200	{object}	responses.Block
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height} [get]
func (handler *BlockHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getBlockRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	block, err := handler.block.ByHeight(c.Request().Context(), req.Height, req.Stats)
	if err != nil {
		return handleError(c, err, handler.block)
	}

	return c.JSON(http.StatusOK, responses.NewBlock(block))
}

type blockListRequest struct {
	Limit  uint64 `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset uint64 `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
	Stats  bool   `query:"stats"  validate:"omitempty"`
}

func (p *blockListRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

// List godoc
//
//	@Summary		List blocks info
//	@Description	List blocks info
//	@Tags			block
//	@ID				list-block
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order"					Enums(asc, desc)
//	@Param			stats	query	boolean	false	"Need join stats for block"
//	@Produce		json
//	@Success		200	{array}		responses.Block
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block [get]
func (handler *BlockHandler) List(c echo.Context) error {
	req, err := bindAndValidate[blockListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	var blocks []*storage.Block
	if req.Stats {
		blocks, err = handler.block.ListWithStats(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	} else {
		blocks, err = handler.block.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	}

	if err != nil {
		return handleError(c, err, handler.block)
	}

	response := make([]responses.Block, len(blocks))
	for i := range blocks {
		response[i] = responses.NewBlock(*blocks[i])
	}

	return returnArray(c, response)
}

type listByHeight struct {
	Height types.Level `param:"height" validate:"min=0"`
	Limit  int         `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int         `query:"offset" validate:"omitempty,min=0"`
}

func (p *listByHeight) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// GetActions godoc
//
//	@Summary		Get actions from begin and end of block
//	@Description	Get actions from begin and end of block
//	@Tags			block
//	@ID				get-block-actions
//	@Param			height	path	integer	true	"Block height"					minimum(1)
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Action
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height}/actions [get]
func (handler *BlockHandler) GetActions(c echo.Context) error {
	req, err := bindAndValidate[listByHeight](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	actions, err := handler.actions.ByBlock(c.Request().Context(), req.Height, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.block)
	}

	response := make([]responses.Action, len(actions))
	for i := range actions {
		response[i] = responses.NewActionWithTx(actions[i])
	}

	return returnArray(c, response)
}

// GetStats godoc
//
//	@Summary		Get block stats by height
//	@Description	Get block stats by height
//	@Tags			block
//	@ID				get-block-stats
//	@Param			height	path	integer	true	"Block height"	minimum(1)
//	@Produce		json
//	@Success		200	{object}	responses.BlockStats
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height}/stats [get]
func (handler *BlockHandler) GetStats(c echo.Context) error {
	req, err := bindAndValidate[getBlockByHeightRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	stats, err := handler.blockStats.ByHeight(c.Request().Context(), req.Height)
	if err != nil {
		return handleError(c, err, handler.block)
	}
	return c.JSON(http.StatusOK, responses.NewBlockStats(&stats))
}

// GetRollupActions godoc
//
//	@Summary		Get rollup actions in the block
//	@Description	Get rollup actions in the block
//	@Tags			block
//	@ID				get-block-rollup-actions
//	@Param			height	path	integer	true	"Block height"					minimum(1)
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.RollupAction
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height}/rollup_actions [get]
func (handler *BlockHandler) GetRollupActions(c echo.Context) error {
	req, err := bindAndValidate[listByHeight](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	actions, err := handler.rollups.ActionsByHeight(c.Request().Context(), req.Height, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.block)
	}
	response := make([]responses.RollupAction, len(actions))
	for i := range response {
		response[i] = responses.NewRollupAction(actions[i])
	}

	return c.JSON(http.StatusOK, response)
}

// GetRollupsActionsCount godoc
//
//	@Summary		Get count of rollup actions
//	@Description	Get count of rollup actions
//	@Tags			block
//	@ID				get-block-rollup-actions-count
//	@Param			height	path	integer	true	"Block height"	minimum(1)
//	@Produce		json
//	@Success		200	{integer}	int64
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height}/rollup_actions/count [get]
func (handler *BlockHandler) GetRollupsActionsCount(c echo.Context) error {
	req, err := bindAndValidate[getBlockByHeightRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	count, err := handler.rollups.CountActionsByHeight(c.Request().Context(), req.Height)
	if err != nil {
		return handleError(c, err, handler.block)
	}

	return c.JSON(http.StatusOK, count)
}

// Count godoc
//
//	@Summary		Get count of blocks in network
//	@Description	Get count of blocks in network
//	@Tags			block
//	@ID				get-block-count
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		500	{object}	Error
//	@Router			/v1/block/count [get]
func (handler *BlockHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.block)
	}
	return c.JSON(http.StatusOK, state.LastHeight+1) // + genesis block
}

// GetTransactions godoc
//
//	@Summary		Get transactions are contained in the block
//	@Description	Get transactions are contained in the block
//	@Tags			block
//	@ID				get-block-transactions
//	@Param			height				path	integer			true	"Block height"					minimum(1)
//	@Param			limit				query	integer			false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset				query	integer			false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Tx
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height}/txs [get]
func (handler *BlockHandler) GetTransactions(c echo.Context) error {
	req, err := bindAndValidate[listByHeight](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	txs, err := handler.txs.ByHeight(c.Request().Context(), req.Height, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.block)
	}
	response := make([]responses.Tx, len(txs))
	for i := range response {
		response[i] = responses.NewTx(txs[i])
	}

	return c.JSON(http.StatusOK, response)
}

// GetPrices godoc
//
//	@Summary		Get prices whuch was published in the block
//	@Description	Get prices whuch was published in the block
//	@Tags			block
//	@ID				get-block-prices
//	@Param			height				path	integer			true	"Block height"					minimum(1)
//	@Param			limit				query	integer			false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset				query	integer			false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Price
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height}/txs [get]
func (handler *BlockHandler) GetPrices(c echo.Context) error {
	req, err := bindAndValidate[listByHeight](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	prices, err := handler.price.ByHeight(c.Request().Context(), req.Height, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.block)
	}
	response := make([]responses.Price, len(prices))
	for i := range response {
		response[i] = responses.NewPrice(prices[i])
	}

	return c.JSON(http.StatusOK, response)
}
