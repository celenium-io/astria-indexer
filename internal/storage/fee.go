// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"time"

	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type IFee interface {
	storage.Table[*Fee]
}

type Fee struct {
	bun.BaseModel `bun:"table:fee" comment:"Table with fees"`

	Id       uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height   pkgTypes.Level  `bun:"height,notnull"              comment:"The number (height) of this block"`
	Time     time.Time       `bun:"time,pk,notnull"             comment:"The time of block"`
	Asset    string          `bun:"asset"                       comment:"Fee asset"`
	Amount   decimal.Decimal `bun:"amount,type:numeric"         comment:"Fee amount"`
	ActionId uint64          `bun:"action_id"                   comment:"Connected action id"`
	TxId     uint64          `bun:"tx_id"                       comment:"Connected transaction id"`
	PayerId  uint64          `bun:"payer_id"                    comment:"Who paid fee"`

	ActionType string   `bun:"-"`
	Payer      *Address `bun:"rel:belongs-to"`
}

func (Fee) TableName() string {
	return "fee"
}
