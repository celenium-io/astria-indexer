// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestBlockByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.ByHeight(ctx, 7965, false)
	s.Require().NoError(err)

	s.Require().EqualValues(7965, block.Height)
	s.Require().EqualValues(2, block.ProposerId)

	appHash, err := hex.DecodeString("0f9c15e37326c29b51e17cb243e621e97c33f4b034b231ffcd9ed5e0a653d584")
	s.Require().NoError(err)
	s.Require().EqualValues(appHash, block.AppHash)

	s.Require().NotNil(block.Proposer)
	s.Require().EqualValues("astria1475jkpuvznd44szgfz8wwdf9w6xh5dx9jwqgvz", block.Proposer.Address)

	s.Require().Nil(block.Stats)
}

func (s *StorageTestSuite) TestBlockByHeightWithStats() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.ByHeight(ctx, 7965, true)
	s.Require().NoError(err)

	s.Require().EqualValues(7965, block.Height)
	s.Require().EqualValues(2, block.ProposerId)

	appHash, err := hex.DecodeString("0f9c15e37326c29b51e17cb243e621e97c33f4b034b231ffcd9ed5e0a653d584")
	s.Require().NoError(err)
	s.Require().EqualValues(appHash, block.AppHash)

	s.Require().NotNil(block.Proposer)
	s.Require().EqualValues("astria1475jkpuvznd44szgfz8wwdf9w6xh5dx9jwqgvz", block.Proposer.Address)

	s.Require().NotNil(block.Stats)
	s.Require().EqualValues(7965, block.Stats.Height)
}

func (s *StorageTestSuite) TestBlockByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := hex.DecodeString("b15d072afc508558b3e962060c701a695af5d6a041d4a25c63240bbff5064b3b")
	s.Require().NoError(err)

	block, err := s.storage.Blocks.ByHash(ctx, hash)
	s.Require().NoError(err)

	s.Require().EqualValues(7965, block.Height)
	s.Require().EqualValues(2, block.ProposerId)

	appHash, err := hex.DecodeString("0f9c15e37326c29b51e17cb243e621e97c33f4b034b231ffcd9ed5e0a653d584")
	s.Require().NoError(err)
	s.Require().EqualValues(appHash, block.AppHash)

	s.Require().NotNil(block.Proposer)
	s.Require().EqualValues("astria1475jkpuvznd44szgfz8wwdf9w6xh5dx9jwqgvz", block.Proposer.Address)
}

func (s *StorageTestSuite) TestBlockLast() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.Last(ctx)
	s.Require().NoError(err)

	s.Require().EqualValues(7965, block.Height)
	s.Require().EqualValues(2, block.ProposerId)

	appHash, err := hex.DecodeString("0f9c15e37326c29b51e17cb243e621e97c33f4b034b231ffcd9ed5e0a653d584")
	s.Require().NoError(err)
	s.Require().EqualValues(appHash, block.AppHash)

	s.Require().NotNil(block.Proposer)
	s.Require().EqualValues("astria1475jkpuvznd44szgfz8wwdf9w6xh5dx9jwqgvz", block.Proposer.Address)
}

func (s *StorageTestSuite) TestBlockListWithStats() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	blocks, err := s.storage.Blocks.ListWithStats(ctx, 1, 0, storage.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(blocks, 1)

	block := blocks[0]
	s.Require().EqualValues(7965, block.Height)
	s.Require().EqualValues(2, block.ProposerId)

	appHash, err := hex.DecodeString("0f9c15e37326c29b51e17cb243e621e97c33f4b034b231ffcd9ed5e0a653d584")
	s.Require().NoError(err)
	s.Require().EqualValues(appHash, block.AppHash)

	s.Require().NotNil(block.Proposer)
	s.Require().EqualValues("astria1475jkpuvznd44szgfz8wwdf9w6xh5dx9jwqgvz", block.Proposer.Address)

	s.Require().NotNil(block.Stats)
}

func (s *StorageTestSuite) TestBlockByProposer() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	blocks, err := s.storage.Blocks.ByProposer(ctx, 2, 1, 0, storage.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(blocks, 1)

	block := blocks[0]
	s.Require().EqualValues(7965, block.Height)
	s.Require().EqualValues(2, block.ProposerId)

	appHash, err := hex.DecodeString("0f9c15e37326c29b51e17cb243e621e97c33f4b034b231ffcd9ed5e0a653d584")
	s.Require().NoError(err)
	s.Require().EqualValues(appHash, block.AppHash)

	s.Require().NotNil(block.Stats)
	s.Require().NotNil(block.Proposer)
}

func (s *StorageTestSuite) TestByIdWithRelations() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	block, err := s.storage.Blocks.ByIdWithRelations(ctx, 7966)
	s.Require().NoError(err)
	s.Require().EqualValues(7965, block.Height)
	s.Require().EqualValues(2, block.ProposerId)

	appHash, err := hex.DecodeString("0f9c15e37326c29b51e17cb243e621e97c33f4b034b231ffcd9ed5e0a653d584")
	s.Require().NoError(err)
	s.Require().EqualValues(appHash, block.AppHash)

	s.Require().NotNil(block.Stats)
	s.Require().NotNil(block.Proposer)
}
