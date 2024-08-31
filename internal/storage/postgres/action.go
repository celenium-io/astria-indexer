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
	query := a.DB().NewSelect().
		Model((*storage.Action)(nil)).
		Where("height = ?", height)

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = a.DB().NewSelect().
		TableExpr("(?) as action", query).
		ColumnExpr("action.*").
		ColumnExpr("fee.asset as fee__asset, fee.amount as fee__amount").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join tx on tx.id = action.tx_id").
		Join("left join fee on fee.action_id = action.id").
		Scan(ctx, &actions)
	return
}

func (a *Action) ByTxId(ctx context.Context, txId uint64, limit, offset int) (actions []storage.Action, err error) {
	query := a.DB().NewSelect().
		Model((*storage.Action)(nil)).
		Where("tx_id = ?", txId)

	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = a.DB().NewSelect().
		TableExpr("(?) as action", query).
		ColumnExpr("fee.asset as fee__asset, fee.amount as fee__amount").
		ColumnExpr("action.*").
		Join("left join fee on fee.action_id = action.id").
		Scan(ctx, &actions)
	return
}

func (a *Action) ByAddress(ctx context.Context, addressId uint64, filters storage.AddressActionsFilter) (actions []storage.AddressAction, err error) {
	subQuery := a.DB().NewSelect().
		Model((*storage.AddressAction)(nil)).
		Where("address_id = ?", addressId)

	if filters.ActionTypes.Bits > 0 {
		subQuery = subQuery.Where("action_type IN (?)", bun.In(filters.ActionTypes.Strings()))
	}

	subQuery = sortScope(subQuery, "action_id", filters.Sort)
	subQuery = limitScope(subQuery, filters.Limit)
	subQuery = offsetScope(subQuery, filters.Offset)

	query := a.DB().NewSelect().
		TableExpr("(?) as address_action", subQuery).
		ColumnExpr("address_action.*").
		ColumnExpr("action.id as action__id, action.height as action__height, action.time as action__time, action.position as action__position, action.type as action__type, action.tx_id as action__tx_id, action.data as action__data").
		ColumnExpr("fee.asset as action__fee__asset, fee.amount as action__fee__amount").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join tx on tx.id = address_action.tx_id").
		Join("left join action on action.id = address_action.action_id").
		Join("left join fee on fee.action_id = address_action.action_id")
	query = sortScope(query, "action_id", filters.Sort)
	err = query.Scan(ctx, &actions)
	return
}

func (a *Action) ByRollup(ctx context.Context, rollupId uint64, limit, offset int, sort sdk.SortOrder) (actions []storage.RollupAction, err error) {
	subQuery := a.DB().NewSelect().
		Model((*storage.RollupAction)(nil)).
		Where("rollup_id = ?", rollupId)

	subQuery = sortScope(subQuery, "action_id", sort)
	subQuery = limitScope(subQuery, limit)
	subQuery = offsetScope(subQuery, offset)

	query := a.DB().NewSelect().
		TableExpr("(?) as rollup_action", subQuery).
		ColumnExpr("rollup_action.*").
		ColumnExpr("fee.asset as action__fee__asset, fee.amount as action__fee__amount").
		ColumnExpr("action.id as action__id, action.height as action__height, action.time as action__time, action.position as action__position, action.type as action__type, action.tx_id as action__tx_id, action.data as action__data").
		ColumnExpr("tx.hash as tx__hash").
		Join("left join tx on tx.id = rollup_action.tx_id").
		Join("left join action on action.id = rollup_action.action_id").
		Join("left join fee on fee.action_id = rollup_action.action_id")
	query = sortScope(query, "action_id", sort)
	err = query.Scan(ctx, &actions)
	return
}
