// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/aopoltorzhicky/astria/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
)

// StatsTestSuite -
type StatsTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *StatsTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer ctxCancel()

	psqlContainer, err := database.NewPostgreSQLContainer(ctx, database.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "timescale/timescaledb-ha:pg15-latest",
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
	}, "../../../database")
	s.Require().NoError(err)
	s.storage = strg

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
func (s *StatsTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *StatsTestSuite) TestSeries() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Stats.Series(ctx, storage.TimeframeHour, storage.SeriesDataSize, storage.NewSeriesRequest(0, 0))
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func (s *StatsTestSuite) TestSeriesWithFrom() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Stats.Series(ctx, storage.TimeframeHour, storage.SeriesDataSize, storage.NewSeriesRequest(1706018798, 0))
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *StatsTestSuite) TestRollupSeries() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Stats.RollupSeries(ctx, 1, storage.TimeframeHour, storage.RollupSeriesActionsCount, storage.NewSeriesRequest(0, 0))
	s.Require().NoError(err)
	s.Require().Len(items, 1)
}

func (s *StatsTestSuite) TestRollupSeriesWithFrom() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	items, err := s.storage.Stats.RollupSeries(ctx, 1, storage.TimeframeHour, storage.RollupSeriesActionsCount, storage.NewSeriesRequest(1706018798, 0))
	s.Require().NoError(err)
	s.Require().Len(items, 0)
}

func (s *StatsTestSuite) TestSummary() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	summary, err := s.storage.Stats.Summary(ctx)
	s.Require().NoError(err)

	s.Require().EqualValues(0.0038194444444444443, summary.BPS)
	s.Require().EqualValues(2327, summary.BlockTime)
	s.Require().EqualValues(330, summary.BytesInBlock)
	s.Require().EqualValues(0, summary.DataSize)
	s.Require().EqualValues("0", summary.Fee.String())
	s.Require().EqualValues(0, summary.RBPS)
	s.Require().EqualValues("0", summary.Supply.String())
	s.Require().EqualValues(1, summary.TxCount)
}

func TestSuiteStats_Run(t *testing.T) {
	suite.Run(t, new(StatsTestSuite))
}
