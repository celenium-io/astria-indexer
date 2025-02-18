// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/pkg/types"
)

func (s *StorageTestSuite) TestBlockSignatureLevels() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	levels, err := s.BlockSignatures.LevelsByValidator(ctx, 2, 7963)
	s.Require().NoError(err)
	s.Require().Len(levels, 2)

	s.Require().Equal([]types.Level{7965, 7964}, levels)
}
