// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestRollupActionsByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	actions, err := s.storage.Rollup.ActionsByHeight(ctx, 7316, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(1, action.ActionId)
	s.Require().EqualValues(1, action.RollupId)
	s.Require().EqualValues(7316, action.Height)
	s.Require().NotNil(action.Rollup)
	s.Require().NotNil(action.Action)
}

func (s *StorageTestSuite) TestRollupCountActionsByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	count, err := s.storage.Rollup.CountActionsByHeight(ctx, 7316)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)
}

func (s *StorageTestSuite) TestRollupActionsByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	actions, err := s.storage.Rollup.ActionsByTxId(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(1, action.ActionId)
	s.Require().EqualValues(1, action.RollupId)
	s.Require().EqualValues(7316, action.Height)
	s.Require().NotNil(action.Rollup)
	s.Require().NotNil(action.Action)
}

func (s *StorageTestSuite) TestRollupCountActionsByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	count, err := s.storage.Rollup.CountActionsByTxId(ctx, 1)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)
}

func (s *StorageTestSuite) TestRollupByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539")
	s.Require().NoError(err)

	rollup, err := s.storage.Rollup.ByHash(ctx, hash)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues(7316, rollup.FirstHeight)
	s.Require().EqualValues(hash, rollup.AstriaId)
	s.Require().EqualValues(112, rollup.Size)
	s.Require().EqualValues(1, rollup.ActionsCount)
}

func (s *StorageTestSuite) TestRollupAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	addresses, err := s.storage.Rollup.Addresses(ctx, 1, 10, 0, storage.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(addresses, 1)

	address := addresses[0]
	s.Require().EqualValues(1, address.AddressId)
	s.Require().EqualValues(1, address.RollupId)
	s.Require().NotNil(address.Address)
	s.Require().EqualValues(1, address.Address.Id)
}

func (s *StorageTestSuite) TestListRollupsByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	rollups, err := s.storage.Rollup.ListRollupsByAddress(ctx, 1, 10, 0, storage.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(rollups, 1)

	rollup := rollups[0]
	s.Require().EqualValues(1, rollup.AddressId)
	s.Require().EqualValues(1, rollup.RollupId)
	s.Require().NotNil(rollup.Rollup)
	s.Require().EqualValues(1, rollup.Rollup.Id)
}
