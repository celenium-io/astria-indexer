// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
)

// StorageTestSuite -
type StorageTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer

	storage *postgres.Storage

	Blocks          storage.IBlock
	BlockStats      storage.IBlockStats
	Bridges         storage.IBridge
	Constants       storage.IConstant
	Tx              storage.ITx
	Transfers       storage.ITransfer
	Fee             storage.IFee
	Deposit         storage.IDeposit
	Action          storage.IAction
	Address         storage.IAddress
	Rollup          storage.IRollup
	BlockSignatures storage.IBlockSignature
	Validator       storage.IValidator
	State           storage.IState
	Search          storage.ISearch
	App             storage.IApp
	Asset           storage.IAsset
}

// SetupSuite -
func (s *StorageTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer ctxCancel()

	psqlContainer, err := database.NewPostgreSQLContainer(ctx, database.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb-ha:pg15.8-ts2.17.0-all",
	})
	s.Require().NoError(err)
	s.psqlContainer = psqlContainer

	strg, err := Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}, "../../../database", false)
	s.Require().NoError(err)
	s.storage = strg

	s.Blocks = NewBlocks(s.storage)
	s.BlockStats = NewBlockStats(s.storage)
	s.Bridges = NewBridge(s.storage)
	s.Constants = NewConstant(s.storage)
	s.Tx = NewTx(s.storage)
	s.Transfers = NewTransfer(s.storage)
	s.Fee = NewFee(s.storage)
	s.Deposit = NewDeposit(s.storage)
	s.Action = NewAction(s.storage)
	s.Address = NewAddress(s.storage)
	s.Rollup = NewRollup(s.storage)
	s.BlockSignatures = NewBlockSignature(s.storage)
	s.Validator = NewValidator(s.storage)
	s.State = NewState(s.storage)
	s.Search = NewSearch(s.storage)
	s.App = NewApp(s.storage)
	s.Asset = NewAsset(s.storage)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("timescaledb"),
		testfixtures.Directory("../../../test/data"),
		testfixtures.UseAlterConstraint(),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

// TearDownSuite -
func (s *StorageTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func TestSuiteStorage_Run(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
