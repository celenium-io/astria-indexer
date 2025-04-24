package storage

import "github.com/uptrace/bun"

type MarketProvider struct {
	bun.BaseModel `bun:"table:market_provider" comment:"Table with market's providers"`

	Pair           string `bun:"pair,pk"          comment:"Trading pair"`
	Provider       string `bun:"provider,pk"      comment:"Provider"`
	OffChainTicker string `bun:"off_chain_ticker" comment:"Off chain ticker"`
}

func (MarketProvider) TableName() string {
	return "market_provider"
}

type MarketProviderUpdate struct {
	MarketProvider
	Type MarketUpdateType
}
