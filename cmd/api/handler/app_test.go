// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/mock"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testApplication = storage.App{
		Id:          1,
		Name:        "test app",
		Description: "loooooooooooooooooong description",
		Website:     "https://website.com",
		Github:      "https://githib.com",
		Twitter:     "https://x.com",
		Logo:        "image.png",
		Slug:        "test-app",
		Category:    types.AppCategoryNft,
	}
	testRollupWithStats = storage.AppWithStats{
		App: testApplication,
		AppStats: storage.AppStats{
			ActionsCount:    100,
			Size:            1000,
			LastActionTime:  testTime,
			FirstActionTime: testTime,
			ActionsCountPct: 0.1,
			SizePct:         0.3,
		},
	}
)

// AppTestSuite -
type AppTestSuite struct {
	suite.Suite
	apps    *mock.MockIApp
	echo    *echo.Echo
	handler *AppHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *AppTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.apps = mock.NewMockIApp(s.ctrl)
	s.handler = NewAppHandler(s.apps)
}

// TearDownSuite -
func (s *AppTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteApp_Run(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func (s *AppTestSuite) TestLeaderboard() {
	for _, sort := range []string{
		"actions_count",
		"time",
		"size",
	} {
		q := make(url.Values)
		q.Add("sort_by", sort)
		q.Add("category", "nft,gaming")

		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)
		c.SetPath("/app")

		s.apps.EXPECT().
			Leaderboard(gomock.Any(), storage.LeaderboardFilters{
				SortField: sort,
				Sort:      sdk.SortOrderDesc,
				Limit:     10,
				Offset:    0,
				Category: []types.AppCategory{
					types.AppCategoryNft,
					types.AppCategoryGaming,
				},
			}).
			Return([]storage.AppWithStats{testRollupWithStats}, nil).
			Times(1)

		s.Require().NoError(s.handler.Leaderboard(c))
		s.Require().Equal(http.StatusOK, rec.Code)

		var rollups []responses.AppWithStats
		err := json.NewDecoder(rec.Body).Decode(&rollups)
		s.Require().NoError(err)
		s.Require().Len(rollups, 1)

		rollup := rollups[0]
		s.Require().EqualValues(1, rollup.Id)
		s.Require().EqualValues("test app", rollup.Name)
		s.Require().EqualValues("image.png", rollup.Logo)
		s.Require().EqualValues("test-app", rollup.Slug)
		s.Require().EqualValues(100, rollup.ActionsCount)
		s.Require().EqualValues(1000, rollup.Size)
		s.Require().EqualValues(testTime, rollup.LastAction)
		s.Require().EqualValues(testTime, rollup.FirstAction)
		s.Require().EqualValues(0.1, rollup.ActionsCountPct)
		s.Require().EqualValues(0.3, rollup.SizePct)
	}
}

func (s *AppTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/app/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("test-app")

	s.apps.EXPECT().
		BySlug(gomock.Any(), "test-app").
		Return(testRollupWithStats, nil).
		Times(1)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var rollup responses.AppWithStats
	err := json.NewDecoder(rec.Body).Decode(&rollup)
	s.Require().NoError(err)
	s.Require().EqualValues(1, rollup.Id)
	s.Require().EqualValues("test app", rollup.Name)
	s.Require().EqualValues("image.png", rollup.Logo)
	s.Require().EqualValues("test-app", rollup.Slug)
	s.Require().EqualValues(100, rollup.ActionsCount)
	s.Require().EqualValues(1000, rollup.Size)
	s.Require().EqualValues(testTime, rollup.LastAction)
	s.Require().EqualValues(testTime, rollup.FirstAction)
	s.Require().EqualValues(0.1, rollup.ActionsCountPct)
	s.Require().EqualValues(0.3, rollup.SizePct)
}

func (s *AppTestSuite) TestActions() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/app/:slug/actions")
	c.SetParamNames("slug")
	c.SetParamValues("test-app")

	s.apps.EXPECT().
		Actions(gomock.Any(), "test-app", 10, 0, sdk.SortOrderAsc).
		Return([]storage.RollupAction{
			testRollupAction,
		}, nil).
		Times(1)

	s.Require().NoError(s.handler.Actions(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var actions []responses.Action
	err := json.NewDecoder(rec.Body).Decode(&actions)
	s.Require().NoError(err)
	s.Require().Len(actions, 1)
	s.Require().EqualValues(1, actions[0].Id)
	s.Require().EqualValues(100, actions[0].Height)
	s.Require().EqualValues(1, actions[0].Position)
	s.Require().Equal(testTime, actions[0].Time)
	s.Require().EqualValues(string(types.ActionTypeRollupDataSubmission), actions[0].Type)
}

func (s *AppTestSuite) TestSeries() {
	for _, name := range []string{
		"size", "actions_count", "size_per_action",
	} {
		for _, tf := range []storage.Timeframe{
			"hour", "day", "month",
		} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := s.echo.NewContext(req, rec)
			c.SetPath("/app/:slug/stats/:name/:timeframe")
			c.SetParamNames("slug", "name", "timeframe")
			c.SetParamValues("test-app", name, string(tf))

			s.apps.EXPECT().
				Series(gomock.Any(), "test-app", tf, name, storage.NewSeriesRequest(0, 0)).
				Return([]storage.SeriesItem{
					{
						Value: "10000",
						Time:  testTime,
					},
				}, nil).
				Times(1)

			s.Require().NoError(s.handler.Series(c))
			s.Require().Equal(http.StatusOK, rec.Code, name, tf)

			var items []responses.SeriesItem
			err := json.NewDecoder(rec.Body).Decode(&items)
			s.Require().NoError(err, name, tf)
			s.Require().Len(items, 1, name, tf)
			s.Require().EqualValues("10000", items[0].Value, name, tf)
			s.Require().EqualValues(testTime.UTC(), items[0].Time.UTC(), name, tf)
		}
	}
}
