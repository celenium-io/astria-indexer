// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IMarket interface {
	List(ctx context.Context, limit, offset int) ([]Market, error)
	Get(ctx context.Context, pair string) (Market, error)
	Decimals(ctx context.Context, pair string) (int, error)
}
type Market struct {
	bun.BaseModel `bun:"table:market" comment:"Table with markets"`

	Pair             string `bun:"pair,pk"            comment:"Trading pair"`
	Decimals         int    `bun:"decimals"           comment:"Decimals"`
	Enabled          bool   `bun:"enabled"            comment:"Is market enabled"`
	MinProviderCount int    `bun:"min_provider_count" comment:"Minimum provider count"`
	Base             string `bun:"base"               comment:"Base asset"`
	Quote            string `bun:"quote"              comment:"Quote asset"`

	Price *Price `bun:"rel:has-one,join:pair=currency_pair"`
}

func (Market) TableName() string {
	return "market"
}

type MarketUpdate struct {
	Market

	Type MarketUpdateType
}

type MarketUpdateType string

const (
	MarketUpdateTypeCreate MarketUpdateType = "create"
	MarketUpdateTypeRemove MarketUpdateType = "remove"
	MarketUpdateTypeUpdate MarketUpdateType = "update"
)
