// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

type Market struct {
	db *database.Bun
}

func NewMarket(conn *postgres.Storage) *Market {
	return &Market{
		db: conn.Connection(),
	}
}
func (m *Market) List(ctx context.Context, limit, offset int) (markets []storage.Market, err error) {
	priceQuery := m.db.DB().NewSelect().
		Model((*storage.Price)(nil)).
		ColumnExpr("max(time) as time").
		ColumnExpr("last(price, time) as price").
		Column("currency_pair").
		Group("currency_pair")

	query := m.db.DB().NewSelect().
		Model((*storage.Market)(nil))

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = m.db.DB().NewSelect().
		With("prices", priceQuery).
		TableExpr("(?) as market", query).
		ColumnExpr("market.*, (p.price/pow(10, market.decimals)) as price__price, p.time as price__time").
		Join("left join prices p on p.currency_pair = market.pair").
		Scan(ctx, &markets)
	return
}
func (m *Market) Get(ctx context.Context, pair string) (market storage.Market, err error) {
	priceQuery := m.db.DB().NewSelect().
		Model((*storage.Price)(nil)).
		ColumnExpr("max(time) as time").
		ColumnExpr("last(price, time) as price").
		Column("currency_pair").
		Group("currency_pair")

	query := m.db.DB().NewSelect().
		Model((*storage.Market)(nil)).
		Where("pair = ?", pair)

	if err = m.db.DB().NewSelect().
		With("prices", priceQuery).
		TableExpr("(?) as market", query).
		ColumnExpr("market.*, (p.price/pow(10, market.decimals)) as price__price, p.time as price__time").
		Join("left join prices p on p.currency_pair = market.pair").
		Scan(ctx, &market); err != nil {
		return
	}

	err = m.db.DB().NewSelect().
		Model(&market.Providers).
		Column("provider", "off_chain_ticker").
		Where("pair = ?", pair).
		Scan(ctx)
	return
}

func (m *Market) Decimals(ctx context.Context, pair string) (decimals int, err error) {
	err = m.db.DB().NewSelect().
		Model((*storage.Market)(nil)).
		Column("decimals").
		Where("pair = ?", pair).
		Limit(1).
		Scan(ctx, &decimals)
	return
}
