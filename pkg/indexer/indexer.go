// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package indexer

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	sdkPg "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"

	internalStorage "github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/indexer/genesis"
	"github.com/celenium-io/astria-indexer/pkg/indexer/parser"
	"github.com/celenium-io/astria-indexer/pkg/indexer/rollback"
	"github.com/celenium-io/astria-indexer/pkg/indexer/storage"
	"github.com/celenium-io/astria-indexer/pkg/node"
	"github.com/celenium-io/astria-indexer/pkg/node/rpc"
	"github.com/pkg/errors"

	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/celenium-io/astria-indexer/pkg/indexer/config"
	"github.com/celenium-io/astria-indexer/pkg/indexer/receiver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Indexer struct {
	cfg      *config.Config
	api      node.Api
	receiver *receiver.Module
	parser   *parser.Module
	storage  *storage.Module
	rollback *rollback.Module
	genesis  *genesis.Module
	stopper  modules.Module
	log      zerolog.Logger

	states  internalStorage.IState
	bridges internalStorage.IBridge
}

func New(cfg *config.Config, pg *sdkPg.Storage, stopperModule modules.Module) (Indexer, error) {
	states := postgres.NewState(pg)
	blocks := postgres.NewBlocks(pg)
	bridges := postgres.NewBridge(pg)
	notificator := postgres.NewNotificator(cfg.Database, pg)

	api, r, err := createReceiver(cfg)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating receiver module")
	}

	rb, err := createRollback(r, pg.Transactable, states, blocks, &api, cfg.Indexer)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating rollback module")
	}

	p, err := createParser(r, &api)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating parser module")
	}

	s, err := createStorage(pg.Transactable, notificator, cfg, p)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating storage module")
	}

	genesisModule, err := createGenesis(pg.Transactable, cfg, r)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating genesis module")
	}

	err = attachStopper(stopperModule, r, p, s, rb, genesisModule)
	if err != nil {
		return Indexer{}, errors.Wrap(err, "while creating stopper module")
	}

	return Indexer{
		cfg:      cfg,
		api:      &api,
		receiver: r,
		parser:   p,
		storage:  s,
		rollback: rb,
		genesis:  genesisModule,
		stopper:  stopperModule,
		log:      log.With().Str("module", "indexer").Logger(),
		states:   states,
		bridges:  bridges,
	}, nil
}

func (i *Indexer) Start(ctx context.Context) {
	i.log.Info().Msg("starting...")

	state, err := loadState(ctx, i.states, i.cfg.Indexer.Name)
	if err != nil {
		log.Err(err).Msg("load state")
	}
	i.receiver.Init(state)

	assets, err := makeBridgeAssetsMap(ctx, i.bridges)
	if err != nil {
		log.Err(err).Msg("make bridge asset map")
	}
	i.parser.Init(ctx, assets)

	i.genesis.Start(ctx)
	i.storage.Start(ctx)
	i.parser.Start(ctx)
	i.receiver.Start(ctx)
}

func (i *Indexer) Close() error {
	i.log.Info().Msg("closing...")

	if err := i.receiver.Close(); err != nil {
		log.Err(err).Msg("closing receiver")
	}
	if err := i.genesis.Close(); err != nil {
		log.Err(err).Msg("closing genesis")
	}
	if err := i.parser.Close(); err != nil {
		log.Err(err).Msg("closing parser")
	}
	if err := i.storage.Close(); err != nil {
		log.Err(err).Msg("closing storage")
	}
	if err := i.rollback.Close(); err != nil {
		log.Err(err).Msg("closing rollback")
	}

	return nil
}

func createReceiver(cfg *config.Config) (rpc.API, *receiver.Module, error) {

	api := rpc.NewAPI(cfg.DataSources["sequencer_rpc"])
	receiverModule := receiver.NewModule(cfg.Indexer, &api)
	return api, &receiverModule, nil
}

func createRollback(receiverModule modules.Module, tx sdk.Transactable, states internalStorage.IState, blocks internalStorage.IBlock, api node.Api, cfg config.Indexer) (*rollback.Module, error) {
	rollbackModule := rollback.NewModule(tx, states, blocks, api, cfg)

	// rollback <- listen signal -- receiver
	if err := rollbackModule.AttachTo(receiverModule, receiver.RollbackOutput, rollback.InputName); err != nil {
		return nil, errors.Wrap(err, "while attaching rollback to receiver")
	}

	// receiver <- listen state -- rollback
	if err := receiverModule.AttachTo(&rollbackModule, rollback.OutputName, receiver.RollbackInput); err != nil {
		return nil, errors.Wrap(err, "while attaching receiver to rollback")
	}

	return &rollbackModule, nil
}

func makeBridgeAssetsMap(ctx context.Context, bridges internalStorage.IBridge) (map[string]string, error) {
	assets := make(map[string]string)
	for end := false; !end; {
		data, err := bridges.ListWithAddress(ctx, 100, len(assets))
		if err != nil {
			return nil, err
		}
		for i := range data {
			assets[data[i].Address.Hash] = data[i].Asset
		}
		end = len(data) < 100
	}
	return assets, nil
}

func createParser(receiverModule modules.Module, api node.Api) (*parser.Module, error) {
	parserModule := parser.NewModule(api)

	if err := parserModule.AttachTo(receiverModule, receiver.BlocksOutput, parser.InputName); err != nil {
		return nil, errors.Wrap(err, "while attaching parser to receiver")
	}

	return &parserModule, nil
}

func createStorage(tx sdk.Transactable, notificator internalStorage.Notificator, cfg *config.Config, parserModule modules.Module) (*storage.Module, error) {
	storageModule := storage.NewModule(tx, notificator, cfg.Indexer)

	if err := storageModule.AttachTo(parserModule, parser.OutputName, storage.InputName); err != nil {
		return nil, errors.Wrap(err, "while attaching storage to parser")
	}

	return &storageModule, nil
}

func createGenesis(tx sdk.Transactable, cfg *config.Config, receiverModule modules.Module) (*genesis.Module, error) {
	genesisModule := genesis.NewModule(tx, cfg.Indexer)

	if err := genesisModule.AttachTo(receiverModule, receiver.GenesisOutput, genesis.InputName); err != nil {
		return nil, errors.Wrap(err, "while attaching genesis to receiver")
	}

	genesisModulePtr := &genesisModule
	if err := receiverModule.AttachTo(genesisModulePtr, genesis.OutputName, receiver.GenesisDoneInput); err != nil {
		return nil, errors.Wrap(err, "while attaching receiver to genesis")
	}

	return genesisModulePtr, nil
}

func attachStopper(stopperModule modules.Module, receiverModule modules.Module, parserModule modules.Module, storageModule modules.Module, rollbackModule modules.Module, genesisModule modules.Module) error {
	if err := stopperModule.AttachTo(receiverModule, receiver.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to receiver")
	}

	if err := stopperModule.AttachTo(parserModule, parser.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to parser")
	}

	if err := stopperModule.AttachTo(storageModule, storage.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to storage")
	}

	if err := stopperModule.AttachTo(rollbackModule, rollback.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to rollback")
	}

	if err := stopperModule.AttachTo(genesisModule, genesis.StopOutput, stopper.InputName); err != nil {
		return errors.Wrap(err, "while attaching stopper to genesis")
	}

	return nil
}

func loadState(ctx context.Context, states internalStorage.IState, indexerName string) (*internalStorage.State, error) {
	state, err := states.ByName(ctx, indexerName)
	if err != nil {
		if states.IsNoRows(err) {
			return nil, nil
		}

		return nil, err
	}

	return &state, nil
}
