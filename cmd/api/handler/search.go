// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	"github.com/labstack/echo/v4"
)

type SearchHandler struct {
	constantCache *cache.ConstantsCache
	search        storage.ISearch
	address       storage.IAddress
	blocks        storage.IBlock
	txs           storage.ITx
	rollups       storage.IRollup
	bridges       storage.IBridge
	validators    storage.IValidator
	celestials    celestials.ICelestial
	app           storage.IApp
}

func NewSearchHandler(
	constantCache *cache.ConstantsCache,
	search storage.ISearch,
	address storage.IAddress,
	blocks storage.IBlock,
	txs storage.ITx,
	rollups storage.IRollup,
	bridges storage.IBridge,
	validators storage.IValidator,
	celestials celestials.ICelestial,
	app storage.IApp,
) *SearchHandler {
	return &SearchHandler{
		constantCache: constantCache,
		search:        search,
		address:       address,
		blocks:        blocks,
		txs:           txs,
		rollups:       rollups,
		bridges:       bridges,
		validators:    validators,
		celestials:    celestials,
		app:           app,
	}
}

var _ Handler = (*SearchHandler)(nil)

func (s *SearchHandler) InitRoutes(srvr *echo.Group) {
	srvr.GET("/search", s.Search)
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

	results, err := s.search.Search(c.Request().Context(), req.Search)
	if err != nil {
		return handleError(c, err, s.address)
	}

	response := make([]responses.SearchResult, len(results))
	for i := range results {

		var body any
		switch results[i].Type {
		case "block":
			block, err := s.blocks.GetByID(c.Request().Context(), results[i].Id)
			if err != nil {
				return handleError(c, err, s.address)
			}
			body = responses.NewBlock(*block)
		case "tx":
			tx, err := s.txs.GetByID(c.Request().Context(), results[i].Id)
			if err != nil {
				return handleError(c, err, s.address)
			}
			body = responses.NewTx(*tx)
		case "rollup":
			rollup, err := s.rollups.GetByID(c.Request().Context(), results[i].Id)
			if err != nil {
				return handleError(c, err, s.address)
			}
			body = responses.NewRollup(rollup)
		case "address":
			address, err := s.address.GetByID(c.Request().Context(), results[i].Id)
			if err != nil {
				return handleError(c, err, s.address)
			}
			sudoAddress, _ := s.constantCache.Get(types.ModuleNameGeneric, "authority_sudo_address")
			ibcSudoAddress, _ := s.constantCache.Get(types.ModuleNameGeneric, "ibc_sudo_address")
			body = responses.NewAddress(*address, nil, sudoAddress, ibcSudoAddress)
		case "validator":
			validator, err := s.validators.GetByID(c.Request().Context(), results[i].Id)
			if err != nil {
				return handleError(c, err, s.address)
			}
			body = responses.NewShortValidator(validator)
		case "bridge":
			bridge, err := s.bridges.ById(c.Request().Context(), results[i].Id)
			if err != nil {
				return handleError(c, err, s.address)
			}
			body = responses.NewBridge(bridge)
		case "app":
			app, err := s.app.GetByID(c.Request().Context(), results[i].Id)
			if err != nil {
				return handleError(c, err, s.address)
			}
			body = responses.NewApp(*app)
		case "celestial":
			address, err := s.address.GetByID(c.Request().Context(), results[i].Id)
			if err != nil {
				return handleError(c, err, s.address)
			}
			sudoAddress, _ := s.constantCache.Get(types.ModuleNameGeneric, "authority_sudo_address")
			ibcSudoAddress, _ := s.constantCache.Get(types.ModuleNameGeneric, "ibc_sudo_address")

			celestials, err := s.celestials.ById(c.Request().Context(), results[i].Value)
			if err != nil {
				return handleError(c, err, s.address)
			}
			address.Celestials = &celestials

			body = responses.NewAddress(*address, nil, sudoAddress, ibcSudoAddress)
			results[i].Type = "address"
		}

		response[i] = responses.NewSearchResult(results[i].Value, results[i].Type, body)
	}
	return returnArray(c, response)
}
