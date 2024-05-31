// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

func saveValidators(
	ctx context.Context,
	tx storage.Transaction,
	validators map[string]*storage.Validator,
) error {
	if len(validators) == 0 {
		return nil
	}

	vals := make([]*storage.Validator, 0)
	for _, val := range validators {
		vals = append(vals, val)
	}

	return tx.UpdateValidators(ctx, vals...)
}
