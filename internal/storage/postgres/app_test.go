// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestLeaderboard() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	for _, column := range []string{
		columnSize, columnActionsCount, columnTime, "",
	} {

		apps, err := s.storage.App.Leaderboard(ctx, storage.LeaderboardFilters{
			SortField: column,
			Sort:      sdk.SortOrderDesc,
			Limit:     10,
			Offset:    0,
		})
		s.Require().NoError(err, column)
		s.Require().Len(apps, 1, column)

		app := apps[0]
		s.Require().EqualValues("App 1", app.Name, column)
		s.Require().EqualValues(34, app.Size, column)
		s.Require().EqualValues(3, app.ActionsCount, column)
		s.Require().False(app.LastActionTime.IsZero())
		s.Require().False(app.FirstActionTime.IsZero())
		s.Require().EqualValues(0.42857142857142855, app.ActionsCountPct)
		s.Require().EqualValues(0.3953488372093023, app.SizePct)
	}
}

func (s *StorageTestSuite) TestLeaderboardWithCategory() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	for _, column := range []string{
		columnSize, columnActionsCount, columnTime, "",
	} {

		apps, err := s.storage.App.Leaderboard(ctx, storage.LeaderboardFilters{
			SortField: column,
			Sort:      sdk.SortOrderDesc,
			Limit:     10,
			Offset:    0,
			Category:  []types.AppCategory{types.AppCategorySocial},
		})
		s.Require().NoError(err, column)
		s.Require().Len(apps, 1, column)

		app := apps[0]
		s.Require().EqualValues("Rollup 3", app.Name, column)
		s.Require().EqualValues("The third", app.Description, column)
		s.Require().EqualValues(34, app.Size, column)
		s.Require().EqualValues(3, app.ActionsCount, column)
		s.Require().False(app.LastActionTime.IsZero())
		s.Require().False(app.FirstActionTime.IsZero())
		s.Require().EqualValues(0.42857142857142855, app.ActionsCountPct)
		s.Require().EqualValues(0.3953488372093023, app.SizePct)
		s.Require().EqualValues("nft", app.Category)
	}
}
