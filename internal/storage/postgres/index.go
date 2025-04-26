// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

func createIndices(ctx context.Context, conn *database.Bun) error {
	log.Info().Msg("creating indexes...")
	return conn.DB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// Address
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Address)(nil)).
			Index("address_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}

		// Block
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Block)(nil)).
			Index("block_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Block)(nil)).
			Index("block_proposer_id_idx").
			Column("proposer_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Block)(nil)).
			Index("block_hash_idx").
			Column("hash").
			Using("HASH").
			Exec(ctx); err != nil {
			return err
		}

		// BlockStats
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.BlockStats)(nil)).
			Index("block_stats_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}

		// Tx
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Tx)(nil)).
			Index("tx_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Tx)(nil)).
			Index("tx_hash_idx").
			Column("hash").
			Using("HASH").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Tx)(nil)).
			Index("tx_status_idx").
			Column("status").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Tx)(nil)).
			Index("tx_signer_id_idx").
			Column("signer_id").
			Exec(ctx); err != nil {
			return err
		}

		// Action
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Action)(nil)).
			Index("action_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Action)(nil)).
			Index("action_tx_id_idx").
			Column("tx_id").
			Where("tx_id IS NOT NULL").
			Exec(ctx); err != nil {
			return err
		}

		// Rollup
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Rollup)(nil)).
			Index("rollup_first_height_idx").
			Column("first_height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Rollup)(nil)).
			Index("rollup_hash_idx").
			Column("astria_id").
			Using("HASH").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Rollup)(nil)).
			Index("rollup_size_idx").
			Column("size").
			Exec(ctx); err != nil {
			return err
		}

		// Rollup actions
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.RollupAction)(nil)).
			Index("rollup_action_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.RollupAction)(nil)).
			Index("rollup_action_tx_id_idx").
			Column("tx_id").
			Where("tx_id IS NOT NULL").
			Exec(ctx); err != nil {
			return err
		}

		// Address actions
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.AddressAction)(nil)).
			Index("address_action_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.AddressAction)(nil)).
			Index("address_action_tx_id_idx").
			Column("tx_id").
			Where("tx_id IS NOT NULL").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.AddressAction)(nil)).
			Index("address_action_type_idx").
			Column("action_type").
			Where("tx_id IS NOT NULL").
			Exec(ctx); err != nil {
			return err
		}

		// Validators
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Validator)(nil)).
			Index("validator_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Validator)(nil)).
			Index("validator_name_idx").
			ColumnExpr("name gin_trgm_ops").
			Using("GIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Validator)(nil)).
			Index("validator_power_idx").
			ColumnExpr("power").
			Exec(ctx); err != nil {
			return err
		}

		// Bridge
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Bridge)(nil)).
			Index("bridge_init_height_idx").
			Column("init_height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Bridge)(nil)).
			Index("bridge_sudo_id_idx").
			Column("sudo_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Bridge)(nil)).
			Index("bridge_withdrawer_id_idx").
			Column("withdrawer_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Bridge)(nil)).
			Index("bridge_rollup_id_idx").
			Column("rollup_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Bridge)(nil)).
			Index("bridge_asset_idx").
			ColumnExpr("asset gin_trgm_ops").
			Using("GIN").
			Exec(ctx); err != nil {
			return err
		}

		// Fee
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Fee)(nil)).
			Index("fee_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Fee)(nil)).
			Index("fee_action_id_idx").
			Column("action_id").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Fee)(nil)).
			Index("fee_tx_id_idx").
			Column("tx_id").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Fee)(nil)).
			Index("fee_asset_idx").
			Column("asset").
			Exec(ctx); err != nil {
			return err
		}

		// Transfer
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Transfer)(nil)).
			Index("transfer_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Transfer)(nil)).
			Index("transfer_src_id_idx").
			Column("src_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Transfer)(nil)).
			Index("transfer_dest_id_idx").
			Column("dest_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Transfer)(nil)).
			Index("transfer_asset_idx").
			Column("asset").
			Exec(ctx); err != nil {
			return err
		}

		// Deposit
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Deposit)(nil)).
			Index("deposit_height_idx").
			Column("height").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Deposit)(nil)).
			Index("deposit_action_id_idx").
			Column("action_id").
			Using("BRIN").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Deposit)(nil)).
			Index("deposit_rollup_id_idx").
			Column("rollup_id").
			Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*storage.Deposit)(nil)).
			Index("deposit_bridge_id_idx").
			Column("bridge_id").
			Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}
