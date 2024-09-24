// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"
)

func (s *StorageTestSuite) TestBridgeByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	bridge, err := s.storage.Bridges.ByAddress(ctx, 1)
	s.Require().NoError(err)

	s.Require().EqualValues(7316, bridge.InitHeight)
	s.Require().EqualValues(1, bridge.AddressId)
	s.Require().EqualValues(1, bridge.SudoId)
	s.Require().EqualValues(1, bridge.WithdrawerId)
	s.Require().EqualValues(1, bridge.RollupId)
	s.Require().EqualValues("nria", bridge.Asset)
	s.Require().EqualValues("nria", bridge.FeeAsset)

	s.Require().NotNil(bridge.Address)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Address.Hash)

	s.Require().NotNil(bridge.Sudo)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Sudo.Hash)

	s.Require().NotNil(bridge.Withdrawer)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Withdrawer.Hash)

	s.Require().NotNil(bridge.Rollup)
	hash, _ := hex.DecodeString("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539")
	s.Require().Equal(hash, bridge.Rollup.AstriaId)
}

func (s *StorageTestSuite) TestBridgeByRollup() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	bridges, err := s.storage.Bridges.ByRollup(ctx, 1, 10, 0)
	s.Require().NoError(err)

	bridge := bridges[0]
	s.Require().EqualValues(7316, bridge.InitHeight)
	s.Require().EqualValues(1, bridge.AddressId)
	s.Require().EqualValues(1, bridge.SudoId)
	s.Require().EqualValues(1, bridge.WithdrawerId)
	s.Require().EqualValues(1, bridge.RollupId)
	s.Require().EqualValues("nria", bridge.Asset)
	s.Require().EqualValues("nria", bridge.FeeAsset)

	s.Require().NotNil(bridge.Address)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Address.Hash)

	s.Require().NotNil(bridge.Sudo)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Sudo.Hash)

	s.Require().NotNil(bridge.Withdrawer)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Withdrawer.Hash)

	s.Require().NotNil(bridge.Rollup)
	hash, _ := hex.DecodeString("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539")
	s.Require().Equal(hash, bridge.Rollup.AstriaId)
}

func (s *StorageTestSuite) TestBridgeByRoles() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	bridges, err := s.storage.Bridges.ByRoles(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(bridges, 1)

	bridge := bridges[0]
	s.Require().EqualValues(7316, bridge.InitHeight)
	s.Require().EqualValues(1, bridge.AddressId)
	s.Require().EqualValues(1, bridge.SudoId)
	s.Require().EqualValues(1, bridge.WithdrawerId)
	s.Require().EqualValues(1, bridge.RollupId)
	s.Require().EqualValues("nria", bridge.Asset)
	s.Require().EqualValues("nria", bridge.FeeAsset)

	s.Require().NotNil(bridge.Address)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Address.Hash)

	s.Require().NotNil(bridge.Sudo)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Sudo.Hash)

	s.Require().NotNil(bridge.Withdrawer)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Withdrawer.Hash)

	s.Require().NotNil(bridge.Rollup)
	hash, _ := hex.DecodeString("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539")
	s.Require().Equal(hash, bridge.Rollup.AstriaId)
}

func (s *StorageTestSuite) TestBridgeListWithAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	bridges, err := s.storage.Bridges.ListWithAddress(ctx, 1, 0)
	s.Require().NoError(err)
	s.Require().Len(bridges, 1)

	bridge := bridges[0]
	s.Require().EqualValues(7316, bridge.InitHeight)
	s.Require().EqualValues(1, bridge.AddressId)
	s.Require().EqualValues(1, bridge.SudoId)
	s.Require().EqualValues(1, bridge.WithdrawerId)
	s.Require().EqualValues(1, bridge.RollupId)
	s.Require().EqualValues("nria", bridge.Asset)
	s.Require().EqualValues("nria", bridge.FeeAsset)

	s.Require().NotNil(bridge.Address)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", bridge.Address.Hash)

	s.Require().Nil(bridge.Sudo)
	s.Require().Nil(bridge.Withdrawer)
	s.Require().Nil(bridge.Rollup)
}
