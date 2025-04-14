// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	module "github.com/celenium-io/celestial-module/pkg/module"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type App struct {
	indexerModule *module.Module
	db            *postgres.Storage
}

func newApp(
	lc fx.Lifecycle,
	db *postgres.Storage,
	indexerModule *module.Module,
) *App {
	app := &App{
		db:            db,
		indexerModule: indexerModule,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			app.indexerModule.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := app.indexerModule.Close(); err != nil {
				return errors.Wrap(err, "closing indexer module")
			}
			if err := app.db.Close(); err != nil {
				return errors.Wrap(err, "closing database")
			}
			return nil
		},
	})
	return app
}
