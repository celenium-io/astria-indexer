// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import "github.com/celenium-io/astria-indexer/internal/storage"

// Bridge model info
//
//	@Description	bridge account information
type Bridge struct {
	Address    string `example:"astria1phym4uktjn6gjle226009ge7u82w0dgtszs8x2" json:"address"    swaggertype:"string"`
	Rollup     []byte `example:"O0Ia+lPYYMf3iFfxBaWXCSdlhphc6d4ZoBXINov6Tjc="  json:"rollup"     swaggertype:"string"`
	Sudo       string `example:"astria1phym4uktjn6gjle226009ge7u82w0dgtszs8x2" json:"sudo"       swaggertype:"string"`
	Withdrawer string `example:"astria1phym4uktjn6gjle226009ge7u82w0dgtszs8x2" json:"withdrawer" swaggertype:"string"`
	Asset      string `example:"nria"                                          json:"asset"      swaggertype:"string"`
	FeeAsset   string `example:"nria"                                          json:"fee_asset"  swaggertype:"string"`
}

func NewBridge(b storage.Bridge) Bridge {
	bridge := Bridge{
		Asset:    b.Asset,
		FeeAsset: b.FeeAsset,
	}

	if b.Address != nil {
		bridge.Address = b.Address.Hash
	}
	if b.Sudo != nil {
		bridge.Sudo = b.Sudo.Hash
	}
	if b.Withdrawer != nil {
		bridge.Withdrawer = b.Withdrawer.Hash
	}
	if b.Rollup != nil {
		bridge.Rollup = b.Rollup.AstriaId
	}

	return bridge
}
