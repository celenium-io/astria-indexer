// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/dipdup-net/go-lib/database"
)

// Constant -
type Constant struct {
	db *database.Bun
}

// NewConstant -
func NewConstant(db *database.Bun) *Constant {
	return &Constant{
		db: db,
	}
}

func (constant *Constant) Get(ctx context.Context, module types.ModuleName, name string) (c storage.Constant, err error) {
	err = constant.db.DB().NewSelect().Model(&c).
		Where("module = ?", module).
		Where("name = ?", name).
		Scan(ctx)
	return
}

func (constant *Constant) ByModule(ctx context.Context, module types.ModuleName) (c []storage.Constant, err error) {
	err = constant.db.DB().NewSelect().Model(&c).
		Where("module = ?", module).
		Scan(ctx)
	return
}

func (constant *Constant) All(ctx context.Context) (c []storage.Constant, err error) {
	err = constant.db.DB().NewSelect().Model(&c).Scan(ctx)
	return
}

func (constant *Constant) IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
