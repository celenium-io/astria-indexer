// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddBridgeTransfer, downAddBridgeTransfer)
}

func upAddBridgeTransfer(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TYPE action_type ADD VALUE ? AFTER ?`, types.ActionTypeRecoverIbcClient.String(), types.ActionTypeBridgeTransfer.String())
	return err
}
func downAddBridgeTransfer(ctx context.Context, db *bun.DB) error {
	return nil
}
