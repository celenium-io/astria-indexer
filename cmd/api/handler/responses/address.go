// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
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
	IsBridge      bool           `example:"false"                                         json:"is_bridge"       swaggertype:"boolean"`
	IsIbcRelayer  bool           `example:"false"                                         json:"is_ibc_relayer"  swaggertype:"boolean"`
	IsSudo        bool           `example:"false"                                         json:"is_sudo"         swaggertype:"boolean"`
	IsIbcSudo     bool           `example:"false"                                         json:"is_ibc_sudo"     swaggertype:"boolean"`

	Balance    []Balance  `json:"balances"`
	Bridge     *Bridge    `json:"bridge,omitempty"`
	Celestials *Celestial `json:"celestials,omitempty"`
}

func NewAddress(addr storage.Address, bridge *storage.Bridge, sudoAddr, ibcSudoAddr string) Address {
	result := Address{
		Id:            addr.Id,
		Height:        addr.Height,
		ActionsCount:  addr.ActionsCount,
		SignedTxCount: addr.SignedTxCount,
		Nonce:         addr.Nonce,
		IsBridge:      addr.IsBridge,
		Hash:          addr.String(),
		Balance:       make([]Balance, 0),
	}

	result.IsSudo = sudoAddr == result.Hash
	result.IsIbcSudo = ibcSudoAddr == result.Hash

	for i := range addr.Balance {
		result.Balance = append(result.Balance, Balance{
			Currency: addr.Balance[i].Currency,
			Value:    addr.Balance[i].Total.String(),
		})
	}

	if bridge != nil {
		b := NewBridge(*bridge)
		result.Bridge = &b
	}
	if addr.IsIbcRelayer != nil {
		result.IsIbcRelayer = *addr.IsIbcRelayer
	}
	result.AddCelestials(addr.Celestials)
	
	return result
}

func (address *Address) AddCelestials(celestials *celestials.Celestial) {
	if celestials != nil {
		address.Celestials = NewCelestial(celestials)
	}
}

// Balance info
//
//	@Description	Balance of address information
type Balance struct {
	Currency string `example:"nria"        json:"currency" swaggertype:"string"`
	Value    string `example:"10000000000" json:"value"    swaggertype:"string"`
}

// Celestial ID
//
//	@Description	Linked celestial id
type Celestial struct {
	Name      string `example:"name"                                                                json:"name"              swaggertype:"string"`
	ImageUrl  string `example:"https://ipfs.io/ipfs/QmUi269vE25fagqhyMCCTNSoiW6x4LHCwwQb3keSrEXAmC" json:"image_url"         swaggertype:"string"`
	Status    string `example:"VERIFIED"                                                            json:"status,omitempty"  swaggertype:"string"`
	IsPrimary bool   `example:"true"                                                                json:"primary,omitempty" swaggertype:"boolean"`
}

func NewCelestial(c *celestials.Celestial) *Celestial {
	return &Celestial{
		ImageUrl:  c.ImageUrl,
		Name:      c.Id,
		IsPrimary: c.Status == celestials.StatusPRIMARY,
		Status:    c.Status.String(),
	}
}

type ShortAddress struct {
	Hash       string     `example:"astria1f0dw5muma062mwfz7g46229adaycpjevnyw9fc" json:"hash" swaggertype:"string"`
	Celestials *Celestial `json:"celestials,omitempty"`
}

func NewShortAddress(address *storage.Address) *ShortAddress {
	if address == nil {
		return nil
	}

	result := new(ShortAddress)
	result.Hash = address.Hash
	if address.Celestials != nil {
		result.Celestials = NewCelestial(address.Celestials)
	}
	return result
}
