// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package testsuite

import (
	"encoding/base64"
	"time"

	"github.com/celenium-io/astria-indexer/pkg/types"
	tmTypes "github.com/cometbft/cometbft/types"
)

func EmptyBlock() (types.BlockData, time.Time) {
	return CreateTestBlock(types.ResponseDeliverTx{}, false)
}

var txs = []string{
	"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
	"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
	// TODO: add tx to test
	// "CkATSNqP3ZBN+DgYi3anxjPcVt14d55cFq5k2hD4bZ+3iblkQ9JDam0MxTFC7osXFM/WFrrVSopirkV50xYGyh8PEiBzAhb1xqe9xHmqHiernziq508EQ6u4fVdyVFwH1Ctz2RpaCjovYXN0cmlhLnByb3RvY29sLnRyYW5zYWN0aW9ucy52MWFscGhhMS5VbnNpZ25lZFRyYW5zYWN0aW9uEhwKCaoDBgoEbnJpYRIPEg1hc3RyaWEtZHVzay04",
}

func CreateTestBlock(tx types.ResponseDeliverTx, withTxs bool) (types.BlockData, time.Time) {
	now := time.Now()
	headerBlock := types.Block{
		Header: types.Header{
			Time: now,
		},
		Data: types.Data{
			Txs: make(tmTypes.Txs, 0),
		},
	}

	var txResults = make([]*types.ResponseDeliverTx, 0)
	if withTxs {
		for i := 0; i < len(txs); i++ {
			raw, _ := base64.StdEncoding.DecodeString(txs[i])
			txResults = append(txResults, &tx)
			headerBlock.Data.Txs = append(headerBlock.Data.Txs, raw)
		}
	}

	block := types.BlockData{
		ResultBlock: types.ResultBlock{
			Block: &headerBlock,
		},
		ResultBlockResults: types.ResultBlockResults{
			TxsResults: txResults,
		},
	}

	return block, now
}

func CreateBlockWithTxs(tx types.ResponseDeliverTx, txData []byte, count int) (types.BlockData, time.Time) {
	now := time.Now()
	headerBlock := types.Block{
		Header: types.Header{
			Time: now,
		},
		Data: types.Data{
			Txs: make(tmTypes.Txs, count),
		},
	}

	var txResults = make([]*types.ResponseDeliverTx, count)
	for i := 0; i < count; i++ {
		txResults[i] = &tx
		headerBlock.Data.Txs[i] = txData
	}

	block := types.BlockData{
		ResultBlock: types.ResultBlock{
			Block: &headerBlock,
		},
		ResultBlockResults: types.ResultBlockResults{
			TxsResults: txResults,
		},
	}

	return block, now
}
