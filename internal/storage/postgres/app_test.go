// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
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
		s.Require().EqualValues(1, app.ActionsCount, column)
		s.Require().False(app.LastActionTime.IsZero())
		s.Require().False(app.FirstActionTime.IsZero())
		s.Require().EqualValues(1, app.ActionsCountPct)
		s.Require().EqualValues(1, app.SizePct)
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
		s.Require().EqualValues("App 1", app.Name, column)
		s.Require().EqualValues(34, app.Size, column)
		s.Require().EqualValues(1, app.ActionsCount, column)
		s.Require().False(app.LastActionTime.IsZero())
		s.Require().False(app.FirstActionTime.IsZero())
		s.Require().EqualValues(1, app.ActionsCountPct)
		s.Require().EqualValues(1, app.SizePct)
	}
}

func (s *StorageTestSuite) TestAppBySlug() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	app, err := s.storage.App.BySlug(ctx, "app-1")
	s.Require().NoError(err)

	s.Require().EqualValues("App 1", app.Name)
	s.Require().EqualValues(34, app.Size)
	s.Require().EqualValues(1, app.ActionsCount)
	s.Require().False(app.LastActionTime.IsZero())
	s.Require().False(app.FirstActionTime.IsZero())
	s.Require().EqualValues(1, app.ActionsCountPct)
	s.Require().EqualValues(1, app.SizePct)
}

func (s *StorageTestSuite) TestAppActions() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	actions, err := s.storage.App.Actions(ctx, "app-1", 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(1, action.RollupId)
	s.Require().EqualValues(1, action.ActionId)
	s.Require().EqualValues(1, action.TxId)
	s.Require().EqualValues(1, action.SenderId)
	s.Require().EqualValues(34, action.Size)
	s.Require().EqualValues(7316, action.Height)
	s.Require().EqualValues("20b0e6310801e7b2a16c69aace7b1a1d550e5c49c80f546941bb1ac747487fe5", hex.EncodeToString(action.Tx.Hash))
	s.Require().EqualValues(types.ActionTypeRollupDataSubmission, action.ActionType)
	s.Require().NotNil(action.Action.Data)
	s.Require().NotNil(action.Action.Fee)
}
