// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	astria "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transaction/v1"
	sequencerblockv1 "buf.build/gen/go/astria/sequencerblock-apis/protocolbuffers/go/astria/sequencerblock/v1"
	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
)

type DecodedTx struct {
	Tx          *astria.Transaction
	UnsignedTx  *astria.TransactionBody
	Actions     []storage.Action
	Signer      *storage.Address
	ActionTypes storageTypes.Bits
	IsDataItem  bool
}

func Tx(b types.BlockData, index int, ctx *Context) (d DecodedTx, err error) {
	raw := b.Block.Txs[index]

	ctx.BytesInBlock += int64(len(raw))

	d.Tx = new(astria.Transaction)
	if err := proto.Unmarshal(raw, d.Tx); err != nil {
		idi, err := isDataItem(b, index, raw, ctx)
		if err != nil {
			return d, errors.Wrap(err, "tx decoding")
		}
		d.IsDataItem = idi
		return d, nil
	}

	body := d.Tx.GetBody()
	if body == nil {
		idi, err := isDataItem(b, index, raw, ctx)
		if err != nil {
			return d, errors.Wrap(err, "tx decoding")
		}
		d.IsDataItem = idi
		return d, nil
	}

	d.UnsignedTx = new(astria.TransactionBody)
	if err := proto.Unmarshal(body.GetValue(), d.UnsignedTx); err != nil {
		return d, errors.Wrap(err, "tx decoding")
	}

	address, err := AddressFromPubKey(d.Tx.GetPublicKey())
	if err != nil {
		return d, errors.Wrapf(err, "decode publick key: %x", d.Tx.GetPublicKey())
	}
	d.Signer = ctx.Addresses.Set(address, b.Height, decimal.Zero, "", 0, 1)
	ctx.Addresses.UpdateNonce(address, d.UnsignedTx.GetParams().GetNonce())

	d.Actions, err = parseActions(b.Height, b.Block.Time, address, &d, ctx)
	if err != nil {
		return d, errors.Wrap(err, "parsing actions")
	}
	ctx.ActionTypes.Set(d.ActionTypes)

	return
}

func isDataItem(b types.BlockData, index int, raw []byte, ctx *Context) (bool, error) {
	if b.Block.Version.App < 3 {
		if len(raw) == 32 && index < 2 {
			return true, nil
		}
	}

	dataItem := new(sequencerblockv1.DataItem)
	err := proto.Unmarshal(raw, dataItem)
	return err == nil, err
}
