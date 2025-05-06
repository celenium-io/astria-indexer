// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/bus"
	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/websocket"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/grafana/pyroscope-go"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type App struct {
	e             *echo.Echo
	db            *postgres.Storage
	wsManager     *websocket.Manager
	dispatcher    *bus.Dispatcher
	constantCache *cache.ConstantsCache
	ttlCache      cache.ICache
	prscp         *pyroscope.Profiler
	constants     storage.IConstant
}

func newApp(
	lc fx.Lifecycle,
	cfg *Config,
	e *echo.Echo,
	db *postgres.Storage,
	wsManager *websocket.Manager,
	dispatcher *bus.Dispatcher,
	constantCache *cache.ConstantsCache,
	ttlCache cache.ICache,
	prscp *pyroscope.Profiler,
	constants storage.IConstant,
) *App {
	app := &App{
		e:             e,
		db:            db,
		wsManager:     wsManager,
		dispatcher:    dispatcher,
		constantCache: constantCache,
		ttlCache:      ttlCache,
		prscp:         prscp,
		constants:     constants,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			dispatcher.Start(ctx)
			wsManager.Start(ctx)
			if err := constantCache.Start(ctx, app.constants); err != nil {
				return errors.Wrap(err, "start constant cache")
			}

			if err := app.e.Start(cfg.ApiConfig.Bind); err != nil && errors.Is(err, http.ErrServerClosed) {
				return errors.Wrap(err, "shutting down the server")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if app.wsManager != nil {
				if err := app.wsManager.Close(); err != nil {
					return errors.Wrap(err, "closing websocket manager")
				}
			}
			if app.ttlCache != nil {
				if err := app.ttlCache.Close(); err != nil {
					return errors.Wrap(err, "closing cache")
				}
			}
			if app.constantCache != nil {
				if err := app.constantCache.Close(); err != nil {
					return errors.Wrap(err, "closing constant cache")
				}
			}
			if app.dispatcher != nil {
				if err := app.dispatcher.Close(); err != nil {
					return errors.Wrap(err, "closing bus dispatcher")
				}
			}

			if err := app.e.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "closing server")
			}

			if app.prscp != nil {
				if err := app.prscp.Stop(); err != nil {
					return errors.Wrap(err, "closing profler")
				}
			}
			if err := app.db.Close(); err != nil {
				return errors.Wrap(err, "closing database")
			}
			return nil
		},
	})
	return app
}
