// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestTxByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("20b0e6310801e7b2a16c69aace7b1a1d550e5c49c80f546941bb1ac747487fe5")
	s.Require().NoError(err)

	tx, err := s.Tx.ByHash(ctx, hash)
	s.Require().NoError(err)

	s.Require().EqualValues(7316, tx.Height)
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(0, tx.Nonce)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().EqualValues(2, tx.Position)
	s.Require().EqualValues(1, tx.SignerId)
	s.Require().EqualValues(types.StatusSuccess, tx.Status)
	s.Require().EqualValues(hash, tx.Hash)
	s.Require().NotNil(tx.Signer)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", tx.Signer.Hash)
	s.Require().NotNil(tx.Signer.Celestials)
	s.Require().EqualValues("name 2", tx.Signer.Celestials.Id)
	s.Require().EqualValues("some_url", tx.Signer.Celestials.ImageUrl)
	s.Require().EqualValues(types.CelestialsStatusVERIFIED, tx.Signer.Celestials.Status)
	s.Require().EqualValues(2, tx.Signer.Celestials.ChangeId)
}

func (s *StorageTestSuite) TestTxByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.Tx.ByHeight(ctx, 7316, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(7316, tx.Height)
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(0, tx.Nonce)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().EqualValues(2, tx.Position)
	s.Require().EqualValues(1, tx.SignerId)
	s.Require().EqualValues(types.StatusSuccess, tx.Status)

	hash, err := hex.DecodeString("20b0e6310801e7b2a16c69aace7b1a1d550e5c49c80f546941bb1ac747487fe5")
	s.Require().NoError(err)
	s.Require().EqualValues(hash, tx.Hash)

	s.Require().NotNil(tx.Signer)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", tx.Signer.Hash)
	s.Require().NotNil(tx.Signer.Celestials)
	s.Require().EqualValues("name 2", tx.Signer.Celestials.Id)
	s.Require().EqualValues("some_url", tx.Signer.Celestials.ImageUrl)
	s.Require().EqualValues(types.CelestialsStatusVERIFIED, tx.Signer.Celestials.Status)
	s.Require().EqualValues(2, tx.Signer.Celestials.ChangeId)
}

func (s *StorageTestSuite) TestTxFilter() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.Tx.Filter(ctx, storage.TxFilter{
		Limit:       10,
		WithActions: true,
		TimeFrom:    time.Date(2023, 11, 30, 23, 52, 23, 0, time.UTC),
		Sort:        sdk.SortOrderAsc,
		ActionTypes: types.NewActionTypeMask(types.ActionTypeRollupDataSubmission.String()),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(7316, tx.Height)
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(0, tx.Nonce)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().EqualValues(2, tx.Position)
	s.Require().EqualValues(1, tx.SignerId)
	s.Require().EqualValues(types.StatusSuccess, tx.Status)

	hash, err := hex.DecodeString("20b0e6310801e7b2a16c69aace7b1a1d550e5c49c80f546941bb1ac747487fe5")
	s.Require().NoError(err)
	s.Require().EqualValues(hash, tx.Hash)

	s.Require().NotNil(tx.Signer)
	s.Require().Equal("astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p", tx.Signer.Hash)
	s.Require().NotNil(tx.Signer.Celestials)
	s.Require().EqualValues("name 2", tx.Signer.Celestials.Id)
	s.Require().EqualValues("some_url", tx.Signer.Celestials.ImageUrl)
	s.Require().EqualValues(types.CelestialsStatusVERIFIED, tx.Signer.Celestials.Status)
	s.Require().EqualValues(2, tx.Signer.Celestials.ChangeId)

	s.Require().Len(tx.Actions, 1)
}

func (s *StorageTestSuite) TestTxByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	txs, err := s.Tx.ByAddress(ctx, 1, storage.TxFilter{
		Limit:       10,
		WithActions: true,
		TimeFrom:    time.Date(2023, 11, 30, 23, 52, 23, 0, time.UTC),
		Sort:        sdk.SortOrderAsc,
		ActionTypes: types.NewActionTypeMask(types.ActionTypeRollupDataSubmission.String()),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(7316, tx.Height)
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(0, tx.Nonce)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().EqualValues(2, tx.Position)
	s.Require().EqualValues(1, tx.SignerId)
	s.Require().EqualValues(types.StatusSuccess, tx.Status)

	hash, err := hex.DecodeString("20b0e6310801e7b2a16c69aace7b1a1d550e5c49c80f546941bb1ac747487fe5")
	s.Require().NoError(err)
	s.Require().EqualValues(hash, tx.Hash)

	s.Require().NotNil(tx.Signer)
	s.Require().NotNil(tx.Signer.Celestials)
	s.Require().EqualValues("name 2", tx.Signer.Celestials.Id)
	s.Require().EqualValues("some_url", tx.Signer.Celestials.ImageUrl)
	s.Require().EqualValues(types.CelestialsStatusVERIFIED, tx.Signer.Celestials.Status)
	s.Require().EqualValues(2, tx.Signer.Celestials.ChangeId)

	s.Require().Len(tx.Actions, 1)
}
