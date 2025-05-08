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
type IPrice interface {
	storage.Table[*Price]

	Series(ctx context.Context, currencyPair string, timeframe Timeframe, fltrs SeriesRequest) ([]Candle, error)
	ByHeight(ctx context.Context, height pkgTypes.Level, limit, offset int) ([]Price, error)
}

type Price struct {
	bun.BaseModel `bun:"price" comment:"Table with currency prices"`

	CurrencyPair string          `bun:"currency_pair,pk,notnull" comment:"Currency pair"`
	Time         time.Time       `bun:"time,pk,notnull"          comment:"The time of price event"`
	Price        decimal.Decimal `bun:"price,type:numeric"       comment:"Price of the asset"`
	Height       pkgTypes.Level  `bun:"height"                   comment:"Block height"`
}

// TableName -
func (Price) TableName() string {
	return "price"
}
