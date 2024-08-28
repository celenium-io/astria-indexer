// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Fee -
type Fee struct {
	*postgres.Table[*storage.Fee]
}

// NewFee -
func NewFee(db *database.Bun) *Fee {
	return &Fee{
		Table: postgres.NewTable[*storage.Fee](db),
	}
}
