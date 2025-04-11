package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

func (s *StorageTestSuite) TestPriceAll() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	prices, err := s.Price.All(ctx, 1, 1)
	s.Require().NoError(err)
	s.Require().Len(prices, 1)

	price := prices[0]
	s.Require().EqualValues("ETH_USDT", price.CurrencyPair)
	s.Require().EqualValues("3000", price.Price.String())
	s.Require().NotEmpty(price.Time)
}

func (s *StorageTestSuite) TestPriceLast() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	price, err := s.Price.Last(ctx, "ETH_USDT")
	s.Require().NoError(err)

	s.Require().EqualValues("ETH_USDT", price.CurrencyPair)
	s.Require().EqualValues("3000", price.Price.String())
	s.Require().NotEmpty(price.Time)
}

func (s *StorageTestSuite) TestPriceSeries() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	prices, err := s.Price.Series(ctx, "ETH_USDT", storage.TimeframeHour)
	s.Require().NoError(err)
	s.Require().Len(prices, 2)

	price := prices[0]
	s.Require().EqualValues("3000", price.High.String())
	s.Require().NotEmpty(price.Time)
}
