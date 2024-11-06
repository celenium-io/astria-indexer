// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/uptrace/bun"
)

type AppBridge struct {
	bun.BaseModel `bun:"app_bridge" comment:"Table with application bridges"`

	AppId    uint64 `bun:"app_id,pk"            comment:"Application id"`
	BridgeId uint64 `bun:"bridge_id,pk"         comment:"Bridge id"`
	Native   bool   `bun:"native,default:false" comment:"Is native bridge for this application"`

	App    *App    `bun:"rel:belongs-to,join:app_id=id"`
	Bridge *Bridge `bun:"rel:belongs-to,join:bridge_id=id"`
}

func (AppBridge) TableName() string {
	return "app_bridge"
}
