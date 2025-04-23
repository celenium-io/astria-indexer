// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/shopspring/decimal"
)

type Price struct {
	CurrencyPair string    `example:"BTC/USDT"                  format:"string"    json:"currency_pair,omitempty" swaggertype:"string"`
	Price        string    `example:"50.00"                     format:"string"    json:"price"                   swaggertype:"string"`
	Time         time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"                    swaggertype:"string"`
}

func NewPrice(price storage.Price) Price {
	return Price{
		CurrencyPair: price.CurrencyPair,
		Price:        price.Price.String(),
		Time:         price.Time,
	}
}

type Market struct {
	Pair             string `example:"BTC/USDT"     format:"string"  json:"pair"               swaggertype:"string"`
	Base             string `example:"BTC"          format:"string"  json:"base"               swaggertype:"string"`
	Quote            string `example:"USDT"         format:"string"  json:"quote"              swaggertype:"string"`
	Decimals         int    `example:"8"            format:"integer" json:"decimals"           swaggertype:"integer"`
	Enabled          bool   `example:"true"         format:"boolean" json:"enabled"            swaggertype:"boolean"`
	MinProviderCount int    `example:"1"            format:"integer" json:"min_provider_count" swaggertype:"integer"`
	Price            *Price `json:"price,omitempty"`
}

func NewMarket(market storage.Market) Market {
	result := Market{
		Pair:             market.Pair,
		Base:             market.Base,
		Quote:            market.Quote,
		Decimals:         market.Decimals,
		Enabled:          market.Enabled,
		MinProviderCount: market.MinProviderCount,
	}

	if market.Price != nil {
		result.Price = &Price{
			Price: decimalPrice(market.Price.Price, market.Decimals),
			Time:  market.Price.Time,
		}
	}

	return result
}

func decimalPrice(price decimal.Decimal, decimals int) string {
	dec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(-decimals)))
	return price.Mul(dec).String()
}
