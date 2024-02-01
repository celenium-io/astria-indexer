// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"fmt"

	"github.com/aopoltorzhicky/astria/pkg/types"
	"github.com/uptrace/bun"
)

type RollupAddress struct {
	bun.BaseModel `bun:"rollup_address" comment:"Table with rollup addresses"`

	RollupId  uint64      `bun:"rollup_id,pk"  comment:"Rollup internal id"`
	AddressId uint64      `bun:"address_id,pk" comment:"Address internal id"`
	Height    types.Level `bun:"height"        comment:"Block height of the first sequence action"`

	Address *Address `bun:"rel:belongs-to,join:address_id=id"`
	Rollup  *Rollup  `bun:"rel:belongs-to,join:rollup_id=id"`
}

func (RollupAddress) TableName() string {
	return "rollup_address"
}

func (r RollupAddress) String() string {
	if r.Address == nil || r.Rollup == nil {
		return ""
	}
	return fmt.Sprintf("%s_%s", r.Address.String(), r.Rollup.String())
}
