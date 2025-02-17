package main

import (
	"context"

	"github.com/celenium-io/astria-indexer/pkg/indexer"
	"github.com/celenium-io/astria-indexer/pkg/indexer/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/grafana/pyroscope-go"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type App struct {
	prscp         *pyroscope.Profiler
	stopperModule modules.Module
	indexerModule indexer.Indexer
	db            *postgres.Storage
}

func newApp(
	lc fx.Lifecycle,
	cfg *config.Config,
	db *postgres.Storage,
	prscp *pyroscope.Profiler,
	stopperModule modules.Module,
	indexerModule indexer.Indexer,
) *App {
	app := &App{
		prscp:         prscp,
		stopperModule: stopperModule,
		indexerModule: indexerModule,
		db:            db,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			app.stopperModule.Start(ctx)
			app.indexerModule.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := app.prscp.Stop(); err != nil {
				return errors.Wrap(err, "closing pyroscrope")
			}
			if err := app.indexerModule.Close(); err != nil {
				return errors.Wrap(err, "closing indexer module")
			}
			if err := app.stopperModule.Close(); err != nil {
				return errors.Wrap(err, "closing stopper module")
			}
			if err := app.db.Close(); err != nil {
				return errors.Wrap(err, "closing database")
			}
			return nil
		},
	})
	return app
}
