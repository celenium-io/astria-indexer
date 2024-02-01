// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"
	"encoding/hex"
	"strings"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/indexer/decode"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (p *Module) parse(ctx context.Context, b types.BlockData) error {
	start := time.Now()
	p.Log.Info().
		Int64("height", b.Block.Height).
		Msg("parsing block...")

	decodeCtx := decode.NewContext()

	txs, err := parseTxs(b, &decodeCtx)
	if err != nil {
		return errors.Wrapf(err, "while parsing block on level=%d", b.Height)
	}

	block := &storage.Block{
		Height:       b.Height,
		Time:         b.Block.Time,
		VersionBlock: b.Block.Version.Block,
		VersionApp:   b.Block.Version.App,

		Hash:               []byte(b.BlockID.Hash),
		ParentHash:         []byte(b.Block.LastBlockID.Hash),
		LastCommitHash:     b.Block.LastCommitHash,
		DataHash:           b.Block.DataHash,
		ValidatorsHash:     b.Block.ValidatorsHash,
		NextValidatorsHash: b.Block.NextValidatorsHash,
		ConsensusHash:      b.Block.ConsensusHash,
		AppHash:            b.Block.AppHash,
		LastResultsHash:    b.Block.LastResultsHash,
		EvidenceHash:       b.Block.EvidenceHash,
		ProposerAddress:    b.Block.ProposerAddress.String(),

		ChainId:       b.Block.ChainID,
		Addresses:     decodeCtx.Addresses,
		Rollups:       decodeCtx.Rollups,
		RollupAddress: decodeCtx.RollupAddress,
		ActionTypes:   decodeCtx.ActionTypes,

		Txs:             txs,
		BlockSignatures: make([]storage.BlockSignature, len(b.Block.LastCommit.Signatures)),
		Stats: &storage.BlockStats{
			Height:       b.Height,
			Time:         b.Block.Time,
			TxCount:      int64(len(txs)),
			Fee:          decimal.Zero,
			SupplyChange: decodeCtx.SupplyChange,
			BytesInBlock: decodeCtx.BytesInBlock,
			GasWanted:    decodeCtx.GasWanted,
			GasUsed:      decodeCtx.GasUsed,
			DataSize:     decodeCtx.DataSize,
		},
	}

	p.parseBlockSignatures(b.Block.LastCommit, block.BlockSignatures)

	p.Log.Info().
		Uint64("height", uint64(block.Height)).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block parsed")

	output := p.MustOutput(OutputName)
	output.Push(block)
	return nil
}

func (p *Module) parseBlockSignatures(commit *types.Commit, signs []storage.BlockSignature) {
	for i := range commit.Signatures {
		signs[i].Height = types.Level(commit.Height)
		signs[i].Time = commit.Signatures[i].Timestamp
		signs[i].Validator = &storage.Validator{
			Address: strings.ToUpper(hex.EncodeToString(commit.Signatures[i].ValidatorAddress)),
		}
	}
}
