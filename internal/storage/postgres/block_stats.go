// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/database"
)

// BlockStats -
type BlockStats struct {
	db *database.Bun
}

// NewBlockStats -
func NewBlockStats(db *database.Bun) *BlockStats {
	return &BlockStats{
		db: db,
	}
}

// ByHeight -
func (b *BlockStats) ByHeight(ctx context.Context, height types.Level) (stats storage.BlockStats, err error) {
	err = b.db.DB().NewSelect().Model(&stats).
		Where("height = ?", height).
		Limit(1).
		Scan(ctx)

	return
}
