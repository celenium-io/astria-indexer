// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IMarket interface {
	List(ctx context.Context, limit, offset int) ([]Market, error)
	Get(ctx context.Context, pair string) (Market, error)
	Decimals(ctx context.Context, pair string) (int, error)
	History(ctx context.Context, pair string, limit, offset int) ([]Market, error)
}
type Market struct {
	bun.BaseModel `bun:"table:market" comment:"Table with markets"`

	Pair             string    `bun:"pair,pk"               comment:"Trading pair"`
	UpdatedAt        time.Time `bun:"updated_at,pk"         comment:"Time when market settings was updated"`
	Decimals         int       `bun:"decimals"              comment:"Decimals"`
	Enabled          bool      `bun:"enabled"               comment:"Is market enabled"`
	MinProviderCount int       `bun:"min_provider_count"    comment:"Minimum provider count"`
	Base             string    `bun:"base"                  comment:"Base asset"`
	Quote            string    `bun:"quote"                 comment:"Quote asset"`
	Removed          bool      `bun:"removed,default:false" comment:"True if market was removed"`

	Price     *Price            `bun:"rel:has-one,join:pair=currency_pair"`
	Providers []*MarketProvider `bun:"rel:has-many,join:pair=pair"`
}

func (Market) TableName() string {
	return "market"
}

func (m Market) String() string {
	var sb = new(strings.Builder)
	sb.WriteString(m.Pair)
	fmt.Fprintf(sb, " was updated at %s", m.UpdatedAt.String())
	fmt.Fprintf(sb, "\r\n\tdecimals: %d | enabled: %v | providers count: %d | removed: %v\r\n", m.Decimals, m.Enabled, m.MinProviderCount, m.Removed)
	return sb.String()
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
