// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"net/http"

	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type App struct {
	e  *echo.Echo
	db *postgres.Storage
}

func newApp(
	lc fx.Lifecycle,
	cfg *Config,
	e *echo.Echo,
	db *postgres.Storage,
) *App {
	app := &App{
		e:  e,
		db: db,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := app.e.Start(cfg.ApiConfig.Bind); err != nil && errors.Is(err, http.ErrServerClosed) {
				return errors.Wrap(err, "shutting down the server")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := app.e.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "closing server")
			}
			if err := app.db.Close(); err != nil {
				return errors.Wrap(err, "closing database")
			}
			return nil
		},
	})
	return app
}
