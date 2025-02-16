// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestBlockStatsByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	stats, err := s.BlockStats.ByHeight(ctx, 7965)
	s.Require().NoError(err)

	s.Require().EqualValues(7965, stats.Height)
	s.Require().EqualValues(7966, stats.Id)
	s.Require().EqualValues(1, stats.TxCount)
	s.Require().EqualValues(2317, stats.BlockTime)
	s.Require().EqualValues(266, stats.BytesInBlock)
	s.Require().EqualValues("0", stats.SupplyChange.String())
}
