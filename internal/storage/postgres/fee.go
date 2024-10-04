// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Fee -
type Fee struct {
	*postgres.Table[*storage.Fee]
}

// NewFee -
func NewFee(db *database.Bun) *Fee {
	return &Fee{
		Table: postgres.NewTable[*storage.Fee](db),
	}
}

func (f *Fee) ByTxId(ctx context.Context, id uint64, limit, offset int) (fees []storage.Fee, err error) {
	query := f.DB().NewSelect().
		Model((*storage.Fee)(nil)).
		Where("tx_id = ?", id)

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = f.DB().NewSelect().
		TableExpr("(?) as fee", query).
		ColumnExpr("fee.*").
		ColumnExpr("address.hash as payer__hash").
		Join("left join address on address.id = fee.payer_id").
		Scan(ctx, &fees)

	return
}

func (f *Fee) ByPayerId(ctx context.Context, id uint64, limit, offset int, sort sdk.SortOrder) (fees []storage.Fee, err error) {
	query := f.DB().NewSelect().
		Model((*storage.Fee)(nil)).
		Where("payer_id = ?", id)

	query = limitScope(query, limit)
	query = offsetScope(query, offset)
	query = sortScope(query, "time", sort)

	err = f.DB().NewSelect().
		TableExpr("(?) as fee", query).
		ColumnExpr("fee.*").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join tx on tx.id = fee.tx_id").
		Scan(ctx, &fees)

	return
}
