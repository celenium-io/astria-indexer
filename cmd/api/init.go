// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/celenium-io/astria-indexer/cmd/api/bus"
	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/websocket"
	"github.com/celenium-io/astria-indexer/internal/profiler"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/grafana/pyroscope-go"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/time/rate"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})
}

func initConfig() (*Config, error) {
	configPath := rootCmd.PersistentFlags().StringP("config", "c", "dipdup.yml", "path to YAML config file")
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err).Msg("command line execute")
		return nil, err
	}

	if err := rootCmd.MarkFlagRequired("config"); err != nil {
		log.Panic().Err(err).Msg("config command line arg is required")
		return nil, err
	}

	var cfg Config
	if err := config.Parse(*configPath, &cfg); err != nil {
		log.Panic().Err(err).Msg("parsing config file")
		return nil, err
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = zerolog.LevelInfoValue
	}

	return &cfg, nil
}

func initLogger(level string) error {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Panic().Err(err).Msg("parsing log level")
		return err
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.Logger.With().Caller().Logger()

	return nil
}

var prscp *pyroscope.Profiler

func initProflier(cfg *profiler.Config) (err error) {
	prscp, err = profiler.New(cfg, "api")
	return
}

func websocketSkipper(c echo.Context) bool {
	return strings.Contains(c.Request().URL.Path, "ws")
}

func postSkipper(c echo.Context) bool {
	if strings.Contains(c.Request().URL.Path, "blob") {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "auth/rollup") {
		return true
	}
	return false
}

func gzipSkipper(c echo.Context) bool {
	if strings.Contains(c.Request().URL.Path, "swagger") {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "metrics") {
		return true
	}
	return websocketSkipper(c)
}

func cacheSkipper(c echo.Context) bool {
	if c.Request().Method != http.MethodGet {
		return true
	}
	if websocketSkipper(c) {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "metrics") {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "head") {
		return true
	}
	return false
}

func initEcho(cfg ApiConfig, db postgres.Storage, env string) *echo.Echo {
	e := echo.New()
	e.Validator = handler.NewApiValidator()

	timeout := 30 * time.Second
	if cfg.RequestTimeout > 0 {
		timeout = time.Duration(cfg.RequestTimeout) * time.Second
	}
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: websocketSkipper,
		Timeout: timeout,
	}))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogMethod:    true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			switch {
			case v.Status == http.StatusOK || v.Status == http.StatusNoContent:
				log.Info().
					Str("uri", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("method", v.Method).
					Str("user-agent", v.UserAgent).
					Str("ip", c.RealIP()).
					Msg("request")
			case v.Status >= 500:
				log.Error().
					Str("uri", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("method", v.Method).
					Str("user-agent", v.UserAgent).
					Str("ip", c.RealIP()).
					Msg("request")
			default:
				log.Warn().
					Str("uri", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("method", v.Method).
					Str("user-agent", v.UserAgent).
					Str("ip", c.RealIP()).
					Msg("request")
			}

			return nil
		},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: gzipSkipper,
	}))
	e.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{
		Skipper: websocketSkipper,
	}))
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CSRFWithConfig(
		middleware.CSRFConfig{
			Skipper: func(c echo.Context) bool {
				return websocketSkipper(c) || postSkipper(c)
			},
		},
	))
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Pre(middleware.RemoveTrailingSlash())

	if cfg.Prometheus {
		e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Namespace: "astria_api",
			Skipper:   websocketSkipper,
		}))
	}
	if cfg.RateLimit > 0 {
		e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Skipper: websocketSkipper,
			Store:   middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit)),
		}))

	}

	if err := initSentry(e, db, cfg.SentryDsn, env); err != nil {
		log.Err(err).Msg("sentry")
	}
	e.Server.IdleTimeout = time.Second * 30

	return e
}

var dispatcher *bus.Dispatcher

func initDispatcher(ctx context.Context, db postgres.Storage) {
	d, err := bus.NewDispatcher(db, db.Blocks)
	if err != nil {
		panic(err)
	}
	dispatcher = d
	dispatcher.Start(ctx)
}

func initDatabase(cfg config.Database, viewsDir string) postgres.Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := postgres.Create(ctx, cfg, viewsDir)
	if err != nil {
		panic(err)
	}
	return db
}

var constantCache *cache.ConstantsCache

func initHandlers(ctx context.Context, e *echo.Echo, cfg Config, db postgres.Storage) {
	v1 := e.Group("v1")

	stateHandlers := handler.NewStateHandler(db.State)
	v1.GET("/head", stateHandlers.Head)
	constantsHandler := handler.NewConstantHandler(db.Constants)
	v1.GET("/constants", constantsHandler.Get)
	v1.GET("/enums", constantsHandler.Enums)

	constantObserver := dispatcher.Observe(storage.ChannelConstant)
	constantCache = cache.NewConstantsCache(constantObserver)
	if err := constantCache.Start(ctx, db.Constants); err != nil {
		panic(err)
	}

	searchHandler := handler.NewSearchHandler(constantCache, db.Search, db.Address, db.Blocks, db.Tx, db.Rollup, db.Bridges, db.Validator)
	v1.GET("/search", searchHandler.Search)

	addressHandler := handler.NewAddressHandler(constantCache, db.Address, db.Tx, db.Action, db.Rollup, db.Fee, db.Bridges, db.Deposit, db.State, cfg.Indexer.Name)
	addressesGroup := v1.Group("/address")
	{
		addressesGroup.GET("", addressHandler.List)
		addressesGroup.GET("/count", addressHandler.Count)
		addressGroup := addressesGroup.Group("/:hash")
		{
			addressGroup.GET("", addressHandler.Get)
			addressGroup.GET("/txs", addressHandler.Transactions)
			addressGroup.GET("/actions", addressHandler.Actions)
			addressGroup.GET("/rollups", addressHandler.Rollups)
			addressGroup.GET("/roles", addressHandler.Roles)
			addressGroup.GET("/fees", addressHandler.Fees)
			addressGroup.GET("/deposits", addressHandler.Deposits)
		}
	}

	blockHandlers := handler.NewBlockHandler(db.Blocks, db.BlockStats, db.Tx, db.Action, db.Rollup, db.State, cfg.Indexer.Name)
	blockGroup := v1.Group("/block")
	{
		blockGroup.GET("", blockHandlers.List)
		blockGroup.GET("/count", blockHandlers.Count)
		heightGroup := blockGroup.Group("/:height")
		{
			heightGroup.GET("", blockHandlers.Get)
			heightGroup.GET("/actions", blockHandlers.GetActions)
			heightGroup.GET("/txs", blockHandlers.GetTransactions)
			heightGroup.GET("/stats", blockHandlers.GetStats)
			heightGroup.GET("/rollup_actions", blockHandlers.GetRollupActions)
			heightGroup.GET("/rollup_actions/count", blockHandlers.GetRollupsActionsCount)
		}
	}

	txHandlers := handler.NewTxHandler(db.Tx, db.Action, db.Rollup, db.Fee, db.State, cfg.Indexer.Name)
	txGroup := v1.Group("/tx")
	{
		txGroup.GET("", txHandlers.List)
		txGroup.GET("/count", txHandlers.Count)
		hashGroup := txGroup.Group("/:hash")
		{
			hashGroup.GET("", txHandlers.Get)
			hashGroup.GET("/actions", txHandlers.GetActions)
			hashGroup.GET("/fees", txHandlers.GetFees)
			hashGroup.GET("/rollup_actions", txHandlers.RollupActions)
			hashGroup.GET("/rollup_actions/count", txHandlers.RollupActionsCount)
		}
	}

	rollupsHandler := handler.NewRollupHandler(constantCache, db.Rollup, db.Action, db.Bridges, db.Deposit, db.State, cfg.Indexer.Name)
	rollupsGroup := v1.Group("/rollup")
	{
		rollupsGroup.GET("", rollupsHandler.List)
		rollupsGroup.GET("/count", rollupsHandler.Count)

		rollupGroup := rollupsGroup.Group("/:hash")
		{
			rollupGroup.GET("", rollupsHandler.Get)
			rollupGroup.GET("/actions", rollupsHandler.Actions)
			rollupGroup.GET("/all_actions", rollupsHandler.AllActions)
			rollupGroup.GET("/addresses", rollupsHandler.Addresses)
			rollupGroup.GET("/bridges", rollupsHandler.Bridges)
			rollupGroup.GET("/deposits", rollupsHandler.Deposits)
		}
	}

	validatorsHandler := handler.NewValidatorHandler(db.Validator, db.Blocks, db.BlockSignatures, db.State, cfg.Indexer.Name)
	validators := v1.Group("/validators")
	{
		validators.GET("", validatorsHandler.List)
		validatorGroup := validators.Group("/:id")
		{
			validatorGroup.GET("", validatorsHandler.Get)
			validatorGroup.GET("/blocks", validatorsHandler.Blocks)
			validatorGroup.GET("/uptime", validatorsHandler.Uptime)
		}
	}

	statsHandler := handler.NewStatsHandler(db.Stats, db.Rollup)
	stats := v1.Group("/stats")
	{
		stats.GET("/summary", statsHandler.Summary)
		stats.GET("/summary/:timeframe", statsHandler.SummaryTimeframe)
		stats.GET("/summary/active_addresses_count", statsHandler.ActiveAddressesCount)
		stats.GET("/series/:name/:timeframe", statsHandler.Series)

		rollup := stats.Group("/rollup")
		{
			rollup.GET("/series/:hash/:name/:timeframe", statsHandler.RollupSeries)
		}

		fee := stats.Group("/fee")
		{
			fee.GET("/summary", statsHandler.FeeSummary)
		}

		token := stats.Group("/token")
		{
			token.GET("/transfer_distribution", statsHandler.TokenTransferDistribution)
		}
	}

	if cfg.ApiConfig.Prometheus {
		e.GET("/metrics", echoprometheus.NewHandler())
	}

	v1.GET("/swagger/*", echoSwagger.WrapHandler)

	if cfg.ApiConfig.Websocket {
		initWebsocket(ctx, v1)
	}

	log.Info().Msg("API routes:")
	for _, route := range e.Routes() {
		log.Info().Msgf("[%s] %s -> %s", route.Method, route.Path, route.Name)
	}
}

func initSentry(e *echo.Echo, db postgres.Storage, dsn, environment string) error {
	if dsn == "" {
		return nil
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		Environment:      environment,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Release:          os.Getenv("TAG"),
	}); err != nil {
		return errors.Wrap(err, "initialization")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	db.SetTracer(tp)

	e.Use(SentryMiddleware())

	return nil
}

var (
	wsManager     *websocket.Manager
	endpointCache *cache.Cache
)

func initWebsocket(ctx context.Context, group *echo.Group) {
	observer := dispatcher.Observe(storage.ChannelHead, storage.ChannelBlock)
	wsManager = websocket.NewManager(observer)
	wsManager.Start(ctx)
	group.GET("/ws", wsManager.Handle)
}

func initCache(ctx context.Context, e *echo.Echo) {
	observer := dispatcher.Observe(storage.ChannelHead)
	endpointCache = cache.NewCache(cache.Config{
		MaxEntitiesCount: 1000,
	}, observer)
	e.Use(cache.Middleware(endpointCache, cacheSkipper))
	endpointCache.Start(ctx)
}
