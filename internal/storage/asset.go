// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type IAsset interface {
	List(ctx context.Context, limit int, offset int, sortBy string, order sdk.SortOrder) ([]Asset, error)
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Asset struct {
	bun.BaseModel `bun:"asset"`

	Asset         string          `bun:"asset"`
	Fee           decimal.Decimal `bun:"fee"`
	FeeCount      int             `bun:"fee_count"`
	Transferred   decimal.Decimal `bun:"transferred"`
	TransferCount int             `bun:"transfer_count"`
	Supply        decimal.Decimal `bun:"supply"`
}
