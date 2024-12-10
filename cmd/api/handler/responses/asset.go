// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import "github.com/celenium-io/astria-indexer/internal/storage"

type Asset struct {
	Fee           string `example:"1000" format:"string" json:"fee"            swaggertype:"string"`
	FeeCount      int    `example:"100"  format:"number" json:"fee_count"      swaggertype:"integer"`
	Transferred   string `example:"1000" format:"string" json:"transferred"    swaggertype:"string"`
	TransferCount int    `example:"100"  format:"number" json:"transfer_count" swaggertype:"integer"`
	Asset         string `example:"nria" format:"string" json:"asset"          swaggertype:"string"`
	Supply        string `example:"1000" format:"string" json:"supply"         swaggertype:"string"`
}

func NewAsset(asset storage.Asset) Asset {
	return Asset{
		Asset:         asset.Asset,
		Fee:           asset.Fee.String(),
		FeeCount:      asset.FeeCount,
		Transferred:   asset.Transferred.String(),
		TransferCount: asset.TransferCount,
		Supply:        asset.Supply.String(),
	}
}
