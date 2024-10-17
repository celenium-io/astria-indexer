// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"
	"testing"

	testsuite "github.com/celenium-io/astria-indexer/internal/test_suite"
	"github.com/celenium-io/astria-indexer/pkg/indexer/decode"
	"github.com/celenium-io/astria-indexer/pkg/node/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestParseTxs_EmptyTxsResults(t *testing.T) {
	block, _ := testsuite.EmptyBlock()

	decodeCtx := decode.NewContext(map[string]string{})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api := mock.NewMockApi(ctrl)
	resultTxs, err := parseTxs(context.Background(), block, &decodeCtx, api)

	assert.NoError(t, err)
	assert.Empty(t, resultTxs)
}

// TODO: uncomment tests after adding a test tx

func TestParseTxs_SuccessTx(t *testing.T) {
	// txRes := types.ResponseDeliverTx{
	// 	Code:      0,
	// 	Data:      []byte{},
	// 	Log:       "[]",
	// 	Info:      "info",
	// 	GasWanted: 12000,
	// 	GasUsed:   1000,
	// 	Events: []types.Event{
	// 		{
	// 			Type: "tx.fees",
	// 			Attributes: []types.EventAttribute{
	// 				{
	// 					Key:   "asset",
	// 					Value: "ibc/704031c868fd3d3c84a1cfa8cb45deba4ea746b44697f7f4a6ed1b8f6c239b82",
	// 				}, {
	// 					Key:   "feeAmount",
	// 					Value: "153",
	// 				}, {
	// 					Key:   "actionType",
	// 					Value: "astria.protocol.transaction.v1alpha1.FeeAssetChange",
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Codespace: "codespace",
	// }
	// block, now := testsuite.CreateTestBlock(txRes, true)
	// ctx := decode.NewContext(map[string]string{})
	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()
	// api := mock.NewMockApi(ctrl)

	// resultTxs, err := parseTxs(context.Background(), block, &ctx, api)

	// assert.NoError(t, err)
	// assert.Len(t, resultTxs, 1)

	// f := resultTxs[0]
	// assert.Equal(t, now, f.Time)
	// assert.Equal(t, storageTypes.StatusSuccess, f.Status)
	// assert.Equal(t, "", f.Error)
	// assert.Equal(t, int64(12000), f.GasWanted)
	// assert.Equal(t, int64(1000), f.GasUsed)
	// assert.Equal(t, "codespace", f.Codespace)
}

func TestParseTxs_FailedTx(t *testing.T) {
	// txRes := types.ResponseDeliverTx{
	// 	Code:      1,
	// 	Data:      []byte{},
	// 	Log:       "something weird happened",
	// 	Info:      "info",
	// 	GasWanted: 12000,
	// 	GasUsed:   1000,
	// 	Events:    nil,
	// 	Codespace: "codespace",
	// }
	// block, now := testsuite.CreateTestBlock(txRes, true)
	// ctx := decode.NewContext(map[string]string{})
	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()
	// api := mock.NewMockApi(ctrl)
	// resultTxs, err := parseTxs(context.Background(), block, &ctx, api)

	// assert.NoError(t, err)
	// assert.Len(t, resultTxs, 1)

	// f := resultTxs[0]
	// assert.Equal(t, now, f.Time)
	// assert.Equal(t, storageTypes.StatusFailed, f.Status)
	// assert.Equal(t, "something weird happened", f.Error)
	// assert.Equal(t, int64(12000), f.GasWanted)
	// assert.Equal(t, int64(1000), f.GasUsed)
	// assert.Equal(t, "codespace", f.Codespace)
}

func TestParseTxs_FailedTxWithNonstandardErrorCode(t *testing.T) {
	// txRes := types.ResponseDeliverTx{
	// 	Code:      300,
	// 	Data:      []byte{},
	// 	Log:       "something unusual happened",
	// 	Info:      "info",
	// 	GasWanted: 12000,
	// 	GasUsed:   1000,
	// 	Events:    nil,
	// 	Codespace: "codespace",
	// }
	// block, now := testsuite.CreateTestBlock(txRes, true)
	// ctx := decode.NewContext(map[string]string{})
	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()
	// api := mock.NewMockApi(ctrl)
	// resultTxs, err := parseTxs(context.Background(), block, &ctx, api)

	// assert.NoError(t, err)
	// assert.Len(t, resultTxs, 1)

	// f := resultTxs[0]
	// assert.Equal(t, now, f.Time)
	// assert.Equal(t, storageTypes.StatusFailed, f.Status)
	// assert.Equal(t, "something unusual happened", f.Error)
	// assert.Equal(t, int64(12000), f.GasWanted)
	// assert.Equal(t, int64(1000), f.GasUsed)
	// assert.Equal(t, "codespace", f.Codespace)
}
