// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

type RollupAction struct {
	Action

	Rollup *Rollup `json:"rollup,omitempty"`
}

func NewRollupAction(action storage.RollupAction) RollupAction {
	result := RollupAction{
		Action: NewAction(*action.Action),
	}

	if action.Tx != nil {
		result.TxHash = hex.EncodeToString(action.Tx.Hash)
	}

	if action.Rollup != nil {
		r := NewRollup(action.Rollup)
		result.Rollup = &r
	}

	return result
}
