// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
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
func (p *Price) Series(ctx context.Context, currencyPair string, timeframe storage.Timeframe) (prices []storage.Candle, err error) {
	query := p.DB().NewSelect()

	switch timeframe {
	case storage.TimeframeHour:
		query = query.Table(storage.ViewPriceByHour)
	case storage.TimeframeDay:
		query = query.Table(storage.ViewPriceByDay)
	default:
		return nil, errors.Errorf("invalid timeframe %s", timeframe)
	}

	err = query.
		Where("currency_pair = ?", currencyPair).
		Order("time DESC").
		Scan(ctx, &prices)
	return
}

func (p *Price) Last(ctx context.Context, currencyPair string) (price storage.Price, err error) {
	err = p.DB().NewSelect().
		Model(&price).
		Where("currency_pair = ?", currencyPair).
		Order("time DESC").
		Limit(1).
		Scan(ctx)
	return
}

func (p *Price) All(ctx context.Context, limit, offset int) (prices []storage.Price, err error) {
	query := p.DB().NewSelect().
		ColumnExpr("max(time) as time").
		ColumnExpr("last(price, time) as price").
		Column("currency_pair").
		Model(&prices).
		Group("currency_pair")

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = query.Scan(ctx)

	return
}
