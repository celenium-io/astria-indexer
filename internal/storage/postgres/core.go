// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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

// Storage -
type Storage struct {
	*postgres.Storage

	cfg        config.Database
	scriptsDir string

	Blocks          models.IBlock
	BlockStats      models.IBlockStats
	Bridges         models.IBridge
	Constants       models.IConstant
	Tx              models.ITx
	Transfers       models.ITransfer
	Fee             models.IFee
	Deposit         models.IDeposit
	Action          models.IAction
	Address         models.IAddress
	Rollup          models.IRollup
	BlockSignatures models.IBlockSignature
	Validator       models.IValidator
	State           models.IState
	Search          models.ISearch
	Stats           models.IStats
	App             models.IApp
	Asset           models.IAsset
	Notificator     *Notificator
}

// Create -
func Create(ctx context.Context, cfg config.Database, scriptsDir string, withMigrations bool) (Storage, error) {
	init := initDatabase
	if withMigrations {
		init = initDatabaseWithMigrations
	}
	strg, err := postgres.Create(ctx, cfg, init)
	if err != nil {
		return Storage{}, err
	}

	s := Storage{
		cfg:             cfg,
		scriptsDir:      scriptsDir,
		Storage:         strg,
		Blocks:          NewBlocks(strg.Connection()),
		BlockStats:      NewBlockStats(strg.Connection()),
		Bridges:         NewBridge(strg.Connection()),
		Constants:       NewConstant(strg.Connection()),
		Action:          NewAction(strg.Connection()),
		Fee:             NewFee(strg.Connection()),
		Deposit:         NewDeposit(strg.Connection()),
		Address:         NewAddress(strg.Connection()),
		BlockSignatures: NewBlockSignature(strg.Connection()),
		Rollup:          NewRollup(strg.Connection()),
		Tx:              NewTx(strg.Connection()),
		Transfers:       NewTransfer(strg.Connection()),
		Validator:       NewValidator(strg.Connection()),
		State:           NewState(strg.Connection()),
		Search:          NewSearch(strg.Connection()),
		App:             NewApp(strg.Connection()),
		Stats:           NewStats(strg.Connection()),
		Asset:           NewAsset(strg.Connection()),
		Notificator:     NewNotificator(cfg, strg.Connection().DB()),
	}

	if err := s.createScripts(ctx, strg.Connection(), "functions", false); err != nil {
		return s, errors.Wrap(err, "creating functions")
	}
	if err := s.createScripts(ctx, strg.Connection(), "views", true); err != nil {
		return s, errors.Wrap(err, "creating views")
	}

	return s, nil
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

func (s Storage) CreateListener() models.Listener {
	return NewNotificator(s.cfg, s.Notificator.db)
}
