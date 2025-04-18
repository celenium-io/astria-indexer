// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Address -
type Address struct {
	*postgres.Table[*storage.Address]
}

// NewAddress -
func NewAddress(db *postgres.Storage) *Address {
	return &Address{
		Table: postgres.NewTable[*storage.Address](db.Connection()),
	}
}

// ByHash -
func (a *Address) ByHash(ctx context.Context, hash string) (address storage.Address, err error) {
	err = a.DB().NewSelect().
		Model(&address).
		Where("hash = ?", hash).
		Join("left join celestial on celestial.address_id = address.id and celestial.status = 'PRIMARY'").
		Relation("Balance").
		Scan(ctx)
	return
}

func (a *Address) ListWithBalance(ctx context.Context, fltrs storage.AddressListFilter) (address []storage.Address, err error) {
	query := a.DB().NewSelect().
		Model(&address).
		Offset(fltrs.Offset)

	if fltrs.Asset != "" {
		query = query.Relation("Balance", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("currency = ?", fltrs.Asset)
		})
	} else {
		query = query.Relation("Balance", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("currency = 'nria'")
		})
	}

	query = addressListFilter(query, fltrs)
	query.Join("left join celestial on celestial.address_id = address.id and celestial.status = 'PRIMARY'")

	err = query.Scan(ctx)
	return
}
