// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Tx -
type Tx struct {
	*postgres.Table[*storage.Tx]
}

// NewTx -
func NewTx(db *database.Bun) *Tx {
	return &Tx{
		Table: postgres.NewTable[*storage.Tx](db),
	}
}

func (tx *Tx) ByHash(ctx context.Context, hash []byte) (transaction storage.Tx, err error) {
	query := tx.DB().NewSelect().Model((*storage.Tx)(nil)).
		Where("hash = ?", hash).
		Limit(1)

	err = tx.DB().NewSelect().
		TableExpr("(?) as tx", query).
		ColumnExpr("tx.*").
		ColumnExpr("address.hash as signer__hash").
		Join("left join address on address.id = tx.signer_id").
		Scan(ctx, &transaction)
	return
}

func (tx *Tx) ByHeight(ctx context.Context, height types.Level, limit, offset int) (txs []storage.Tx, err error) {
	query := tx.DB().NewSelect().Model(&txs).
		Where("tx.height = ?", height).
		Relation("Signer")

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	err = query.Scan(ctx)
	return
}

func (tx *Tx) Filter(ctx context.Context, fltrs storage.TxFilter) (txs []storage.Tx, err error) {
	query := tx.DB().NewSelect().Model(&txs).Offset(fltrs.Offset).Relation("Signer")
	query = txFilter(query, fltrs)

	err = query.Scan(ctx)
	return
}

func (tx *Tx) ByAddress(ctx context.Context, addressId uint64, fltrs storage.TxFilter) (txs []storage.Tx, err error) {
	query := tx.DB().NewSelect().
		Model(&txs).
		Where("signer_id = ?", addressId).
		Relation("Signer").
		Offset(fltrs.Offset)

	query = txFilter(query, fltrs)

	err = query.Scan(ctx)
	return txs, err
}
