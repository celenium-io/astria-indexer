// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
)

func updateState(block *storage.Block, totalAccounts, totalRollups, totalBridges, totalBytes int64, state *storage.State) {
	if types.Level(block.Id) <= state.LastHeight {
		return
	}

	state.LastHeight = block.Height
	state.LastHash = block.Hash
	state.LastTime = block.Time
	state.TotalTx += block.Stats.TxCount
	state.TotalAccounts += totalAccounts
	state.TotalRollups += totalRollups
	state.TotalBridges += totalBridges
	state.TotalBytes += totalBytes
	state.TotalSupply = state.TotalSupply.Add(block.Stats.SupplyChange)
	state.ChainId = block.ChainId
}
