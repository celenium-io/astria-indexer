// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"github.com/aopoltorzhicky/astria/internal/storage"
	"github.com/aopoltorzhicky/astria/pkg/types"
)

func updateState(block *storage.Block, totalAccounts, totalRollups int64, state *storage.State) {
	if types.Level(block.Id) <= state.LastHeight {
		return
	}

	state.LastHeight = block.Height
	state.LastHash = block.Hash
	state.LastTime = block.Time
	state.TotalTx += block.Stats.TxCount
	state.TotalAccounts += totalAccounts
	state.TotalRollups += totalRollups
	state.TotalFee = state.TotalFee.Add(block.Stats.Fee)
	state.TotalSupply = state.TotalSupply.Add(block.Stats.SupplyChange)
	state.ChainId = block.ChainId
}
