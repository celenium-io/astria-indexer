// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

func (s *StorageTestSuite) TestConstantGet() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	c, err := s.Constants.Get(ctx, "block", "block_max_bytes")
	s.Require().NoError(err)

	s.Require().EqualValues("block", c.Module)
	s.Require().EqualValues("block_max_bytes", c.Name)
	s.Require().EqualValues("22020096", c.Value)
}

func (s *StorageTestSuite) TestConstantByModule() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	consts, err := s.Constants.ByModule(ctx, "block")
	s.Require().NoError(err)
	s.Require().Len(consts, 2)

	s.Require().EqualValues("block", consts[0].Module)
	s.Require().EqualValues("block_max_bytes", consts[0].Name)
	s.Require().EqualValues("22020096", consts[0].Value)

	s.Require().EqualValues("block", consts[1].Module)
	s.Require().EqualValues("block_max_gas", consts[1].Name)
	s.Require().EqualValues("-1", consts[1].Value)
}

func (s *StorageTestSuite) TestConstantAll() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	consts, err := s.Constants.All(ctx)
	s.Require().NoError(err)
	s.Require().Len(consts, 9)
}

func (s *StorageTestSuite) TestConstantIsNoRows() {
	s.Require().True(s.Constants.IsNoRows(sql.ErrNoRows))
	s.Require().True(s.Constants.IsNoRows(errors.Wrap(sql.ErrNoRows, "some text")))
	s.Require().False(s.Constants.IsNoRows(errors.New("test")))
}
