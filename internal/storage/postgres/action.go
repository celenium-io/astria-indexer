// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/aopoltorzhicky/astria/internal/storage"
	"github.com/aopoltorzhicky/astria/pkg/types"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Action -
type Action struct {
	*postgres.Table[*storage.Action]
}

// NewAction -
func NewAction(db *database.Bun) *Action {
	return &Action{
		Table: postgres.NewTable[*storage.Action](db),
	}
}

func (a *Action) ByBlock(ctx context.Context, height types.Level, limit, offset int) (actions []storage.ActionWithTx, err error) {
	query := a.DB().NewSelect().Model(&actions).
		Where("action_with_tx.height = ?", height).
		Relation("Tx", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("hash")
		})

	query = limitScope(query, limit)

	if offset > 0 {
		query = query.Offset(offset)
	}
	err = query.Scan(ctx)
	return
}

func (a *Action) ByTxId(ctx context.Context, txId uint64, limit, offset int) (actions []storage.Action, err error) {
	query := a.DB().NewSelect().Model(&actions).
		Where("tx_id = ?", txId)
	query = limitScope(query, limit)

	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx)
	return
}

func (a *Action) ByAddress(ctx context.Context, addressId uint64, filters storage.AddressActionsFilter) (actions []storage.AddressAction, err error) {
	query := a.DB().NewSelect().Model(&actions).
		Where("address_id = ?", addressId).
		Relation("Action").
		Relation("Tx", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("hash")
		})

	if filters.ActionTypes.Bits > 0 {
		query = query.Where("action_type IN (?)", bun.In(filters.ActionTypes.Strings()))
	}

	query = sortScope(query, "action_id", filters.Sort)
	query = limitScope(query, filters.Limit)

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err = query.Scan(ctx)
	return
}

func (a *Action) ByRollup(ctx context.Context, rollupId uint64, limit, offset int, sort sdk.SortOrder) (actions []storage.RollupAction, err error) {
	query := a.DB().NewSelect().Model(&actions).
		Where("rollup_id = ?", rollupId).
		Relation("Action").
		Relation("Tx", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("hash")
		})

	query = sortScope(query, "action_id", sort)
	query = limitScope(query, limit)

	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx)
	return
}
