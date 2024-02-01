// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/hex"

	"github.com/aopoltorzhicky/astria/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IRollup interface {
	storage.Table[*Rollup]

	ActionsByHeight(ctx context.Context, height types.Level, limit, offset int) ([]RollupAction, error)
	CountActionsByHeight(ctx context.Context, height types.Level) (int64, error)
	ActionsByTxId(ctx context.Context, txId uint64, limit, offset int) ([]RollupAction, error)
	CountActionsByTxId(ctx context.Context, txId uint64) (int64, error)
	ByHash(ctx context.Context, hash []byte) (Rollup, error)
	Addresses(ctx context.Context, rollupId uint64, limit, offset int, sort sdk.SortOrder) ([]RollupAddress, error)
	ListRollupsByAddress(ctx context.Context, addressId uint64, limit, offset int, sort sdk.SortOrder) ([]RollupAddress, error)
}

type Rollup struct {
	bun.BaseModel `bun:"rollup" comment:"Table with rollups"`

	Id           uint64      `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	AstriaId     []byte      `bun:"astria_id,unique:rollup_id"  comment:"Astria rollup identity"`
	FirstHeight  types.Level `bun:"first_height"                comment:"Block number of the first rollup occurance"`
	ActionsCount int64       `bun:"actions_count"               comment:"Count of actions in which the rollup was involved"`
	Size         int64       `bun:"size"                        comment:"Count bytes which was saved in the rollup"`
}

// TableName -
func (Rollup) TableName() string {
	return "rollup"
}

func (r Rollup) String() string {
	return hex.EncodeToString(r.AstriaId)
}
