// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/aopoltorzhicky/astria/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestAddressByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("3fff1c39b9d163bfb9bcbf9dfea78675f1b4bc2c")
	s.Require().NoError(err)

	address, err := s.storage.Address.ByHash(ctx, hash)
	s.Require().NoError(err)

	s.Require().EqualValues(0, address.Height)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(1, address.Nonce)
	s.Require().EqualValues(1, address.ActionsCount)
	s.Require().EqualValues(2, address.SignedTxCount)
	s.Require().EqualValues(hash, address.Hash)
}

func (s *StorageTestSuite) TestAddressListWithBalances() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	addresses, err := s.storage.Address.ListWithBalance(ctx, storage.AddressListFilter{
		Sort:  sdk.SortOrderAsc,
		Limit: 1,
	})
	s.Require().NoError(err)
	s.Require().Len(addresses, 1)

	address := addresses[0]
	s.Require().EqualValues(0, address.Height)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(1, address.Nonce)
	s.Require().EqualValues(1, address.ActionsCount)
	s.Require().EqualValues(2, address.SignedTxCount)

	hash, err := hex.DecodeString("3fff1c39b9d163bfb9bcbf9dfea78675f1b4bc2c")
	s.Require().NoError(err)
	s.Require().EqualValues(hash, address.Hash)
}
