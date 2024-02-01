// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	pkgTypes "github.com/aopoltorzhicky/astria/pkg/types"
	"github.com/shopspring/decimal"

	"github.com/aopoltorzhicky/astria/cmd/api/handler/responses"
	"github.com/aopoltorzhicky/astria/internal/currency"
	"github.com/aopoltorzhicky/astria/internal/storage"
	"github.com/aopoltorzhicky/astria/internal/storage/mock"
	"github.com/aopoltorzhicky/astria/internal/storage/types"
	testsuite "github.com/aopoltorzhicky/astria/internal/test_suite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testAddress = storage.Address{
		Hash:          testsuite.RandomHash(20),
		Id:            1,
		Nonce:         10,
		ActionsCount:  1,
		SignedTxCount: 1,
		Balance: &storage.Balance{
			Currency: currency.DefaultCurrency,
			Total:    decimal.RequireFromString("1000"),
			Id:       1,
		},
	}
	testAddressHash = hex.EncodeToString(testAddress.Hash)
	testBlock       = storage.Block{
		Id:           1,
		Hash:         []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		Height:       100,
		VersionBlock: 11,
		VersionApp:   1,
		Time:         testTime,
	}
	testBlockHash  = hex.EncodeToString(testBlock.Hash)
	testBlockStats = storage.BlockStats{
		TxCount:   1,
		Time:      testTime,
		Height:    100,
		BlockTime: 11043,
	}
	testValidator = storage.Validator{
		Id:         1,
		Name:       "name",
		Address:    "012345",
		PubkeyType: "tendermint/PubKeyEd25519",
		Power:      decimal.RequireFromString("1"),
	}
	testBlockWithStats = storage.Block{
		Id:           1,
		Hash:         []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		Height:       100,
		VersionBlock: 11,
		VersionApp:   1,
		Time:         testTime,
		Stats:        &testBlockStats,
		Proposer:     &testValidator,
	}
	testRollup = storage.Rollup{
		Id:           1,
		FirstHeight:  100,
		AstriaId:     testsuite.RandomHash(32),
		ActionsCount: 1,
		Size:         10,
	}
	testRollupHash = hex.EncodeToString(testRollup.AstriaId)

	testRollupAction = storage.RollupAction{
		Action: &storage.Action{
			Id:       1,
			Height:   100,
			Time:     testTime,
			Position: 1,
			Type:     types.ActionTypeSequence,
			TxId:     1,
			Data: map[string]any{
				"rollup_id": hex.EncodeToString(testRollup.AstriaId),
				"data":      testsuite.MustHexDecode("deadbeaf"),
			},
		},
		Rollup:   &testRollup,
		RollupId: testRollup.Id,
		ActionId: 1,
	}

	testTx = storage.Tx{
		Id:           1,
		Height:       100,
		Time:         testTime,
		Position:     1,
		GasWanted:    10,
		GasUsed:      8,
		ActionsCount: 1,
		Status:       types.StatusSuccess,
		Nonce:        10,
		Hash:         testsuite.RandomHash(32),
		Codespace:    "codespace",
		Signature:    testsuite.RandomHash(32),
		Signer:       &testAddress,
		SignerId:     testAddress.Id,
		ActionTypes:  types.ActionTypeSequenceBits,
		Actions: []storage.Action{
			*testRollupAction.Action,
		},
	}
	testTxHash = hex.EncodeToString(testTx.Hash)
)

// BlockTestSuite -
type BlockTestSuite struct {
	suite.Suite
	blocks     *mock.MockIBlock
	blockStats *mock.MockIBlockStats
	txs        *mock.MockITx
	actions    *mock.MockIAction
	rollups    *mock.MockIRollup
	state      *mock.MockIState
	echo       *echo.Echo
	handler    *BlockHandler
	ctrl       *gomock.Controller
}

// SetupSuite -
func (s *BlockTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.blocks = mock.NewMockIBlock(s.ctrl)
	s.blockStats = mock.NewMockIBlockStats(s.ctrl)
	s.txs = mock.NewMockITx(s.ctrl)
	s.rollups = mock.NewMockIRollup(s.ctrl)
	s.actions = mock.NewMockIAction(s.ctrl)
	s.state = mock.NewMockIState(s.ctrl)
	s.handler = NewBlockHandler(s.blocks, s.blockStats, s.txs, s.actions, s.rollups, s.state, testIndexerName)
}

// TearDownSuite -
func (s *BlockTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteBlock_Run(t *testing.T) {
	suite.Run(t, new(BlockTestSuite))
}

func (s *BlockTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100), false).
		Return(testBlock, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var block responses.Block
	err := json.NewDecoder(rec.Body).Decode(&block)
	s.Require().NoError(err)
	s.Require().EqualValues(1, block.Id)
	s.Require().EqualValues(100, block.Height)
	s.Require().Equal("1", block.VersionApp)
	s.Require().Equal("11", block.VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", block.Hash.String())
	s.Require().Equal(testTime, block.Time)
	s.Require().Nil(block.Stats)
}

func (s *BlockTestSuite) TestGetNoContent() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100), false).
		Return(storage.Block{}, sql.ErrNoRows)

	s.blocks.EXPECT().
		IsNoRows(gomock.Any()).
		Return(true)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusNoContent, rec.Code)
}

func (s *BlockTestSuite) TestGetWithoutStats() {
	q := make(url.Values)
	q.Set("stats", "false")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100), false).
		Return(testBlock, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var block responses.Block
	err := json.NewDecoder(rec.Body).Decode(&block)
	s.Require().NoError(err)
	s.Require().EqualValues(1, block.Id)
	s.Require().EqualValues(100, block.Height)
	s.Require().Equal("1", block.VersionApp)
	s.Require().Equal("11", block.VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", block.Hash.String())
	s.Require().Equal(testTime, block.Time)
	s.Require().Nil(block.Stats)
}

func (s *BlockTestSuite) TestGetWithStats() {
	q := make(url.Values)
	q.Set("stats", "true")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100), true).
		Return(testBlockWithStats, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var block responses.Block
	err := json.NewDecoder(rec.Body).Decode(&block)
	s.Require().NoError(err)
	s.Require().EqualValues(1, block.Id)
	s.Require().EqualValues(100, block.Height)
	s.Require().Equal("1", block.VersionApp)
	s.Require().Equal("11", block.VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", block.Hash.String())
	s.Require().Equal(testTime, block.Time)
	s.Require().NotNil(block.Stats)
	s.Require().EqualValues(1, block.Stats.TxCount)
}

func (s *BlockTestSuite) TestGetInvalidBlockHeight() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("invalid")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "parsing")
}

func (s *BlockTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block")

	s.blocks.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.Block{
			&testBlock,
		}, nil).
		MaxTimes(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blocks []responses.Block
	err := json.NewDecoder(rec.Body).Decode(&blocks)
	s.Require().NoError(err)
	s.Require().Len(blocks, 1)
	s.Require().EqualValues(1, blocks[0].Id)
	s.Require().EqualValues(100, blocks[0].Height)
	s.Require().Equal("1", blocks[0].VersionApp)
	s.Require().Equal("11", blocks[0].VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", blocks[0].Hash.String())
	s.Require().Equal(testTime, blocks[0].Time)
}

func (s *BlockTestSuite) TestListWithStats() {
	q := make(url.Values)
	q.Set("stats", "true")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block")

	s.blocks.EXPECT().
		ListWithStats(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.Block{
			&testBlockWithStats,
		}, nil).
		MaxTimes(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blocks []responses.Block
	err := json.NewDecoder(rec.Body).Decode(&blocks)
	s.Require().NoError(err)
	s.Require().Len(blocks, 1)
	s.Require().EqualValues(1, blocks[0].Id)
	s.Require().EqualValues(100, blocks[0].Height)
	s.Require().Equal("1", blocks[0].VersionApp)
	s.Require().Equal("11", blocks[0].VersionBlock)
	s.Require().Equal("000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F", blocks[0].Hash.String())
	s.Require().Equal(testTime, blocks[0].Time)
}

func (s *BlockTestSuite) TestGetActions() {
	q := make(url.Values)
	q.Set("limit", "2")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/actions")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.actions.EXPECT().
		ByBlock(gomock.Any(), pkgTypes.Level(100), 2, 0).
		Return([]storage.ActionWithTx{
			{
				Action: storage.Action{
					Id:       1,
					Height:   100,
					Time:     testTime,
					Position: 2,
					Type:     types.ActionTypeSequence,
					TxId:     10,
					Data: map[string]any{
						"test": "value",
					},
				},
				Tx: &testTx,
			},
		}, nil)

	s.Require().NoError(s.handler.GetActions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var actions []responses.Action
	err := json.NewDecoder(rec.Body).Decode(&actions)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)
	s.Require().EqualValues(1, actions[0].Id)
	s.Require().EqualValues(100, actions[0].Height)
	s.Require().EqualValues(2, actions[0].Position)
	s.Require().Equal(testTime, actions[0].Time)
	s.Require().Equal(types.ActionTypeSequence, actions[0].Type)
	s.Require().Equal(hex.EncodeToString(testTx.Hash), actions[0].TxHash)
}

func (s *BlockTestSuite) TestGetStats() {
	req := httptest.NewRequest(http.MethodGet, "/?", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/stats")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blockStats.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100)).
		Return(testBlockStats, nil)

	s.Require().NoError(s.handler.GetStats(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var stats responses.BlockStats
	err := json.NewDecoder(rec.Body).Decode(&stats)
	s.Require().NoError(err)
	s.Require().EqualValues(1, stats.TxCount)
	s.Require().EqualValues(11043, stats.BlockTime)
}

func (s *BlockTestSuite) TestGetRollupActions() {
	req := httptest.NewRequest(http.MethodGet, "/?", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/rollup_actions")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.rollups.EXPECT().
		ActionsByHeight(gomock.Any(), pkgTypes.Level(100), int(10), int(0)).
		Return([]storage.RollupAction{testRollupAction}, nil)

	s.Require().NoError(s.handler.GetRollupActions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var actions []responses.RollupAction
	err := json.NewDecoder(rec.Body).Decode(&actions)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)

	action := actions[0]
	s.Require().EqualValues(1, action.Id)
	s.Require().EqualValues(100, action.Height)
	s.Require().EqualValues(1, action.Position)
	s.Require().Equal(testTime, action.Time)
	s.Require().EqualValues(string(types.ActionTypeSequence), action.Type)
}

func (s *BlockTestSuite) TestGetRollupActionsCount() {
	req := httptest.NewRequest(http.MethodGet, "/?", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/rollup_actions/count")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.rollups.EXPECT().
		CountActionsByHeight(gomock.Any(), pkgTypes.Level(100)).
		Return(12, nil)

	s.Require().NoError(s.handler.GetRollupsActionsCount(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count int
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(count, 12)
}

func (s *BlockTestSuite) TestCount() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/count")

	s.state.EXPECT().
		ByName(gomock.Any(), testIndexerName).
		Return(testState, nil)

	s.Require().NoError(s.handler.Count(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var count uint64
	err := json.NewDecoder(rec.Body).Decode(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(101, count)
}

func (s *BlockTestSuite) TestGetTransactions() {
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("offset", "0")

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/txs")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.txs.EXPECT().
		ByHeight(gomock.Any(), pkgTypes.Level(100), 10, 0).
		Return([]storage.Tx{
			testTx,
		}, nil)

	s.Require().NoError(s.handler.GetTransactions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var txs []responses.Tx
	err := json.NewDecoder(rec.Body).Decode(&txs)
	s.Require().NoError(err)
	s.Require().Len(txs, 1)

	tx := txs[0]
	s.Require().EqualValues(1, tx.Id)
	s.Require().EqualValues(100, tx.Height)
	s.Require().Equal(testTime, tx.Time)
	s.Require().Equal(testTxHash, tx.Hash)
	s.Require().EqualValues(1, tx.Position)
	s.Require().EqualValues(10, tx.GasWanted)
	s.Require().EqualValues(8, tx.GasUsed)
	s.Require().EqualValues(1, tx.ActionsCount)
	s.Require().EqualValues(10, tx.Nonce)
	s.Require().EqualValues(hex.EncodeToString(testAddress.Hash), tx.Signer)
	s.Require().Equal("codespace", tx.Codespace)
	s.Require().Equal(types.StatusSuccess, tx.Status)
}
