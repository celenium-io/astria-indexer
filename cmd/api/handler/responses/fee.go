// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
)

type FullFee struct {
	Amount string         `example:"1000"                                                             format:"string"    json:"amount"            swaggertype:"string"`
	Asset  string         `example:"nria"                                                             format:"string"    json:"asset"             swaggertype:"string"`
	Payer  string         `example:"astria1phym4uktjn6gjle226009ge7u82w0dgtszs8x2"                    format:"string"    json:"payer,omitempty"   swaggertype:"string"`
	TxHash string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty" swaggertype:"string"`
	Height pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"            swaggertype:"integer"`
	Time   time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"              swaggertype:"string"`
}

func NewFullFee(fee storage.Fee) FullFee {
	ff := FullFee{
		Time:   fee.Time,
		Height: fee.Height,
		Asset:  fee.Asset,
		Amount: fee.Amount.String(),
	}

	if fee.Payer != nil {
		ff.Payer = fee.Payer.Hash
	}

	if fee.Tx != nil {
		ff.TxHash = hex.EncodeToString(fee.Tx.Hash)
	}

	return ff
}