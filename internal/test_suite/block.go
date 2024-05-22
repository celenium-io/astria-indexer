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
	"CkDNosMUJAdxeMWfzZPWO0f1FbF4ez8ITp8OLEtjukT+V7R74OY6qqhBv46Tq4LMQBJW2qV/eSJa3rIhNijxR+kKEiCyrIGvc6WftEDvMbvCIN/RE/QJmS3CVxPz1hdXvW1VtBq7AQo6L2FzdHJpYS5wcm90b2NvbC50cmFuc2FjdGlvbnMudjFhbHBoYTEuVW5zaWduZWRUcmFuc2FjdGlvbhJ9CmpaaAoiCiBD3GftSSqiCwsXSvGXMO/RH4PVA73gqzIUzcWUSbo2OBIgcEAxyGj9PTyEoc+oy0Xeuk6nRrRGl/f0pu0bj2wjm4IaIHBAMcho/T08hKHPqMtF3rpOp0a0Rpf39KbtG49sI5uCEg8SDWFzdHJpYS1kdXNrLTY=",
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
