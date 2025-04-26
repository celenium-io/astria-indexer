// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type LeaderboardFilters struct {
	SortField string
	Sort      storage.SortOrder
	Limit     int
	Offset    int
	Category  []types.AppCategory
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IApp interface {
	storage.Table[*App]

	Leaderboard(ctx context.Context, fltrs LeaderboardFilters) ([]AppWithStats, error)
	BySlug(ctx context.Context, slug string) (AppWithStats, error)
	ByRollupId(ctx context.Context, rollupId uint64) (AppWithStats, error)
}

type App struct {
	bun.BaseModel `bun:"app" comment:"Table with applications."`

	Id             uint64            `bun:"id,pk,notnull,autoincrement"            comment:"Unique internal identity"`
	Group          string            `bun:"group"                                  comment:"Application group"`
	Name           string            `bun:"name"                                   comment:"Application name"`
	Slug           string            `bun:"slug,unique:app_slug"                   comment:"Application slug"`
	Github         string            `bun:"github"                                 comment:"Application github link"`
	Twitter        string            `bun:"twitter"                                comment:"Application twitter account link"`
	Website        string            `bun:"website"                                comment:"Application website link"`
	Logo           string            `bun:"logo"                                   comment:"Application logo link"`
	Description    string            `bun:"description"                            comment:"Application description"`
	Explorer       string            `bun:"explorer"                               comment:"Application explorer link"`
	L2Beat         string            `bun:"l2beat"                                 comment:"Link to L2Beat"`
	Links          []string          `bun:"links,array"                            comment:"Additional links"`
	Stack          string            `bun:"stack"                                  comment:"Using stack"`
	VM             string            `bun:"vm"                                     comment:"Virtual machine"`
	Provider       string            `bun:"provider"                               comment:"RaaS"`
	Type           types.AppType     `bun:"type,type:app_type"                     comment:"Type of application: settled or sovereign"`
	Category       types.AppCategory `bun:"category,type:app_category"             comment:"Category of applications"`
	RollupId       uint64            `bun:"rollup_id,notnull,unique:app_rollup_id" comment:"Rollup internal identity"`
	NativeBridgeId uint64            `bun:"native_bridge_id"                       comment:"Native bridge internal id"`

	Bridge *Address `bun:"rel:belongs-to"`
	Rollup *Rollup  `bun:"rel:belongs-to"`
}

func (App) TableName() string {
	return "app"
}

func (app App) IsEmpty() bool {
	return app.Group == "" &&
		app.Name == "" &&
		app.Slug == "" &&
		app.Github == "" &&
		app.Twitter == "" &&
		app.Website == "" &&
		app.Logo == "" &&
		app.Description == "" &&
		app.Explorer == "" &&
		app.L2Beat == "" &&
		app.Links == nil &&
		app.Stack == "" &&
		app.VM == "" &&
		app.Provider == "" &&
		app.Type == "" &&
		app.Category == "" &&
		app.RollupId == 0 &&
		app.NativeBridgeId == 0
}

type AppWithStats struct {
	App
	AppStats
}

type AppStats struct {
	Size            int64     `bun:"size"`
	MinSize         int64     `bun:"min_size"`
	MaxSize         int64     `bun:"max_size"`
	AvgSize         float64   `bun:"avg_size"`
	ActionsCount    int64     `bun:"actions_count"`
	LastActionTime  time.Time `bun:"last_time"`
	FirstActionTime time.Time `bun:"first_time"`
}
