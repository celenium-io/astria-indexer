// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"

	"github.com/aopoltorzhicky/astria/internal/storage"
	"github.com/aopoltorzhicky/astria/pkg/types"
)

type Rollup struct {
	Id           uint64      `example:"321"                                                              json:"id"            swaggertype:"integer"`
	FirstHeight  types.Level `example:"100"                                                              json:"first_height"  swaggertype:"integer"`
	AstriaId     string      `example:"19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539" json:"hash"          swaggertype:"string"`
	ActionsCount int64       `example:"100"                                                              json:"actions_count" swaggertype:"integer"`
	Size         int64       `example:"100"                                                              json:"size"          swaggertype:"integer"`
}

func NewRollup(rollup *storage.Rollup) Rollup {
	return Rollup{
		Id:           rollup.Id,
		AstriaId:     hex.EncodeToString(rollup.AstriaId),
		FirstHeight:  rollup.FirstHeight,
		ActionsCount: rollup.ActionsCount,
		Size:         rollup.Size,
	}
}
