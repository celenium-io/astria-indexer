// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

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

func (app *App) Leaderboard(ctx context.Context, fltrs storage.LeaderboardFilters) (rollups []storage.RollupWithStats, err error) {
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
		ColumnExpr("*").
		Offset(fltrs.Offset)

	if len(fltrs.Category) > 0 {
		query = query.Where("category IN (?)", bun.In(fltrs.Category))
	}

	query = sortScope(query, fltrs.SortField, fltrs.Sort)
	query = limitScope(query, fltrs.Limit)
	err = query.Scan(ctx, &rollups)
	return
}
