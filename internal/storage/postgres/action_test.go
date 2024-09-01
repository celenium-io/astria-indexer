// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
)

func (s *StorageTestSuite) TestActionByBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	actions, err := s.storage.Action.ByBlock(ctx, 7316, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(7316, action.Height)
	s.Require().EqualValues(1, action.Id)
	s.Require().EqualValues(0, action.Position)
	s.Require().EqualValues(1, action.TxId)
	s.Require().EqualValues(types.ActionTypeSequence, action.Type)
	s.Require().NotNil(action.Data)
	s.Require().NotNil(action.Fee)
	s.Require().NotNil(action.Tx)
	s.Require().NotEmpty(action.Tx.Hash)
}

func (s *StorageTestSuite) TestActionByTxId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	actions, err := s.storage.Action.ByTxId(ctx, 1, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(7316, action.Height)
	s.Require().EqualValues(1, action.Id)
	s.Require().EqualValues(0, action.Position)
	s.Require().EqualValues(1, action.TxId)
	s.Require().EqualValues(types.ActionTypeSequence, action.Type)
	s.Require().NotNil(action.Data)
	s.Require().NotNil(action.Fee)
}

func (s *StorageTestSuite) TestActionByAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	actions, err := s.storage.Action.ByAddress(ctx, 1, storage.AddressActionsFilter{
		Sort:        sdk.SortOrderAsc,
		Limit:       1,
		ActionTypes: types.NewActionTypeMask(types.ActionTypeSequence.String()),
	})
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(7316, action.Height)
	s.Require().EqualValues(1, action.ActionId)
	s.Require().EqualValues(1, action.ActionId)
	s.Require().NotNil(action.Tx)
	s.Require().NotNil(action.Action)
	s.Require().EqualValues(1, action.Action.TxId)
	s.Require().EqualValues(types.ActionTypeSequence, action.Action.Type)
	s.Require().NotNil(action.Action.Data)
	s.Require().NotNil(action.Action.Fee)
}

func (s *StorageTestSuite) TestActionByRollup() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	actions, err := s.storage.Action.ByRollup(ctx, 1, 1, 0, sdk.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(7316, action.Height)
	s.Require().EqualValues(1, action.ActionId)
	s.Require().EqualValues(1, action.ActionId)
	s.Require().NotNil(action.Tx)
	s.Require().NotNil(action.Action)
	s.Require().EqualValues(1, action.Action.TxId)
	s.Require().EqualValues(types.ActionTypeSequence, action.Action.Type)
	s.Require().NotNil(action.Action.Data)
	s.Require().NotNil(action.Action.Fee)
}

func (s *StorageTestSuite) TestActionByRollupAndBridge() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	actions, err := s.storage.Action.ByRollupAndBridge(ctx, 1, storage.RollupAndBridgeActionsFilter{
		Sort:          sdk.SortOrderAsc,
		Limit:         10,
		Offset:        0,
		RollupActions: true,
		BridgeActions: true,
	})
	s.Require().NoError(err)
	s.Require().Len(actions, 2)

	action := actions[0]
	s.Require().EqualValues(7316, action.Height)
	s.Require().NotNil(action.Tx)
	s.Require().NotNil(action.Action)
	s.Require().EqualValues(1, action.Action.TxId)
	s.Require().EqualValues(types.ActionTypeSequence, action.Action.Type)
	s.Require().NotNil(action.Action.Data)
	s.Require().NotNil(action.Action.Fee)
}

func (s *StorageTestSuite) TestActionByRollupAndBridgeWithoutRollupActions() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	actions, err := s.storage.Action.ByRollupAndBridge(ctx, 1, storage.RollupAndBridgeActionsFilter{
		Sort:          sdk.SortOrderAsc,
		Limit:         10,
		Offset:        0,
		RollupActions: false,
		BridgeActions: true,
	})
	s.Require().NoError(err)
	s.Require().Len(actions, 2)

	action := actions[0]
	s.Require().EqualValues(7316, action.Height)
	s.Require().NotNil(action.Tx)
	s.Require().NotNil(action.Action)
	s.Require().EqualValues(1, action.Action.TxId)
	s.Require().EqualValues(types.ActionTypeSequence, action.Action.Type)
	s.Require().NotNil(action.Action.Data)
	s.Require().NotNil(action.Action.Fee)
}
