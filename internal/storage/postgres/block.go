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

// Blocks -
type Blocks struct {
	*postgres.Table[*storage.Block]
}

// NewBlocks -
func NewBlocks(db *database.Bun) *Blocks {
	return &Blocks{
		Table: postgres.NewTable[*storage.Block](db),
	}
}

// ByHeight -
func (b *Blocks) ByHeight(ctx context.Context, height types.Level, withStats bool) (block storage.Block, err error) {
	query := b.DB().NewSelect().Model(&block).
		Where("block.height = ?", height).
		Relation("Proposer", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("id", "address", "name")
		}).
		Limit(1)

	if withStats {
		query = query.Relation("Stats")
	}

	err = query.Scan(ctx)
	return
}

// Last -
func (b *Blocks) Last(ctx context.Context) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).
		Relation("Proposer", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("id", "address", "name")
		}).
		Order("id desc").
		Limit(1).
		Scan(ctx)
	return
}

// ByHash -
func (b *Blocks) ByHash(ctx context.Context, hash []byte) (block storage.Block, err error) {
	err = b.DB().NewSelect().
		Model(&block).
		Where("hash = ?", hash).
		Relation("Stats").
		Relation("Proposer", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column("id", "address", "name")
		}).
		Limit(1).
		Scan(ctx)
	return
}

// ListWithStats -
func (b *Blocks) ListWithStats(ctx context.Context, limit, offset uint64, order sdk.SortOrder) (blocks []*storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&blocks)
	subQuery = postgres.Pagination(subQuery, limit, offset, order)

	query := b.DB().NewSelect().
		ColumnExpr("block.*").
		ColumnExpr("v.id AS proposer__id, v.address as proposer__address, v.name as proposer__name").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count").
		ColumnExpr("stats.block_time AS stats__block_time, stats.bytes_in_block AS stats__bytes_in_block").
		ColumnExpr("stats.supply_change AS stats__supply_change, stats.fee AS stats__fee, stats.gas_used AS stats__gas_used, stats.gas_wanted AS stats__gas_wanted").
		TableExpr("(?) as block", subQuery).
		Join("LEFT JOIN block_stats as stats ON stats.height = block.height").
		Join("LEFT JOIN validator as v ON v.id = block.proposer_id")
	query = sortScope(query, "block.id", order)
	err = query.Scan(ctx, &blocks)

	return
}

func (b *Blocks) ByProposer(ctx context.Context, proposerId uint64, limit, offset int, order sdk.SortOrder) (blocks []storage.Block, err error) {
	query := b.DB().NewSelect().Model(&blocks).
		Where("proposer_id = ?", proposerId).
		Relation("Stats")

	query = sortScope(query, "id", order)
	query = limitScope(query, limit)
	query = offsetScope(query, offset)

	err = query.Scan(ctx)
	return
}

func (b *Blocks) ByIdWithRelations(ctx context.Context, id uint64) (block storage.Block, err error) {
	err = b.DB().NewSelect().Model(&block).
		Where("block.id = ?", id).
		Limit(1).
		Relation("Stats").
		Relation("Proposer").
		Scan(ctx)
	return
}
