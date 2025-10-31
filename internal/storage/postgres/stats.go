// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Stats struct {
	db *database.Bun
}

func NewStats(conn *postgres.Storage) Stats {
	return Stats{
		db: conn.Connection(),
	}
}

func (s Stats) Series(ctx context.Context, timeframe storage.Timeframe, name string, req storage.SeriesRequest) (response []storage.SeriesItem, err error) {
	var view string
	switch timeframe {
	case storage.TimeframeHour:
		view = storage.ViewBlockStatsByHour
	case storage.TimeframeDay:
		view = storage.ViewBlockStatsByDay
	case storage.TimeframeMonth:
		view = storage.ViewBlockStatsByMonth
	default:
		return nil, errors.Errorf("unexpected timeframe %s", timeframe)
	}

	query := s.db.DB().NewSelect().Table(view)

	switch name {
	case storage.SeriesDataSize:
		query.ColumnExpr("ts, data_size as value")
	case storage.SeriesTPS:
		query.ColumnExpr("ts, tps as value, tps_max as max, tps_min as min")
	case storage.SeriesBPS:
		query.ColumnExpr("ts, bps as value, bps_max as max, bps_min as min")
	case storage.SeriesRBPS:
		query.ColumnExpr("ts, rbps as value, rbps_max as max, rbps_min as min")
	case storage.SeriesSupplyChange:
		query.ColumnExpr("ts, supply_change as value")
	case storage.SeriesBlockTime:
		query.ColumnExpr("ts, block_time as value")
	case storage.SeriesTxCount:
		query.ColumnExpr("ts, tx_count as value")
	case storage.SeriesBytesInBlock:
		query.ColumnExpr("ts, bytes_in_block as value")
	default:
		return nil, errors.Errorf("unexpected series name: %s", name)
	}

	if !req.From.IsZero() {
		query = query.Where("ts >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("ts < ?", req.To)
	}

	err = query.Limit(100).Scan(ctx, &response)
	return
}

func (s Stats) RollupSeries(ctx context.Context, rollupId uint64, timeframe storage.Timeframe, name string, req storage.SeriesRequest) (response []storage.SeriesItem, err error) {
	var view string
	switch timeframe {
	case storage.TimeframeHour:
		view = storage.ViewRollupStatsByHour
	case storage.TimeframeDay:
		view = storage.ViewRollupStatsByDay
	case storage.TimeframeMonth:
		view = storage.ViewRollupStatsByMonth
	default:
		return nil, errors.Errorf("unexpected timeframe %s", timeframe)
	}

	query := s.db.DB().NewSelect().Table(view).
		Where("rollup_id = ?", rollupId)

	switch name {
	case storage.RollupSeriesActionsCount:
		query.ColumnExpr("ts, actions_count as value")
	case storage.RollupSeriesAvgSize:
		query.ColumnExpr("ts, avg_size as value")
	case storage.RollupSeriesMaxSize:
		query.ColumnExpr("ts, max_size as value")
	case storage.RollupSeriesMinSize:
		query.ColumnExpr("ts, min_size as value")
	case storage.RollupSeriesSize:
		query.ColumnExpr("ts, size as value")
	default:
		return nil, errors.Errorf("unexpected series name: %s", name)
	}

	if !req.From.IsZero() {
		query = query.Where("ts >= ?", req.From)
	}
	if !req.To.IsZero() {
		query = query.Where("ts < ?", req.To)
	}

	err = query.Limit(100).Scan(ctx, &response)
	return
}

func (s Stats) Summary(ctx context.Context) (summary storage.NetworkSummary, err error) {
	err = s.db.DB().NewSelect().Table(storage.ViewBlockStatsByMonth).
		ColumnExpr("sum(data_size) as data_size, sum(supply_change) as supply, sum(tx_count) as tx_count, sum(bytes_in_block) as bytes_in_block").
		ColumnExpr("avg(tps) as tps, avg(bps) as bps, avg(rbps) as rbps, avg(block_time) as block_time").
		Scan(ctx, &summary)
	return
}

func (s Stats) buildSummaryQuery(table string, prevDate, currDate time.Time) (*bun.SelectQuery, *bun.SelectQuery) {
	curr := s.db.DB().NewSelect().
		Table(table).
		Where("ts >= ?", currDate).
		ColumnExpr("sum(data_size) as data_size, sum(tx_count) as tx_count, sum(bytes_in_block) as bytes_in_block").
		ColumnExpr("avg(tps) as tps, avg(bps) as bps, avg(rbps) as rbps, avg(block_time) as block_time")

	prev := s.db.DB().NewSelect().
		Table(table).
		Where("ts < ?", currDate).
		Where("ts >= ?", prevDate).
		ColumnExpr("sum(data_size) as data_size, sum(tx_count) as tx_count, sum(bytes_in_block) as bytes_in_block").
		ColumnExpr("avg(tps) as tps, avg(bps) as bps, avg(rbps) as rbps, avg(block_time) as block_time")

	return curr, prev
}

func (s Stats) SummaryTimeframe(ctx context.Context, timeframe storage.Timeframe) (summary storage.NetworkSummaryWithChange, err error) {
	var (
		currDate, prevDate time.Time
		table              string
	)

	switch timeframe {
	case storage.TimeframeDay:
		currDate = time.Now().AddDate(0, 0, -1).UTC()
		prevDate = currDate.AddDate(0, 0, -1).UTC()
		table = storage.ViewBlockStatsByHour

	case storage.TimeframeWeek:
		currDate = time.Now().AddDate(0, 0, -7).UTC()
		prevDate = currDate.AddDate(0, 0, -7).UTC()
		table = storage.ViewBlockStatsByHour

	case storage.TimeframeMonth:
		currDate = time.Now().AddDate(0, -1, 0).UTC()
		prevDate = currDate.AddDate(0, -1, 0).UTC()
		table = storage.ViewBlockStatsByMonth

	default:
		return summary, errors.Errorf("unknown timeframe: %s", timeframe)
	}

	curr, prev := s.buildSummaryQuery(table, prevDate, currDate)

	err = s.db.DB().NewSelect().
		With("curr", curr).
		With("prev", prev).
		Table("curr", "prev").
		ColumnExpr("curr.data_size as data_size, curr.tx_count as tx_count, curr.bytes_in_block as bytes_in_block, curr.tps as tps, curr.bps as bps, curr.rbps as rbps, curr.block_time as block_time").
		ColumnExpr("case when prev.data_size = 0 then 100 else (curr.data_size - prev.data_size) * 100 / prev.data_size end as data_size_pct").
		ColumnExpr("case when prev.tx_count = 0 then 100 else (curr.tx_count - prev.tx_count)* 100 / prev.tx_count end as tx_count_pct").
		ColumnExpr("case when prev.bytes_in_block = 0 then 100 else (curr.bytes_in_block - prev.bytes_in_block)* 100 / prev.bytes_in_block end as bytes_in_block_pct").
		ColumnExpr("case when prev.tps = 0 then 100 else (curr.tps - prev.tps)* 100 / prev.tps end as tps_pct").
		ColumnExpr("case when prev.bps = 0 then 100 else (curr.bps - prev.bps)* 100 / prev.bps end as bps_pct").
		ColumnExpr("case when prev.rbps = 0 then 100 else (curr.rbps - prev.rbps)* 100 / prev.rbps end as rbps_pct").
		ColumnExpr("case when prev.block_time = 0 then 100 else (curr.block_time - prev.block_time)* 100 / prev.block_time end as block_time_pct").
		Scan(ctx, &summary)

	return
}

func (s Stats) FeeSummary(ctx context.Context) (response []storage.FeeSummary, err error) {
	err = s.db.DB().NewSelect().
		Table(storage.ViewFeeStatsByMonth).
		ColumnExpr("sum(amount) as amount").
		ColumnExpr("min(min_amount) as min_amount").
		ColumnExpr("max(max_amount) as max_amount").
		ColumnExpr("sum(fee_count) as fee_count").
		Column("asset").
		Group("asset").
		Order("amount desc").
		Scan(ctx, &response)
	return
}

func (s Stats) TokenTransferDistribution(ctx context.Context, limit int) (items []storage.TokenTransferDistributionItem, err error) {
	query := s.db.DB().NewSelect().
		Table(storage.ViewTransferStatsByMonth).
		ColumnExpr("sum(transfers_count) as transfers_count").
		ColumnExpr("sum(amount) as amount").
		Column("asset").
		Group("asset").
		Order("amount desc")

	query = limitScope(query, limit)
	err = query.Scan(ctx, &items)
	return
}

func (s Stats) ActiveAddressesCount(ctx context.Context) (val int64, err error) {
	err = s.db.DB().NewSelect().
		Model((*storage.Tx)(nil)).
		ColumnExpr("count(distinct signer_id)").
		Where("time > now() - '1 month'::interval").
		Scan(ctx, &val)
	return
}
