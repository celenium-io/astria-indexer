// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestMarketList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	markets, err := s.Market.List(ctx, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(markets, 4)

	s.Require().EqualValues("TIA_BTC", markets[3].Pair)
	s.Require().EqualValues("TIA", markets[3].Base)
	s.Require().EqualValues("BTC", markets[3].Quote)
	s.Require().EqualValues(8, markets[3].Decimals)
	s.Require().EqualValues(1, markets[2].MinProviderCount)
	s.Require().True(markets[3].Enabled)
	s.Require().NotNil(markets[3].Price)
	s.Require().EqualValues("0.00000001", markets[3].Price.Price.String())
	s.Require().EqualValues("2025-01-01 00:00:00 +0000 UTC", markets[3].UpdatedAt.String())

	s.Require().EqualValues("TIA_USD", markets[0].Pair)
	s.Require().EqualValues("TIA", markets[0].Base)
	s.Require().EqualValues("USD", markets[0].Quote)
	s.Require().EqualValues(6, markets[0].Decimals)
	s.Require().EqualValues(1, markets[0].MinProviderCount)
	s.Require().True(markets[0].Enabled)
	s.Require().NotNil(markets[0].Price)
	s.Require().EqualValues("0.000025", markets[0].Price.Price.String())
	s.Require().EqualValues("2025-01-01 00:00:00 +0000 UTC", markets[0].UpdatedAt.String())
}

func (s *StorageTestSuite) TestMarketGet() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	market, err := s.Market.Get(ctx, "TIA_USD")
	s.Require().NoError(err)

	s.Require().EqualValues("TIA_USD", market.Pair)
	s.Require().EqualValues("TIA", market.Base)
	s.Require().EqualValues("USD", market.Quote)
	s.Require().EqualValues(6, market.Decimals)
	s.Require().EqualValues(1, market.MinProviderCount)
	s.Require().True(market.Enabled)
	s.Require().NotNil(market.Price)
	s.Require().EqualValues("0.000025", market.Price.Price.String())
	s.Require().Len(market.Providers, 1)
	s.Require().EqualValues("coingecko", market.Providers[0].Provider)
	s.Require().EqualValues("TIA/USD", market.Providers[0].OffChainTicker)
	s.Require().EqualValues("2025-01-01 00:00:00 +0000 UTC", market.UpdatedAt.String())
}

func (s *StorageTestSuite) TestMarketDecimals() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	decimals, err := s.Market.Decimals(ctx, "TIA_BTC")
	s.Require().NoError(err)
	s.Require().EqualValues(8, decimals)

	decimals, err = s.Market.Decimals(ctx, "TIA_USD")
	s.Require().NoError(err)
	s.Require().EqualValues(6, decimals)
}

func (s *StorageTestSuite) TestHistoryList() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	markets, err := s.Market.History(ctx, "TIA_USD", 10, 0)
	s.Require().NoError(err)
	s.Require().Len(markets, 2)

	s.Require().EqualValues("TIA_USD", markets[0].Pair)
	s.Require().EqualValues("TIA", markets[0].Base)
	s.Require().EqualValues("USD", markets[0].Quote)
	s.Require().EqualValues(6, markets[0].Decimals)
	s.Require().EqualValues(1, markets[0].MinProviderCount)
	s.Require().True(markets[0].Enabled)
	s.Require().Nil(markets[0].Price)
	s.Require().EqualValues("2025-01-01 00:00:00 +0000 UTC", markets[0].UpdatedAt.String())

	s.Require().EqualValues("TIA_USD", markets[1].Pair)
	s.Require().EqualValues("TIA", markets[1].Base)
	s.Require().EqualValues("USD", markets[1].Quote)
	s.Require().EqualValues(4, markets[1].Decimals)
	s.Require().EqualValues(1, markets[1].MinProviderCount)
	s.Require().True(markets[1].Enabled)
	s.Require().Nil(markets[1].Price)
	s.Require().EqualValues("2024-01-01 00:00:00 +0000 UTC", markets[1].UpdatedAt.String())
}
