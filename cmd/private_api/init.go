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

	"cosmossdk.io/errors"
	"github.com/celenium-io/astria-indexer/cmd/private_api/handler"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func newServer(cfg *Config, handlers []handler.Handler) (*echo.Echo, error) {
	e := echo.New()
	e.Validator = handler.NewApiValidator()

	timeout := 30 * time.Second
	if cfg.ApiConfig.RequestTimeout > 0 {
		timeout = time.Duration(cfg.ApiConfig.RequestTimeout) * time.Second
	}
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: middleware.DefaultSkipper,
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
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Pre(middleware.RemoveTrailingSlash())

	if cfg.ApiConfig.RateLimit > 0 {
		e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store:   middleware.NewRateLimiterMemoryStore(rate.Limit(cfg.ApiConfig.RateLimit)),
		}))

	}

	e.Server.IdleTimeout = time.Second * 30

	v1 := e.Group("v1")
	for _, handler := range handlers {
		handler.InitRoutes(v1)
	}

	log.Info().Msg("API routes:")
	for _, route := range e.Routes() {
		log.Info().Msgf("[%s] %s -> %s", route.Method, route.Path, route.Name)
	}

	return e, nil
}

func newDatabase(cfg *Config) (*sdk.Storage, error) {
	return postgres.Create(context.Background(), cfg.Database, cfg.Indexer.ScriptsDir, false)
}

func newTransactable(db *sdk.Storage) storage.Transactable {
	return db.Transactable
}
