// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/aopoltorzhicky/astria/internal/storage"
)

func saveAction(
	ctx context.Context,
	tx storage.Transaction,
	actions []*storage.Action,
	addrToId map[string]uint64,
) error {
	if len(actions) == 0 {
		return nil
	}

	if err := tx.SaveActions(ctx, actions...); err != nil {
		return err
	}

	var (
		rollupActions  = make([]*storage.RollupAction, 0)
		addrActions    = make([]*storage.AddressAction, 0)
		balanceUpdates = make([]storage.BalanceUpdate, 0)
	)
	for i := range actions {
		if actions[i].RollupAction != nil {
			actions[i].RollupAction.ActionId = actions[i].Id
			actions[i].RollupAction.RollupId = actions[i].RollupAction.Rollup.Id
			actions[i].RollupAction.TxId = actions[i].TxId
			rollupActions = append(rollupActions, actions[i].RollupAction)
		}

		for j := range actions[i].Addresses {
			actions[i].Addresses[j].ActionId = actions[i].Id
			actions[i].Addresses[j].AddressId = actions[i].Addresses[j].Address.Id
			actions[i].Addresses[j].TxId = actions[i].TxId
		}
		addrActions = append(addrActions, actions[i].Addresses...)

		for j := range actions[i].BalanceUpdates {
			actions[i].BalanceUpdates[j].AddressId = actions[i].BalanceUpdates[j].Address.Id
		}
		balanceUpdates = append(balanceUpdates, actions[i].BalanceUpdates...)
	}

	if err := tx.SaveRollupActions(ctx, rollupActions...); err != nil {
		return err
	}
	if err := tx.SaveAddressActions(ctx, addrActions...); err != nil {
		return err
	}
	if err := tx.SaveBalanceUpdates(ctx, balanceUpdates...); err != nil {
		return err
	}

	return nil
}
