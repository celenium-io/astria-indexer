// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Timeframe string

const (
	TimeframeHour  Timeframe = "hour"
	TimeframeDay   Timeframe = "day"
	TimeframeWeek  Timeframe = "week"
	TimeframeMonth Timeframe = "month"
)

type TPS struct {
	Low               float64
	High              float64
	Current           float64
	ChangeLastHourPct float64
}

type TxCountForLast24hItem struct {
	Time    time.Time `bun:"ts"`
	TxCount int64     `bun:"tx_count"`
	TPS     float64   `bun:"tps"`
}

type SeriesRequest struct {
	From time.Time
	To   time.Time
}

func NewSeriesRequest(from, to int64) (sr SeriesRequest) {
	if from > 0 {
		sr.From = time.Unix(from, 0).UTC()
	}
	if to > 0 {
		sr.To = time.Unix(to, 0).UTC()
	}
	return
}

type SeriesItem struct {
	Time  time.Time `bun:"ts"`
	Value string    `bun:"value"`
	Max   string    `bun:"max"`
	Min   string    `bun:"min"`
}

const (
	SeriesDataSize     = "data_size"
	SeriesTPS          = "tps"
	SeriesBPS          = "bps"
	SeriesRBPS         = "rbps"
	SeriesSupplyChange = "supply_change"
	SeriesBlockTime    = "block_time"
	SeriesTxCount      = "tx_count"
	SeriesBytesInBlock = "bytes_in_block"

	RollupSeriesActionsCount = "actions_count"
	RollupSeriesSize         = "size"
	RollupSeriesAvgSize      = "avg_size"
	RollupSeriesMinSize      = "min_size"
	RollupSeriesMaxSize      = "max_size"
)

type NetworkSummary struct {
	DataSize     int64           `bun:"data_size"`
	TPS          float64         `bun:"tps"`
	BPS          float64         `bun:"bps"`
	RBPS         float64         `bun:"rbps"`
	Fee          decimal.Decimal `bun:"fee"`
	Supply       decimal.Decimal `bun:"supply"`
	BlockTime    float64         `bun:"block_time"`
	TxCount      int64           `bun:"tx_count"`
	BytesInBlock int64           `bun:"bytes_in_block"`
}

type NetworkSummaryWithChange struct {
	DataSize        int64   `bun:"data_size"`
	DataSizePct     float64 `bun:"data_size_pct"`
	TPS             float64 `bun:"tps"`
	TPSPct          float64 `bun:"tps_pct"`
	BPS             float64 `bun:"bps"`
	BPSPct          float64 `bun:"bps_pct"`
	RBPS            float64 `bun:"rbps"`
	RBPSPct         float64 `bun:"rbps_pct"`
	BlockTime       float64 `bun:"block_time"`
	BlockTimePct    float64 `bun:"block_time_pct"`
	TxCount         int64   `bun:"tx_count"`
	TxCountPct      float64 `bun:"tx_count_pct"`
	BytesInBlock    int64   `bun:"bytes_in_block"`
	BytesInBlockPct float64 `bun:"bytes_in_block_pct"`
}

type RollupSummary struct {
	ActionsCount int64 `bun:"actions_count"`
	Size         int64 `bun:"size"`
	AvgSize      int64 `bun:"avg_size"`
	MinSize      int64 `bun:"min_size"`
	MaxSize      int64 `bun:"max_size"`
}

type FeeSummary struct {
	Asset     string `bun:"asset"`
	Amount    string `bun:"amount"`
	MinAmount string `bun:"min_amount"`
	MaxAmount string `bun:"max_amount"`
	FeeCount  int64  `bun:"fee_count"`
}

type TokenTransferDistributionItem struct {
	Asset          string `bun:"asset"`
	Amount         string `bun:"amount"`
	TransfersCount int64  `bun:"transfers_count"`
}

type Candle struct {
	Time         time.Time       `bun:"time"`
	Open         decimal.Decimal `bun:"open"`
	Close        decimal.Decimal `bun:"close"`
	High         decimal.Decimal `bun:"high"`
	Low          decimal.Decimal `bun:"low"`
	CurrencyPair string          `bun:"currency_pair"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IStats interface {
	Summary(ctx context.Context) (NetworkSummary, error)
	SummaryTimeframe(ctx context.Context, timeframe Timeframe) (NetworkSummaryWithChange, error)
	Series(ctx context.Context, timeframe Timeframe, name string, req SeriesRequest) ([]SeriesItem, error)
	RollupSeries(ctx context.Context, rollupId uint64, timeframe Timeframe, name string, req SeriesRequest) ([]SeriesItem, error)
	FeeSummary(ctx context.Context) ([]FeeSummary, error)
	TokenTransferDistribution(ctx context.Context, limit int) ([]TokenTransferDistributionItem, error)
	ActiveAddressesCount(ctx context.Context) (int64, error)
}
