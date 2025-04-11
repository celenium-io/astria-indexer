// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddPriceFeed, downPriceFeed)
}

func upAddPriceFeed(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TYPE action_type ADD VALUE ? AFTER ?`, types.ActionTypeRecoverIbcClient.String(), types.ActionTypeBridgeTransfer.String())
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `ALTER TYPE action_type ADD VALUE ? AFTER ?`, types.ActionTypePriceFeed.String(), types.ActionTypeRecoverIbcClient.String())
	return err
}
func downPriceFeed(ctx context.Context, db *bun.DB) error {
	return nil
}
