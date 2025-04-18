// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/celenium-io/astria-indexer/pkg/indexer/config"
	conf "github.com/dipdup-net/go-lib/config"
)

type CelestialsConfig struct {
	ChainId string `validate:"required" yaml:"chain_id"`
}

type Config struct {
	*config.Config `yaml:",inline"`

	Celestials CelestialsConfig `validate:"required" yaml:"celestials"`
}

func datasourceConfig(cfg *Config) conf.DataSource {
	return cfg.DataSources["celestials"]
}

func indexerName(cfg *Config) string {
	return cfg.Indexer.Name
}

func networkName(cfg *Config) string {
	return cfg.Celestials.ChainId
}
