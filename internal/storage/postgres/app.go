// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"fmt"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

// App -
type App struct {
	*postgres.Table[*storage.App]
}

// NewApp -
func NewApp(db *database.Bun) *App {
	return &App{
		Table: postgres.NewTable[*storage.App](db),
	}
}

func (app *App) Leaderboard(ctx context.Context, fltrs storage.LeaderboardFilters) (rollups []storage.AppWithStats, err error) {
	switch fltrs.SortField {
	case columnTime:
		fltrs.SortField = "last_time"
	case columnSize, columnActionsCount:
	case "":
		fltrs.SortField = columnSize
	default:
		return nil, errors.Errorf("unknown sort field: %s", fltrs.SortField)
	}

	query := app.DB().NewSelect().
		Table(storage.ViewLeaderboard).
		Offset(fltrs.Offset)

	if len(fltrs.Category) > 0 {
		query = query.Where("category IN (?)", bun.In(fltrs.Category))
	}

	query = sortScope(query, fmt.Sprintf("%s.%s", storage.ViewLeaderboard, fltrs.SortField), fltrs.Sort)
	query = limitScope(query, fltrs.Limit)

	query = query.
		ColumnExpr("leaderboard.*").
		ColumnExpr("address.hash as bridge__hash").
		ColumnExpr("rollup.astria_id as rollup__astria_id").
		Join("left join address on native_bridge_id = address.id").
		Join("left join rollup on rollup.id = rollup_id")
	err = query.Scan(ctx, &rollups)
	return
}

func (app *App) BySlug(ctx context.Context, slug string) (result storage.AppWithStats, err error) {
	err = app.DB().NewSelect().
		Table(storage.ViewLeaderboard).
		ColumnExpr("leaderboard.*").
		ColumnExpr("address.hash as bridge__hash").
		ColumnExpr("rollup.astria_id as rollup__astria_id").
		Join("left join address on native_bridge_id = address.id").
		Join("left join rollup on rollup.id = rollup_id").
		Where("slug = ?", slug).
		Limit(1).
		Scan(ctx, &result)
	return
}

func (app *App) ByRollupId(ctx context.Context, rollupId uint64) (result storage.AppWithStats, err error) {
	err = app.DB().NewSelect().
		Table(storage.ViewLeaderboard).
		ColumnExpr("leaderboard.*").
		ColumnExpr("address.hash as bridge__hash").
		Where("rollup_id = ?", rollupId).
		Join("left join address on native_bridge_id = address.id").
		Limit(1).
		Scan(ctx, &result)
	return
}
