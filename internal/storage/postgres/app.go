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

	if len(appIds) == 0 {
		return
	}

	subQuery := app.DB().NewSelect().
		Model((*storage.RollupAction)(nil))

	subQuery.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
		for i := range appIds {
			sq = sq.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("rollup_id = ?", appIds[i].RolllupId).Where("sender_id = ?", appIds[i].AddressId)
			})
		}

		return sq
	})

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

func (app *App) Series(ctx context.Context, slug string, timeframe storage.Timeframe, column string, req storage.SeriesRequest) (items []storage.SeriesItem, err error) {
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

	if len(appIds) == 0 {
		return
	}

	query := app.DB().NewSelect().Order("time desc").Limit(100).Group("time")

	switch timeframe {
	case storage.TimeframeHour:
		query = query.Table("app_stats_by_hour")
	case storage.TimeframeDay:
		query = query.Table("app_stats_by_day")
	case storage.TimeframeMonth:
		query = query.Table("app_stats_by_month")
	default:
		return nil, errors.Errorf("invalid timeframe: %s", timeframe)
	}

	switch column {
	case "actions_count":
		query = query.ColumnExpr("sum(actions_count) as value, time as ts")
	case "size":
		query = query.ColumnExpr("sum(size) as value, time as ts")
	case "size_per_action":
		query = query.ColumnExpr("(sum(size) / sum(actions_count)) as value, time as ts")
	default:
		return nil, errors.Errorf("invalid column: %s", column)
	}

	if !req.From.IsZero() {
		query = query.Where("time >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("time < ?", req.To)
	}

	query.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
		for i := range appIds {
			sq = sq.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("rollup_id = ?", appIds[i].RolllupId).Where("sender_id = ?", appIds[i].AddressId)
			})
		}

		return sq
	})

	err = query.Scan(ctx, &items)

	return
}
