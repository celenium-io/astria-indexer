// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/celenium-io/astria-indexer/pkg/indexer/decode"
	"github.com/celenium-io/astria-indexer/pkg/node"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
)

func parseTxs(ctx context.Context, b types.BlockData, decodeCtx *decode.Context, api node.Api) ([]*storage.Tx, error) {
	count := len(b.Block.Txs)
	index := 0
	if count == 0 {
		return []*storage.Tx{}, nil
	}

	if len(b.Block.Txs) >= 2 && len(b.Block.Txs[0]) == 32 && len(b.Block.Txs[1]) == 32 {
		count -= 2
		index = 2
		decodeCtx.BytesInBlock += 64
	}

	txs := make([]*storage.Tx, count)

	for i := index; i < len(b.TxsResults); i++ {
		if err := parseEvents(ctx, b.TxsResults[i].Events, decodeCtx, api); err != nil {
			return nil, errors.Wrap(err, "parse events")
		}

		t, err := parseTx(b, i, b.TxsResults[i], decodeCtx)
		if err != nil {
			return nil, err
		}
		txs[i-index] = &t

		decodeCtx.GasWanted += b.TxsResults[i].GasWanted
		decodeCtx.GasUsed += b.TxsResults[i].GasUsed
	}

	return txs, nil
}

func parseTx(b types.BlockData, index int, txRes *types.ResponseDeliverTx, ctx *decode.Context) (storage.Tx, error) {
	d, err := decode.Tx(b, index, ctx)
	if err != nil {
		return storage.Tx{}, errors.Wrapf(err, "while parsing Tx on index %d", index)
	}

	t := storage.Tx{
		Height:       b.Height,
		Time:         b.Block.Time,
		Position:     int64(index),
		GasWanted:    txRes.GasWanted,
		GasUsed:      txRes.GasUsed,
		ActionsCount: int64(len(d.Actions)),
		Status:       storageTypes.StatusSuccess,
		Codespace:    txRes.Codespace,
		Hash:         b.Block.Txs[index].Hash(),
		Signature:    d.Tx.GetSignature(),
		Nonce:        d.UnsignedTx.GetParams().GetNonce(),
		Signer:       d.Signer,
		ActionTypes:  d.ActionTypes,

		Actions:   d.Actions,
		BytesSize: int64(len(txRes.Data)),
	}

	if txRes.IsFailed() {
		t.Status = storageTypes.StatusFailed
		t.Error = txRes.Log
	}

	return t, nil
}
