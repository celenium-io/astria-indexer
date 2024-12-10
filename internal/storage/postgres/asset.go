// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
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

var validSortFieldsForAssetList = map[string]struct{}{
	"fee":            {},
	"transferred":    {},
	"transfer_count": {},
	"fee_count":      {},
	"supply":         {},
}

func (a *Asset) List(ctx context.Context, limit int, offset int, sortBy string, order sdk.SortOrder) (assets []storage.Asset, err error) {
	transferredQuery := a.db.DB().NewSelect().
		Model((*storage.Transfer)(nil)).
		ColumnExpr("asset, count(*) as c, sum(amount) as amount").
		Group("asset")

	feesQuery := a.db.DB().NewSelect().
		Model((*storage.Fee)(nil)).
		ColumnExpr("asset, count(*) as c, sum(amount) as amount").
		Group("asset")

	supplyQuery := a.db.DB().NewSelect().
		Model((*storage.Balance)(nil)).
		ColumnExpr("currency, sum(total) as amount").
		Group("currency")

	query := a.db.DB().NewSelect().
		With("fees", feesQuery).
		With("transferred", transferredQuery).
		With("supply", supplyQuery).
		Table("supply").
		ColumnExpr("(case when fees.asset is not NULL then fees.asset when supply.currency is not NULL then supply.currency else transferred.asset end) as asset").
		ColumnExpr("(case when fees.amount is NULL then 0 else fees.amount end) as fee").
		ColumnExpr("(case when transferred.amount is NULL then 0 else transferred.amount end) as transferred").
		ColumnExpr("(case when supply.amount is NULL then 0 else supply.amount end) as supply").
		ColumnExpr("(case when fees.c is NULL then 0 else fees.c end) as fee_count").
		ColumnExpr("(case when transferred.c is NULL then 0 else transferred.c end) as transfer_count").
		Join("left join transferred on supply.currency = transferred.asset").
		Join("left join fees on supply.currency = fees.asset")

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	if _, ok := validSortFieldsForAssetList[sortBy]; ok {
		query = sortScope(query, sortBy, order)
	} else {
		query = sortScope(query, "supply", order)
	}

	err = query.Scan(ctx, &assets)
	return
}
