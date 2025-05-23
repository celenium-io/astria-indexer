// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/astria-indexer/internal/test_suite"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

// TransactionTestSuite -
type TransactionTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       *postgres.Storage

	Address        storage.IAddress
	Rollup         storage.IRollup
	BlockSignature storage.IBlockSignature
	Blocks         storage.IBlock
	Constants      storage.IConstant
	Validator      storage.IValidator
	Markets        storage.IMarket
	Prices         storage.IPrice
}

// SetupSuite -
func (s *TransactionTestSuite) SetupSuite() {
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
	s.Address = NewAddress(s.storage)
	s.Rollup = NewRollup(s.storage)
	s.BlockSignature = NewBlockSignature(s.storage)
	s.Constants = NewConstant(s.storage)
	s.Validator = NewValidator(s.storage)
	s.Blocks = NewBlocks(s.storage)
	s.Markets = NewMarket(s.storage)
	s.Prices = NewPrice(s.storage)
}

// TearDownSuite -
func (s *TransactionTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *TransactionTestSuite) BeforeTest(suiteName, testName string) {
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

func TestSuiteTransaction_Run(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}

func (s *TransactionTestSuite) TestSaveAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	addresses := []*storage.Address{
		{
			Height:       1000,
			Hash:         testsuite.RandomAddress(),
			IsIbcRelayer: testsuite.Ptr(true),
		}, {
			Height: 1000,
			Hash:   testsuite.RandomAddress(),
		}, {
			Height:       1000,
			Hash:         testsuite.RandomAddress(),
			IsIbcRelayer: testsuite.Ptr(false),
		}, {
			Height: 1000,
			Hash:   testsuite.RandomAddress(),
		}, {
			Height: 1000,
			Hash:   testsuite.RandomAddress(),
		},
	}

	count1, err := tx.SaveAddresses(ctx, addresses...)
	s.Require().NoError(err)
	s.Require().EqualValues(5, count1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	s.Require().Greater(addresses[0].Id, uint64(0))
	s.Require().Greater(addresses[1].Id, uint64(0))

	replyAddress := storage.Address{
		Height:       1000,
		Hash:         addresses[1].Hash,
		Id:           2,
		IsIbcRelayer: testsuite.Ptr(true),
	}

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	count2, err := tx2.SaveAddresses(ctx, &replyAddress)
	s.Require().NoError(err)
	s.Require().EqualValues(0, count2)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
	s.Require().Equal(replyAddress.Id, addresses[1].Id)

	response, err := s.Address.List(ctx, 10, 6, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(response, 5)

	s.Require().Equal(addresses[0], response[0])
	s.Require().Equal(replyAddress, *response[1])
}

func (s *TransactionTestSuite) TestSaveConstants() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	constants := make([]storage.Constant, 5)
	for i := 0; i < 5; i++ {
		constants[i].Module = types.ModuleNameValues()[i]
		constants[i].Name = fmt.Sprintf("constant_%d", i)
		constants[i].Value = strconv.FormatInt(int64(i), 10)
	}

	err = tx.SaveConstants(ctx, constants...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveTransactions() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	txs := make([]*storage.Tx, 5)
	for i := 0; i < 5; i++ {
		txs[i] = &storage.Tx{
			Height:   pkgTypes.Level(10000),
			Time:     time.Now(),
			Position: int64(i),
			Status:   types.StatusSuccess,
			Hash:     testsuite.RandomHash(32),
			SignerId: uint64(i),
		}
	}

	err = tx.SaveTransactions(ctx, txs...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveActions() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	actions := make([]*storage.Action, 5)
	for i := 0; i < 5; i++ {
		actions[i] = &storage.Action{
			Height:   pkgTypes.Level(10000),
			Time:     time.Now(),
			Position: int64(i),
			Type:     types.ActionTypeValues()[i],
			TxId:     uint64(i),
			Data: map[string]any{
				"rollup_id": i,
			},
		}
	}

	err = tx.SaveActions(ctx, actions...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveValidators() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	validators := make([]*storage.Validator, 5)
	for i := 0; i < 5; i++ {
		pubkey := testsuite.RandomHash(32)
		address := hex.EncodeToString(pubkey[:20])
		validators[i] = &storage.Validator{
			Height:     pkgTypes.Level(10000),
			Address:    address,
			PubkeyType: "tendermint/PubKeyEd25519",
			PubKey:     pubkey,
			Power:      decimal.New(1, 0),
		}
		if i%2 == 0 {
			validators[i].Name = fmt.Sprintf("validator_%d", i)
		}
	}

	err = tx.SaveValidators(ctx, validators...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveRollups() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	rollups := make([]*storage.Rollup, 5)
	for i := 0; i < 5; i++ {
		rollups[i] = &storage.Rollup{
			AstriaId:     testsuite.RandomHash(32),
			FirstHeight:  10000,
			ActionsCount: 1,
			Size:         10,
		}
	}

	count, err := tx.SaveRollups(ctx, rollups...)
	s.Require().NoError(err)
	s.Require().EqualValues(5, count)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	ret, err := s.Rollup.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(ret, 7)
}

func (s *TransactionTestSuite) TestSaveRollupActions() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	ra := make([]*storage.RollupAction, 5)
	for i := 0; i < 5; i++ {
		ra[i] = &storage.RollupAction{
			RollupId:   uint64(i + 1),
			ActionId:   uint64(5 - i),
			Height:     10000,
			ActionType: types.ActionTypeBridgeLock,
		}
	}

	err = tx.SaveRollupActions(ctx, ra...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveRollupAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	ra := make([]*storage.RollupAddress, 5)
	for i := 0; i < 5; i++ {
		ra[i] = &storage.RollupAddress{
			RollupId:  uint64(i + 1),
			AddressId: uint64(5 - i),
			Height:    10000,
		}
	}

	err = tx.SaveRollupAddresses(ctx, ra...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveFees() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	fees := make([]*storage.Fee, 5)
	for i := 0; i < 5; i++ {
		fees[i] = &storage.Fee{
			PayerId:  uint64(i + 1),
			TxId:     uint64(5 - i),
			ActionId: uint64(5 - i),
			Height:   10000,
		}
	}

	err = tx.SaveFees(ctx, fees...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveDeposits() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	deposits := make([]*storage.Deposit, 5)
	for i := 0; i < 5; i++ {
		deposits[i] = &storage.Deposit{
			BridgeId: uint64(i + 1),
			RollupId: uint64(5 - i),
			ActionId: uint64(5 - i),
			Height:   10000,
		}
	}

	err = tx.SaveDeposits(ctx, deposits...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveMsgAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	addresses := make([]*storage.AddressAction, 5)
	for i := 0; i < 5; i++ {
		addresses[i] = &storage.AddressAction{
			AddressId:  uint64(i + 1),
			ActionId:   uint64(5 - i),
			ActionType: types.ActionTypeValues()[i],
		}
	}

	err = tx.SaveAddressActions(ctx, addresses...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveBlockSignatures() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	bs := make([]storage.BlockSignature, 5)
	for i := 0; i < 5; i++ {
		bs[i].ValidatorId = uint64(i + 1)
		bs[i].Height = 10000
		bs[i].Time = time.Now()
	}

	err = tx.SaveBlockSignatures(ctx, bs...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveBalances() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	balances := make([]storage.Balance, 5)
	for i := 0; i < 5; i++ {
		balances[i].Id = uint64(i + 1)
		balances[i].Total = decimal.RequireFromString("1000")
		balances[i].Currency = string(currency.Nria)
	}

	err = tx.SaveBalances(ctx, balances...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveBalanceUpdates() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	balanceUpdates := make([]storage.BalanceUpdate, 5)
	for i := 0; i < 5; i++ {
		balanceUpdates[i].Height = 1000
		balanceUpdates[i].AddressId = uint64(i + 1000)
		balanceUpdates[i].Currency = string(currency.Nria)
		balanceUpdates[i].Update = decimal.RequireFromString("1000")
	}

	err = tx.SaveBalanceUpdates(ctx, balanceUpdates...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveBridges() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	bridges := make([]*storage.Bridge, 5)
	for i := 0; i < 5; i++ {
		bridges[i] = new(storage.Bridge)
		bridges[i].AddressId = uint64(i + 1000)
		bridges[i].WithdrawerId = uint64(i + 100)
		bridges[i].SudoId = uint64(i + 10)
		bridges[i].RollupId = uint64(i + 50)
	}

	count, err := tx.SaveBridges(ctx, bridges...)
	s.Require().NoError(err)
	s.Require().EqualValues(5, count)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveTransfers() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	transfers := make([]*storage.Transfer, 5)
	for i := 0; i < 5; i++ {
		transfers[i] = new(storage.Transfer)
		transfers[i].SourceId = uint64(i + 1000)
		transfers[i].DestinationId = uint64(i + 100)
		transfers[i].Amount = decimal.NewFromInt(int64(i))
		transfers[i].Asset = string(currency.Nria)
	}

	err = tx.SaveTransfers(ctx, transfers...)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestGetProposerId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	id, err := tx.GetProposerId(ctx, "astria1c220qfmjrwqlk939ca5a5z2rjxryyr9m3ah8gl")
	s.Require().NoError(err)
	s.Require().EqualValues(3, id)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestLastBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	block, err := tx.LastBlock(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(7965, block.Height)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestState() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	state, err := tx.State(ctx, testDipdupName)
	s.Require().NoError(err)
	s.Require().EqualValues(7965, state.LastHeight)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestLastNonce() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	nonce, err := tx.LastNonce(ctx, 1)
	s.Require().NoError(err)
	s.Require().EqualValues(1, nonce)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestValidators() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	validators, err := tx.Validators(ctx)
	s.Require().NoError(err)
	s.Require().Len(validators, 3)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackBlock() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackBlock(ctx, 7965)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	block, err := s.Blocks.Last(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(7964, block.Height)
}

func (s *TransactionTestSuite) TestRollbackBridge() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	count, err := tx.RollbackBridges(ctx, 7316)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackFees() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackFees(ctx, 7316)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackDeposits() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackDeposits(ctx, 7965)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackBlockStats() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	stats, err := tx.RollbackBlockStats(ctx, 7965)
	s.Require().NoError(err)
	s.Require().EqualValues(2317, stats.BlockTime)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	addresses, err := tx.RollbackAddresses(ctx, 7965)
	s.Require().NoError(err)
	s.Require().Len(addresses, 1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackTxs() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	txs, err := tx.RollbackTxs(ctx, 7965)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackActions() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	actions, err := tx.RollbackActions(ctx, 7965)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackValidators() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackValidators(ctx, 0)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackBlockSignatures() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackBlockSignatures(ctx, 7965)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackBalanceUpdates() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	updates, err := tx.RollbackBalanceUpdates(ctx, 7965)
	s.Require().NoError(err)
	s.Require().Len(updates, 2)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackAddressActions() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	actions, err := tx.RollbackAddressActions(ctx, 7965)
	s.Require().NoError(err)
	s.Require().Len(actions, 2)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackRollupActions() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	actions, err := tx.RollbackRollupActions(ctx, 7316)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackRollupAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackRollupAddresses(ctx, 7316)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackTransfers() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackTransfers(ctx, 7965)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackPrices() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackPrices(ctx, 7965)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	prices, err := s.Prices.ByHeight(ctx, 7965, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(prices, 0)
}

func (s *TransactionTestSuite) TestRollbackRollups() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	rollups, err := tx.RollbackRollups(ctx, 7316)
	s.Require().NoError(err)
	s.Require().Len(rollups, 1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestRollbackBalances() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RollbackBalances(ctx, []uint64{3, 4})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestUpdateAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.UpdateAddresses(ctx, &storage.Address{
		Id:            1,
		ActionsCount:  1,
		Nonce:         10,
		SignedTxCount: 1,
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	address, err := s.Address.GetByID(ctx, 1)
	s.Require().NoError(err)
	s.Require().EqualValues(10, address.Nonce)
	s.Require().EqualValues(2, address.ActionsCount)
	s.Require().EqualValues(3, address.SignedTxCount)
}

func (s *TransactionTestSuite) TestUpdateRollup() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.UpdateRollups(ctx, &storage.Rollup{
		Id:           1,
		ActionsCount: 1,
		Size:         100,
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	rollup, err := s.Rollup.GetByID(ctx, 1)
	s.Require().NoError(err)
	s.Require().EqualValues(212, rollup.Size)
	s.Require().EqualValues(2, rollup.ActionsCount)
}

func (s *TransactionTestSuite) TestRetentionBlockSignatures() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.RetentionBlockSignatures(ctx, 7964)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	signs, err := s.BlockSignature.List(ctx, 10, 0, sdk.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(signs, 3)
}

func (s *TransactionTestSuite) TestCreateValidator() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	pk, err := hex.DecodeString("52415f09dbee4297cc9a841c2c2312bf903fc53c48860d788ae66097355a5851")
	s.Require().NoError(err)

	val := &storage.Validator{
		PubKey: pk,
		Power:  decimal.NewFromInt(10000),
	}
	err = tx.SaveValidators(ctx, val)
	s.Require().NoError(err)
	s.Require().Greater(val.Id, uint64(0))

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	validator, err := s.Validator.GetByID(ctx, val.Id)
	s.Require().NoError(err)
	s.Require().EqualValues("10000", validator.Power.String())
}

func (s *TransactionTestSuite) TestUpdateConstants() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx.UpdateConstants(ctx, &storage.Constant{
		Module: types.ModuleNameGeneric,
		Name:   "authority_sudo_key",
		Value:  "100",
	})
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	c, err := s.Constants.Get(ctx, types.ModuleNameGeneric, "authority_sudo_key")
	s.Require().NoError(err)
	s.Require().EqualValues("100", c.Value)
}

func (s *TransactionTestSuite) TestGetBridgeIdByAddressId() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	id, err := tx.GetBridgeIdByAddressId(ctx, 1)
	s.Require().NoError(err)
	s.Require().EqualValues(1, id)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))
}

func (s *TransactionTestSuite) TestSaveMarkets() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	ts := time.Now().UTC()

	create := storage.MarketUpdate{
		Market: storage.Market{
			Pair:             "ETH_USD",
			Base:             "ETH",
			Quote:            "USD",
			Decimals:         18,
			MinProviderCount: 1,
			Enabled:          true,
			UpdatedAt:        ts,
		},
		Type: storage.MarketUpdateTypeCreate,
	}

	update := storage.MarketUpdate{
		Market: storage.Market{
			Pair:             "TIA_USD",
			Base:             "TIA",
			Quote:            "USD",
			Decimals:         6,
			MinProviderCount: 2,
			Enabled:          false,
			UpdatedAt:        ts,
		},
		Type: storage.MarketUpdateTypeUpdate,
	}

	delete := storage.MarketUpdate{
		Market: storage.Market{
			Pair:             "TIA_BTC",
			Base:             "TIA",
			Quote:            "BTC",
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          false,
			UpdatedAt:        ts,
		},
		Type: storage.MarketUpdateTypeRemove,
	}

	err = tx.SaveMarkets(ctx, create, update, delete)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	markets, err := s.Markets.List(ctx, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(markets, 5)
}

func (s *TransactionTestSuite) TestSaveMarketProviders() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	create := storage.MarketProviderUpdate{
		MarketProvider: storage.MarketProvider{
			Pair:           "ETH_USD",
			Provider:       "binance",
			OffChainTicker: "ETH/USD",
		},
		Type: storage.MarketUpdateTypeCreate,
	}

	update := storage.MarketProviderUpdate{
		MarketProvider: storage.MarketProvider{
			Pair:           "TIA_USD",
			Provider:       "coingecko",
			OffChainTicker: "TIA_USD",
		},
		Type: storage.MarketUpdateTypeUpdate,
	}

	delete := storage.MarketProviderUpdate{
		MarketProvider: storage.MarketProvider{
			Pair:           "TIA_BTC",
			Provider:       "binance",
			OffChainTicker: "TIA/BTC",
		},
		Type: storage.MarketUpdateTypeRemove,
	}

	err = tx.SaveMarketProviders(ctx, create, update, delete)
	s.Require().NoError(err)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	markets, err := s.Markets.List(ctx, 10, 0)
	s.Require().NoError(err)
	s.Require().Len(markets, 4)
}
