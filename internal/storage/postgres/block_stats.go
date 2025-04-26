// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// BlockStats -
type BlockStats struct {
	db *postgres.Storage
}

// NewBlockStats -
func NewBlockStats(db *postgres.Storage) *BlockStats {
	return &BlockStats{
		db: db,
	}
}

// ByHeight -
func (b *BlockStats) ByHeight(ctx context.Context, height types.Level) (stats storage.BlockStats, err error) {
	err = b.db.Connection().DB().NewSelect().Model(&stats).
		Where("height = ?", height).
		Limit(1).
		Scan(ctx)

	return
}
