// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"strconv"

	"github.com/celenium-io/astria-indexer/internal/astria"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
)

// Search -
type Search struct {
	db *database.Bun
}

// NewSearch -
func NewSearch(db *database.Bun) *Search {
	return &Search{
		db: db,
	}
}

func (s *Search) Search(ctx context.Context, query string) (results []storage.SearchResult, err error) {
	text := "%" + query + "%"
	searchQuery := s.db.DB().NewSelect().
		Model((*storage.Validator)(nil)).
		ColumnExpr("id, name as value, 'validator' as type").
		Where("name ILIKE ?", text)

	bridgeQuery := s.db.DB().NewSelect().
		Model((*storage.Bridge)(nil)).
		ColumnExpr("id, asset as value, 'bridge' as type").
		Where("asset ILIKE ?", text)

	appQuery := s.db.DB().NewSelect().
		Model((*storage.App)(nil)).
		ColumnExpr("id, name as value, 'app' as type").
		Where("name ILIKE ?", text)

	searchQuery = searchQuery.
		UnionAll(bridgeQuery).
		UnionAll(appQuery)

	if height, err := strconv.ParseInt(query, 10, 64); err == nil {
		heightQuery := s.db.DB().NewSelect().
			Model((*storage.Block)(nil)).
			ColumnExpr("id, encode(hash, 'hex') as value, 'block' as type").
			Where("height = ?", height)

		searchQuery = searchQuery.UnionAll(heightQuery)
	}

	if hash, err := hex.DecodeString(query); err == nil {
		blockQuery := s.db.DB().NewSelect().
			Model((*storage.Block)(nil)).
			ColumnExpr("id, encode(hash, 'hex') as value, 'block' as type").
			Where("hash = ?", hash)
		txQuery := s.db.DB().NewSelect().
			Model((*storage.Tx)(nil)).
			ColumnExpr("id, encode(hash, 'hex') as value, 'tx' as type").
			Where("hash = ?", hash)
		rollupQuery := s.db.DB().NewSelect().
			Model((*storage.Rollup)(nil)).
			ColumnExpr("id, encode(astria_id, 'hex') as value, 'rollup' as type").
			Where("astria_id = ?", hash)

		searchQuery = searchQuery.
			UnionAll(blockQuery).
			UnionAll(txQuery).
			UnionAll(rollupQuery)
	}

	if astria.IsAddress(query) {
		addressQuery := s.db.DB().NewSelect().
			Model((*storage.Address)(nil)).
			ColumnExpr("id, hash as value, 'address' as type").
			Where("hash = ?", query)

		validatorQuery := s.db.DB().NewSelect().
			Model((*storage.Validator)(nil)).
			ColumnExpr("id, name as value, 'validator' as type").
			Where("address = ?", query)

		searchQuery = searchQuery.
			UnionAll(addressQuery).
			UnionAll(validatorQuery)
	}

	if decoded, err := base64.StdEncoding.DecodeString(query); err == nil {
		rollupQuery := s.db.DB().NewSelect().
			Model((*storage.Rollup)(nil)).
			ColumnExpr("id, encode(astria_id, 'hex') as value, 'rollup' as type").
			Where("astria_id = ?", decoded)

		searchQuery = searchQuery.
			UnionAll(rollupQuery)
	}

	err = s.db.DB().NewSelect().
		TableExpr("(?) as search", searchQuery).
		Limit(10).
		Offset(0).
		Scan(ctx, &results)

	return
}
