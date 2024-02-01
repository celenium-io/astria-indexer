// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"
)

func (s *StorageTestSuite) TestSearchBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("b15d072afc508558b3e962060c701a695af5d6a041d4a25c63240bbff5064b3b")
	s.Require().NoError(err)

	results, err := s.storage.Search.Search(ctx, hash)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("b15d072afc508558b3e962060c701a695af5d6a041d4a25c63240bbff5064b3b", result.Value)
	s.Require().EqualValues("block", result.Type)
}

func (s *StorageTestSuite) TestSearchTx() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("20b0e6310801e7b2a16c69aace7b1a1d550e5c49c80f546941bb1ac747487fe5")
	s.Require().NoError(err)

	results, err := s.storage.Search.Search(ctx, hash)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("20b0e6310801e7b2a16c69aace7b1a1d550e5c49c80f546941bb1ac747487fe5", result.Value)
	s.Require().EqualValues("tx", result.Type)
}

func (s *StorageTestSuite) TestSearchRollup() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539")
	s.Require().NoError(err)

	results, err := s.storage.Search.Search(ctx, hash)
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539", result.Value)
	s.Require().EqualValues("rollup", result.Type)
}

func (s *StorageTestSuite) TestSearchTextValidator() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.storage.Search.SearchText(ctx, "node")
	s.Require().NoError(err)
	s.Require().Len(results, 3)

	result := results[0]
	s.Require().EqualValues("node0", result.Value)
	s.Require().EqualValues("validator", result.Type)
}
