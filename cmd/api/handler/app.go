// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
)

type AppHandler struct {
	apps storage.IApp
}

func NewAppHandler(
	apps storage.IApp,
) *AppHandler {
	return &AppHandler{
		apps: apps,
	}
}

var _ Handler = (*AppHandler)(nil)

func (handler *AppHandler) InitRoutes(srvr *echo.Group) {
	apps := srvr.Group("/app")
	{
		apps.GET("", handler.Leaderboard)
		app := apps.Group("/:slug")
		{
			app.GET("", handler.Get)
		}
	}
}

type leaderboardRequest struct {
	Limit    int         `query:"limit"    validate:"omitempty,min=1,max=100"`
	Offset   int         `query:"offset"   validate:"omitempty,min=0"`
	Sort     string      `query:"sort"     validate:"omitempty,oneof=asc desc"`
	SortBy   string      `query:"sort_by"  validate:"omitempty,oneof=time actions_count size"`
	Category StringArray `query:"category" validate:"omitempty,dive,app_category"`
}

func (p *leaderboardRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = desc
	}
	if p.SortBy == "" {
		p.SortBy = "size"
	}
}

// Leaderboard godoc
//
//		@Summary		List applications info
//		@Description	List applications info
//		@Tags			applications
//		@ID				list-applications
//		@Param			limit	 query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//		@Param			offset	 query	integer	false	"Offset"						mininum(1)
//		@Param			sort	 query	string	false	"Sort order. Default: desc"		Enums(asc, desc)
//		@Param			sort_by	 query	string	false	"Sort field. Default: size"		Enums(time, actions_count, size)
//	    @Param          category query  string  false   "Comma-separated application category list"
//		@Produce		json
//		@Success		200	{array}		responses.AppWithStats
//		@Failure		400	{object}	Error
//		@Failure		500	{object}	Error
//		@Router			/app [get]
func (handler AppHandler) Leaderboard(c echo.Context) error {
	req, err := bindAndValidate[leaderboardRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	categories := make([]types.AppCategory, len(req.Category))
	for i := range categories {
		categories[i] = types.AppCategory(req.Category[i])
	}

	apps, err := handler.apps.Leaderboard(c.Request().Context(), storage.LeaderboardFilters{
		SortField: req.SortBy,
		Sort:      pgSort(req.Sort),
		Limit:     req.Limit,
		Offset:    req.Offset,
		Category:  categories,
	})
	if err != nil {
		return handleError(c, err, handler.apps)
	}
	response := make([]responses.AppWithStats, len(apps))
	for i := range apps {
		response[i] = responses.NewAppWithStats(apps[i])
	}
	return returnArray(c, response)
}

type getAppRequest struct {
	Slug string `param:"slug" validate:"required"`
}

// Get godoc
//
//	@Summary		Get application info
//	@Description	Get application info
//	@Tags			applications
//	@ID				get-application
//	@Param			slug	path	string	true	"Slug"
//	@Produce		json
//	@Success		200	{object}	responses.AppWithStats
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/app/{slug} [get]
func (handler AppHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getAppRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	rollup, err := handler.apps.BySlug(c.Request().Context(), req.Slug)
	if err != nil {
		return handleError(c, err, handler.apps)
	}

	return c.JSON(http.StatusOK, responses.NewAppWithStats(rollup))
}
