// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Constant -
type Constant struct {
	db *postgres.Storage
}

// NewConstant -
func NewConstant(db *postgres.Storage) *Constant {
	return &Constant{
		db: db,
	}
}

func (constant *Constant) Get(ctx context.Context, module types.ModuleName, name string) (c storage.Constant, err error) {
	err = constant.db.Connection().DB().NewSelect().Model(&c).
		Where("module = ?", module).
		Where("name = ?", name).
		Scan(ctx)
	return
}

func (constant *Constant) ByModule(ctx context.Context, module types.ModuleName) (c []storage.Constant, err error) {
	err = constant.db.Connection().DB().NewSelect().Model(&c).
		Where("module = ?", module).
		Scan(ctx)
	return
}

func (constant *Constant) All(ctx context.Context) (c []storage.Constant, err error) {
	err = constant.db.Connection().DB().NewSelect().Model(&c).Scan(ctx)
	return
}

func (constant *Constant) IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
