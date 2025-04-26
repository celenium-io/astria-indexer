// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/uptrace/bun"
)

type RollupAction struct {
	bun.BaseModel `bun:"rollup_action" comment:"Table with rollup actions"`

	RollupId   uint64           `bun:"rollup_id,pk"                 comment:"Rollup internal id"`
	ActionId   uint64           `bun:"action_id,pk"                 comment:"Action internal id"`
	Time       time.Time        `bun:"time,notnull,pk"              comment:"Action time"`
	ActionType types.ActionType `bun:"action_type,type:action_type" comment:"Action type"`
	Height     pkgTypes.Level   `bun:"height"                       comment:"Action block height"`
	TxId       uint64           `bun:"tx_id"                        comment:"Transaction internal id"`
	Size       int64            `bun:"size"                         comment:"Count bytes which was pushed to the rollup"`

	Action *Action `bun:"rel:belongs-to,join:action_id=id"`
	Rollup *Rollup `bun:"rel:belongs-to,join:rollup_id=id"`
	Tx     *Tx     `bun:"rel:belongs-to,join:tx_id=id"`
}

func (RollupAction) TableName() string {
	return "rollup_action"
}
