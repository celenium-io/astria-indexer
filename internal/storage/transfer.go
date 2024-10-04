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

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ITransfer interface {
	storage.Table[*Transfer]
}

type Transfer struct {
	bun.BaseModel `bun:"transfer" comment:"Table with asset transfers"`

	Id            uint64          `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height        pkgTypes.Level  `bun:"height,notnull"              comment:"The number (height) of this block"`
	Time          time.Time       `bun:"time,pk,notnull"             comment:"The time of block"`
	Asset         string          `bun:"asset"                       comment:"Transfer asset"`
	Amount        decimal.Decimal `bun:"amount,type:numeric"         comment:"Transfer amount"`
	SourceId      uint64          `bun:"src_id"                      comment:"Who made transfer"`
	DestinationId uint64          `bun:"dest_id"                     comment:"Who receive transfer"`

	Source      *Address `bun:"rel:belongs-to"`
	Destination *Address `bun:"rel:belongs-to"`
}

func (Transfer) TableName() string {
	return "transfer"
}
