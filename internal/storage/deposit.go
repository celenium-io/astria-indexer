// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IDeposit interface {
	storage.Table[*Deposit]

	ByBridgeId(ctx context.Context, bridgeId uint64, limit, offset int, sort storage.SortOrder) ([]Deposit, error)
	ByRollupId(ctx context.Context, rollupId uint64, limit, offset int, sort storage.SortOrder) ([]Deposit, error)
}

type Deposit struct {
	bun.BaseModel `bun:"deposit" comment:"Table with deposits"`

	Id                      uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Height                  pkgTypes.Level  `bun:",notnull"                    comment:"The number (height) of this block"`
	Time                    time.Time       `bun:"time,pk,notnull"             comment:"The time of block"`
	BridgeId                uint64          `bun:"bridge_id"                   comment:"Bridge id"`
	RollupId                uint64          `bun:"rollup_id"                   comment:"Rollup id"`
	Asset                   string          `bun:"asset"                       comment:"Deposit asset"`
	Amount                  decimal.Decimal `bun:"amount,type:numeric"         comment:"Deposit amount"`
	DestinationChainAddress string          `bun:"destination_chain_address"   comment:"Destination chain address"`
	ActionId                uint64          `bun:"action_id"                   comment:"Internal action id"`
	TxId                    uint64          `bun:"tx_id"                       comment:"Internal transaction id"`

	Bridge *Bridge `bun:"rel:belongs-to"`
	Rollup *Rollup `bun:"rel:belongs-to"`
	Action *Action `bun:"rel:belongs-to"`
	Tx     *Tx     `bun:"rel:belongs-to"`
}

func (*Deposit) TableName() string {
	return "deposit"
}
