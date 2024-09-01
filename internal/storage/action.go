// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"

	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IAction interface {
	storage.Table[*Action]

	ByTxId(ctx context.Context, txId uint64, limit, offset int) ([]Action, error)
	ByBlock(ctx context.Context, height pkgTypes.Level, limit, offset int) ([]ActionWithTx, error)
	ByAddress(ctx context.Context, addressId uint64, filters AddressActionsFilter) ([]AddressAction, error)
	ByRollup(ctx context.Context, rollupId uint64, limit, offset int, sort storage.SortOrder) ([]RollupAction, error)
	ByRollupAndBridge(ctx context.Context, rollupId uint64, fltrs RollupAndBridgeActionsFilter) ([]ActionWithTx, error)
}

type AddressActionsFilter struct {
	Limit       int
	Offset      int
	Sort        storage.SortOrder
	ActionTypes types.ActionTypeMask
}

type RollupAndBridgeActionsFilter struct {
	Limit         int
	Offset        int
	Sort          storage.SortOrder
	RollupActions bool
	BridgeActions bool
}

type ActionWithTx struct {
	bun.BaseModel `bun:"action"`

	Action
	Tx *Tx `bun:"rel:belongs-to"`
}

// Action -
type Action struct {
	bun.BaseModel `bun:"action" comment:"Table with actions"`

	Id       uint64           `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Height   pkgTypes.Level   `bun:",notnull"                    comment:"The number (height) of this block"`
	Time     time.Time        `bun:"time,pk,notnull"             comment:"The time of block"`
	Position int64            `bun:"position"                    comment:"Position in transaction"`
	Type     types.ActionType `bun:",type:action_type"           comment:"Action type"`
	TxId     uint64           `bun:"tx_id"                       comment:"Parent transaction id"`
	Data     map[string]any   `bun:"data,type:jsonb"             comment:"Action data"`

	// Rollup         *Rollup          `bun:"-"`
	Addresses      []*AddressAction `bun:"-"`
	BalanceUpdates []BalanceUpdate  `bun:"-"`
	RollupAction   *RollupAction    `bun:"-"`
	Fee            *Fee             `bun:"rel:has-one,join:id=action_id"`
}

// TableName -
func (Action) TableName() string {
	return "action"
}
