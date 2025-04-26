// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/pkg/errors"
)

func (module *Module) saveTransactions(
	ctx context.Context,
	tx storage.Transaction,
	addrToId map[string]uint64,
	txs ...*storage.Tx,
) error {
	if len(txs) == 0 {
		return nil
	}

	for i := range txs {
		if signerId, ok := addrToId[txs[i].Signer.String()]; ok {
			txs[i].SignerId = signerId
		} else {
			return errors.Errorf("unknown signer id")
		}
	}

	return tx.SaveTransactions(ctx, txs...)
}
