// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import "github.com/uptrace/bun"

type AppId struct {
	bun.BaseModel `bun:"app_id" comment:"Table with application bridges"`

	AppId     uint64 `bun:"app_id,pk"     comment:"Application id"`
	RolllupId uint64 `bun:"rollup_id,pk"  comment:"Rollup id"`
	AddressId uint64 `bun:"address_id,pk" comment:"Address id"`

	App     *App     `bun:"rel:belongs-to,join:app_id=id"`
	Rollup  *Rollup  `bun:"rel:belongs-to,join:rollup_id=id"`
	Address *Address `bun:"rel:belongs-to,join:address_id=id"`
}

func (AppId) TableName() string {
	return "app_id"
}
