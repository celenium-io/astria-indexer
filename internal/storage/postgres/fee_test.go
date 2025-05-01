// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestFeeByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	fees, err := s.Fee.ByTxId(ctx, 1, 10, 0)
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

	s.Require().NotNil(fee.Payer.Celestials)
	s.Require().EqualValues("name 2", fee.Payer.Celestials.Id)
	s.Require().EqualValues("some_url", fee.Payer.Celestials.ImageUrl)
	s.Require().EqualValues(types.CelestialsStatusVERIFIED, fee.Payer.Celestials.Status)
	s.Require().EqualValues(2, fee.Payer.Celestials.ChangeId)
}

func (s *StorageTestSuite) TestFeeByPayerId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	fees, err := s.Fee.ByPayerId(ctx, 1, 10, 0, storage.SortOrderDesc)
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

func (s *StorageTestSuite) TestFullTxFee() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	fees, err := s.Fee.FullTxFee(ctx, 1)
	s.Require().NoError(err)
	s.Require().Len(fees, 1)

	fee := fees[0]
	s.Require().EqualValues("100", fee.Amount.String())
	s.Require().EqualValues("nria", fee.Asset)
}
