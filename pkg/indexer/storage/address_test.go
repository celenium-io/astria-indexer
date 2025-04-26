// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"testing"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_saveAddresses(t *testing.T) {
	tests := []struct {
		name      string
		addresses map[string]*storage.Address
		addr      map[string]uint64
		total     int64
		wantErr   bool
	}{
		{
			name:      "test 1",
			addresses: map[string]*storage.Address{},
			addr:      nil,
			total:     0,
			wantErr:   false,
		}, {
			name: "test 2",
			addresses: map[string]*storage.Address{
				"deadbeaf": {
					Hash:   "deadbeaf",
					Height: 100,
					Balance: []*storage.Balance{
						{
							Currency: "nria",
							Total:    decimal.RequireFromString("1"),
						},
					},
				},
			},
			addr: map[string]uint64{
				"deadbeaf": 1,
			},
			total:   1,
			wantErr: false,
		}, {
			name: "test 3",
			addresses: map[string]*storage.Address{
				"deadbeaf": {
					Hash:   "deadbeaf",
					Height: 100,
					Balance: []*storage.Balance{
						{
							Currency: "nria",
							Total:    decimal.RequireFromString("1"),
						},
					},
				},
			},
			addr: map[string]uint64{
				"deadbeaf": 1,
			},
			total:   0,
			wantErr: false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	for _, tt := range tests {
		tx := mock.NewMockTransaction(ctrl)

		tx.EXPECT().
			SaveAddresses(ctx, gomock.Any()).
			MaxTimes(1).
			MinTimes(0).
			DoAndReturn(func(_ context.Context, addresses ...*storage.Address) (int64, error) {
				require.Equal(t, len(tt.addresses), len(addresses))
				for i := range addresses {
					addresses[i].Id = uint64(i + 1)
				}
				return tt.total, nil
			})

		tx.EXPECT().
			SaveBalances(ctx, gomock.Any()).
			MaxTimes(1).
			MinTimes(0).
			DoAndReturn(func(_ context.Context, balances ...storage.Balance) error {
				require.Equal(t, len(tt.addresses), len(balances))
				return nil
			})

		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := saveAddresses(ctx, tx, tt.addresses)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.addr, got)
			require.Equal(t, tt.total, got1)
		})
	}
}
