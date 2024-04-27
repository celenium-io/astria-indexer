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
	"VgeosXgFzGcWbA+kQcQLNLiu8YP1J7qwUgp/e99TIXM=",
	"NEsEB484/VUytEr01o0tdtmPOpWqLC8cEP6ASq2LW1w=",
	"CkDwtsuHHqKfR4VeCIU6FxjJUWlvSzkEFBL+f8gSZIvtRK80YmkfvnrrtiNi0DBIDbRISsIoeCF4a9OaA8JUXCYAEiCyrIGvc6WftEDvMbvCIN/RE/QJmS3CVxPz1hdXvW1VtBp9CmpaaAoiCiA7Qhr6U9hgx/eIV/EFpZcJJ2WGmFzp3hmgFcg2i/pONxIgcEAxyGj9PTyEoc+oy0Xeuk6nRrRGl/f0pu0bj2wjm4IaIHBAMcho/T08hKHPqMtF3rpOp0a0Rpf39KbtG49sI5uCEg8SDWFzdHJpYS1kdXNrLTU=",
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
