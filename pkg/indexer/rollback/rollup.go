// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rollback

import (
	"context"
	"encoding/base64"

	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
)

func rollbackRollups(
	ctx context.Context,
	tx storage.Transaction,
	height types.Level,
	actions []storage.Action,
) (int64, error) {
	rollups, err := tx.RollbackRollups(ctx, height)
	if err != nil {
		return 0, err
	}

	rollbackActions, err := tx.RollbackRollupActions(ctx, height)
	if err != nil {
		return 0, err
	}

	if err := tx.RollbackRollupAddresses(ctx, height); err != nil {
		return 0, err
	}

	m := make(map[uint64]*storage.Rollup)
	for i := range rollups {
		m[rollups[i].Id] = &rollups[i]
	}

	mapActions := make(map[uint64]storage.Action)
	for i := range actions {
		mapActions[actions[i].Id] = actions[i]
	}

	updates := make(map[uint64]*storage.Rollup)
	for i := range rollbackActions {
		if _, ok := m[rollbackActions[i].RollupId]; ok {
			continue
		}
		action, ok := mapActions[rollbackActions[i].ActionId]
		if !ok {
			return 0, errors.Errorf("can't find action with id: %d", rollbackActions[i].ActionId)
		}
		if err := updateRollups(updates, rollbackActions[i].RollupId, action); err != nil {
			return 0, err
		}
	}

	arr := make([]*storage.Rollup, 0)
	for i := range updates {
		arr = append(arr, updates[i])
	}

	if err := tx.UpdateRollups(ctx, arr...); err != nil {
		return 0, err
	}

	return int64(len(rollups)), nil
}

func updateRollups(updates map[uint64]*storage.Rollup, rollupId uint64, action storage.Action) error {
	if action.Type != storageTypes.ActionTypeRollupDataSubmission {
		return errors.Errorf("invalid action type: %s", action.Type)
	}
	size, err := getActionSize(action)
	if err != nil {
		return err
	}
	if update, ok := updates[rollupId]; ok {
		update.ActionsCount -= 1
		update.Size -= size
	} else {
		updates[rollupId] = &storage.Rollup{
			Id:           rollupId,
			Size:         -size,
			ActionsCount: -1,
		}
	}

	return nil
}

func getActionSize(action storage.Action) (int64, error) {
	data, ok := action.Data["data"]
	if !ok {
		return 0, errors.Errorf("can't find 'data' in (%d) %##v", action.Id, action.Data)
	}
	str, ok := data.(string)
	if !ok {
		return 0, errors.Errorf("invalid 'data' type in (%d) %##v", action.Id, action.Data)
	}
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return 0, err
	}
	return int64(len(bytes)), nil
}
