// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/node/mock"
	"github.com/celenium-io/astria-indexer/pkg/types"
	cometTypes "github.com/cometbft/cometbft/types"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var testTime = time.Now()

func createModules(t *testing.T, ctrl *gomock.Controller) (modules.BaseModule, string, Module) {
	writerModule := modules.New("writer-module")
	outputName := "write"
	writerModule.CreateOutput(outputName)

	api := mock.NewMockApi(ctrl)
	parserModule := NewModule(api)

	err := parserModule.AttachTo(&writerModule, outputName, InputName)
	assert.NoError(t, err)

	return writerModule, outputName, parserModule
}

func getExpectedBlock() storage.Block {
	return storage.Block{
		Id:                 0,
		Height:             100,
		Time:               testTime,
		VersionBlock:       1,
		VersionApp:         3,
		Hash:               types.Hex{0x0, 0x0, 0x0, 0x2},
		ParentHash:         types.Hex{0x0, 0x0, 0x0, 0x1},
		LastCommitHash:     types.Hex{0x0, 0x0, 0x1, 0x1},
		DataHash:           types.Hex{0x0, 0x0, 0x1, 0x2},
		ValidatorsHash:     types.Hex{0x0, 0x0, 0x1, 0x3},
		NextValidatorsHash: types.Hex{0x0, 0x0, 0x1, 0x4},
		ConsensusHash:      types.Hex{0x0, 0x0, 0x1, 0x5},
		AppHash:            types.Hex{0x0, 0x0, 0x1, 0x6},
		LastResultsHash:    types.Hex{0x0, 0x0, 0x1, 0x7},
		EvidenceHash:       types.Hex{0x0, 0x0, 0x1, 0x8},
		ProposerAddress:    "astria1qqqqzzgxcftkc",
		ChainId:            "explorer-test",
		Txs:                make([]*storage.Tx, 0),
		Stats: &storage.BlockStats{
			Id:           0,
			Height:       100,
			Time:         testTime,
			TxCount:      0,
			Fee:          decimal.Zero,
			SupplyChange: decimal.Zero,
		},
		Addresses:       make(map[string]*storage.Address),
		Rollups:         make(map[string]*storage.Rollup),
		RollupAddress:   make(map[string]*storage.RollupAddress),
		Validators:      make(map[string]*storage.Validator),
		BlockSignatures: []storage.BlockSignature{},
		Constants:       make([]*storage.Constant, 0),
		Bridges:         make([]*storage.Bridge, 0),
		Transfers:       make([]*storage.Transfer, 0),
		Prices:          make([]storage.Price, 0),
		MarketUpdates:   make([]storage.MarketUpdate, 0),
		MarketProviders: make([]storage.MarketProviderUpdate, 0),
	}
}

func getBlock() types.BlockData {
	return types.BlockData{
		ResultBlock: types.ResultBlock{
			BlockID: types.BlockId{
				Hash: types.Hex{0x0, 0x0, 0x0, 0x2},
			},
			Block: &types.Block{
				Header: types.Header{
					Version: types.Consensus{
						Block: 1,
						App:   3,
					},
					ChainID: "explorer-test",
					Height:  1000,
					Time:    testTime,
					LastBlockID: types.BlockId{
						Hash: types.Hex{0x0, 0x0, 0x0, 0x1},
					},
					LastCommitHash:     types.Hex{0x0, 0x0, 0x1, 0x1},
					DataHash:           types.Hex{0x0, 0x0, 0x1, 0x2},
					ValidatorsHash:     types.Hex{0x0, 0x0, 0x1, 0x3},
					NextValidatorsHash: types.Hex{0x0, 0x0, 0x1, 0x4},
					ConsensusHash:      types.Hex{0x0, 0x0, 0x1, 0x5},
					AppHash:            types.Hex{0x0, 0x0, 0x1, 0x6},
					LastResultsHash:    types.Hex{0x0, 0x0, 0x1, 0x7},
					EvidenceHash:       types.Hex{0x0, 0x0, 0x1, 0x8},
					ProposerAddress:    types.Hex{0x0, 0x0, 0x1, 0x9},
				},
				Data: types.Data{
					Txs: nil,
				},
				LastCommit: &types.Commit{
					Height:     999,
					Round:      0,
					Signatures: []cometTypes.CommitSig{},
				},
			},
		},
		ResultBlockResults: types.ResultBlockResults{
			Height:              100,
			TxsResults:          nil,
			BeginBlockEvents:    nil,
			FinalizeBlockEvents: nil,
			ValidatorUpdates:    nil,
			ConsensusParamUpdates: &types.ConsensusParams{
				Block: &types.BlockParams{
					MaxBytes: 0,
					MaxGas:   0,
				},
				Evidence: &types.EvidenceParams{
					MaxAgeNumBlocks: 0,
					MaxAgeDuration:  0,
					MaxBytes:        0,
				},
				Validator: &types.ValidatorParams{
					PubKeyTypes: nil,
				},
				Version: &types.VersionParams{
					AppVersion: 0,
				},
			},
		},
	}
}

func TestParserModule_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writerModule, outputName, parserModule := createModules(t, ctrl)

	readerModule := modules.New("reader-module")
	readerInputName := "read"
	readerModule.CreateInput(readerInputName)

	err := readerModule.AttachTo(&parserModule, OutputName, readerInputName)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()

	parserModule.Start(ctx)

	block := getBlock()
	writerModule.MustOutput(outputName).Push(block)

	for {
		select {
		case <-ctx.Done():
			t.Error("stop by cancelled context")
		case msg, ok := <-readerModule.MustInput(readerInputName).Listen():
			assert.True(t, ok, "received value should be delivered by successful send operation")

			parsedBlock, ok := msg.(*storage.Block)
			assert.Truef(t, ok, "invalid message type: %T", msg)

			expectedBlock := getExpectedBlock()
			assert.Equal(t, &expectedBlock, parsedBlock)
			return
		}
	}
}

func TestModule_OnClosedChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, _, parserModule := createModules(t, ctrl)

	stopperModule := modules.New("stopper-module")
	stopInputName := "stop-signal"
	stopperModule.CreateInput(stopInputName)

	err := stopperModule.AttachTo(&parserModule, StopOutput, stopInputName)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*1)
	defer cancel()

	parserModule.Start(ctx)

	err = parserModule.MustInput(InputName).Close()
	assert.NoError(t, err)

	for {
		select {
		case <-ctx.Done():
			t.Error("stop by cancelled context")
		case msg := <-stopperModule.MustInput(stopInputName).Listen():
			assert.Equal(t, struct{}{}, msg)
			return
		}
	}
}

func TestModule_OnParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writerModule, writerOutputName, parserModule := createModules(t, ctrl)

	stopperModule := modules.New("stopper-module")
	stopInputName := "stop-signal"
	stopperModule.CreateInput(stopInputName)

	err := stopperModule.AttachTo(&parserModule, StopOutput, stopInputName)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*1)
	defer cancel()

	parserModule.Start(ctx)

	block := getBlock()
	block.Block.Data.Txs = []cometTypes.Tx{
		// unfinished sequence of tx bytes
		{10, 171, 1, 10, 168, 1, 10, 35, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98},
	}
	block.ResultBlockResults.TxsResults = []*types.ResponseDeliverTx{
		{
			Code:      0,
			Data:      []byte{18, 45, 10, 43, 47, 99, 111, 115, 109, 111, 115, 46, 115, 116, 97, 107, 105, 110, 103, 46, 118, 49, 98, 101, 116, 97},
			Log:       "",
			Info:      "",
			Events:    nil,
			Codespace: "",
		},
	}
	writerModule.MustOutput(writerOutputName).Push(block)

	for {
		select {
		case <-ctx.Done():
			t.Error("stop by cancelled context")
		case msg := <-stopperModule.MustInput(stopInputName).Listen():
			assert.Equal(t, struct{}{}, msg)
			return
		}
	}
}
