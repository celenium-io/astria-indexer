// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Transfer -
type Transfer struct {
	*postgres.Table[*storage.Transfer]
}

// NewTransfer -
func NewTransfer(db *postgres.Storage) *Transfer {
	return &Transfer{
		Table: postgres.NewTable[*storage.Transfer](db.Connection()),
	}
}
