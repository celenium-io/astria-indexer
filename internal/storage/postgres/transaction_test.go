// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

// TransactionTestSuite -
type TransactionTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
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

	replyAddress := storage.Address{}
	addresses := make([]*storage.Address, 0, 5)
	for i := 0; i < 5; i++ {
		addresses = append(addresses, &storage.Address{
			Height: pkgTypes.Level(10000 + i),
			Hash:   testsuite.RandomAddress(),
			Id:     uint64(i),
		})

		if i == 2 {
			replyAddress.Hash = addresses[i].Hash
			replyAddress.Height = addresses[i].Height + 1
		}
	}

	count1, err := tx.SaveAddresses(ctx, addresses...)
	s.Require().NoError(err)
	s.Require().EqualValues(5, count1)

	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	s.Require().Greater(addresses[0].Id, uint64(0))
	s.Require().Greater(addresses[1].Id, uint64(0))

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	count2, err := tx2.SaveAddresses(ctx, &replyAddress)
	s.Require().NoError(err)
	s.Require().EqualValues(0, count2)

	s.Require().NoError(tx2.Flush(ctx))
	s.Require().NoError(tx2.Close(ctx))
	s.Require().Equal(replyAddress.Id, addresses[2].Id)
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

	ret, err := s.storage.Rollup.List(ctx, 10, 0, sdk.SortOrderAsc)
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

	err = tx.SaveBridges(ctx, bridges...)
	s.Require().NoError(err)

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

	block, err := s.storage.Blocks.Last(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(7964, block.Height)
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

	address, err := s.storage.Address.GetByID(ctx, 1)
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

	rollup, err := s.storage.Rollup.GetByID(ctx, 1)
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

	signs, err := s.storage.BlockSignatures.List(ctx, 10, 0, sdk.SortOrderDesc)
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

	validator, err := s.storage.Validator.GetByID(ctx, val.Id)
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

	c, err := s.storage.Constants.Get(ctx, types.ModuleNameGeneric, "authority_sudo_key")
	s.Require().NoError(err)
	s.Require().EqualValues("100", c.Value)
}
