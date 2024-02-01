// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
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
