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
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// AssetTestSuite -
type AssetTestSuite struct {
	suite.Suite
	asset   *mock.MockIAsset
	block   *mock.MockIBlock
	echo    *echo.Echo
	handler *AssetHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *AssetTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.asset = mock.NewMockIAsset(s.ctrl)
	s.block = mock.NewMockIBlock(s.ctrl)
	s.handler = NewAssetHandler(s.asset, s.block)
}

// TearDownSuite -
func (s *AssetTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteAsset_Run(t *testing.T) {
	suite.Run(t, new(AssetTestSuite))
}

func (s *AssetTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/asset")

	s.asset.EXPECT().
		List(gomock.Any(), 10, 0, "", sdk.SortOrderDesc).
		Return([]storage.Asset{
			{
				Asset:         "asset",
				Transferred:   decimal.NewFromInt(10),
				Fee:           decimal.NewFromInt(20),
				TransferCount: 2,
				FeeCount:      3,
			},
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var assets []responses.Asset
	err := json.NewDecoder(rec.Body).Decode(&assets)
	s.Require().NoError(err)
	s.Require().Len(assets, 1)

	s.Require().EqualValues("asset", assets[0].Asset)
	s.Require().EqualValues("10", assets[0].Transferred)
	s.Require().EqualValues("20", assets[0].Fee)
	s.Require().EqualValues(2, assets[0].TransferCount)
	s.Require().EqualValues(3, assets[0].FeeCount)
}
