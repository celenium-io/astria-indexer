// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/shopspring/decimal"
)

func Test_updateState(t *testing.T) {
	type args struct {
		block         *storage.Block
		totalAccounts int64
		totalRollups  int64
		totalBytes    int64
		state         *storage.State
	}

	now := time.Now()
	after := time.Now().Add(time.Minute)

	tests := []struct {
		name string
		args args
		want storage.State
	}{
		{
			name: "test 1",
			args: args{
				block: &storage.Block{
					Height:  101,
					Time:    after,
					ChainId: "test",
					Stats: &storage.BlockStats{
						TxCount:      10,
						SupplyChange: decimal.RequireFromString("100"),
						Fee:          decimal.RequireFromString("10"),
					},
				},
				totalAccounts: 10,
				totalRollups:  11,
				totalBytes:    12,
				state: &storage.State{
					Id:            1,
					Name:          "test",
					LastHeight:    100,
					LastTime:      now,
					ChainId:       "chain_id",
					TotalTx:       10,
					TotalAccounts: 2,
					TotalRollups:  12,
					TotalSupply:   decimal.RequireFromString("1000"),
					TotalBridges:  10,
					TotalBytes:    12,
				},
			},
			want: storage.State{
				Id:            1,
				Name:          "test",
				LastHeight:    101,
				LastTime:      after,
				ChainId:       "chain_id",
				TotalTx:       20,
				TotalAccounts: 12,
				TotalRollups:  23,
				TotalSupply:   decimal.RequireFromString("1100"),
				TotalBridges:  20,
				TotalBytes:    24,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateState(tt.args.block, tt.args.totalAccounts, tt.args.totalRollups, 10, tt.args.totalBytes, tt.args.state)
		})
	}
}
