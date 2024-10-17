// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rollback

import (
	"bytes"
	"context"

	"github.com/celenium-io/astria-indexer/pkg/node"

	"github.com/celenium-io/astria-indexer/pkg/indexer/config"
	"github.com/celenium-io/astria-indexer/pkg/types"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	InputName  = "signal"
	OutputName = "state"
	StopOutput = "stop"
)

// Module - executes rollback on signal from input and notify all subscribers about new state after rollback operation.
//
//	                |----------------|
//	                |                |
//	-- struct{} ->  |     MODULE     |  -- storage.State ->
//	                |                |
//	                |----------------|
type Module struct {
	modules.BaseModule
	tx        sdk.Transactable
	state     storage.IState
	blocks    storage.IBlock
	node      node.Api
	indexName string
}

var _ modules.Module = (*Module)(nil)

func NewModule(
	tx sdk.Transactable,
	state storage.IState,
	blocks storage.IBlock,
	node node.Api,
	cfg config.Indexer,
) Module {
	module := Module{
		BaseModule: modules.New("rollback"),
		tx:         tx,
		state:      state,
		blocks:     blocks,
		node:       node,
		indexName:  cfg.Name,
	}

	module.CreateInput(InputName)
	module.CreateOutput(OutputName)
	module.CreateOutput(StopOutput)

	return module
}

// Start -
func (module *Module) Start(ctx context.Context) {
	module.G.GoCtx(ctx, module.listen)
}

func (module *Module) listen(ctx context.Context) {
	module.Log.Info().Msg("module started")
	input := module.MustInput(InputName)

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-input.Listen():
			if !ok {
				module.Log.Warn().Msg("can't read message from input, channel was dried and closed")
				module.MustOutput(StopOutput).Push(struct{}{})
				return
			}

			if err := module.rollback(ctx); err != nil {
				module.Log.Err(err).Msgf("error occurred")
			}
		}
	}
}

// Close -
func (module *Module) Close() error {
	module.Log.Info().Msg("closing module...")
	module.G.Wait()
	return nil
}

func (module *Module) rollback(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
		default:
			lastBlock, err := module.blocks.Last(ctx)
			if err != nil {
				return errors.Wrap(err, "receive last block from database")
			}

			nodeBlock, err := module.node.Block(ctx, lastBlock.Height)
			if err != nil {
				return errors.Wrapf(err, "receive block from node by height: %d", lastBlock.Height)
			}

			log.Debug().
				Uint64("height", uint64(lastBlock.Height)).
				Hex("db_block_hash", lastBlock.Hash).
				Hex("node_block_hash", nodeBlock.BlockID.Hash).
				Msg("comparing hash...")

			if bytes.Equal(lastBlock.Hash, nodeBlock.BlockID.Hash) {
				return module.finish(ctx)
			}

			log.Warn().
				Uint64("height", uint64(lastBlock.Height)).
				Hex("db_block_hash", lastBlock.Hash).
				Hex("node_block_hash", nodeBlock.BlockID.Hash).
				Msg("need rollback")

			if err := module.rollbackBlock(ctx, lastBlock.Height); err != nil {
				return errors.Wrapf(err, "rollback block: %d", lastBlock.Height)
			}
		}
	}
}

func (module *Module) finish(ctx context.Context) error {
	newState, err := module.state.ByName(ctx, module.indexName)
	if err != nil {
		return err
	}
	module.MustOutput(OutputName).Push(newState)

	log.Info().
		Uint64("new_height", uint64(newState.LastHeight)).
		Msg("roll backed to new height")

	return nil
}

func (module *Module) rollbackBlock(ctx context.Context, height types.Level) error {
	tx, err := postgres.BeginTransaction(ctx, module.tx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	if err := rollbackBlock(ctx, tx, height, module.indexName); err != nil {
		return tx.HandleError(ctx, err)
	}

	return nil
}

func rollbackBlock(ctx context.Context, tx storage.Transaction, height types.Level, indexName string) error {
	if err := tx.RollbackBlock(ctx, height); err != nil {
		return err
	}

	blockStats, err := tx.RollbackBlockStats(ctx, height)
	if err != nil {
		return err
	}

	txs, err := tx.RollbackTxs(ctx, height)
	if err != nil {
		return err
	}

	actions, err := tx.RollbackActions(ctx, height)
	if err != nil {
		return err
	}

	addressActions, err := tx.RollbackAddressActions(ctx, height)
	if err != nil {
		return err
	}

	countDeletedAddresses, err := rollbackAddress(ctx, tx, height, addressActions, txs)
	if err != nil {
		return errors.Wrap(err, "address")
	}

	countDeletedRollups, err := rollbackRollups(ctx, tx, height, actions)
	if err != nil {
		return errors.Wrap(err, "rollups")
	}

	if err := tx.RollbackValidators(ctx, height); err != nil {
		return err
	}

	if err := tx.RollbackFees(ctx, height); err != nil {
		return err
	}

	if err := tx.RollbackDeposits(ctx, height); err != nil {
		return err
	}

	if err := tx.RollbackTransfers(ctx, height); err != nil {
		return err
	}

	if err := tx.RollbackBlockSignatures(ctx, height); err != nil {
		return err
	}

	deletedBridges, err := tx.RollbackBridges(ctx, height)
	if err != nil {
		return errors.Wrap(err, "bridges")
	}

	newBlock, err := tx.LastBlock(ctx)
	if err != nil {
		return err
	}
	state, err := tx.State(ctx, indexName)
	if err != nil {
		return err
	}

	state.LastHeight = newBlock.Height
	state.LastHash = newBlock.Hash
	state.LastTime = newBlock.Time
	state.TotalTx -= blockStats.TxCount
	state.TotalAccounts -= int64(countDeletedAddresses)
	state.TotalRollups -= countDeletedRollups
	state.TotalSupply = state.TotalSupply.Sub(blockStats.SupplyChange)
	state.TotalBridges -= int64(deletedBridges)

	if err := tx.Update(ctx, &state); err != nil {
		return err
	}

	return tx.Flush(ctx)
}
