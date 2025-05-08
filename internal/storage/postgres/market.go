// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
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
		Model((*storage.Market)(nil)).
		ColumnExpr("max(updated_at) as updated_at").
		ColumnExpr("last(decimals, updated_at) as decimals").
		ColumnExpr("last(enabled, updated_at) as enabled").
		ColumnExpr("last(min_provider_count, updated_at) as min_provider_count").
		ColumnExpr("last(base, updated_at) as base").
		ColumnExpr("last(quote, updated_at) as quote").
		ColumnExpr("max(updated_at) as updated_at").
		Column("pair").
		Group("pair")

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

	pairQuery := m.db.DB().NewSelect().
		Model((*storage.Market)(nil)).
		Where("pair = ?", pair)

	query := m.db.DB().NewSelect().
		TableExpr("(?) as pairs", pairQuery).
		ColumnExpr("max(updated_at) as updated_at").
		ColumnExpr("last(decimals, updated_at) as decimals").
		ColumnExpr("last(enabled, updated_at) as enabled").
		ColumnExpr("last(min_provider_count, updated_at) as min_provider_count").
		ColumnExpr("last(base, updated_at) as base").
		ColumnExpr("last(quote, updated_at) as quote").
		ColumnExpr("max(updated_at) as updated_at").
		Column("pair").
		Group("pair")

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
		ColumnExpr("last(decimals, updated_at) as decimals").
		Where("pair = ?", pair).
		Limit(1).
		Scan(ctx, &decimals)
	return
}

func (m *Market) History(ctx context.Context, pair string, limit, offset int) (markets []storage.Market, err error) {
	query := m.db.DB().NewSelect().
		Model(&markets).
		Where("pair = ?", pair)
	query = limitScope(query, limit)
	query = offsetScope(query, offset)
	query = sortScope(query, "updated_at", sdk.SortOrderDesc)
	err = query.Scan(ctx)
	return
}
