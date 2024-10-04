// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestFeeByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	fees, err := s.storage.Fee.ByTxId(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(fees, 1)

	fee := fees[0]
	s.Require().EqualValues(7316, fee.Height)
	s.Require().EqualValues(1, fee.Id)
	s.Require().EqualValues(1, fee.ActionId)
	s.Require().EqualValues(1, fee.TxId)
	s.Require().EqualValues("100", fee.Amount.String())
	s.Require().EqualValues("nria", fee.Asset)
	s.Require().NotNil(fee.Payer)
	s.Require().NotEmpty(fee.Payer.Hash)
}

func (s *StorageTestSuite) TestFeeByPayerId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	fees, err := s.storage.Fee.ByPayerId(ctx, 1, 10, 0, storage.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(fees, 1)

	fee := fees[0]
	s.Require().EqualValues(7316, fee.Height)
	s.Require().EqualValues(1, fee.Id)
	s.Require().EqualValues(1, fee.ActionId)
	s.Require().EqualValues(1, fee.TxId)
	s.Require().EqualValues("100", fee.Amount.String())
	s.Require().EqualValues("nria", fee.Asset)
	s.Require().NotNil(fee.Tx)
	s.Require().NotEmpty(fee.Tx.Hash)
}
