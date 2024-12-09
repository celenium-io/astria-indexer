// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/shopspring/decimal"
)

func (s *StorageTestSuite) TestAssetList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	assets, err := s.storage.Asset.List(ctx, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(assets, 3)

	m := map[string]storage.Asset{
		"asset-1": {
			Asset:         "asset-1",
			FeeCount:      0,
			Fee:           decimal.Zero,
			Transferred:   decimal.NewFromInt(1),
			TransferCount: 1,
		},
		"asset-2": {
			Asset:         "asset-2",
			FeeCount:      1,
			Fee:           decimal.NewFromInt(100),
			Transferred:   decimal.Zero,
			TransferCount: 0,
		},
		"nria": {
			Asset:         "nria",
			FeeCount:      1,
			Fee:           decimal.NewFromInt(100),
			Transferred:   decimal.NewFromInt(1),
			TransferCount: 1,
		},
	}

	for i := range assets {
		s.Require().Contains(m, assets[i].Asset)

		a := m[assets[i].Asset]
		s.Require().Equal(a.Asset, assets[i].Asset)
		s.Require().Equal(a.TransferCount, assets[i].TransferCount)
		s.Require().Equal(a.FeeCount, assets[i].FeeCount)
		s.Require().Equal(a.Transferred.String(), assets[i].Transferred.String())
		s.Require().Equal(a.Fee.String(), assets[i].Fee.String())
	}
}
