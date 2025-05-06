// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celenium-io/astria-indexer/cmd/api/bus"
	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	_ "github.com/celenium-io/astria-indexer/cmd/api/docs"
	"github.com/celenium-io/astria-indexer/cmd/api/handler"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	celestialsStorage "github.com/celenium-io/celestial-module/pkg/storage"
	"github.com/ipfans/fxlogger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "DipDup Verticals | Astria explorer API",
}

// @title					Swagger Astria Explorer API
// @version					1.0
// @description				This is docs of Astria Explorer API.
// @host					api.astrotrek.io
//
// @query.collection.format	multi
func main() {
	app := fx.New(
		fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		fx.StartTimeout(5*time.Minute),
		fx.Provide(
			loadConfig,
			databaseConfig,
			indexerName,
			fx.Annotate(
				cacheUrl,
				fx.ResultTags(`name:"cache_url"`),
			),
			newProflier,
			fx.Annotate(
				newServer,
				fx.ParamTags("", "", `group:"handlers"`),
			),
			bus.NewDispatcher,
			fx.Annotate(
				cache.InitCache,
				fx.ParamTags(`name:"cache_url"`),
			),
			newConstantCache,
			newWebsocket,
			newApp,

			newDatabase,
			fx.Annotate(
				postgres.NewListenerFactory,
				fx.As(new(storage.ListenerFactory)),
			),
			fx.Annotate(
				postgres.NewAction,
				fx.As(new(storage.IAction)),
			),
			fx.Annotate(
				postgres.NewAddress,
				fx.As(new(storage.IAddress)),
			),
			fx.Annotate(
				postgres.NewApp,
				fx.As(new(storage.IApp)),
			),
			fx.Annotate(
				postgres.NewAsset,
				fx.As(new(storage.IAsset)),
			),
			fx.Annotate(
				postgres.NewBlockSignature,
				fx.As(new(storage.IBlockSignature)),
			),
			fx.Annotate(
				postgres.NewBlockStats,
				fx.As(new(storage.IBlockStats)),
			),
			fx.Annotate(
				postgres.NewBlocks,
				fx.As(new(storage.IBlock)),
			),
			fx.Annotate(
				postgres.NewBridge,
				fx.As(new(storage.IBridge)),
			),
			fx.Annotate(
				postgres.NewConstant,
				fx.As(new(storage.IConstant)),
			),
			fx.Annotate(
				postgres.NewDeposit,
				fx.As(new(storage.IDeposit)),
			),
			fx.Annotate(
				postgres.NewFee,
				fx.As(new(storage.IFee)),
			),
			fx.Annotate(
				postgres.NewRollup,
				fx.As(new(storage.IRollup)),
			),
			fx.Annotate(
				postgres.NewSearch,
				fx.As(new(storage.ISearch)),
			),
			fx.Annotate(
				postgres.NewState,
				fx.As(new(storage.IState)),
			),
			fx.Annotate(
				postgres.NewStats,
				fx.As(new(storage.IStats)),
			),
			fx.Annotate(
				postgres.NewTransfer,
				fx.As(new(storage.ITransfer)),
			),
			fx.Annotate(
				postgres.NewTx,
				fx.As(new(storage.ITx)),
			),
			fx.Annotate(
				postgres.NewValidator,
				fx.As(new(storage.IValidator)),
			),
			fx.Annotate(
				postgres.NewPrice,
				fx.As(new(storage.IPrice)),
			),
			fx.Annotate(
				postgres.NewMarket,
				fx.As(new(storage.IMarket)),
			),
			fx.Annotate(
				newCelestials,
				fx.As(new(celestialsStorage.ICelestial)),
			),

			AsHandler(handler.NewAddressHandler),
			AsHandler(handler.NewAppHandler),
			AsHandler(handler.NewAssetHandler),
			AsHandler(handler.NewRollupHandler),
			AsHandler(handler.NewSearchHandler),
			AsHandler(handler.NewStateHandler),
			AsHandler(handler.NewValidatorHandler),
			AsHandler(handler.NewBlockHandler),
			AsHandler(handler.NewConstantHandler),
			AsHandler(handler.NewStatsHandler),
			AsHandler(handler.NewTxHandler),
			AsHandler(handler.NewPriceHandler),
			AsHandler(handler.NewActionHandler),
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
