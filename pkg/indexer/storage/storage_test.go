// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	internalPg "github.com/celenium-io/astria-indexer/internal/storage/postgres"
	indexerCfg "github.com/celenium-io/astria-indexer/pkg/indexer/config"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
)

const testIndexerName = "dipdup_astria_indexer"

// ModuleTestSuite -
type ModuleTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       *postgres.Storage
	notificator   storage.Notificator
	blocks        storage.IBlock
	state         storage.IState
}

// SetupSuite -
func (s *ModuleTestSuite) SetupSuite() {
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

	dbCfg := config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}

	strg, err := internalPg.Create(ctx, dbCfg, "../../../database", false)
	s.Require().NoError(err)
	s.storage = strg
	s.notificator = internalPg.NewNotificator(dbCfg, strg)
	s.blocks = internalPg.NewBlocks(strg)
	s.state = internalPg.NewState(strg)
}

// TearDownSuite -
func (s *ModuleTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *ModuleTestSuite) BeforeTest(suiteName, testName string) {
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

func (s *ModuleTestSuite) TestBlockLast() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	module := NewModule(s.storage.Transactable, s.notificator, indexerCfg.Indexer{Name: testIndexerName})
	module.Start(ctx)

	hash, err := hex.DecodeString("F44BC94BF7D064ADF82618F2691D2353161DE232ECB3091B7E5C89B453C79456")
	s.Require().NoError(err)

	module.MustInput(InputName).Push(&storage.Block{
		Height:          10001,
		Hash:            hash,
		VersionBlock:    11,
		VersionApp:      1,
		ProposerAddress: "astria16rgmx2s86kk2r69rhjnvs9y44ujfhadc7yav9a",
		Time:            time.Now(),
		Stats:           &storage.BlockStats{},
	})
	time.Sleep(time.Second)

	block, err := s.blocks.Last(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(10001, block.Height)
	s.Require().EqualValues(1, block.VersionApp)
	s.Require().EqualValues(11, block.VersionBlock)
	s.Require().Equal(hash, block.Hash.Bytes())

	state, err := s.state.ByName(ctx, testIndexerName)
	s.Require().NoError(err)
	s.Require().Equal(testIndexerName, state.Name)
	s.Require().EqualValues(10001, state.LastHeight)

	s.Require().NoError(module.Close())
}

func TestSuiteModule_Run(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
