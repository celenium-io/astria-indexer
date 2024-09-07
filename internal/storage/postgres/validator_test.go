// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestValidatorListByPower() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	validators, err := s.storage.Validator.ListByPower(ctx, 1, 0, storage.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(validators, 1)

	val := validators[0]
	s.Require().EqualValues(3, val.Id)
	s.Require().EqualValues("2", val.Power.String())
	s.Require().EqualValues("node2", val.Name)
	s.Require().EqualValues("astria1c220qfmjrwqlk939ca5a5z2rjxryyr9m3ah8gl", val.Address)
}
