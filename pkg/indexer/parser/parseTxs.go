// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"github.com/aopoltorzhicky/astria/internal/storage"
	storageTypes "github.com/aopoltorzhicky/astria/internal/storage/types"
	"github.com/aopoltorzhicky/astria/pkg/indexer/decode"
	"github.com/aopoltorzhicky/astria/pkg/types"
	"github.com/pkg/errors"
)

func parseTxs(b types.BlockData, ctx *decode.Context) ([]*storage.Tx, error) {
	count := len(b.Block.Txs)
	index := 0
	if count == 0 {
		return []*storage.Tx{}, nil
	}

	if len(b.Block.Txs) >= 2 && len(b.Block.Txs[0]) == 32 && len(b.Block.Txs[1]) == 32 {
		count -= 2
		index = 2
		ctx.BytesInBlock += 64
	}

	txs := make([]*storage.Tx, count)

	for i := index; i < len(b.TxsResults); i++ {
		t, err := parseTx(b, i, b.TxsResults[i], ctx)
		if err != nil {
			return nil, err
		}
		txs[i-index] = &t

		ctx.GasWanted += b.TxsResults[i].GasWanted
		ctx.GasUsed += b.TxsResults[i].GasUsed
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
		Signature:    d.Tx.Signature,
		Nonce:        d.Tx.Transaction.Nonce,
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
