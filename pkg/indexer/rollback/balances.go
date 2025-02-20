// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
)

func rollbackBalances(
	ctx context.Context,
	tx storage.Transaction,
	deletedAddresses []storage.Address,
	height types.Level,
) error {
	ids := make([]uint64, len(deletedAddresses))
	for i := range deletedAddresses {
		ids[i] = deletedAddresses[i].Id
	}

	if err := tx.RollbackBalances(ctx, ids); err != nil {
		return err
	}

	if err := getBalanceUpdates(ctx, tx, height); err != nil {
		return err
	}

	return nil
}

func getBalanceUpdates(
	ctx context.Context,
	tx storage.Transaction,
	height types.Level,
) error {
	balances := make(map[uint64]*storage.Balance)

	updates, err := tx.RollbackBalanceUpdates(ctx, height)
	if err != nil {
		return err
	}

	if len(updates) == 0 {
		return nil
	}

	for i := range updates {
		updateBalances(balances, updates[i])
	}

	arr := make([]storage.Balance, 0)
	for _, b := range balances {
		arr = append(arr, *b)
	}

	return tx.SaveBalances(ctx, arr...)
}

func updateBalances(m map[uint64]*storage.Balance, update storage.BalanceUpdate) {
	if addr, ok := m[update.AddressId]; ok {
		addr.Total = addr.Total.Sub(update.Update)
	} else {
		m[update.AddressId] = &storage.Balance{
			Total:    update.Update.Neg(),
			Id:       update.AddressId,
			Currency: update.Currency,
		}
	}
}
