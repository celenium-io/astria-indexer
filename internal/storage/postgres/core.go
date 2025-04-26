// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	models "github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres/migrations"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

// Create -
func Create(ctx context.Context, cfg config.Database, scriptsDir string, withMigrations bool) (*postgres.Storage, error) {
	init := initDatabase
	if withMigrations {
		init = initDatabaseWithMigrations
	}
	strg, err := postgres.Create(ctx, cfg, init)
	if err != nil {
		return nil, errors.Wrap(err, "create database")
	}

	if err := createScripts(ctx, strg.Connection(), scriptsDir, "functions", false); err != nil {
		return nil, errors.Wrap(err, "creating functions")
	}
	if err := createScripts(ctx, strg.Connection(), scriptsDir, "views", true); err != nil {
		return nil, errors.Wrap(err, "creating views")
	}

	return strg, nil
}

func initDatabase(ctx context.Context, conn *database.Bun) error {
	if err := createExtensions(ctx, conn); err != nil {
		return errors.Wrap(err, "create extensions")
	}

	if err := createTypes(ctx, conn); err != nil {
		return errors.Wrap(err, "creating custom types")
	}

	// register many-to-many relationships
	conn.DB().RegisterModel(
		(*models.RollupAction)(nil),
		(*models.RollupAddress)(nil),
		(*models.AddressAction)(nil),
	)

	if err := database.CreateTables(ctx, conn, models.Models...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return err
	}

	if err := database.MakeComments(ctx, conn, models.Models...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return errors.Wrap(err, "make comments")
	}

	if err := createHypertables(ctx, conn); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return errors.Wrap(err, "create hypertables")
	}

	return createIndices(ctx, conn)
}

func initDatabaseWithMigrations(ctx context.Context, conn *database.Bun) error {
	if err := initDatabase(ctx, conn); err != nil {
		return err
	}
	return migrateDatabase(ctx, conn)
}

func migrateDatabase(ctx context.Context, db *database.Bun) error {
	migrator := migrate.NewMigrator(db.DB(), migrations.Migrations)
	if err := migrator.Init(ctx); err != nil {
		return err
	}
	if err := migrator.Lock(ctx); err != nil {
		return err
	}
	defer migrator.Unlock(ctx) //nolint:errcheck

	_, err := migrator.Migrate(ctx)
	return err
}

func createHypertables(ctx context.Context, conn *database.Bun) error {
	return conn.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		for _, model := range []storage.Model{
			&models.Block{},
			&models.BlockStats{},
			&models.Tx{},
			&models.Action{},
			&models.BlockSignature{},
			&models.RollupAction{},
			&models.Fee{},
			&models.Transfer{},
			&models.Deposit{},
			&models.Price{},
		} {
			if _, err := tx.ExecContext(ctx,
				`SELECT create_hypertable(?, 'time', chunk_time_interval => INTERVAL '1 month', if_not_exists => TRUE);`,
				model.TableName(),
			); err != nil {
				return err
			}
		}
		return nil
	})
}

func createExtensions(ctx context.Context, conn *database.Bun) error {
	return conn.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pg_trgm;")
		return err
	})
}
