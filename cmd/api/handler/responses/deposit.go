// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
)

type Deposit struct {
	Id                      uint64         `example:"321"                                                              format:"int64"     json:"id"                        swaggertype:"integer"`
	Height                  pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"                    swaggertype:"integer"`
	Time                    time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                      swaggertype:"string"`
	Amount                  string         `example:"1000"                                                             format:"string"    json:"amount"                    swaggertype:"string"`
	Asset                   string         `example:"nria"                                                             format:"string"    json:"asset"                     swaggertype:"string"`
	DestinationChainAddress string         `example:"0x8bAec8896775DDa83796eda3e7E67217b5E3C5dA"                       format:"string"    json:"destination_chain_address" swaggertype:"string"`
	TxHash                  string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty"         swaggertype:"string"`
	Rollup                  []byte         `example:"O0Ia+lPYYMf3iFfxBaWXCSdlhphc6d4ZoBXINov6Tjc="                     format:"string"    json:"rollup,omitempty"          swaggertype:"string"`

	Bridge *ShortAddress `json:"bridge,omitempty"`
}

func NewDeposit(d storage.Deposit) Deposit {
	deposit := Deposit{
		Id:                      d.Id,
		Height:                  d.Height,
		Time:                    d.Time,
		Amount:                  d.Amount.String(),
		Asset:                   d.Asset,
		DestinationChainAddress: d.DestinationChainAddress,
	}

	if d.Tx != nil {
		deposit.TxHash = hex.EncodeToString(d.Tx.Hash)
	}
	if d.Rollup != nil {
		deposit.Rollup = d.Rollup.AstriaId
	}
	if d.Bridge != nil && d.Bridge.Address != nil {
		NewShortAddress(d.Bridge.Address)
	}

	return deposit
}
