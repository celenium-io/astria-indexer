// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/hex"

	"github.com/celenium-io/astria-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IRollup interface {
	sdk.Table[*Rollup]

	ActionsByHeight(ctx context.Context, height types.Level, limit, offset int) ([]RollupAction, error)
	CountActionsByHeight(ctx context.Context, height types.Level) (int64, error)
	ActionsByTxId(ctx context.Context, txId uint64, limit, offset int) ([]RollupAction, error)
	CountActionsByTxId(ctx context.Context, txId uint64) (int64, error)
	ByHash(ctx context.Context, hash []byte) (Rollup, error)
	Addresses(ctx context.Context, rollupId uint64, limit, offset int, sort sdk.SortOrder) ([]RollupAddress, error)
	ListRollupsByAddress(ctx context.Context, addressId uint64, limit, offset int, sort sdk.SortOrder) ([]RollupAddress, error)
	ListExt(ctx context.Context, fltrs RollupListFilter) ([]Rollup, error)
}

type Rollup struct {
	bun.BaseModel `bun:"rollup" comment:"Table with rollups"`

	Id           uint64      `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	AstriaId     []byte      `bun:"astria_id,unique:rollup_id"  comment:"Astria rollup identity"`
	FirstHeight  types.Level `bun:"first_height"                comment:"Block number of the first rollup occurrence"`
	ActionsCount int64       `bun:"actions_count"               comment:"Count of actions in which the rollup was involved"`
	BridgeCount  int64       `bun:"bridge_count"                comment:"Count of connected bridges"`
	Size         int64       `bun:"size"                        comment:"Count bytes which was saved in the rollup"`
}

// TableName -
func (Rollup) TableName() string {
	return "rollup"
}

func (r Rollup) String() string {
	return hex.EncodeToString(r.AstriaId)
}

type RollupListFilter struct {
	Limit     int
	Offset    int
	SortField string
	SortOrder sdk.SortOrder
}
