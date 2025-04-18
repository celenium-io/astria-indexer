// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
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
	TxHash string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty" swaggertype:"string"`
	Height pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"            swaggertype:"integer"`
	Time   time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"              swaggertype:"string"`

	Payer *ShortAddress `json:"payer,omitempty"`
}

func NewFullFee(fee storage.Fee) FullFee {
	ff := FullFee{
		Time:   fee.Time,
		Height: fee.Height,
		Asset:  fee.Asset,
		Amount: fee.Amount.String(),
		Payer:  NewShortAddress(fee.Payer),
	}

	if fee.Tx != nil {
		ff.TxHash = hex.EncodeToString(fee.Tx.Hash)
	}

	return ff
}

type TxFee struct {
	Amount string `example:"1000" format:"string" json:"amount" swaggertype:"string"`
	Asset  string `example:"nria" format:"string" json:"asset"  swaggertype:"string"`
}

func NewTxFee(fee storage.Fee) TxFee {
	return TxFee{
		Asset:  fee.Asset,
		Amount: fee.Amount.String(),
	}
}
