// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/celenium-io/astria-indexer/cmd/api/bus"
	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/celenium-io/astria-indexer/cmd/api/handler"
	"github.com/celenium-io/astria-indexer/cmd/api/handler/websocket"
	"github.com/celenium-io/astria-indexer/internal/profiler"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/getsentry/sentry-go"
	"github.com/grafana/pyroscope-go"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/time/rate"
)

func init() {
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
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		}).
		With().Caller().
		Logger()
}

func loadConfig() (*Config, error) {
	configPath := rootCmd.PersistentFlags().StringP("config", "c", "dipdup.yml", "path to YAML config file")
	if err := rootCmd.Execute(); err != nil {
		return nil, errors.Wrap(err, "command line execute")
	}

	if err := rootCmd.MarkFlagRequired("config"); err != nil {
		return nil, errors.Wrap(err, "config command line arg is required")
	}

	var cfg Config
	if err := config.Parse(*configPath, &cfg); err != nil {
		return nil, errors.Wrap(err, "parsing config file")
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = zerolog.LevelInfoValue
	}

	if err := setLoggerLevel(cfg.LogLevel); err != nil {
		return nil, errors.Wrap(err, "set log level")
	}

	return &cfg, nil
}

func setLoggerLevel(level string) error {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return errors.Wrap(err, "parsing log level")
	}
	zerolog.SetGlobalLevel(logLevel)
	return nil
}

func newProflier(cfg *Config) (*pyroscope.Profiler, error) {
	return profiler.New(cfg.Profiler, "api")
}

func newServer(cfg *Config, wsManager *websocket.Manager, handlers []handler.Handler) (*echo.Echo, error) {
	e := echo.New()
	e.Validator = handler.NewApiValidator()
	e.Server.IdleTimeout = time.Second * 30

	e.Pre(middleware.RemoveTrailingSlash())

	timeout := 30 * time.Second
	if cfg.ApiConfig.RequestTimeout > 0 {
		timeout = time.Duration(cfg.ApiConfig.RequestTimeout) * time.Second
	}
	timeoutConfig := middleware.TimeoutConfig{
		Skipper: websocketSkipper,
		Timeout: timeout,
	}
	loggerConfig := middleware.RequestLoggerConfig{
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
	}
	gzipConfig := middleware.GzipConfig{
		Skipper: gzipSkipper,
	}
	decompressConfig := middleware.DecompressConfig{
		Skipper: websocketSkipper,
	}
	csrfConfig := middleware.CSRFConfig{
		Skipper: func(c echo.Context) bool {
			return websocketSkipper(c) || postSkipper(c)
		},
	}

	middlewares := []echo.MiddlewareFunc{
		middleware.TimeoutWithConfig(timeoutConfig),
		middleware.RequestLoggerWithConfig(loggerConfig),
		middleware.GzipWithConfig(gzipConfig),
		middleware.DecompressWithConfig(decompressConfig),
		middleware.BodyLimit("2M"),
		middleware.CSRFWithConfig(csrfConfig),
		middleware.CORS(),
		middleware.Recover(),
		middleware.Secure(),
	}

	if cfg.ApiConfig.Prometheus {
		middlewares = append(middlewares, echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Namespace: "astria_api",
			Skipper:   websocketSkipper,
		}))
	}
	if cfg.ApiConfig.RateLimit > 0 {
		middlewares = append(middlewares, middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Skipper: websocketSkipper,
			Store:   middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.ApiConfig.RateLimit)),
		}))
	}

	sentryMiddleware, err := initSentry(cfg.ApiConfig.SentryDsn, cfg.Environment)
	if err != nil {
		return nil, errors.Wrap(err, "init sentry")
	}
	if sentryMiddleware != nil {
		middlewares = append(middlewares, sentryMiddleware)
	}

	e.Use(middlewares...)

	v1 := e.Group("v1")
	for _, handler := range handlers {
		handler.InitRoutes(v1)
	}

	if cfg.ApiConfig.Websocket {
		wsManager.InitRoutes(v1)
	}

	if cfg.ApiConfig.Prometheus {
		e.GET("/metrics", echoprometheus.NewHandler())
	}

	v1.GET("/swagger/*", echoSwagger.WrapHandler)

	log.Info().Msg("API routes:")
	for _, route := range e.Routes() {
		log.Info().Msgf("[%s] %s -> %s", route.Method, route.Path, route.Name)
	}

	return e, nil
}

func newDatabase(cfg *Config) (*sdk.Storage, error) {
	return postgres.Create(context.Background(), cfg.Database, cfg.Indexer.ScriptsDir, false)
}

func newConstantCache(dispatcher *bus.Dispatcher) *cache.ConstantsCache {
	constantObserver := dispatcher.Observe(storage.ChannelConstant)
	return cache.NewConstantsCache(constantObserver)
}

func initSentry(dsn, environment string) (echo.MiddlewareFunc, error) {
	if dsn == "" {
		return nil, nil
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		Environment:      environment,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Release:          os.Getenv("TAG"),
	}); err != nil {
		return nil, errors.Wrap(err, "initialization")
	}

	return SentryMiddleware(), nil
}

func newWebsocket(dispatcher *bus.Dispatcher) *websocket.Manager {
	observer := dispatcher.Observe(storage.ChannelHead, storage.ChannelBlock)
	wsManager := websocket.NewManager(observer)
	return wsManager
}

func newEndpointCache(e *echo.Echo, dispatcher *bus.Dispatcher) *cache.Cache {
	observer := dispatcher.Observe(storage.ChannelHead)
	endpointCache := cache.NewCache(cache.Config{
		MaxEntitiesCount: 1000,
	}, observer)
	e.Use(cache.Middleware(endpointCache, cacheSkipper))
	return endpointCache
}
