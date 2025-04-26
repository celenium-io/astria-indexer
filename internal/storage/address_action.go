// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/uptrace/bun"
)

type AddressAction struct {
	bun.BaseModel `bun:"address_action" comment:"Table with address actions"`

	AddressId  uint64           `bun:"address_id,pk"                comment:"Address internal id"`
	ActionId   uint64           `bun:"action_id,pk"                 comment:"Action internal id"`
	TxId       uint64           `bun:"tx_id"                        comment:"Tx internal id"`
	ActionType types.ActionType `bun:"action_type,type:action_type" comment:"Action type"`
	Time       time.Time        `bun:"time,notnull,pk"              comment:"Action time"`
	Height     pkgTypes.Level   `bun:"height"                       comment:"Action block height"`

	Address *Address `bun:"rel:belongs-to,join:address_id=id"`
	Action  *Action  `bun:"rel:belongs-to,join:action_id=id"`
	Tx      *Tx      `bun:"rel:belongs-to,join:tx_id=id"`
}

func (AddressAction) TableName() string {
	return "address_action"
}
