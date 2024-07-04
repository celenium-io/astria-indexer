// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
)

// Address model info
//
//	@Description	address information
type Address struct {
	Id            uint64         `example:"321"                                           json:"id"              swaggertype:"integer"`
	Height        pkgTypes.Level `example:"100"                                           json:"first_height"    swaggertype:"integer"`
	ActionsCount  int64          `example:"10"                                            json:"actions_count"   swaggertype:"integer"`
	SignedTxCount int64          `example:"10"                                            json:"signed_tx_count" swaggertype:"integer"`
	Nonce         uint32         `example:"10"                                            json:"nonce"           swaggertype:"integer"`
	Hash          string         `example:"astria1phym4uktjn6gjle226009ge7u82w0dgtszs8x2" json:"hash"            swaggertype:"string"`

	Balance *Balance `json:"balance,omitempty"`
	Bridge  *Bridge  `json:"bridge,omitempty"`
}

func NewAddress(addr storage.Address, bridge *storage.Bridge) Address {
	result := Address{
		Id:            addr.Id,
		Height:        addr.Height,
		ActionsCount:  addr.ActionsCount,
		SignedTxCount: addr.SignedTxCount,
		Nonce:         addr.Nonce,
		Hash:          addr.String(),
	}

	if addr.Balance != nil {
		result.Balance = &Balance{
			Currency: addr.Balance.Currency,
			Value:    addr.Balance.Total.String(),
		}
	}

	if bridge != nil {
		b := NewBridge(*bridge)
		result.Bridge = &b
	}

	return result
}

// Balance info
//
//	@Description	Balance of address information
type Balance struct {
	Currency string `example:"nria"        json:"currency" swaggertype:"string"`
	Value    string `example:"10000000000" json:"value"    swaggertype:"string"`
}
