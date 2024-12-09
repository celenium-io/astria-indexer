// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type AssetHandler struct {
	asset  storage.IAsset
	blocks storage.IBlock
}

func NewAssetHandler(
	asset storage.IAsset,
	blocks storage.IBlock,
) *AssetHandler {
	return &AssetHandler{
		asset:  asset,
		blocks: blocks,
	}
}

type assetListRequest struct {
	Limit  uint64 `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset uint64 `query:"offset" validate:"omitempty,min=0"`
}

func (p *assetListRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// List godoc
//
//	@Summary		Get assets info
//	@Description	Get assets info
//	@Tags			assets
//	@ID				get-asset
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{object}	responses.Asset
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/asset [get]
func (handler *AssetHandler) List(c echo.Context) error {
	req, err := bindAndValidate[assetListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	assets, err := handler.asset.List(c.Request().Context(), int(req.Limit), int(req.Offset))
	if err != nil {
		return handleError(c, err, handler.blocks)
	}

	response := make([]responses.Asset, len(assets))
	for i := range assets {
		response[i] = responses.NewAsset(assets[i])
	}

	return c.JSON(http.StatusOK, response)
}
