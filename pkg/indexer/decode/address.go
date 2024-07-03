// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/shopspring/decimal"
)

type Addresses map[string]*storage.Address

func NewAddress() Addresses {
	return make(map[string]*storage.Address)
}

func (a Addresses) Set(address string, height types.Level, change decimal.Decimal, actionCount int, signedTxCount int) *storage.Address {
	if addr, ok := a[address]; ok {
		addr.Balance.Total = addr.Balance.Total.Add(change)
		addr.ActionsCount += int64(actionCount)
		addr.SignedTxCount += int64(signedTxCount)
		return addr
	}
	addr := &storage.Address{
		Height:        height,
		Hash:          address,
		ActionsCount:  int64(actionCount),
		SignedTxCount: int64(signedTxCount),
		Balance: &storage.Balance{
			Total:    change,
			Currency: currency.DefaultCurrency,
		},
	}
	a[address] = addr
	return addr
}

func (a Addresses) UpdateNonce(address string, nonce uint32) {
	if address, ok := a[address]; ok {
		address.Nonce = nonce
	}
}

func (a Addresses) Get(address string) (*storage.Address, bool) {
	addr, ok := a[address]
	return addr, ok
}
