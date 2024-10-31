// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlockStats interface {
	ByHeight(ctx context.Context, height pkgTypes.Level) (stats BlockStats, err error)
}

type BlockStats struct {
	bun.BaseModel `bun:"table:block_stats" comment:"Table with block stats"`

	Id     uint64         `bun:",pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height pkgTypes.Level `bun:"height"                    comment:"The number (height) of this block"`
	Time   time.Time      `bun:"time,pk,notnull"           comment:"The time of block"`

	TxCount      int64           `bun:"tx_count"         comment:"Count of transactions in block"`
	BlockTime    uint64          `bun:"block_time"       comment:"Time in milliseconds between current and previous block"`
	SupplyChange decimal.Decimal `bun:",type:numeric"    comment:"Change of total supply in the block"`
	Fee          decimal.Decimal `bun:"fee,type:numeric" comment:"Summary block fee"`
	BytesInBlock int64           `bun:"bytes_in_block"   comment:"Size of all transactions in bytes"`
	DataSize     int64           `bun:"data_size"        comment:"Size of all rollup data in block"`
}

func (BlockStats) TableName() string {
	return "block_stats"
}
