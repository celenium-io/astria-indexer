// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestAddressByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash := "astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p"

	address, err := s.storage.Address.ByHash(ctx, hash)
	s.Require().NoError(err)

	s.Require().EqualValues(0, address.Height)
	s.Require().EqualValues(1, address.Id)
	s.Require().EqualValues(1, address.Nonce)
	s.Require().EqualValues(1, address.ActionsCount)
	s.Require().EqualValues(2, address.SignedTxCount)
	s.Require().EqualValues(hash, address.Hash)
	s.Require().Len(address.Balance, 2)
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
	s.Require().Len(address.Balance, 1)

	s.Require().EqualValues("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", address.Hash)
}
