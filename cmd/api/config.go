// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/celenium-io/astria-indexer/internal/profiler"
	indexerConfig "github.com/celenium-io/astria-indexer/pkg/indexer/config"
	"github.com/dipdup-net/go-lib/config"
)

type Config struct {
	*config.Config `yaml:",inline"`
	LogLevel       string                `validate:"omitempty,oneof=debug trace info warn error fatal panic" yaml:"log_level"`
	ApiConfig      ApiConfig             `validate:"required"                                                yaml:"api"`
	Profiler       *profiler.Config      `validate:"omitempty"                                               yaml:"profiler"`
	Indexer        indexerConfig.Indexer `validate:"required"                                                yaml:"indexer"`
	Environment    string                `validate:"omitempty,oneof=development production"                  yaml:"environment"`
}

type ApiConfig struct {
	Bind           string  `validate:"required,hostname_port" yaml:"bind"`
	RateLimit      float64 `validate:"omitempty,min=0"        yaml:"rate_limit"`
	Prometheus     bool    `validate:"omitempty"              yaml:"prometheus"`
	RequestTimeout int     `validate:"omitempty,min=1"        yaml:"request_timeout"`
	BlobReceiver   string  `validate:"required"               yaml:"blob_receiver"`
	SentryDsn      string  `validate:"omitempty"              yaml:"sentry_dsn"`
	Websocket      bool    `validate:"omitempty"              yaml:"websocket"`
	Cache          string  `validate:"omitempty,url"          yaml:"cache"`
}

func indexerName(cfg *Config) string {
	return cfg.Indexer.Name
}

func cacheUrl(cfg *Config) string {
	return cfg.ApiConfig.Cache
}

func databaseConfig(cfg *Config) config.Database {
	return cfg.Database
}
