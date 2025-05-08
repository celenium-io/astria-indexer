// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
)

type Price struct {
	*postgres.Table[*storage.Price]
}

func NewPrice(db *postgres.Storage) *Price {
	return &Price{
		Table: postgres.NewTable[*storage.Price](db.Connection()),
	}
}
func (p *Price) Series(ctx context.Context, currencyPair string, timeframe storage.Timeframe, filters storage.SeriesRequest) (prices []storage.Candle, err error) {
	query := p.DB().NewSelect()

	switch timeframe {
	case storage.TimeframeHour:
		query = query.Table(storage.ViewPriceByHour)
	case storage.TimeframeDay:
		query = query.Table(storage.ViewPriceByDay)
	default:
		return nil, errors.Errorf("invalid timeframe %s", timeframe)
	}

	if !filters.From.IsZero() {
		query = query.Where("time >= ?", filters.From)
	}
	if !filters.To.IsZero() {
		query = query.Where("time < ?", filters.To)
	}

	err = query.
		Where("currency_pair = ?", currencyPair).
		Order("time DESC").
		Scan(ctx, &prices)
	return
}

func (p *Price) ByHeight(ctx context.Context, height pkgTypes.Level, limit, offset int) (prices []storage.Price, err error) {
	query := p.DB().NewSelect().
		Model((*storage.Price)(nil)).
		Where("height = ?", height)

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = p.DB().NewSelect().
		TableExpr("(?) as price", query).
		ColumnExpr("price.currency_pair as currency_pair, price.time as time, (price.price/pow(10, coalesce(pair.decimals, 0))) as price").
		Join("left join lateral (select * from market where market.pair = price.currency_pair and market.updated_at <= price.time order by updated_at desc limit 1) pair on true").
		Scan(ctx, &prices)
	return
}
