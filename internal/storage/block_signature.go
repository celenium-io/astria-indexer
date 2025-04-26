// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlockSignature interface {
	storage.Table[*BlockSignature]

	LevelsByValidator(ctx context.Context, validatorId uint64, startHeight pkgTypes.Level) ([]pkgTypes.Level, error)
}

type BlockSignature struct {
	bun.BaseModel `bun:"block_signature" comment:"Table with block signatures"`

	Id          uint64         `bun:"id,pk,notnull,autoincrement" comment:"Unique internal id"`
	Height      pkgTypes.Level `bun:",notnull"                    comment:"The number (height) of this block"`
	Time        time.Time      `bun:"time,pk,notnull"             comment:"The time of block"`
	ValidatorId uint64         `bun:"validator_id"                comment:"Validator's internal identity"`

	Validator *Validator `bun:"rel:belongs-to"`
}

func (BlockSignature) TableName() string {
	return "block_signature"
}
