// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
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

func (app *App) BySlug(ctx context.Context, slug string) (result storage.AppWithStats, err error) {
	err = app.DB().NewSelect().
		Table(storage.ViewLeaderboard).
		Where("slug = ?", slug).
		Limit(1).
		Scan(ctx, &result)
	return
}

func (app *App) Actions(ctx context.Context, slug string, limit, offset int, sort sdk.SortOrder) (result []storage.RollupAction, err error) {
	var id uint64
	if err = app.DB().NewSelect().
		Column("id").
		Model((*storage.App)(nil)).
		Where("slug = ?", slug).
		Limit(1).
		Scan(ctx, &id); err != nil {
		return
	}

	var appIds []storage.AppId
	if err = app.DB().NewSelect().
		Model(&appIds).
		Where("app_id = ?", id).
		Scan(ctx); err != nil {
		return
	}

	subQuery := app.DB().NewSelect().
		Model((*storage.RollupAction)(nil))

	subQuery = limitScope(subQuery, limit)
	subQuery = offsetScope(subQuery, offset)
	subQuery = sortScope(subQuery, "action_id", sort)

	err = app.DB().NewSelect().
		TableExpr("(?) as rollup_action", subQuery).
		ColumnExpr("rollup_action.*").
		ColumnExpr("action.data as action__data, action.position as action__position").
		ColumnExpr("tx.hash as tx__hash").
		ColumnExpr("fee.asset as action__fee__asset, fee.amount as action__fee__amount").
		Join("left join fee on fee.action_id = rollup_action.action_id").
		Join("left join action on action.id = rollup_action.action_id").
		Join("left join tx on tx.id = rollup_action.tx_id").
		Scan(ctx, &result)

	return
}
