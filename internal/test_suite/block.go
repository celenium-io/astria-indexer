// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package testsuite

import (
	"encoding/base64"
	"time"

	"github.com/aopoltorzhicky/astria/pkg/types"
	tmTypes "github.com/cometbft/cometbft/types"
)

func EmptyBlock() (types.BlockData, time.Time) {
	return CreateTestBlock(types.ResponseDeliverTx{}, false)
}

var txs = []string{
	"VgeosXgFzGcWbA+kQcQLNLiu8YP1J7qwUgp/e99TIXM=",
	"NEsEB484/VUytEr01o0tdtmPOpWqLC8cEP6ASq2LW1w=",
	"CkCfKZ6MNVHTNOY39V/JfReAS413lp+9m18mwZwmzqM0/aKsfMIqTLsJdzIKNV+m9fd48ITUwWOLU17kcZx6URkNEiDsYdXD5V4Xza1Y2UcI61HZVIjxmIXFBSEIu9lkT8BLNxrKBQjCNBKiBRKfBQogGbqKuz5LVqMJ32dWxHuX4pjjpy2IRJ02oPrbHKc2ZTkS+gQC+QJ2gw3sr4IEUYSy0F4AhLLQXg6DAtyslAmg65pJxLDfmTp0rxalNWWGsSmIgLkCBHw50TAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMYAAAAAAAAAAAAAAABJF6l0antuClcVnMt/WmdEJH8tDapmVFBOz0u8NriFCQWtTZPZM7U5k1NtwUPKVdDfjLE4AAAHfOgoeuJ6rlP3k1cuCBXKRi0ltGJDjbVskLry4I2phvVyG6nzcK5Phsa0Erk3IHcJUpZ7oVn8+HuqEj9kMjL91xQbjBcicjPQz4nhZLeMugf4N06BYc+FyGzIv0sPKRKwixAWVddZvD+Xuky42pgcBrqZB01ho7sF4bn1PfmzDfLV8hwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACNAwAAB5UAqjanAAAAAAAAAAAAAAAAJUOBeypNWzTeiwAz1W9c7LF2H7QADeyvAAAAAAAAAAAAAAAAFVHN0OTfMW4bBEqEKDyEvJJM2r8AAAAAAAAAAAAAAAA4vfk0aAkIcb+8o0ORpOTXEVEkswAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACxorwuxQAAAAAAAAAAAAAAAAAAAAAAAAAAAMABoNaizXg5BZDQ6WNvMDtkVxjyXDeq7W7GuG/15jVqN+TqoHmMf4525t9BWdSjfSshbLppfvsDOQzE62o8SFqnIOhjGiBwQDHIaP09PIShz6jLRd66TqdGtEaX9/Sm7RuPbCObgg==",
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
