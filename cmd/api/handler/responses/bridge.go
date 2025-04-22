// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
)

// Bridge model info
//
//	@Description	bridge account information
type Bridge struct {
	Rollup   []byte `example:"O0Ia+lPYYMf3iFfxBaWXCSdlhphc6d4ZoBXINov6Tjc=" json:"rollup"    swaggertype:"string"`
	Asset    string `example:"nria"                                         json:"asset"     swaggertype:"string"`
	FeeAsset string `example:"nria"                                         json:"fee_asset" swaggertype:"string"`

	Address    *ShortAddress `json:"address,omitempty"`
	Sudo       *ShortAddress `json:"sudo,omitempty"`
	Withdrawer *ShortAddress `json:"withdrawer,omitempty"`
}

func NewBridge(b storage.Bridge) Bridge {
	bridge := Bridge{
		Asset:      b.Asset,
		FeeAsset:   b.FeeAsset,
		Address:    NewShortAddress(b.Address),
		Sudo:       NewShortAddress(b.Sudo),
		Withdrawer: NewShortAddress(b.Withdrawer),
	}

	if b.Rollup != nil {
		bridge.Rollup = b.Rollup.AstriaId
	}

	return bridge
}
