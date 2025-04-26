// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

func saveAddresses(
	ctx context.Context,
	tx storage.Transaction,
	addresses map[string]*storage.Address,
) (map[string]uint64, int64, error) {
	if len(addresses) == 0 {
		return nil, 0, nil
	}

	data := make([]*storage.Address, 0, len(addresses))
	for key := range addresses {
		data = append(data, addresses[key])
	}

	totalAccounts, err := tx.SaveAddresses(ctx, data...)
	if err != nil {
		return nil, 0, err
	}

	addToId := make(map[string]uint64)
	balances := make([]storage.Balance, 0)
	for i := range data {
		addToId[data[i].String()] = data[i].Id
		for j := range data[i].Balance {
			data[i].Balance[j].Id = data[i].Id
			balances = append(balances, *data[i].Balance[j])
		}
	}
	err = tx.SaveBalances(ctx, balances...)
	return addToId, totalAccounts, err
}
