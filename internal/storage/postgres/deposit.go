// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Fee -
type Deposit struct {
	*postgres.Table[*storage.Deposit]
}

// NewDeposit -
func NewDeposit(db *postgres.Storage) *Deposit {
	return &Deposit{
		Table: postgres.NewTable[*storage.Deposit](db.Connection()),
	}
}

func (d *Deposit) ByBridgeId(ctx context.Context, bridgeId uint64, limit, offset int, sort sdk.SortOrder) (deposits []storage.Deposit, err error) {
	query := d.DB().NewSelect().
		Model((*storage.Deposit)(nil)).
		Where("bridge_id = ?", bridgeId)

	query = limitScope(query, limit)
	query = offsetScope(query, offset)
	query = sortScope(query, "time", sort)

	err = d.DB().NewSelect().
		TableExpr("(?) as deposit", query).
		ColumnExpr("deposit.*").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join tx on tx.id = tx_id").
		Scan(ctx, &deposits)
	return
}

func (d *Deposit) ByRollupId(ctx context.Context, rollupId uint64, limit, offset int, sort sdk.SortOrder) (deposits []storage.Deposit, err error) {
	query := d.DB().NewSelect().
		Model((*storage.Deposit)(nil)).
		Where("rollup_id = ?", rollupId)

	query = limitScope(query, limit)
	query = offsetScope(query, offset)
	query = sortScope(query, "time", sort)

	err = d.DB().NewSelect().
		TableExpr("(?) as deposit", query).
		ColumnExpr("deposit.*").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("bridge.address_id as bridge__address_id").
		ColumnExpr("address.hash as bridge__address__hash").
		Join("left join tx on tx.id = tx_id").
		Join("left join bridge on bridge_id = bridge.id").
		Join("left join address on address_id = address.id").
		Scan(ctx, &deposits)

	return
}
