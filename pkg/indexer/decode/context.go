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
	Addresses        Addresses
	Rollups          Rollups
	Validators       Validators
	RollupAddress    map[string]*storage.RollupAddress
	AddressActions   map[string]*storage.AddressAction
	SupplyChange     decimal.Decimal
	BytesInBlock     int64
	GasUsed          int64
	GasWanted        int64
	DataSize         int64
	ActionTypes      storageTypes.Bits
	Constants        map[string]*storage.Constant
	Bridges          map[string]*storage.Bridge
	Fees             map[int64]*storage.Fee
	Transfers        []*storage.Transfer
	Deposits         map[int64]*storage.Deposit
	HasWriteAckError bool
	Proposer         string

	bridgeAssets map[string]string
}

func NewContext(bridgeAssets map[string]string) Context {
	return Context{
		Addresses:     NewAddress(),
		Rollups:       NewRollups(),
		RollupAddress: make(map[string]*storage.RollupAddress),
		SupplyChange:  decimal.Zero,
		Validators:    NewValidators(),
		Constants:     make(map[string]*storage.Constant),
		Bridges:       make(map[string]*storage.Bridge),
		Fees:          make(map[int64]*storage.Fee),
		Transfers:     make([]*storage.Transfer, 0),
		Deposits:      make(map[int64]*storage.Deposit),

		bridgeAssets: bridgeAssets,
	}
}

func (ctx *Context) AddGenericConstant(key, value string) {
	k := fmt.Sprintf("%s-%s", key, storageTypes.ModuleNameGeneric)
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

func (ctx *Context) AddBridge(b *storage.Bridge) {
	ctx.Bridges[b.Address.Hash] = b
}

func (ctx *Context) BridgesArray() []*storage.Bridge {
	arr := make([]*storage.Bridge, 0)
	for _, val := range ctx.Bridges {
		arr = append(arr, val)
	}
	return arr
}

func (ctx *Context) AddFee(idx int64, fee *storage.Fee) {
	ctx.Fees[idx] = fee
}

func (ctx *Context) ClearTx() {
	clear(ctx.Fees)
	ctx.HasWriteAckError = false
}

func (ctx *Context) AddBridgeAsset(bridge, asset string) {
	ctx.bridgeAssets[bridge] = asset
}

func (ctx *Context) AddTransfer(transfer *storage.Transfer) {
	ctx.Transfers = append(ctx.Transfers, transfer)
}

func (ctx *Context) AddDeposit(idx int64, deposit *storage.Deposit) {
	ctx.Deposits[idx] = deposit
}
