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
type ITx interface {
	storage.Table[*Tx]

	ByHash(ctx context.Context, hash []byte) (Tx, error)
	ByHeight(ctx context.Context, height pkgTypes.Level, limit, offset int) ([]Tx, error)
	ByAddress(ctx context.Context, addressId uint64, fltrs TxFilter) ([]Tx, error)
	Filter(ctx context.Context, fltrs TxFilter) ([]Tx, error)
}

type TxFilter struct {
	Limit       int
	Offset      int
	Sort        storage.SortOrder
	Status      []string
	ActionTypes types.ActionTypeMask
	Height      uint64
	TimeFrom    time.Time
	TimeTo      time.Time
	WithActions bool
}

// Tx -
type Tx struct {
	bun.BaseModel `bun:"tx" comment:"Table with transactions"`

	Id           uint64         `bun:"id,autoincrement,pk,notnull" comment:"Unique internal id"`
	Height       pkgTypes.Level `bun:",notnull"                    comment:"The number (height) of this block"`
	Time         time.Time      `bun:"time,pk,notnull"             comment:"The time of block"`
	Position     int64          `bun:"position"                    comment:"Position in block"`
	GasWanted    int64          `bun:"gas_wanted"                  comment:"Gas wanted"`
	GasUsed      int64          `bun:"gas_used"                    comment:"Gas used"`
	ActionsCount int64          `bun:"actions_count"               comment:"Actions count in transaction"`
	Status       types.Status   `bun:"status,type:status"          comment:"Transaction status"`

	Error       string     `bun:"error,type:text"     comment:"Error string if failed"`
	Codespace   string     `bun:"codespace,type:text" comment:"Codespace"`
	SignerId    uint64     `bun:"signer_id"           comment:"Signer internal identity"`
	ActionTypes types.Bits `bun:"action_types"        comment:"Bit mask for action types contained in tx"`
	Nonce       uint32     `bun:"nonce"               comment:"Nonce"`
	Hash        []byte     `bun:"hash"                comment:"Transaction hash"`
	Signature   []byte     `bun:"signature"           comment:"Signature"`

	Actions   []Action `bun:"rel:has-many,join:id=tx_id"`
	Signer    *Address `bun:"rel:belongs-to"`
	BytesSize int64    `bun:"-"`
}

// TableName -
func (Tx) TableName() string {
	return "tx"
}
