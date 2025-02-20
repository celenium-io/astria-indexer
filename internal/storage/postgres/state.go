// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// State -
type State struct {
	*postgres.Table[*storage.State]
}

// NewState -
func NewState(db *postgres.Storage) *State {
	return &State{
		Table: postgres.NewTable[*storage.State](db.Connection()),
	}
}

// ByName -
func (s *State) ByName(ctx context.Context, name string) (state storage.State, err error) {
	err = s.DB().NewSelect().Model(&state).Where("name = ?", name).Scan(ctx)
	return
}
