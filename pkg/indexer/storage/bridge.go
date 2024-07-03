// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/pkg/errors"
)

func saveBridges(
	ctx context.Context,
	tx storage.Transaction,
	addrToId map[string]uint64,
	bridges []*storage.Bridge,
) error {
	if len(bridges) == 0 {
		return nil
	}

	for i := range bridges {
		if id, ok := addrToId[bridges[i].Address.Hash]; ok {
			bridges[i].AddressId = id
		} else {
			return errors.Errorf("unknown bridge address: %s", bridges[i].Address.Hash)
		}

		if id, ok := addrToId[bridges[i].Sudo.Hash]; ok {
			bridges[i].SudoId = id
		} else {
			return errors.Errorf("unknown sudo bridge address: %s", bridges[i].Sudo.Hash)
		}

		if id, ok := addrToId[bridges[i].Withdrawer.Hash]; ok {
			bridges[i].WithdrawerId = id
		} else {
			return errors.Errorf("unknown withdrawer bridge address: %s", bridges[i].Withdrawer.Hash)
		}

		if bridges[i].Rollup.Id == 0 {
			rollup, err := tx.GetRollup(ctx, bridges[i].Rollup.AstriaId)
			if err != nil {
				return err
			}
			bridges[i].RollupId = rollup.Id
		} else {
			bridges[i].RollupId = bridges[i].Rollup.Id
		}
	}

	return tx.SaveBridges(ctx, bridges...)
}
