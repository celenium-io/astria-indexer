// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package genesis

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
)

func (module *Module) save(ctx context.Context, data parsedData) error {
	start := time.Now()
	module.Log.Info().Uint64("height", uint64(data.block.Height)).Msg("saving block...")
	tx, err := postgres.BeginTransaction(ctx, module.storage)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	if err := tx.SaveConstants(ctx, data.constants...); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, &data.block); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, data.block.Stats); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.SaveTransactions(ctx, data.block.Txs...); err != nil {
		return tx.HandleError(ctx, err)
	}

	var totalAccounts int64
	if len(data.addresses) > 0 {
		entities := make([]*storage.Address, 0, len(data.addresses))
		for key := range data.addresses {
			entities = append(entities, data.addresses[key])
		}

		totalAccounts, err = tx.SaveAddresses(ctx, entities...)
		if err != nil {
			return tx.HandleError(ctx, err)
		}

		balances := make([]storage.Balance, 0)
		for i := range entities {
			for _, balance := range entities[i].Balance {
				balances = append(balances, *balance)
			}
		}
		if err := tx.SaveBalances(ctx, balances...); err != nil {
			return tx.HandleError(ctx, err)
		}

		for i := range data.balanceUpdates {
			if addr, ok := data.addresses[data.balanceUpdates[i].Address.String()]; ok {
				data.balanceUpdates[i].AddressId = addr.Id
			}
		}
		if err := tx.SaveBalanceUpdates(ctx, data.balanceUpdates...); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.SaveValidators(ctx, data.validators...); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Add(ctx, &storage.State{
		Name:            module.indexerName,
		LastHeight:      data.block.Height,
		LastTime:        data.block.Time,
		LastHash:        data.block.Hash,
		ChainId:         data.block.ChainId,
		TotalTx:         data.block.Stats.TxCount,
		TotalSupply:     data.block.Stats.SupplyChange,
		TotalAccounts:   totalAccounts,
		TotalValidators: len(data.validators),
	}); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	module.Log.Info().
		Uint64("height", data.block.Id).
		Str("block_fee", data.block.Stats.Fee.String()).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block saved")
	return nil
}
