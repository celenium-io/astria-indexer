// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"strconv"

	"github.com/celenium-io/astria-indexer/internal/profiler"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/celenium-io/astria-indexer/pkg/indexer/config"
	dipdupCfg "github.com/dipdup-net/go-lib/config"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/grafana/pyroscope-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func loadConfig() (*config.Config, error) {
	configPath := rootCmd.PersistentFlags().StringP("config", "c", "dipdup.yml", "path to YAML config file")
	if err := rootCmd.Execute(); err != nil {
		return nil, errors.Wrap(err, "command line execute")
	}

	if err := rootCmd.MarkFlagRequired("config"); err != nil {
		return nil, errors.Wrap(err, "config command line arg is required")
	}

	var cfg config.Config
	if err := dipdupCfg.Parse(*configPath, &cfg); err != nil {
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

func newProflier(cfg *config.Config) (*pyroscope.Profiler, error) {
	return profiler.New(cfg.Profiler, "api")
}

func newDatabase(cfg *config.Config) (*sdk.Storage, error) {
	return postgres.Create(context.Background(), cfg.Database, cfg.Indexer.ScriptsDir, false)
}
