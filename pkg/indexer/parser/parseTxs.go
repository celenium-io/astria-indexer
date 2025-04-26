// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
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
	txs := make([]*storage.Tx, 0)

	for i := 0; i < len(b.TxsResults); i++ {
		if err := parseEvents(ctx, b.TxsResults[i].Events, b.Height, decodeCtx, api); err != nil {
			return nil, errors.Wrap(err, "parse events")
		}

		t, err := parseTx(b, i, decodeCtx)
		if err != nil {
			return nil, err
		}
		if len(t.Hash) == 0 {
			continue
		}
		txs = append(txs, &t)

		decodeCtx.ClearTx()
	}

	return txs, nil
}

func parseTx(b types.BlockData, index int, ctx *decode.Context) (storage.Tx, error) {
	d, err := decode.Tx(b, index, ctx)
	if err != nil {
		return storage.Tx{}, errors.Wrapf(err, "while parsing Tx on index %d", index)
	}

	if d.IsDataItem {
		return storage.Tx{}, nil
	}

	result := b.TxsResults[index]

	t := storage.Tx{
		Height:       b.Height,
		Time:         b.Block.Time,
		Position:     int64(index),
		ActionsCount: int64(len(d.Actions)),
		Status:       storageTypes.StatusSuccess,
		Codespace:    result.Codespace,
		Hash:         b.Block.Txs[index].Hash(),
		Signature:    d.Tx.GetSignature(),
		Nonce:        d.UnsignedTx.GetParams().GetNonce(),
		Signer:       d.Signer,
		ActionTypes:  d.ActionTypes,

		Actions:   d.Actions,
		BytesSize: int64(len(result.Data)),
	}

	if result.IsFailed() {
		t.Status = storageTypes.StatusFailed
		t.Error = result.Log
	}

	return t, nil
}
