// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/pkg/errors"
)

func (module *Module) saveRollup(
	ctx context.Context,
	tx storage.Transaction,
	addrToId map[string]uint64,
	rollups map[string]*storage.Rollup,
	rollupAddress map[string]*storage.RollupAddress,
) (int64, error) {
	if len(rollups) == 0 {
		return 0, nil
	}

	data := make([]*storage.Rollup, 0)
	for _, value := range rollups {
		if value.BridgeAddress != nil {
			if id, ok := addrToId[value.BridgeAddress.String()]; ok {
				value.BridgeAddressId = id
			} else {
				return 0, errors.Errorf("unknown bridge address id: %s", value.BridgeAddress.String())
			}
		}
		data = append(data, value)
	}

	count, err := tx.SaveRollups(ctx, data...)
	if err != nil {
		return count, err
	}

	ra := make([]*storage.RollupAddress, 0)
	for _, value := range rollupAddress {
		value.RollupId = value.Rollup.Id
		value.AddressId = value.Address.Id
		ra = append(ra, value)
	}
	if err := tx.SaveRollupAddresses(ctx, ra...); err != nil {
		return 0, err
	}

	return count, nil
}
