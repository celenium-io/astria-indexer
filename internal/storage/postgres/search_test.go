// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

func (s *StorageTestSuite) TestSearchBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "b15d072afc508558b3e962060c701a695af5d6a041d4a25c63240bbff5064b3b")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("b15d072afc508558b3e962060c701a695af5d6a041d4a25c63240bbff5064b3b", result.Value)
	s.Require().EqualValues("block", result.Type)
}

func (s *StorageTestSuite) TestSearchBlockByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "7965")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("b15d072afc508558b3e962060c701a695af5d6a041d4a25c63240bbff5064b3b", result.Value)
	s.Require().EqualValues("block", result.Type)
}

func (s *StorageTestSuite) TestSearchTx() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "20b0e6310801e7b2a16c69aace7b1a1d550e5c49c80f546941bb1ac747487fe5")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("20b0e6310801e7b2a16c69aace7b1a1d550e5c49c80f546941bb1ac747487fe5", result.Value)
	s.Require().EqualValues("tx", result.Type)
}

func (s *StorageTestSuite) TestSearchRollup() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539", result.Value)
	s.Require().EqualValues("rollup", result.Type)
}

func (s *StorageTestSuite) TestSearchAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "astria1e9q7egqgz8rz6aej8nr57swqgaeujhz04vd9q5")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("astria1e9q7egqgz8rz6aej8nr57swqgaeujhz04vd9q5", result.Value)
	s.Require().EqualValues("address", result.Type)
	s.Require().EqualValues(8, result.Id)
}

func (s *StorageTestSuite) TestSearchValidator() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "node0")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("node0", result.Value)
	s.Require().EqualValues("validator", result.Type)
	s.Require().EqualValues(1, result.Id)
}

func (s *StorageTestSuite) TestSearchValidatorByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "astria16rgmx2s86kk2r69rhjnvs9y44ujfhadc7yav9a")
	s.Require().NoError(err)
	s.Require().Len(results, 2)

	result := results[0]
	s.Require().EqualValues("astria16rgmx2s86kk2r69rhjnvs9y44ujfhadc7yav9a", result.Value)
	s.Require().EqualValues("address", result.Type)
	s.Require().EqualValues(3, result.Id)

	result1 := results[1]
	s.Require().EqualValues("node0", result1.Value)
	s.Require().EqualValues("validator", result1.Type)
	s.Require().EqualValues(1, result1.Id)
}

func (s *StorageTestSuite) TestSearchBridge() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "nri")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("nria", result.Value)
	s.Require().EqualValues("bridge", result.Type)
	s.Require().EqualValues(1, result.Id)
}

func (s *StorageTestSuite) TestSearchRollupByHash() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "GbqKuz5LVqMJ32dWxHuX4pjjpy2IRJ02oPrbHKc2ZTk=")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("19ba8abb3e4b56a309df6756c47b97e298e3a72d88449d36a0fadb1ca7366539", result.Value)
	s.Require().EqualValues("rollup", result.Type)
}

func (s *StorageTestSuite) TestSearchApp() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	results, err := s.Search.Search(ctx, "p 1")
	s.Require().NoError(err)
	s.Require().Len(results, 1)

	result := results[0]
	s.Require().EqualValues("App 1", result.Value)
	s.Require().EqualValues("app", result.Type)
	s.Require().EqualValues(1, result.Id)
}
