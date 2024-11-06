// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

type AppWithStats struct {
	Id          uint64 `example:"321"                                       format:"integer" json:"id"                    swaggertype:"integer"`
	Name        string `example:"Rollup name"                               format:"string"  json:"name"                  swaggertype:"string"`
	Description string `example:"Long rollup description"                   format:"string"  json:"description,omitempty" swaggertype:"string"`
	Website     string `example:"https://website.com"                       format:"string"  json:"website,omitempty"     swaggertype:"string"`
	Twitter     string `example:"https://x.com/account"                     format:"string"  json:"twitter,omitempty"     swaggertype:"string"`
	Github      string `example:"https://github.com/account"                format:"string"  json:"github,omitempty"      swaggertype:"string"`
	Logo        string `example:"https://some_link.com/image.png"           format:"string"  json:"logo,omitempty"        swaggertype:"string"`
	Slug        string `example:"rollup_slug"                               format:"string"  json:"slug"                  swaggertype:"string"`
	L2Beat      string `example:"https://l2beat.com/scaling/projects/karak" format:"string"  json:"l2_beat,omitempty"     swaggertype:"string"`
	Explorer    string `example:"https://explorer.karak.network/"           format:"string"  json:"explorer,omitempty"    swaggertype:"string"`
	Stack       string `example:"op_stack"                                  format:"string"  json:"stack,omitempty"       swaggertype:"string"`
	Type        string `example:"settled"                                   format:"string"  json:"type,omitempty"        swaggertype:"string"`
	Category    string `example:"nft"                                       format:"string"  json:"category,omitempty"    swaggertype:"string"`
	VM          string `example:"evm"                                       format:"string"  json:"vm,omitempty"          swaggertype:"string"`
	Provider    string `example:"name"                                      format:"string"  json:"provider,omitempty"    swaggertype:"string"`

	ActionsCount    int64     `example:"2"                         format:"integer"   json:"actions_count"      swaggertype:"integer"`
	Size            int64     `example:"1000"                      format:"integer"   json:"size"               swaggertype:"integer"`
	LastAction      time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"last_message_time"  swaggertype:"string"`
	FirstAction     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"first_message_time" swaggertype:"string"`
	SizePct         float64   `example:"0.9876"                    format:"float"     json:"size_pct"           swaggertype:"number"`
	ActionsCountPct float64   `example:"0.9876"                    format:"float"     json:"actions_count_pct"  swaggertype:"number"`

	Links []string `json:"links,omitempty"`
}

func NewAppWithStats(r storage.AppWithStats) AppWithStats {
	return AppWithStats{
		Id:              r.Id,
		Name:            r.Name,
		Description:     r.Description,
		Github:          r.Github,
		Twitter:         r.Twitter,
		Website:         r.Website,
		Logo:            r.Logo,
		L2Beat:          r.L2Beat,
		Explorer:        r.Explorer,
		Links:           r.Links,
		Stack:           r.Stack,
		Slug:            r.Slug,
		ActionsCount:    r.ActionsCount,
		Size:            r.Size,
		SizePct:         r.SizePct,
		ActionsCountPct: r.ActionsCountPct,
		LastAction:      r.LastActionTime,
		FirstAction:     r.FirstActionTime,
		Category:        r.Category.String(),
		Type:            r.Type.String(),
		Provider:        r.Provider,
		VM:              r.VM,
	}
}

// type App struct {
// 	Id          uint64 `example:"321"                             format:"integer" json:"id"                    swaggertype:"integer"`
// 	Name        string `example:"Rollup name"                     format:"string"  json:"name"                  swaggertype:"string"`
// 	Description string `example:"Long rollup description"         format:"string"  json:"description,omitempty" swaggertype:"string"`
// 	Website     string `example:"https://website.com"             format:"string"  json:"website,omitempty"     swaggertype:"string"`
// 	Twitter     string `example:"https://x.com/account"           format:"string"  json:"twitter,omitempty"     swaggertype:"string"`
// 	Github      string `example:"https://github.com/account"      format:"string"  json:"github,omitempty"      swaggertype:"string"`
// 	Logo        string `example:"https://some_link.com/image.png" format:"string"  json:"logo,omitempty"        swaggertype:"string"`
// 	Slug        string `example:"rollup_slug"                     format:"string"  json:"slug"                  swaggertype:"string"`
// 	L2Beat      string `example:"https://github.com/account"      format:"string"  json:"l2_beat,omitempty"     swaggertype:"string"`
// 	Explorer    string `example:"https://explorer.karak.network/" format:"string"  json:"explorer,omitempty"    swaggertype:"string"`
// 	Stack       string `example:"op_stack"                        format:"string"  json:"stack,omitempty"       swaggertype:"string"`
// 	Type        string `example:"settled"                         format:"string"  json:"type,omitempty"        swaggertype:"string"`
// 	Category    string `example:"nft"                             format:"string"  json:"category,omitempty"    swaggertype:"string"`
// 	Provider    string `example:"name"                            format:"string"  json:"provider,omitempty"    swaggertype:"string"`
// 	VM          string `example:"evm"                             format:"string"  json:"vm,omitempty"          swaggertype:"string"`

// 	Links []string `json:"links,omitempty"`
// }

// func NewApp(r *storage.App) App {
// 	return App{
// 		Id:          r.Id,
// 		Name:        r.Name,
// 		Description: r.Description,
// 		Github:      r.Github,
// 		Twitter:     r.Twitter,
// 		Website:     r.Website,
// 		Logo:        r.Logo,
// 		Slug:        r.Slug,
// 		L2Beat:      r.L2Beat,
// 		Stack:       r.Stack,
// 		Explorer:    r.Explorer,
// 		Links:       r.Links,
// 		Category:    r.Category.String(),
// 		Type:        r.Type.String(),
// 		Provider:    r.Provider,
// 		VM:          r.VM,
// 	}
// }
