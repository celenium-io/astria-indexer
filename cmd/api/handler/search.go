// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/hex"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type SearchHandler struct {
	search     storage.ISearch
	address    storage.IAddress
	blocks     storage.IBlock
	txs        storage.ITx
	rollups    storage.IRollup
	validators storage.IValidator
}

func NewSearchHandler(
	search storage.ISearch,
	address storage.IAddress,
	blocks storage.IBlock,
	txs storage.ITx,
	rollups storage.IRollup,
	validators storage.IValidator,
) *SearchHandler {
	return &SearchHandler{
		search:     search,
		address:    address,
		blocks:     blocks,
		txs:        txs,
		rollups:    rollups,
		validators: validators,
	}
}

type searchRequest struct {
	Search string `query:"query" validate:"required"`
}

// Search godoc
//
//	@Summary				Search by hash or text
//	@Tags					search
//	@ID						search
//	@Param					query	query	string	true	"Search string"
//	@Produce				json
//	@Success				200	{array}		responses.SearchResult
//	@Failure				400	{object}	Error
//	@Failure				500	{object}	Error
//	@Router					/v1/search [get]
func (s *SearchHandler) Search(c echo.Context) error {
	req, err := bindAndValidate[searchRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if isAddress(req.Search) {
		hash, err := hex.DecodeString(req.Search)
		if err != nil {
			return badRequestError(c, err)
		}
		address, err := s.address.ByHash(c.Request().Context(), hash)
		if err != nil {
			return internalServerError(c, err)
		}
		results := []responses.SearchResult{
			responses.NewSearchResult(address.String(), "address", responses.NewAddress(address, nil)),
		}
		return returnArray(c, results)
	}

	if isHash(req.Search) {
		hash, err := hex.DecodeString(req.Search)
		if err != nil {
			return badRequestError(c, err)
		}
		results, err := s.search.Search(c.Request().Context(), hash)
		if err != nil {
			return internalServerError(c, err)
		}

		response := make([]responses.SearchResult, len(results))
		for i := range results {

			var body any
			switch results[i].Type {
			case "block":
				block, err := s.blocks.GetByID(c.Request().Context(), results[i].Id)
				if err != nil {
					return internalServerError(c, err)
				}
				body = responses.NewBlock(*block)
			case "tx":
				tx, err := s.txs.GetByID(c.Request().Context(), results[i].Id)
				if err != nil {
					return internalServerError(c, err)
				}
				body = responses.NewTx(*tx)
			case "rollup":
				rollup, err := s.rollups.GetByID(c.Request().Context(), results[i].Id)
				if err != nil {
					return internalServerError(c, err)
				}
				body = responses.NewRollup(rollup)
			}

			response[i] = responses.NewSearchResult(results[i].Value, results[i].Type, body)
		}
		return returnArray(c, response)
	}

	results, err := s.search.SearchText(c.Request().Context(), req.Search)
	if err != nil {
		return internalServerError(c, err)
	}
	response := make([]responses.SearchResult, len(results))
	for i := range results {
		var body any
		if results[i].Type == "validator" {
			validator, err := s.validators.GetByID(c.Request().Context(), results[i].Id)
			if err != nil {
				return internalServerError(c, err)
			}
			body = responses.NewValidator(validator)
		}

		response[i] = responses.NewSearchResult(results[i].Value, results[i].Type, body)
	}
	return returnArray(c, response)
}
