// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Validator -
type Validator struct {
	*postgres.Table[*storage.Validator]
}

// NewValidator -
func NewValidator(db *database.Bun) *Validator {
	return &Validator{
		Table: postgres.NewTable[*storage.Validator](db),
	}
}

func (v *Validator) ListByPower(ctx context.Context, limit, offset int, order sdk.SortOrder) (validators []storage.Validator, err error) {
	query := v.DB().NewSelect().
		Model(&validators).
		Offset(offset)
	query = limitScope(query, limit)
	query = sortScope(query, "power", order)

	err = query.Scan(ctx)
	return
}
