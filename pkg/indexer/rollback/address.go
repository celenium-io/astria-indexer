// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
)

func rollbackAddress(
	ctx context.Context,
	tx storage.Transaction,
	height types.Level,
	addressActions []storage.AddressAction,
	txs []storage.Tx,
) (int, error) {
	deletedAddresses, err := tx.RollbackAddresses(ctx, height)
	if err != nil {
		return 0, err
	}

	if err := rollbackBalances(ctx, tx, deletedAddresses, height); err != nil {
		return 0, err
	}

	addresses := make(map[uint64]*storage.Address, 0)
	for i := range addressActions {
		if addr, ok := addresses[addressActions[i].AddressId]; ok {
			addr.ActionsCount -= 1
		} else {
			nonce, err := tx.LastNonce(ctx, addressActions[i].AddressId)
			if err != nil {
				return 0, err
			}
			addresses[addressActions[i].AddressId] = &storage.Address{
				Id:           addressActions[i].AddressId,
				ActionsCount: -1,
				Nonce:        nonce,
			}
		}
	}

	for i := range txs {
		if addr, ok := addresses[txs[i].SignerId]; ok {
			addr.SignedTxCount -= 1
			addr.ActionsCount -= 1
		} else {
			addresses[addressActions[i].AddressId] = &storage.Address{
				Id:            txs[i].SignerId,
				SignedTxCount: -1,
				ActionsCount:  -1,
			}
		}
	}

	arr := make([]*storage.Address, 0)
	for i := range addresses {
		arr = append(arr, addresses[i])
	}

	if err := tx.UpdateAddresses(ctx, arr...); err != nil {
		return 0, err
	}

	return len(deletedAddresses), nil
}
