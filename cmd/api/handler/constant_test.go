// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// ConstantTestSuite -
type ConstantTestSuite struct {
	suite.Suite
	constants *mock.MockIConstant
	echo      *echo.Echo
	handler   *ConstantHandler
	ctrl      *gomock.Controller
}

// SetupSuite -
func (s *ConstantTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.constants = mock.NewMockIConstant(s.ctrl)
	s.handler = NewConstantHandler(s.constants)
}

// TearDownSuite -
func (s *ConstantTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteConstant_Run(t *testing.T) {
	suite.Run(t, new(ConstantTestSuite))
}

func (s *ConstantTestSuite) TestEnums() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/enums")

	s.Require().NoError(s.handler.Enums(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var enums responses.Enums
	err := json.NewDecoder(rec.Body).Decode(&enums)
	s.Require().NoError(err)
	s.Require().Len(enums.ActionType, 18)
	s.Require().Len(enums.Status, 2)
}

func (s *ConstantTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/constants")

	s.constants.EXPECT().
		All(gomock.Any()).
		Return([]storage.Constant{
			{
				Module: "module",
				Name:   "name",
				Value:  "value",
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var consts responses.Constants
	err := json.NewDecoder(rec.Body).Decode(&consts)
	s.Require().NoError(err)
	s.Require().NotEmpty(consts.Module)
	for _, module := range consts.Module {
		s.Require().NotEmpty(module)
	}
}
