// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/celenium-io/astria-indexer/pkg/indexer/config"
	"github.com/pkg/errors"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

const (
	InputName  = "data"
	StopOutput = "stop"
)

// Module - saves received from input block to storage.
//
//	                     |----------------|
//	                     |                |
//	-- storage.Block ->  |     MODULE     |
//	                     |                |
//	                     |----------------|
type Module struct {
	modules.BaseModule
	storage     sdk.Transactable
	notificator storage.Notificator
	indexerName string
	validators  map[string]uint64
}

var _ modules.Module = (*Module)(nil)

// NewModule -
func NewModule(
	storage sdk.Transactable,
	notificator storage.Notificator,
	cfg config.Indexer,
) Module {
	m := Module{
		BaseModule:  modules.New("storage"),
		storage:     storage,
		indexerName: cfg.Name,
		notificator: notificator,
		validators:  make(map[string]uint64),
	}

	m.CreateInputWithCapacity(InputName, 16)
	m.CreateOutput(StopOutput)

	return m
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
		case msg, ok := <-input.Listen():
			if !ok {
				module.Log.Warn().Msg("can't read message from input")
				module.MustOutput(StopOutput).Push(struct{}{})
				continue
			}
			block, ok := msg.(*storage.Block)
			if !ok {
				module.Log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			state, err := module.saveBlock(ctx, block)
			if err != nil {
				module.Log.Err(err).
					Uint64("height", uint64(block.Height)).
					Msg("block saving error")
				module.MustOutput(StopOutput).Push(struct{}{})
				continue
			}

			if err := module.notify(ctx, state, block); err != nil {
				module.Log.Err(err).Msg("block notification error")
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

func (module *Module) saveBlock(ctx context.Context, block *storage.Block) (storage.State, error) {
	start := time.Now()
	module.Log.Info().Uint64("height", uint64(block.Height)).Msg("saving block...")
	tx, err := postgres.BeginTransaction(ctx, module.storage)
	if err != nil {
		return storage.State{}, err
	}
	defer tx.Close(ctx)

	state, err := module.processBlockInTransaction(ctx, tx, block)
	if err != nil {
		return state, tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return state, tx.HandleError(ctx, err)
	}
	module.Log.Info().
		Uint64("height", uint64(block.Height)).
		Time("block_time", block.Time).
		Int64("ms", time.Since(start).Milliseconds()).
		Int("tx_count", len(block.Txs)).
		Msg("block saved")
	return state, nil
}

func (module *Module) processBlockInTransaction(ctx context.Context, tx storage.Transaction, block *storage.Block) (storage.State, error) {
	state, err := tx.State(ctx, module.indexerName)
	if err != nil {
		return state, err
	}
	block.Stats.BlockTime = uint64(block.Time.Sub(state.LastTime).Milliseconds())

	if len(module.validators) > 0 {
		if id, ok := module.validators[block.ProposerAddress]; ok {
			block.ProposerId = id
		} else {
			return state, errors.Errorf("unknown block proposer: %s", block.ProposerAddress)
		}
	} else {
		proposerId, err := tx.GetProposerId(ctx, block.ProposerAddress)
		if err != nil {
			return state, errors.Wrap(err, "can't find block proposer")
		}
		block.ProposerId = proposerId
	}

	if err := tx.Add(ctx, block); err != nil {
		return state, err
	}

	if err := tx.Add(ctx, block.Stats); err != nil {
		return state, err
	}

	addrToId, totalAccounts, err := saveAddresses(ctx, tx, block.Addresses)
	if err != nil {
		return state, err
	}

	if err := module.saveTransactions(ctx, tx, addrToId, block.Txs...); err != nil {
		return state, err
	}

	totalRollups, err := module.saveRollup(ctx, tx, addrToId, block.Rollups, block.RollupAddress)
	if err != nil {
		return state, err
	}

	var actions = make([]*storage.Action, 0)

	for i := range block.Txs {
		for j := range block.Txs[i].Actions {
			block.Txs[i].Actions[j].TxId = block.Txs[i].Id
			actions = append(actions, &block.Txs[i].Actions[j])
		}
	}

	if err := saveAction(ctx, tx, actions); err != nil {
		return state, err
	}

	if err := module.saveBlockSignatures(ctx, tx, block.BlockSignatures, block.Height); err != nil {
		return state, err
	}

	if err := saveValidators(ctx, tx, block.Validators); err != nil {
		return state, err
	}

	updateState(block, totalAccounts, totalRollups, &state)
	if err := tx.Update(ctx, &state); err != nil {
		return state, err
	}

	return state, nil
}

func (module *Module) notify(ctx context.Context, state storage.State, block *storage.Block) error {
	if time.Since(block.Time) > time.Hour {
		// do not notify all about events if initial indexing is in progress
		return nil
	}

	rawState, err := json.Marshal(state)
	if err != nil {
		return err
	}
	if err := module.notificator.Notify(ctx, storage.ChannelHead, string(rawState)); err != nil {
		return err
	}

	blockId := strconv.FormatUint(block.Id, 10)
	if err := module.notificator.Notify(ctx, storage.ChannelBlock, blockId); err != nil {
		return err
	}

	return nil
}
