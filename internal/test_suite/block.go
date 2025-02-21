// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

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
	"CkCBcBs4ojg/IV1xLhFs5RNGSMwrlPQTuNSZB2WpfR1HD8AQF72Vdr/VhCKn2ppL4KBGgiKp+FqzchZdUIqTCy8IEiD1A4vkUMma1i1j6Mq95eIjQWn290avEMhbIEaK5psMmRq7AQovL2FzdHJpYS5wcm90b2NvbC50cmFuc2FjdGlvbi52MS5UcmFuc2FjdGlvbkJvZHkShwEKEBIOYXN0cmlhLWR1c2stMTESc2JxCi8SLWFzdHJpYTF5cWRqbm5tcnA3dzV5Z3dqMGRrbGRzZ3pqaHY1dmNha3A3eWV1ORIGCICU69wDGgRucmlhIgRucmlhKioweDlkMENFQzdCRUI5NDhBYjA0NmU4YjY0RTlhYTZDYzliNzMxQTk2MTM=",
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
