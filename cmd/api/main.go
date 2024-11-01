// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/celenium-io/astria-indexer/cmd/api/docs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "DipDup Verticals | Astria explorer API",
}

// @title					Swagger Astria Explorer API
// @version					1.0
// @description				This is docs of Astria Explorer API.
// @host					api-dusk.astrotrek.io
//
// @query.collection.format	multi
func main() {
	cfg, err := initConfig()
	if err != nil {
		return
	}

	if err = initLogger(cfg.LogLevel); err != nil {
		return
	}

	if err := initProflier(cfg.Profiler); err != nil {
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	db := initDatabase(cfg.Database, cfg.Indexer.ScriptsDir)
	e := initEcho(cfg.ApiConfig, db, cfg.Environment)
	initDispatcher(ctx, db)
	initHandlers(ctx, e, *cfg, db)
	initCache(ctx, e)

	go func() {
		if err := e.Start(cfg.ApiConfig.Bind); err != nil && errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	cancel()

	if wsManager != nil {
		if err := wsManager.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}
	if endpointCache != nil {
		if err := endpointCache.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}
	if constantCache != nil {
		if err := constantCache.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}
	if dispatcher != nil {
		if err := dispatcher.Close(); err != nil {
			e.Logger.Fatal(err)
		}
	}

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	if prscp != nil {
		if err := prscp.Stop(); err != nil {
			e.Logger.Fatal(err)
		}
	}
}
