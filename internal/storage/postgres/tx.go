// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Tx -
type Tx struct {
	*postgres.Table[*storage.Tx]
}

// NewTx -
func NewTx(db *postgres.Storage) *Tx {
	return &Tx{
		Table: postgres.NewTable[*storage.Tx](db.Connection()),
	}
}

func (tx *Tx) ByHash(ctx context.Context, hash []byte) (transaction storage.Tx, err error) {
	query := tx.DB().NewSelect().Model((*storage.Tx)(nil)).
		Where("hash = ?", hash).
		Limit(1)

	q := tx.DB().NewSelect().
		TableExpr("(?) as tx", query).
		ColumnExpr("tx.*").
		ColumnExpr("address.hash as signer__hash").
		Join("left join address on address.id = tx.signer_id")

	q = joinCelestials(q, "signer__", "tx.signer_id")
	err = q.Scan(ctx, &transaction)
	return
}

func (tx *Tx) ByHeight(ctx context.Context, height types.Level, limit, offset int) (txs []storage.Tx, err error) {
	query := tx.DB().NewSelect().Model((*storage.Tx)(nil)).
		Where("tx.height = ?", height)

	query = limitScope(query, limit)
	if offset > 0 {
		query = query.Offset(offset)
	}

	q := tx.DB().NewSelect().
		TableExpr("(?) as tx", query).
		ColumnExpr("tx.*").
		ColumnExpr("address.hash as signer__hash").
		Join("left join address on address.id = tx.signer_id")

	q = joinCelestials(q, "signer__", "tx.signer_id")
	err = q.Scan(ctx, &txs)
	return
}

func (tx *Tx) Filter(ctx context.Context, fltrs storage.TxFilter) (txs []storage.Tx, err error) {
	query := tx.DB().NewSelect().
		Model(&txs).
		ColumnExpr("tx.*").
		ColumnExpr("address.hash as signer__hash").
		Join("left join address on address.id = tx.signer_id")

	query = txFilter(query, fltrs)
	query = joinCelestials(query, "signer__", "tx.signer_id")
	err = query.Scan(ctx)
	return
}

func (tx *Tx) ByAddress(ctx context.Context, addressId uint64, fltrs storage.TxFilter) (txs []storage.Tx, err error) {
	query := tx.DB().NewSelect().
		Model(&txs).
		Where("signer_id = ?", addressId).
		ColumnExpr("tx.*").
		ColumnExpr("address.hash as signer__hash").
		Join("left join address on address.id = tx.signer_id")

	query = txFilter(query, fltrs)
	query = joinCelestials(query, "signer__", "tx.signer_id")
	err = query.Scan(ctx)
	return txs, err
}
