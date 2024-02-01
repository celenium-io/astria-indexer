// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/aopoltorzhicky/astria/internal/storage"
	"github.com/aopoltorzhicky/astria/pkg/types"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// BlockSignature -
type BlockSignature struct {
	*postgres.Table[*storage.BlockSignature]
}

// NewBlockSignature -
func NewBlockSignature(db *database.Bun) *BlockSignature {
	return &BlockSignature{
		Table: postgres.NewTable[*storage.BlockSignature](db),
	}
}

func (bs *BlockSignature) LevelsByValidator(ctx context.Context, validatorId uint64, startHeight types.Level) (levels []types.Level, err error) {
	err = bs.DB().NewSelect().
		Model((*storage.BlockSignature)(nil)).
		Column("height").
		Where("validator_id = ?", validatorId).
		Where("height > ?", startHeight).
		Order("id desc").
		Scan(ctx, &levels)
	return
}
