// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
)

type Rollup struct {
	Id           uint64      `example:"321"                                          json:"id"            swaggertype:"integer"`
	FirstHeight  types.Level `example:"100"                                          json:"first_height"  swaggertype:"integer"`
	AstriaId     []byte      `example:"O0Ia+lPYYMf3iFfxBaWXCSdlhphc6d4ZoBXINov6Tjc=" json:"hash"          swaggertype:"string"`
	ActionsCount int64       `example:"101"                                          json:"actions_count" swaggertype:"integer"`
	BridgeCount  int64       `example:"2"                                            json:"bridge_count"  swaggertype:"integer"`
	Size         int64       `example:"100"                                          json:"size"          swaggertype:"integer"`

	App *AppWithStats `json:"app,omitempty"`
}

func NewRollup(rollup *storage.Rollup) Rollup {
	r := Rollup{
		Id:           rollup.Id,
		AstriaId:     rollup.AstriaId,
		FirstHeight:  rollup.FirstHeight,
		ActionsCount: rollup.ActionsCount,
		BridgeCount:  rollup.BridgeCount,
		Size:         rollup.Size,
	}

	return r
}
