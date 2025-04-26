// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package testsuite

import (
	"encoding/base64"
	"time"

	"github.com/celenium-io/astria-indexer/pkg/types"
	tmTypes "github.com/cometbft/cometbft/types"
)

func EmptyBlock() (types.BlockData, time.Time) {
	return createTestBlock(types.ResponseDeliverTx{}, 0, []string{}, false)
}

func CreateTestBlockV3(tx types.ResponseDeliverTx) (types.BlockData, time.Time) {
	return createTestBlock(tx, 3, txsV3, true)
}

func CreateTestBlockV0(tx types.ResponseDeliverTx) (types.BlockData, time.Time) {
	return createTestBlock(tx, 0, txsV0, true)
}

func CreateTestBlockV0ToV3(tx types.ResponseDeliverTx) (types.BlockData, time.Time) {
	return createTestBlock(tx, 0, txsV3, true)
}

var txsV3 = []string{
	"CiDjsMRCmPwcFJr79MiZb7kkJ65B5GSbk0yklZkbeFK4VQ==",
	"EiDjsMRCmPwcFJr79MiZb7kkJ65B5GSbk0yklZkbeFK4VQ==",
	"GmYKIFzfzCuQeY5A3ruKkHkUwtvLBr5FtNTSA5RaeD1h/MuuCiAuH0tPTyRhm6lcyjNE3BRnGPaHrNbdEnY6NZJ1+Q/L9wogLh9LT08kYZupXMozRNwUZxj2h6zW3RJ2OjWSdfkPy/c=",
	"CkCBcBs4ojg/IV1xLhFs5RNGSMwrlPQTuNSZB2WpfR1HD8AQF72Vdr/VhCKn2ppL4KBGgiKp+FqzchZdUIqTCy8IEiD1A4vkUMma1i1j6Mq95eIjQWn290avEMhbIEaK5psMmRq7AQovL2FzdHJpYS5wcm90b2NvbC50cmFuc2FjdGlvbi52MS5UcmFuc2FjdGlvbkJvZHkShwEKEBIOYXN0cmlhLWR1c2stMTESc2JxCi8SLWFzdHJpYTF5cWRqbm5tcnA3dzV5Z3dqMGRrbGRzZ3pqaHY1dmNha3A3eWV1ORIGCICU69wDGgRucmlhIgRucmlhKioweDlkMENFQzdCRUI5NDhBYjA0NmU4YjY0RTlhYTZDYzliNzMxQTk2MTM=",
}

var txsV0 = []string{
	"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
	"47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
	"CkDdioSOXz6qMrInWE3y+PMYsD9vmYXWwBoitzMylh3fi5Bru6SwanaybOflDioXHOoZDZFoY4QDtdxZu70KpboMEiBHWSfO5m65QYhB6ytVGUZ7HUgfR+jsMItBj+vdOHzccBr+AQovL2FzdHJpYS5wcm90b2NvbC50cmFuc2FjdGlvbi52MS5UcmFuc2FjdGlvbkJvZHkSygEKCggEEgZhc3RyaWESCboDBioECgASABIJugMGYgQKABIAEgm6AwYiBAoAEgASCboDBnIECgASABIJugMGQgQKABIAEgm6AwYyBAoAEgASCboDBjoECgASABIOugMLUgkKBQiAreIEEgASDboDChIICgQIwIQ9EgASDboDCloICgIIIBICCAMSC7oDCAoGCgIIZBIAEgy6AwlqBwoDCOIJEgASDLoDCRoHCgMIxBMSABIMugMJSgcKAwjEExIAEgq6AweCAQQKABIA",
}

func createTestBlock(tx types.ResponseDeliverTx, version uint64, txs []string, withTxs bool) (types.BlockData, time.Time) {
	now := time.Now()
	headerBlock := types.Block{
		Header: types.Header{
			Version: types.Consensus{
				App: version,
			},
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
