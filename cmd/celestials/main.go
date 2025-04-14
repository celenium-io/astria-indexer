// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	indexer "github.com/celenium-io/celestial-module/pkg/module"
	celestialsStorage "github.com/celenium-io/celestial-module/pkg/storage"
	"github.com/ipfans/fxlogger"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "celestials",
	Short: "DipDup Verticals | Celestials indexer for Astria",
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	app := fx.New(
		fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		fx.StartTimeout(5*time.Minute),
		fx.Provide(
			loadConfig,
			datasourceConfig,
			newProfiler,
			setAddressHandler,
			newTransactable,
			newDatabase,
			fx.Annotate(
				indexerName,
				fx.ResultTags(`name:"indexerName"`),
			),
			fx.Annotate(
				networkName,
				fx.ResultTags(`name:"network"`),
			),
			fx.Annotate(
				postgres.NewAddress,
				fx.As(new(storage.IAddress)),
			),
			fx.Annotate(
				newCelestials,
				fx.As(new(celestialsStorage.ICelestial)),
			),
			fx.Annotate(
				newCelestialState,
				fx.As(new(celestialsStorage.ICelestialState)),
			),
			fx.Annotate(
				indexer.New,
				fx.ParamTags(
					``,
					``,
					``,
					``,
					``,
					`name:"indexerName"`,
					`name:"network"`,
				),
			),
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
