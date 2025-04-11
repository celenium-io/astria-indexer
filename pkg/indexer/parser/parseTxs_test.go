// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"
	"time"

	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/astria-indexer/internal/test_suite"
	"github.com/celenium-io/astria-indexer/pkg/indexer/decode"
	"github.com/celenium-io/astria-indexer/pkg/node/mock"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestParseTxs_EmptyTxsResults(t *testing.T) {
	block, _ := testsuite.EmptyBlock()

	decodeCtx := decode.NewContext(map[string]string{}, time.Now())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api := mock.NewMockApi(ctrl)
	resultTxs, err := parseTxs(t.Context(), block, &decodeCtx, api)

	assert.NoError(t, err)
	assert.Empty(t, resultTxs)
}

func TestParseTxs_SuccessTx(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code: 0,
		Data: []byte{},
		Log:  "[]",
		Info: "info",
		Events: []types.Event{
			{
				Type: "tx.deposit",
				Attributes: []types.EventAttribute{
					{
						Key:   "bridgeAddress",
						Value: "astria1yqdjnnmrp7w5ygwj0dkldsgzjhv5vcakp7yeu9",
					}, {
						Key:   "rollupId",
						Value: "xuzPaQ1qQbu1wl5oBdEiP-Xl977xEmJOVjVpqNrktnU=",
					}, {
						Key:   "amount",
						Value: "1000000000",
					}, {
						Key:   "asset",
						Value: "nria",
					}, {
						Key:   "destinationChainAddress",
						Value: "0x9d0CEC7BEB948Ab046e8b64E9aa6Cc9b731A9613",
					}, {
						Key:   "sourceTransactionId",
						Value: "c8c58bfd5d4745c08149da81c24df6738659634507c3eb535fa803c69d6a2b46",
					}, {
						Key:   "sourceActionIndex",
						Value: "0",
					},
				},
			},
		},
		Codespace: "codespace",
	}
	block, now := testsuite.CreateTestBlock(txRes, true)
	ctx := decode.NewContext(map[string]string{}, time.Now())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api := mock.NewMockApi(ctrl)

	resultTxs, err := parseTxs(t.Context(), block, &ctx, api)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusSuccess, f.Status)
	assert.Equal(t, "", f.Error)
	assert.Equal(t, "codespace", f.Codespace)
}

func TestParseTxs_FailedTx(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code:      1,
		Data:      []byte{},
		Log:       "something weird happened",
		Info:      "info",
		Events:    nil,
		Codespace: "codespace",
	}
	block, now := testsuite.CreateTestBlock(txRes, true)
	ctx := decode.NewContext(map[string]string{}, time.Now())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api := mock.NewMockApi(ctrl)
	resultTxs, err := parseTxs(t.Context(), block, &ctx, api)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusFailed, f.Status)
	assert.Equal(t, "something weird happened", f.Error)
	assert.Equal(t, "codespace", f.Codespace)
}

func TestParseTxs_FailedTxWithNonstandardErrorCode(t *testing.T) {
	txRes := types.ResponseDeliverTx{
		Code:      300,
		Data:      []byte{},
		Log:       "something unusual happened",
		Info:      "info",
		Events:    nil,
		Codespace: "codespace",
	}
	block, now := testsuite.CreateTestBlock(txRes, true)
	ctx := decode.NewContext(map[string]string{}, time.Now())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api := mock.NewMockApi(ctrl)
	resultTxs, err := parseTxs(t.Context(), block, &ctx, api)

	assert.NoError(t, err)
	assert.Len(t, resultTxs, 1)

	f := resultTxs[0]
	assert.Equal(t, now, f.Time)
	assert.Equal(t, storageTypes.StatusFailed, f.Status)
	assert.Equal(t, "something unusual happened", f.Error)
	assert.Equal(t, "codespace", f.Codespace)
}
