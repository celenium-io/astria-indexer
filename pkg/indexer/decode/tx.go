// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	astria "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
)

type Context struct {
	Addresses      Addresses
	Rollups        Rollups
	RollupAddress  map[string]*storage.RollupAddress
	AddressActions map[string]*storage.AddressAction
	SupplyChange   decimal.Decimal
	BytesInBlock   int64
	GasUsed        int64
	GasWanted      int64
	DataSize       int64
	ActionTypes    storageTypes.Bits
}

func NewContext() Context {
	return Context{
		Addresses:     NewAddress(),
		Rollups:       NewRollups(),
		RollupAddress: make(map[string]*storage.RollupAddress),
		SupplyChange:  decimal.Zero,
	}
}

type DecodedTx struct {
	Tx          *astria.SignedTransaction
	UnsignedTx  *astria.UnsignedTransaction
	Actions     []storage.Action
	Signer      *storage.Address
	ActionTypes storageTypes.Bits
}

func Tx(b types.BlockData, index int, ctx *Context) (d DecodedTx, err error) {
	raw := b.Block.Txs[index]

	ctx.BytesInBlock += int64(len(raw))

	d.Tx = new(astria.SignedTransaction)
	if err := proto.Unmarshal(raw, d.Tx); err != nil {
		return d, errors.Wrap(err, "tx decoding")
	}

	if d.Tx.GetTransaction() == nil {
		return d, errors.Wrap(err, "nil decoded tx")
	}

	d.UnsignedTx = new(astria.UnsignedTransaction)
	if err := proto.Unmarshal(d.Tx.GetTransaction().GetValue(), d.UnsignedTx); err != nil {
		return d, errors.Wrap(err, "tx decoding")
	}

	address := AddressFromPubKey(d.Tx.GetPublicKey())
	d.Signer = ctx.Addresses.Set(address, b.Height, decimal.Zero, 0, 1)
	ctx.Addresses.UpdateNonce(address, d.UnsignedTx.GetParams().GetNonce())

	d.Actions, err = parseActions(b.Height, b.Block.Time, address, &d, ctx)
	if err != nil {
		return d, errors.Wrap(err, "parsing actions")
	}
	ctx.ActionTypes.Set(d.ActionTypes)

	return
}
