// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upTotalBytes, downTotalBytes)
}

func upTotalBytes(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE IF EXISTS public.state ADD COLUMN IF NOT EXISTS total_bytes int8 NULL`)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `COMMENT ON COLUMN public.state.total_bytes IS 'Total rollup bytes'`)
	if err != nil {
		return err
	}

	var totalBytes int64
	err = db.NewSelect().Model((*storage.Rollup)(nil)).ColumnExpr("sum(size)").Scan(ctx, &totalBytes)
	if err != nil {
		return err
	}

	_, err = db.NewUpdate().Set("total_bytes = ?", totalBytes).Model((*storage.State)(nil)).Where("name = 'dipdup_astria_indexer'").Exec(ctx)
	return err
}
func downTotalBytes(ctx context.Context, db *bun.DB) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE IF EXISTS public.state DROP COLUMN IF NOT EXISTS total_bytes`)
	return err
}
