// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"fmt"

	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/shopspring/decimal"
)

type Context struct {
	Addresses      Addresses
	Rollups        Rollups
	Validators     Validators
	RollupAddress  map[string]*storage.RollupAddress
	AddressActions map[string]*storage.AddressAction
	SupplyChange   decimal.Decimal
	BytesInBlock   int64
	GasUsed        int64
	GasWanted      int64
	DataSize       int64
	ActionTypes    storageTypes.Bits
	Constants      map[string]*storage.Constant
}

func NewContext() Context {
	return Context{
		Addresses:     NewAddress(),
		Rollups:       NewRollups(),
		RollupAddress: make(map[string]*storage.RollupAddress),
		SupplyChange:  decimal.Zero,
		Validators:    NewValidators(),
		Constants:     make(map[string]*storage.Constant),
	}
}

func (ctx *Context) AddGenericConstant(key, value string) {
	k := fmt.Sprintf("%s-%s", key, value)
	ctx.Constants[k] = &storage.Constant{
		Module: storageTypes.ModuleNameGeneric,
		Name:   key,
		Value:  value,
	}
}

func (ctx *Context) ConstantsArray() []*storage.Constant {
	arr := make([]*storage.Constant, 0)
	for _, val := range ctx.Constants {
		arr = append(arr, val)
	}
	return arr
}
