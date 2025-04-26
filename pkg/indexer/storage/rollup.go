// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

func (module *Module) saveRollup(
	ctx context.Context,
	tx storage.Transaction,
	rollups map[string]*storage.Rollup,
	rollupAddress map[string]*storage.RollupAddress,
) (int64, int64, error) {
	if len(rollups) == 0 {
		return 0, 0, nil
	}

	var totalSize int64

	data := make([]*storage.Rollup, 0)
	for _, value := range rollups {
		data = append(data, value)
		totalSize += value.Size
	}

	count, err := tx.SaveRollups(ctx, data...)
	if err != nil {
		return count, totalSize, err
	}

	ra := make([]*storage.RollupAddress, 0)
	for _, value := range rollupAddress {
		value.RollupId = value.Rollup.Id
		value.AddressId = value.Address.Id
		ra = append(ra, value)
	}
	if err := tx.SaveRollupAddresses(ctx, ra...); err != nil {
		return 0, 0, err
	}

	return count, totalSize, nil
}
