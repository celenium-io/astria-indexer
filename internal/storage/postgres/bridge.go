// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Bridge -
type Bridge struct {
	*postgres.Table[*storage.Bridge]
}

// NewBridge -
func NewBridge(db *postgres.Storage) *Bridge {
	return &Bridge{
		Table: postgres.NewTable[*storage.Bridge](db.Connection()),
	}
}

func (b *Bridge) ByAddress(ctx context.Context, addressId uint64) (bridge storage.Bridge, err error) {
	query := b.DB().NewSelect().
		Model((*storage.Bridge)(nil)).
		Where("address_id = ?", addressId)

	err = b.DB().NewSelect().
		TableExpr("(?) as bridge", query).
		ColumnExpr("bridge.*").
		ColumnExpr("address.hash as address__hash").
		ColumnExpr("sudo.hash as sudo__hash").
		ColumnExpr("withdrawer.hash as withdrawer__hash").
		ColumnExpr("rollup.astria_id as rollup__astria_id").
		Join("left join address as address on address.id = bridge.address_id").
		Join("left join address as sudo on sudo.id = bridge.sudo_id").
		Join("left join address as withdrawer on withdrawer.id = bridge.withdrawer_id").
		Join("left join rollup on rollup.id = bridge.rollup_id").
		Scan(ctx, &bridge)
	return
}

func (b *Bridge) ByRollup(ctx context.Context, rollupId uint64, limit, offset int) (bridge []storage.Bridge, err error) {
	query := b.DB().NewSelect().
		Model((*storage.Bridge)(nil)).
		Where("rollup_id = ?", rollupId).
		Offset(offset)

	query = limitScope(query, limit)

	err = b.DB().NewSelect().
		TableExpr("(?) as bridge", query).
		ColumnExpr("bridge.*").
		ColumnExpr("address.hash as address__hash").
		ColumnExpr("sudo.hash as sudo__hash").
		ColumnExpr("withdrawer.hash as withdrawer__hash").
		ColumnExpr("rollup.astria_id as rollup__astria_id").
		Join("left join address as address on address.id = bridge.address_id").
		Join("left join address as sudo on sudo.id = bridge.sudo_id").
		Join("left join address as withdrawer on withdrawer.id = bridge.withdrawer_id").
		Join("left join rollup on rollup.id = bridge.rollup_id").
		Scan(ctx, &bridge)
	return
}

func (b *Bridge) ByRoles(ctx context.Context, addressId uint64, limit, offset int) (result []storage.Bridge, err error) {
	query := b.DB().NewSelect().
		Model((*storage.Bridge)(nil)).
		Where("sudo_id = ?", addressId).
		WhereOr("withdrawer_id = ?", addressId).
		Offset(offset)

	query = limitScope(query, limit)

	err = b.DB().NewSelect().
		TableExpr("(?) as bridge", query).
		ColumnExpr("bridge.*").
		ColumnExpr("address.hash as address__hash").
		ColumnExpr("sudo.hash as sudo__hash").
		ColumnExpr("withdrawer.hash as withdrawer__hash").
		ColumnExpr("rollup.astria_id as rollup__astria_id").
		Join("left join address as address on address.id = bridge.address_id").
		Join("left join address as sudo on sudo.id = bridge.sudo_id").
		Join("left join address as withdrawer on withdrawer.id = bridge.withdrawer_id").
		Join("left join rollup on rollup.id = bridge.rollup_id").
		Scan(ctx, &result)
	return
}

func (b *Bridge) ListWithAddress(ctx context.Context, limit, offset int) (result []storage.Bridge, err error) {
	query := b.DB().NewSelect().
		Model((*storage.Bridge)(nil)).
		Offset(offset)

	query = limitScope(query, limit)
	query = sortScope(query, "id", sdk.SortOrderAsc)

	err = b.DB().NewSelect().
		TableExpr("(?) as bridge", query).
		ColumnExpr("bridge.*").
		ColumnExpr("address.hash as address__hash").
		Join("left join address as address on address.id = bridge.address_id").
		Scan(ctx, &result)
	return
}

func (b *Bridge) ById(ctx context.Context, id uint64) (bridge storage.Bridge, err error) {
	query := b.DB().NewSelect().
		Model((*storage.Bridge)(nil)).
		Where("id = ?", id)

	err = b.DB().NewSelect().
		TableExpr("(?) as bridge", query).
		ColumnExpr("bridge.*").
		ColumnExpr("address.hash as address__hash").
		ColumnExpr("rollup.astria_id as rollup__astria_id").
		Join("left join address as address on address.id = bridge.address_id").
		Join("left join rollup on rollup.id = bridge.rollup_id").
		Scan(ctx, &bridge)

	return
}
