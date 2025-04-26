// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"
)

const testDipdupName = "dipdup_astria_indexer"

func (s *StorageTestSuite) TestStateByName() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	state, err := s.State.ByName(ctx, testDipdupName)
	s.Require().NoError(err)

	s.Require().EqualValues(testDipdupName, state.Name)
	s.Require().EqualValues(1, state.Id)
	s.Require().EqualValues(7965, state.LastHeight)
	s.Require().EqualValues("astria-dusk-2-final", state.ChainId)
	s.Require().EqualValues(2, state.TotalTx)
	s.Require().EqualValues(6, state.TotalAccounts)
	s.Require().EqualValues(1, state.TotalRollups)
	s.Require().EqualValues("1000000000000000000000", state.TotalSupply.String())
}
