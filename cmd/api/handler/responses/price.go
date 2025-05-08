// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

type Price struct {
	Price string    `example:"50.00"                     format:"string"    json:"value"          swaggertype:"string"`
	Time  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"           swaggertype:"string"`
	Pair  string    `example:"BTC/USDT"                  format:"string"    json:"pair,omitempty" swaggertype:"string"`
}

func NewPrice(price storage.Price) Price {
	return Price{
		Price: price.Price.String(),
		Time:  price.Time,
		Pair:  price.CurrencyPair,
	}
}

type Market struct {
	Pair             string     `example:"BTC/USDT"                  format:"string"            json:"pair"               swaggertype:"string"`
	Base             string     `example:"BTC"                       format:"string"            json:"base"               swaggertype:"string"`
	Quote            string     `example:"USDT"                      format:"string"            json:"quote"              swaggertype:"string"`
	Decimals         int        `example:"8"                         format:"integer"           json:"decimals"           swaggertype:"integer"`
	Enabled          bool       `example:"true"                      format:"boolean"           json:"enabled"            swaggertype:"boolean"`
	MinProviderCount int        `example:"1"                         format:"integer"           json:"min_provider_count" swaggertype:"integer"`
	UpdatedAt        time.Time  `example:"2023-07-04T03:10:57+00:00" format:"date-time"         json:"updated_at"         swaggertype:"string"`
	Price            *Price     `json:"price,omitempty"`
	Providers        []Provider `json:"providers,omitempty"          swaggertype:"array,object"`
}

func NewMarket(market storage.Market) Market {
	result := Market{
		Pair:             market.Pair,
		Base:             market.Base,
		Quote:            market.Quote,
		Decimals:         market.Decimals,
		Enabled:          market.Enabled,
		MinProviderCount: market.MinProviderCount,
		UpdatedAt:        market.UpdatedAt,
		Providers:        make([]Provider, len(market.Providers)),
	}

	if market.Price != nil {
		p := NewPrice(*market.Price)
		result.Price = &p
	}

	for i := range market.Providers {
		result.Providers[i] = NewProvider(market.Providers[i])
	}

	return result
}

type Provider struct {
	Provider       string `example:"binance"  format:"string" json:"provider"         swaggertype:"string"`
	OffChainTicker string `example:"BTC/USDT" format:"string" json:"off_chain_ticker" swaggertype:"string"`
}

func NewProvider(provider *storage.MarketProvider) Provider {
	return Provider{
		Provider:       provider.Provider,
		OffChainTicker: provider.OffChainTicker,
	}
}
