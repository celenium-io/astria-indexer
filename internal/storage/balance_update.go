// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBalanceUpdate interface {
	storage.Table[*BalanceUpdate]
}

type BalanceUpdate struct {
	bun.BaseModel `bun:"balance_update" comment:"Table with account balance updates"`

	Id        uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height    pkgTypes.Level  `bun:",notnull"                    comment:"The number (height) of this block"`
	AddressId uint64          `bun:"address_id"                  comment:"Address internal identity"`
	Update    decimal.Decimal `bun:"update,type:numeric"         comment:"Balance update"`
	Currency  string          `bun:"currency"                    comment:"Currency"`

	Address *Address `bun:"rel:belongs-to"`
}

func (BalanceUpdate) TableName() string {
	return "balance_update"
}
