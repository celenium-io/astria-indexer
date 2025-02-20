// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/astria-indexer/internal/test_suite"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_rollbackBlock(t *testing.T) {
	indexName := "test"
	height := pkgTypes.Level(10000)
	blockTime := time.Now()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mock.NewMockTransaction(ctrl)

	t.Run("rollback block", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		tx.EXPECT().
			RollbackBlock(ctx, height).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackBlockStats(ctx, height).
			Return(storage.BlockStats{
				Height:       height,
				Time:         blockTime,
				TxCount:      1,
				BlockTime:    2,
				BytesInBlock: 100,
			}, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackTxs(ctx, height).
			Return([]storage.Tx{
				{
					Id:           1,
					Height:       height,
					Time:         blockTime,
					Position:     0,
					ActionsCount: 3,
					Status:       types.StatusSuccess,
					SignerId:     2,
					Nonce:        2,
					Hash:         testsuite.RandomHash(32),
				},
			}, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackActions(ctx, height).
			Return([]storage.Action{
				{
					Id:       1,
					Height:   height,
					Time:     blockTime,
					Position: 0,
					Type:     types.ActionTypeRollupDataSubmission,
					TxId:     1,
					Data: map[string]any{
						"rollup_id": "deadbeaf",
						"data":      "AAAAAA==",
					},
				}, {
					Id:       2,
					Height:   height,
					Time:     blockTime,
					Position: 1,
					Type:     types.ActionTypeTransfer,
					TxId:     1,
					Data: map[string]any{
						"to":     "deadbeaf",
						"amount": "100",
					},
				}, {
					Id:       3,
					Height:   height,
					Time:     blockTime,
					Position: 2,
					Type:     types.ActionTypeRollupDataSubmission,
					TxId:     1,
					Data: map[string]any{
						"data":      "+G6AhDuaygeCUgiUaN0ig7sPHLWZae8gW9rtKb4FEKSIiscjBInoAACAgxvZgqDlaFLJ2rb9OUtQRsM/meiHSoW2nSkIGJiW6fhUti+v16Ani2wgQDfXhYkgZylMwLhCXtawIhnoA8eVSnnsg/7jGQ==",
						"rollup_id": "GbqKuz5LVqMJ32dWxHuX4pjjpy2IRJ02oPrbHKc2ZTk=",
					},
				},
			}, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackAddressActions(ctx, height).
			Return([]storage.AddressAction{
				{
					Height:     height,
					Time:       blockTime,
					ActionType: types.ActionTypeRollupDataSubmission,
					ActionId:   1,
					AddressId:  1,
				}, {
					Height:     height,
					Time:       blockTime,
					ActionType: types.ActionTypeTransfer,
					ActionId:   2,
					AddressId:  1,
				}, {
					Height:     height,
					Time:       blockTime,
					ActionType: types.ActionTypeTransfer,
					ActionId:   2,
					AddressId:  2,
				}, {
					Height:     height,
					Time:       blockTime,
					ActionType: types.ActionTypeRollupDataSubmission,
					ActionId:   3,
					AddressId:  1,
				},
			}, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackAddresses(ctx, height).
			Return([]storage.Address{}, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackBalances(ctx, []uint64{}).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackBalanceUpdates(ctx, height).
			Return([]storage.BalanceUpdate{
				{
					Height:    height,
					AddressId: 1,
					Update:    decimal.RequireFromString("100"),
					Currency:  currency.DefaultCurrency,
				}, {
					Height:    height,
					AddressId: 2,
					Update:    decimal.RequireFromString("-100"),
					Currency:  currency.DefaultCurrency,
				},
			}, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			SaveBalances(ctx, gomock.Any()).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			LastNonce(ctx, uint64(1)).
			Return(1, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			LastNonce(ctx, uint64(2)).
			Return(0, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			UpdateAddresses(ctx, gomock.Any()).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackRollups(ctx, height).
			Return([]storage.Rollup{
				{
					Id:           1,
					ActionsCount: 1,
					AstriaId:     []byte("deadbeaf"),
					FirstHeight:  height,
					Size:         10,
				},
			}, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackRollupActions(ctx, height).
			Return([]storage.RollupAction{
				{
					RollupId: 1,
					ActionId: 1,
					Time:     blockTime,
					Height:   height,
				}, {
					RollupId: 2,
					ActionId: 3,
					Time:     blockTime,
					Height:   height,
				},
			}, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackRollupAddresses(ctx, height).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			UpdateRollups(ctx, gomock.Any()).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackValidators(ctx, height).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackBlockSignatures(ctx, height).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			RollbackFees(ctx, height).
			Return(nil).
			Times(1)

		tx.EXPECT().
			RollbackDeposits(ctx, height).
			Return(nil).
			Times(1)

		tx.EXPECT().
			RollbackTransfers(ctx, height).
			Return(nil).
			Times(1)

		tx.EXPECT().
			RollbackBridges(ctx, height).
			Return(0, nil).
			Times(1)

		lastBlock := storage.Block{
			Height:         height - 1,
			Time:           blockTime.Add(-time.Minute),
			VersionBlock:   1,
			VersionApp:     11,
			Hash:           testsuite.RandomHash(32),
			ParentHash:     testsuite.RandomHash(32),
			LastCommitHash: testsuite.RandomHash(32),
			AppHash:        testsuite.RandomHash(32),
			ConsensusHash:  testsuite.RandomHash(32),
			ProposerId:     1,
		}
		tx.EXPECT().
			LastBlock(ctx).
			Return(lastBlock, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			State(ctx, indexName).
			Return(storage.State{
				Name:          indexName,
				LastHeight:    height,
				LastTime:      blockTime,
				LastHash:      testsuite.RandomHash(32),
				TotalTx:       10,
				TotalAccounts: 10,
				TotalRollups:  10,
				TotalSupply:   decimal.RequireFromString("1000"),
			}, nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			Update(ctx, gomock.Any()).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		tx.EXPECT().
			Flush(ctx).
			Return(nil).
			MaxTimes(1).
			MinTimes(1)

		err := rollbackBlock(ctx, tx, height, indexName)
		require.NoError(t, err)
	})
}
