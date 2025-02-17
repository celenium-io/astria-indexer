// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestDepositByBridgeId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	deposits, err := s.Deposit.ByBridgeId(ctx, 1, 10, 0, storage.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(deposits, 1)

	deposit := deposits[0]
	s.Require().EqualValues(7965, deposit.Height)
	s.Require().EqualValues(1, deposit.RollupId)
	s.Require().EqualValues(1, deposit.BridgeId)
	s.Require().EqualValues("100", deposit.Amount.String())
	s.Require().EqualValues("destination_chain_address", deposit.DestinationChainAddress)
	s.Require().EqualValues("nria", deposit.Asset)

	s.Require().NotNil(deposit.Tx)
	hash := hex.EncodeToString(deposit.Tx.Hash)
	s.Require().Equal("a7bc8121a38725bd33e5d66b80817a2ba39e517fb6b9244a7081ad2fb210bfcc", hash)

	s.Require().Nil(deposit.Action)
	s.Require().Nil(deposit.Bridge)
	s.Require().Nil(deposit.Rollup)
}

func (s *StorageTestSuite) TestDepositByRollupId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	deposits, err := s.Deposit.ByRollupId(ctx, 1, 10, 0, storage.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(deposits, 1)

	deposit := deposits[0]
	s.Require().EqualValues(7965, deposit.Height)
	s.Require().EqualValues(1, deposit.RollupId)
	s.Require().EqualValues(1, deposit.BridgeId)
	s.Require().EqualValues("100", deposit.Amount.String())
	s.Require().EqualValues("destination_chain_address", deposit.DestinationChainAddress)
	s.Require().EqualValues("nria", deposit.Asset)

	s.Require().NotNil(deposit.Tx)
	hash := hex.EncodeToString(deposit.Tx.Hash)
	s.Require().Equal("a7bc8121a38725bd33e5d66b80817a2ba39e517fb6b9244a7081ad2fb210bfcc", hash)

	s.Require().NotNil(deposit.Bridge)
	s.Require().NotNil(deposit.Bridge.Address)
	s.Require().EqualValues("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", deposit.Bridge.Address.Hash)

	s.Require().Nil(deposit.Action)
	s.Require().Nil(deposit.Rollup)
}
