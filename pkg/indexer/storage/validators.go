// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

func (module *Module) saveValidators(
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

	if err := tx.SaveValidators(ctx, vals...); err != nil {
		return err
	}

	for i := range vals {
		if _, ok := module.validators[vals[i].Address]; !ok {
			module.validators[vals[i].Address] = vals[i].Id
		}
	}

	return nil
}
