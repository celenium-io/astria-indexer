// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celenium-io/astria-indexer/cmd/private_api/handler"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/ipfans/fxlogger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var rootCmd = &cobra.Command{
	Use: "private_api",
}

func main() {
	app := fx.New(
		fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		fx.StartTimeout(5*time.Minute),
		fx.Provide(
			loadConfig,

			fx.Annotate(
				newServer,
				fx.ParamTags("", `group:"handlers"`),
			),
			newApp,

			newDatabase,
			newTransactable,
			fx.Annotate(
				postgres.NewAddress,
				fx.As(new(storage.IAddress)),
			),
			fx.Annotate(
				postgres.NewApp,
				fx.As(new(storage.IApp)),
			),
			fx.Annotate(
				postgres.NewRollup,
				fx.As(new(storage.IRollup)),
			),

			AsHandler(handler.NewAppHandler),
		),
		fx.Invoke(func(*App) {}),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	if err := app.Start(ctx); err != nil {
		log.Err(err).Msg("start app")
		os.Exit(1)
	}

	<-ctx.Done()
	cancel()

	if err := app.Stop(ctx); err != nil {
		log.Err(err).Msg("stop app")
		os.Exit(1)
	}
}

func AsHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(handler.Handler)),
		fx.ResultTags(`group:"handlers"`),
	)
}
