// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

type State struct {
	Id            uint64         `example:"321"                                                              format:"int64"     json:"id"             swaggertype:"integer"`
	Name          string         `example:"indexer"                                                          format:"string"    json:"name"           swaggertype:"string"`
	ChainID       string         `example:"astria-dusk-7"                                                    format:"string"    json:"chain_id"       swaggertype:"string"`
	LastHeight    pkgTypes.Level `example:"100"                                                              format:"int64"     json:"last_height"    swaggertype:"integer"`
	LastHash      string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"string"    json:"hash"           swaggertype:"string"`
	LastTime      time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"last_time"      swaggertype:"string"`
	TotalTx       int64          `example:"23456"                                                            format:"int64"     json:"total_tx"       swaggertype:"integer"`
	TotalAccounts int64          `example:"43"                                                               format:"int64"     json:"total_accounts" swaggertype:"integer"`
	TotalRollups  int64          `example:"312"                                                              format:"int64"     json:"total_rollups"  swaggertype:"integer"`
	TotalBridges  int64          `example:"312"                                                              format:"int64"     json:"total_bridges"  swaggertype:"integer"`
	TotalBytes    int64          `example:"312"                                                              format:"int64"     json:"total_bytes"    swaggertype:"integer"`
	TotalSupply   string         `example:"312"                                                              format:"string"    json:"total_supply"   swaggertype:"string"`
	Synced        bool           `example:"true"                                                             format:"boolean"   json:"synced"         swaggertype:"boolean"`
}

func NewState(state storage.State) State {
	return State{
		Id:            state.Id,
		Name:          state.Name,
		ChainID:       state.ChainId,
		LastHeight:    state.LastHeight,
		LastHash:      hex.EncodeToString(state.LastHash),
		LastTime:      state.LastTime,
		TotalTx:       state.TotalTx,
		TotalAccounts: state.TotalAccounts,
		TotalRollups:  state.TotalRollups,
		TotalBridges:  state.TotalBridges,
		TotalBytes:    state.TotalBytes,
		TotalSupply:   state.TotalSupply.String(),
		Synced:        !state.LastTime.UTC().Add(2 * time.Minute).Before(time.Now().UTC()),
	}
}
