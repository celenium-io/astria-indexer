// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Rollup -
type Rollup struct {
	*postgres.Table[*storage.Rollup]
}

// NewRollup -
func NewRollup(db *database.Bun) *Rollup {
	return &Rollup{
		Table: postgres.NewTable[*storage.Rollup](db),
	}
}

func (r *Rollup) ActionsByHeight(ctx context.Context, height types.Level, limit, offset int) (actions []storage.RollupAction, err error) {
	query := r.DB().NewSelect().Model(&actions).
		Where("rollup_action.height = ?", height).
		Relation("Rollup").
		Relation("Action")

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = query.Scan(ctx)
	return
}

func (r *Rollup) CountActionsByHeight(ctx context.Context, height types.Level) (int64, error) {
	count, err := r.DB().NewSelect().Model((*storage.RollupAction)(nil)).
		Where("height = ?", height).
		Count(ctx)
	return int64(count), err
}

func (r *Rollup) ActionsByTxId(ctx context.Context, txId uint64, limit, offset int) (actions []storage.RollupAction, err error) {
	query := r.DB().NewSelect().Model(&actions).
		Where("rollup_action.tx_id = ?", txId).
		Relation("Rollup").
		Relation("Action")

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = query.Scan(ctx)
	return
}

func (r *Rollup) CountActionsByTxId(ctx context.Context, txId uint64) (int64, error) {
	count, err := r.DB().NewSelect().Model((*storage.RollupAction)(nil)).
		Where("tx_id = ?", txId).
		Count(ctx)
	return int64(count), err
}

func (r *Rollup) ByHash(ctx context.Context, hash []byte) (rollup storage.Rollup, err error) {
	err = r.DB().NewSelect().Model(&rollup).
		Where("astria_id = ?", hash).
		Scan(ctx)
	return
}

func (r *Rollup) Addresses(ctx context.Context, rollupId uint64, limit, offset int, sort sdk.SortOrder) (addresses []storage.RollupAddress, err error) {
	query := r.DB().NewSelect().Model(&addresses).
		Where("rollup_id = ?", rollupId).
		Relation("Address")

	query = limitScope(query, limit)
	query = sortScope(query, "address_id", sort)
	query = offsetScope(query, offset)

	err = query.Scan(ctx)
	return
}

func (r *Rollup) ListRollupsByAddress(ctx context.Context, addressId uint64, limit, offset int, sort sdk.SortOrder) (addresses []storage.RollupAddress, err error) {
	query := r.DB().NewSelect().Model(&addresses).
		Where("address_id = ?", addressId).
		Relation("Rollup")

	query = limitScope(query, limit)
	query = sortScope(query, "rollup_id", sort)
	query = offsetScope(query, offset)

	err = query.Scan(ctx)
	return
}
