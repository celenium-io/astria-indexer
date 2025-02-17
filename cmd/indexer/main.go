// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/stopper"
	"github.com/ipfans/fxlogger"
	"go.uber.org/fx"

	"github.com/celenium-io/astria-indexer/pkg/indexer"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "indexer",
	Short: "DipDup Verticals | Astria Indexer",
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	app := fx.New(
		fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		fx.Supply(cancel),
		fx.Provide(
			loadConfig,
			newProflier,
			newDatabase,
			fx.Annotate(
				stopper.NewModule,
				fx.As(new(modules.Module)),
			),
			indexer.New,
			newApp,
		),
		fx.Invoke(func(*App) {}),
	)

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
