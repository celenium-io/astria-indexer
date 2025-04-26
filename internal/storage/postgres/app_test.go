// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
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

		apps, err := s.App.Leaderboard(ctx, storage.LeaderboardFilters{
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
		s.Require().NotNil(app.Bridge)
		s.Require().EqualValues("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", app.Bridge.Hash)
		s.Require().NotNil(app.Rollup)
		s.Require().EqualValues("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539", hex.EncodeToString(app.Rollup.AstriaId))
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

		apps, err := s.App.Leaderboard(ctx, storage.LeaderboardFilters{
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
		s.Require().NotNil(app.Bridge)
		s.Require().EqualValues("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", app.Bridge.Hash)
		s.Require().NotNil(app.Rollup)
		s.Require().EqualValues("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539", hex.EncodeToString(app.Rollup.AstriaId))
	}
}

func (s *StorageTestSuite) TestAppBySlug() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	app, err := s.App.BySlug(ctx, "app-1")
	s.Require().NoError(err)

	s.Require().EqualValues("App 1", app.Name)
	s.Require().EqualValues(34, app.Size)
	s.Require().EqualValues(1, app.ActionsCount)
	s.Require().False(app.LastActionTime.IsZero())
	s.Require().False(app.FirstActionTime.IsZero())
	s.Require().NotNil(app.Bridge)
	s.Require().EqualValues("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", app.Bridge.Hash)
	s.Require().NotNil(app.Rollup)
	s.Require().EqualValues("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539", hex.EncodeToString(app.Rollup.AstriaId))
}

func (s *StorageTestSuite) TestAppByRollupId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err := s.storage.Connection().Exec(ctx, "REFRESH MATERIALIZED VIEW leaderboard;")
	s.Require().NoError(err)

	app, err := s.App.ByRollupId(ctx, 1)
	s.Require().NoError(err)

	s.Require().EqualValues("App 1", app.Name)
	s.Require().EqualValues(34, app.Size)
	s.Require().EqualValues(1, app.ActionsCount)
	s.Require().False(app.LastActionTime.IsZero())
	s.Require().False(app.FirstActionTime.IsZero())
	s.Require().NotNil(app.Bridge)
	s.Require().EqualValues("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", app.Bridge.Hash)
	s.Require().Nil(app.Rollup)
}
