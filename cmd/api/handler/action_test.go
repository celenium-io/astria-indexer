// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// ActionTestSuite -
type ActionTestSuite struct {
	suite.Suite
	actions *mock.MockIAction
	echo    *echo.Echo
	handler *ActionHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *ActionTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.actions = mock.NewMockIAction(s.ctrl)
	s.handler = NewActionHandler(s.actions)
}

// TearDownSuite -
func (s *ActionTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteAction_Run(t *testing.T) {
	suite.Run(t, new(ActionTestSuite))
}

func (s *ActionTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/action/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	s.actions.EXPECT().
		ById(gomock.Any(), uint64(1)).
		Return(storage.ActionWithTx{
			Tx:     &testTx,
			Action: *testRollupAction.Action,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var action responses.Action
	err := json.NewDecoder(rec.Body).Decode(&action)
	s.Require().NoError(err)

	s.Require().EqualValues(1, action.Id)
	s.Require().EqualValues(100, action.Height)
	s.Require().EqualValues(1, action.Position)
	s.Require().Equal(testTime, action.Time)
	s.Require().EqualValues(string(types.ActionTypeRollupDataSubmission), action.Type)
	s.Require().NotEmpty(action.TxHash)
}
