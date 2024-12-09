// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

type Asset struct {
	db *database.Bun
}

// NewAsset -
func NewAsset(db *database.Bun) *Asset {
	return &Asset{
		db: db,
	}
}

func (a *Asset) List(ctx context.Context, limit int, offset int) (assets []storage.Asset, err error) {
	transferredQuery := a.db.DB().NewSelect().
		Model((*storage.Transfer)(nil)).
		ColumnExpr("asset, count(*) as c, sum(amount) as amount").
		Group("asset")

	feesQuery := a.db.DB().NewSelect().
		Model((*storage.Fee)(nil)).
		ColumnExpr("asset, count(*) as c, sum(amount) as amount").
		Group("asset")

	query := a.db.DB().NewSelect().
		With("fees", feesQuery).
		With("transferred", transferredQuery).
		Table("fees").
		ColumnExpr("(case when fees.asset is NULL then transferred.asset else fees.asset end) as asset").
		ColumnExpr("(case when fees.amount is NULL then 0 else fees.amount end) as fee").
		ColumnExpr("(case when transferred.amount is NULL then 0 else transferred.amount end) as transferred").
		ColumnExpr("fees.c as fee_count, transferred.c as transfer_count").
		Join("full outer join transferred on transferred.asset = fees.asset")

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = query.Scan(ctx, &assets)
	return
}
