// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Blocks -
type Blocks struct {
	*postgres.Table[*storage.Block]
}

// NewBlocks -
func NewBlocks(db *postgres.Storage) *Blocks {
	return &Blocks{
		Table: postgres.NewTable[*storage.Block](db.Connection()),
	}
}

// ByHeight -
func (b *Blocks) ByHeight(ctx context.Context, height types.Level, withStats bool) (block storage.Block, err error) {
	subQuery := b.DB().NewSelect().
		Model((*storage.Block)(nil)).
		Where("height = ?", height).
		Limit(1)

	query := b.DB().NewSelect().
		TableExpr("(?) as block", subQuery).
		ColumnExpr("block.*").
		ColumnExpr("validator.id as proposer__id, validator.address as proposer__address, validator.name as proposer__name").
		Join("left join validator on block.proposer_id = validator.id")

	if withStats {
		query = query.
			ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count").
			ColumnExpr("stats.block_time AS stats__block_time, stats.bytes_in_block AS stats__bytes_in_block").
			ColumnExpr("stats.supply_change AS stats__supply_change").
			Join("left join block_stats as stats ON stats.height = block.height AND stats.time = block.time")
	}

	err = query.Scan(ctx, &block)
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
	subQuery := b.DB().NewSelect().
		Model((*storage.Block)(nil)).
		Where("hash = ?", hash).
		Limit(1)

	err = b.DB().NewSelect().
		TableExpr("(?) as block", subQuery).
		ColumnExpr("block.*").
		ColumnExpr("validator.id as proposer__id, validator.address as proposer__address, validator.name as proposer__name").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count").
		ColumnExpr("stats.block_time AS stats__block_time, stats.bytes_in_block AS stats__bytes_in_block").
		ColumnExpr("stats.supply_change AS stats__supply_change").
		Join("left join validator on block.proposer_id = validator.id").
		Join("left join block_stats as stats ON stats.height = block.height AND stats.time = block.time").
		Scan(ctx, &block)
	return
}

// ListWithStats -
func (b *Blocks) ListWithStats(ctx context.Context, limit, offset uint64, order sdk.SortOrder) (blocks []*storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&blocks)
	subQuery = sortScope(subQuery, "block.time", order)
	subQuery = limitScope(subQuery, int(limit))
	subQuery = offsetScope(subQuery, int(offset))

	query := b.DB().NewSelect().
		ColumnExpr("block.*").
		ColumnExpr("v.id AS proposer__id, v.address as proposer__address, v.name as proposer__name").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count").
		ColumnExpr("stats.block_time AS stats__block_time, stats.bytes_in_block AS stats__bytes_in_block").
		ColumnExpr("stats.supply_change AS stats__supply_change").
		TableExpr("(?) as block", subQuery).
		Join("LEFT JOIN block_stats as stats ON stats.height = block.height AND stats.time = block.time").
		Join("LEFT JOIN validator as v ON v.id = block.proposer_id")
	query = sortScope(query, "block.time", order)
	err = query.Scan(ctx, &blocks)

	return
}

func (b *Blocks) ByProposer(ctx context.Context, proposerId uint64, limit, offset int, order sdk.SortOrder) (blocks []storage.Block, err error) {
	subQuery := b.DB().NewSelect().Model(&blocks).
		Where("proposer_id = ?", proposerId)

	subQuery = sortScope(subQuery, "id", order)
	subQuery = limitScope(subQuery, limit)
	subQuery = offsetScope(subQuery, offset)

	query := b.DB().NewSelect().
		ColumnExpr("block.*").
		ColumnExpr("v.id AS proposer__id, v.address as proposer__address, v.name as proposer__name").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count").
		ColumnExpr("stats.block_time AS stats__block_time, stats.bytes_in_block AS stats__bytes_in_block").
		ColumnExpr("stats.supply_change AS stats__supply_change").
		TableExpr("(?) as block", subQuery).
		Join("LEFT JOIN block_stats as stats ON stats.height = block.height AND stats.time = block.time").
		Join("LEFT JOIN validator as v ON v.id = block.proposer_id")
	query = sortScope(query, "block.id", order)
	err = query.Scan(ctx, &blocks)
	return
}

func (b *Blocks) ByIdWithRelations(ctx context.Context, id uint64) (block storage.Block, err error) {
	query := b.DB().NewSelect().
		Model((*storage.Block)(nil)).
		Where("id = ?", id).
		Limit(1)

	err = b.DB().NewSelect().
		TableExpr("(?) as block", query).
		ColumnExpr("block.*").
		ColumnExpr("validator.id as proposer__id, validator.address as proposer__address, validator.name as proposer__name, validator.pubkey as proposer__pubkey, validator.pubkey_type as proposer__pubkey_type, validator.power as proposer__power, validator.height as proposer__height").
		ColumnExpr("stats.id AS stats__id, stats.height AS stats__height, stats.time AS stats__time, stats.tx_count AS stats__tx_count").
		ColumnExpr("stats.block_time AS stats__block_time, stats.bytes_in_block AS stats__bytes_in_block").
		ColumnExpr("stats.supply_change AS stats__supply_change").
		Join("left join validator on block.proposer_id = validator.id").
		Join("left join block_stats as stats ON stats.height = block.height AND stats.time = block.time").
		Scan(ctx, &block)

	return
}
