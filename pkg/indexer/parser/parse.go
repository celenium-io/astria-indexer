// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"
	"encoding/hex"
	"strings"
	"time"

	"github.com/celenium-io/astria-indexer/internal/astria"
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

	proposer, err := astria.EncodeFromHex(b.Block.ProposerAddress.String())
	if err != nil {
		return errors.Wrap(err, "decoding block proposer address")
	}

	decodeCtx := decode.NewContext(p.bridgeAssets, b.Block.Time)
	decodeCtx.Proposer = proposer

	if err := parseEvents(ctx, b.FinalizeBlockEvents, b.Height, &decodeCtx, p.api); err != nil {
		return errors.Wrap(err, "parse finalize events")
	}

	txs, err := parseTxs(ctx, b, &decodeCtx, p.api)
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
		ProposerAddress:    proposer,

		ChainId:       b.Block.ChainID,
		Addresses:     decodeCtx.Addresses,
		Rollups:       decodeCtx.Rollups,
		RollupAddress: decodeCtx.RollupAddress,
		Validators:    decodeCtx.Validators,
		ActionTypes:   decodeCtx.ActionTypes,
		Constants:     decodeCtx.ConstantsArray(),
		Bridges:       decodeCtx.BridgesArray(),
		Transfers:     decodeCtx.Transfers,

		Txs: txs,
		Stats: &storage.BlockStats{
			Height:       b.Height,
			Time:         b.Block.Time,
			TxCount:      int64(len(txs)),
			Fee:          decimal.Zero,
			SupplyChange: decodeCtx.SupplyChange,
			BytesInBlock: decodeCtx.BytesInBlock,
			DataSize:     decodeCtx.DataSize,
		},
		Prices: decodeCtx.Prices,
	}

	block.BlockSignatures = p.parseBlockSignatures(b.Block.LastCommit)

	p.Log.Info().
		Uint64("height", uint64(block.Height)).
		Int64("ms", time.Since(start).Milliseconds()).
		Msg("block parsed")

	output := p.MustOutput(OutputName)
	output.Push(block)
	return nil
}

func (p *Module) parseBlockSignatures(commit *types.Commit) []storage.BlockSignature {
	signs := make([]storage.BlockSignature, 0)
	for i := range commit.Signatures {
		if commit.Signatures[i].BlockIDFlag != 2 {
			continue
		}
		signs = append(signs, storage.BlockSignature{
			Height: types.Level(commit.Height),
			Time:   commit.Signatures[i].Timestamp,
			Validator: &storage.Validator{
				Address: strings.ToUpper(hex.EncodeToString(commit.Signatures[i].ValidatorAddress)),
			},
		})
	}
	return signs
}
