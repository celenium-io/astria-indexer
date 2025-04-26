// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/labstack/echo/v4"
)

type ValidatorHandler struct {
	validators      storage.IValidator
	blocks          storage.IBlock
	blockSignatures storage.IBlockSignature
	state           storage.IState
	indexerName     string
}

func NewValidatorHandler(
	validators storage.IValidator,
	blocks storage.IBlock,
	blockSignatures storage.IBlockSignature,
	state storage.IState,
	indexerName string,
) *ValidatorHandler {
	return &ValidatorHandler{
		validators:      validators,
		blocks:          blocks,
		blockSignatures: blockSignatures,
		state:           state,
		indexerName:     indexerName,
	}
}

var _ Handler = (*ValidatorHandler)(nil)

func (handler *ValidatorHandler) InitRoutes(srvr *echo.Group) {
	validators := srvr.Group("/validators")
	{
		validators.GET("", handler.List)
		validatorGroup := validators.Group("/:id")
		{
			validatorGroup.GET("", handler.Get)
			validatorGroup.GET("/blocks", handler.Blocks)
			validatorGroup.GET("/uptime", handler.Uptime)
		}
	}
}

type validatorRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`
}

// Get godoc
//
//	@Summary		Get validator info
//	@Description	Get validator info
//	@Tags			validator
//	@ID				get-validator
//	@Param			id	path	integer	true	"Internal validator id"
//	@Produce		json
//	@Success		200	{object}	responses.Validator
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/validators/{id} [get]
func (handler *ValidatorHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[validatorRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	validator, err := handler.validators.GetByID(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.validators)
	}

	return c.JSON(http.StatusOK, responses.NewValidator(validator))
}

type validatorsRequest struct {
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
}

func (p *validatorsRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

// List godoc
//
//	@Summary		List validators
//	@Description	List validators sorted by power
//	@Tags			validator
//	@ID				list-validator
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order"					Enums(asc, desc)
//	@Produce		json
//	@Produce		json
//	@Success		200	{array}		responses.Validator
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/validators [get]
func (handler *ValidatorHandler) List(c echo.Context) error {
	req, err := bindAndValidate[validatorsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	validators, err := handler.validators.ListByPower(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.validators)
	}

	response := make([]responses.Validator, len(validators))
	for i := range validators {
		response[i] = *responses.NewValidator(&validators[i])
	}

	return returnArray(c, response)
}

type listValidatorRequest struct {
	Id     uint64 `param:"id"     validate:"required,min=1"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
}

func (p *listValidatorRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = desc
	}
}

// Blocks godoc
//
//	@Summary		List blocks which was proposed by validator
//	@Description	List blocks which was proposed by validator
//	@Tags			validator
//	@ID				get-validator-blocks
//	@Param			id		path	integer	true	"Internal validator id"
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order"					Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}	responses.Block
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/validators/{id}/blocks [get]
func (handler *ValidatorHandler) Blocks(c echo.Context) error {
	req, err := bindAndValidate[listValidatorRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	blocks, err := handler.blocks.ByProposer(c.Request().Context(), req.Id, req.Limit, req.Offset, pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.blocks)
	}

	response := make([]responses.Block, len(blocks))
	for i := range blocks {
		response[i] = responses.NewBlock(blocks[i])
	}
	return returnArray(c, response)
}

type validatorUptimeRequest struct {
	Id    uint64      `param:"id"    validate:"required,min=1"`
	Limit types.Level `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (r *validatorUptimeRequest) SetDefault() {
	if r.Limit == 0 {
		r.Limit = 10
	}
}

// Uptime godoc
//
//	@Summary		Get validator's uptime and history of signed block
//	@Description	Get validator's uptime and history of signed block
//	@Tags			validator
//	@ID				get-validator-uptime
//	@Param			id		path	integer	true	"Internal validator id"
//	@Param			limit	query	integer	false	"Count of requested blocks"	mininum(1)	maximum(100)
//	@Produce		json
//	@Success		200	{object}	responses.ValidatorUptime
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/validators/{id}/uptime [get]
func (handler *ValidatorHandler) Uptime(c echo.Context) error {
	req, err := bindAndValidate[validatorUptimeRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.blockSignatures)
	}

	startHeight := state.LastHeight - req.Limit - 1
	levels, err := handler.blockSignatures.LevelsByValidator(c.Request().Context(), req.Id, startHeight)
	if err != nil {
		return handleError(c, err, handler.blockSignatures)
	}

	response := responses.NewValidatorUptime(levels, state.LastHeight-1, req.Limit)
	return c.JSON(http.StatusOK, response)
}
