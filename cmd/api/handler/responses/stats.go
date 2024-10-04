// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

type SeriesItem struct {
	Time  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"          swaggertype:"string"`
	Value string    `example:"0.17632"                   format:"string"    json:"value"         swaggertype:"string"`
	Max   string    `example:"0.17632"                   format:"string"    json:"max,omitempty" swaggertype:"string"`
	Min   string    `example:"0.17632"                   format:"string"    json:"min,omitempty" swaggertype:"string"`
}

func NewSeriesItem(item storage.SeriesItem) SeriesItem {
	return SeriesItem{
		Time:  item.Time,
		Value: item.Value,
		Max:   item.Max,
		Min:   item.Min,
	}
}

type RollupSeriesItem struct {
	Time  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"  swaggertype:"string"`
	Value string    `example:"0.17632"                   format:"string"    json:"value" swaggertype:"string"`
}

func NewRollupSeriesItem(item storage.SeriesItem) RollupSeriesItem {
	return RollupSeriesItem{
		Time:  item.Time,
		Value: item.Value,
	}
}

type NetworkSummary struct {
	DataSize     int64   `example:"1000000" format:"integer" json:"data_size"`
	TxCount      int64   `example:"100"     format:"integer" json:"tx_count"`
	BytesInBlock int64   `example:"1024"    format:"integer" json:"bytes_in_block"`
	TPS          float64 `example:"0.17632" format:"number"  json:"tps"`
	BPS          float64 `example:"0.17632" format:"number"  json:"bps"`
	RBPS         float64 `example:"0.17632" format:"number"  json:"rbps"`
	BlockTime    float64 `example:"2345"    format:"number"  json:"block_time"`
	Fee          string  `example:"1012012" format:"string"  json:"fee"`
	Supply       string  `example:"1029129" format:"string"  json:"supply"`
}

func NewNetworkSummary(summary storage.NetworkSummary) NetworkSummary {
	return NetworkSummary{
		DataSize:     summary.DataSize,
		TxCount:      summary.TxCount,
		BytesInBlock: summary.BytesInBlock,
		TPS:          summary.TPS,
		BPS:          summary.BPS,
		RBPS:         summary.RBPS,
		BlockTime:    summary.BlockTime,
		Fee:          summary.Fee.String(),
		Supply:       summary.Supply.String(),
	}
}

type NetworkSummaryWithChange struct {
	DataSize        int64   `example:"1000000" format:"integer" json:"data_size"`
	DataSizePct     float64 `example:"17.632"  format:"number"  json:"data_size_pct"`
	TxCount         int64   `example:"100"     format:"integer" json:"tx_count"`
	TxCountPct      float64 `example:"17.632"  format:"number"  json:"tx_count_pct"`
	BytesInBlock    int64   `example:"1024"    format:"integer" json:"bytes_in_block"`
	BytesInBlockPct float64 `example:"17.632"  format:"number"  json:"bytes_in_block_pct"`
	TPS             float64 `example:"0.17632" format:"number"  json:"tps"`
	TPSPct          float64 `example:"17.632"  format:"number"  json:"tps_pct"`
	BPS             float64 `example:"0.17632" format:"number"  json:"bps"`
	BPSPct          float64 `example:"17.632"  format:"number"  json:"bps_pct"`
	RBPS            float64 `example:"0.17632" format:"number"  json:"rbps"`
	RBPSPct         float64 `example:"17.632"  format:"number"  json:"rbps_pct"`
	BlockTime       float64 `example:"2345"    format:"number"  json:"block_time"`
	BlockTimePct    float64 `example:"17.632"  format:"number"  json:"block_time_pct"`
}

func NewNetworkSummaryWithChange(summary storage.NetworkSummaryWithChange) NetworkSummaryWithChange {
	return NetworkSummaryWithChange{
		DataSize:        summary.DataSize,
		DataSizePct:     summary.DataSizePct,
		TxCount:         summary.TxCount,
		TxCountPct:      summary.TxCountPct,
		BytesInBlock:    summary.BytesInBlock,
		BytesInBlockPct: summary.BytesInBlockPct,
		TPS:             summary.TPS,
		TPSPct:          summary.TPSPct,
		BPS:             summary.BPS,
		BPSPct:          summary.BPSPct,
		RBPS:            summary.RBPS,
		RBPSPct:         summary.RBPSPct,
		BlockTime:       summary.BlockTime,
		BlockTimePct:    summary.BlockTimePct,
	}
}

type FeeSummary struct {
	Asset     string `example:"nria"    format:"string"  json:"asset"`
	Amount    string `example:"1000000" format:"integer" json:"amount"`
	MinAmount string `example:"1000000" format:"integer" json:"min_amount"`
	MaxAmount string `example:"1000000" format:"integer" json:"max_amount"`
	FeeCount  int64  `example:"1000000" format:"integer" json:"fee_count"`
}

func NewFeeSummary(summary storage.FeeSummary) FeeSummary {
	return FeeSummary{
		Asset:     summary.Asset,
		Amount:    summary.Amount,
		MinAmount: summary.MinAmount,
		MaxAmount: summary.MaxAmount,
		FeeCount:  summary.FeeCount,
	}
}
