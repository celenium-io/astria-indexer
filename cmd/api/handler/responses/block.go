// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"strconv"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
)

type Block struct {
	Id                 uint64          `example:"321"                                                              json:"id"                   swaggertype:"integer"`
	Height             uint64          `example:"100"                                                              json:"height"               swaggertype:"integer"`
	Time               time.Time       `example:"2023-07-04T03:10:57+00:00"                                        json:"time"                 swaggertype:"string"`
	VersionBlock       string          `example:"11"                                                               json:"version_block"        swaggertype:"string"`
	VersionApp         string          `example:"1"                                                                json:"version_app"          swaggertype:"string"`
	Hash               pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"hash"                 swaggertype:"string"`
	ParentHash         pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"parent_hash"          swaggertype:"string"`
	LastCommitHash     pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"last_commit_hash"     swaggertype:"string"`
	DataHash           pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"data_hash"            swaggertype:"string"`
	ValidatorsHash     pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"validators_hash"      swaggertype:"string"`
	NextValidatorsHash pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"next_validators_hash" swaggertype:"string"`
	ConsensusHash      pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"consensus_hash"       swaggertype:"string"`
	AppHash            pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"app_hash"             swaggertype:"string"`
	LastResultsHash    pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"last_results_hash"    swaggertype:"string"`
	EvidenceHash       pkgTypes.Hex    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"evidence_hash"        swaggertype:"string"`
	ActionTypes        []string        `example:"rollup_data_submission,transfer"                                  json:"action_types"         swaggertype:"string"`
	Proposer           *ShortValidator `json:"proposer,omitempty"`

	Stats *BlockStats `json:"stats,omitempty"`
}

func NewBlock(block storage.Block) Block {
	result := Block{
		Id:                 block.Id,
		Height:             uint64(block.Height),
		Time:               block.Time,
		VersionBlock:       strconv.FormatUint(block.VersionBlock, 10),
		VersionApp:         strconv.FormatUint(block.VersionApp, 10),
		Hash:               block.Hash,
		ParentHash:         block.ParentHash,
		LastCommitHash:     block.LastCommitHash,
		DataHash:           block.DataHash,
		ValidatorsHash:     block.ValidatorsHash,
		NextValidatorsHash: block.NextValidatorsHash,
		ConsensusHash:      block.ConsensusHash,
		AppHash:            block.AppHash,
		LastResultsHash:    block.LastResultsHash,
		EvidenceHash:       block.EvidenceHash,
		ActionTypes:        types.NewActionTypeMaskBits(block.ActionTypes).Strings(),
	}
	result.Proposer = NewShortValidator(block.Proposer)

	if block.Stats != nil {
		result.Stats = NewBlockStats(block.Stats)
	}
	return result
}

type BlockStats struct {
	TxCount      int64  `example:"12"          json:"tx_count"       swaggertype:"integer"`
	Fee          string `example:"28347628346" json:"fee"            swaggertype:"string"`
	SupplyChange string `example:"8635234"     json:"supply_change"  swaggertype:"string"`
	BlockTime    uint64 `example:"12354"       json:"block_time"     swaggertype:"integer"`
	GasWanted    int64  `example:"1234"        json:"gas_wanted"     swaggertype:"integer"`
	GasUsed      int64  `example:"1234"        json:"gas_used"       swaggertype:"integer"`
	BytesInBlock int64  `example:"1234"        json:"bytes_in_block" swaggertype:"integer"`
}

func NewBlockStats(stats *storage.BlockStats) *BlockStats {
	return &BlockStats{
		TxCount:      stats.TxCount,
		Fee:          stats.Fee.String(),
		SupplyChange: stats.SupplyChange.String(),
		BlockTime:    stats.BlockTime,
		GasUsed:      stats.GasUsed,
		GasWanted:    stats.GasWanted,
		BytesInBlock: stats.BytesInBlock,
	}
}
