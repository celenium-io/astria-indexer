// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"github.com/celenium-io/astria-indexer/internal/astria"
	"github.com/celenium-io/astria-indexer/internal/profiler"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestial-module/pkg/module"
	celestialsPg "github.com/celenium-io/celestial-module/pkg/storage/postgres"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	sdkPg "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/grafana/pyroscope-go"
	"github.com/pkg/errors"
	"os"
	"strconv"

	dipdupCfg "github.com/dipdup-net/go-lib/config"
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

func loadConfig() (*Config, error) {
	configPath := rootCmd.PersistentFlags().StringP("config", "c", "dipdup.yml", "path to YAML config file")
	if err := rootCmd.Execute(); err != nil {
		return nil, errors.Wrap(err, "command line execute")
	}

	if err := rootCmd.MarkFlagRequired("config"); err != nil {
		return nil, errors.Wrap(err, "config command line arg is required")
	}

	var cfg Config
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

func newCelestials(db *sdkPg.Storage) *celestialsPg.Celestials {
	return &celestialsPg.Celestials{
		Bun: db.Connection(),
	}
}

func newCelestialState(db *sdkPg.Storage) *celestialsPg.CelestialState {
	return celestialsPg.NewCelestialState(db.Connection())
}

func newProfiler(cfg *Config) (*pyroscope.Profiler, error) {
	return profiler.New(cfg.Profiler, "celestials")
}

func newDatabase(cfg *Config) (*sdkPg.Storage, error) {
	return postgres.Create(context.Background(), cfg.Database, cfg.Indexer.ScriptsDir, true)
}

func newTransactable(db *sdkPg.Storage) sdk.Transactable {
	return db.Transactable
}

func setAddressHandler(repo storage.IAddress) module.AddressHandler {
	return func(ctx context.Context, address string) (uint64, error) {
		return addressHandler(ctx, repo, address)
	}
}

func addressHandler(ctx context.Context, repo storage.IAddress, address string) (uint64, error) {
	prefix, hash, err := astria.DecodeAddress(address)
	if err != nil {
		return 0, errors.Wrap(err, "decoding address")
	}
	if prefix != astria.Prefix && prefix != astria.PrefixCompat {
		return 0, errors.Errorf("invalid prefix: %s", prefix)
	}

	addr, err := repo.ByHash(ctx, string(hash))
	if err != nil {
		return 0, errors.Errorf("can't find address %s in database", address)
	}

	return addr.Id, nil
}
