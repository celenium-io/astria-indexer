// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

func saveTransfers(
	ctx context.Context,
	tx storage.Transaction,
	transfers []*storage.Transfer,
	addrToId map[string]uint64,
) error {
	for i := range transfers {
		if transfers[i].Source != nil {
			if id, ok := addrToId[transfers[i].Source.Hash]; ok {
				transfers[i].SourceId = id
			}
		}

		if transfers[i].Destination != nil {
			if id, ok := addrToId[transfers[i].Destination.Hash]; ok {
				transfers[i].DestinationId = id
			}
		}
	}

	return tx.SaveTransfers(ctx, transfers...)
}
