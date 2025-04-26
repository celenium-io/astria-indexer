// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/hex"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
)

type Rollups map[string]*storage.Rollup

func NewRollups() Rollups {
	return make(map[string]*storage.Rollup)
}

func (r Rollups) Set(rollupId []byte, height types.Level, size int) *storage.Rollup {
	sRollupId := hex.EncodeToString(rollupId)

	if rollup, ok := r[sRollupId]; ok {
		rollup.ActionsCount += 1
		rollup.Size += int64(size)
		return rollup
	}

	rollup := &storage.Rollup{
		FirstHeight:  height,
		AstriaId:     rollupId,
		ActionsCount: 1,
		Size:         int64(size),
	}
	r[sRollupId] = rollup
	return rollup
}
