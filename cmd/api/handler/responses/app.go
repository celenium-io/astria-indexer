// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/base64"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

type AppWithStats struct {
	Id           uint64 `example:"321"                                           format:"integer" json:"id"                    swaggertype:"integer"`
	Name         string `example:"Rollup name"                                   format:"string"  json:"name"                  swaggertype:"string"`
	Description  string `example:"Long rollup description"                       format:"string"  json:"description,omitempty" swaggertype:"string"`
	Website      string `example:"https://website.com"                           format:"string"  json:"website,omitempty"     swaggertype:"string"`
	Twitter      string `example:"https://x.com/account"                         format:"string"  json:"twitter,omitempty"     swaggertype:"string"`
	Github       string `example:"https://github.com/account"                    format:"string"  json:"github,omitempty"      swaggertype:"string"`
	Logo         string `example:"https://some_link.com/image.png"               format:"string"  json:"logo,omitempty"        swaggertype:"string"`
	Slug         string `example:"rollup_slug"                                   format:"string"  json:"slug"                  swaggertype:"string"`
	L2Beat       string `example:"https://l2beat.com/scaling/projects/karak"     format:"string"  json:"l2_beat,omitempty"     swaggertype:"string"`
	Explorer     string `example:"https://explorer.karak.network/"               format:"string"  json:"explorer,omitempty"    swaggertype:"string"`
	Stack        string `example:"op_stack"                                      format:"string"  json:"stack,omitempty"       swaggertype:"string"`
	Type         string `example:"settled"                                       format:"string"  json:"type,omitempty"        swaggertype:"string"`
	Category     string `example:"nft"                                           format:"string"  json:"category,omitempty"    swaggertype:"string"`
	VM           string `example:"evm"                                           format:"string"  json:"vm,omitempty"          swaggertype:"string"`
	Provider     string `example:"name"                                          format:"string"  json:"provider,omitempty"    swaggertype:"string"`
	NativeBridge string `example:"astria1phym4uktjn6gjle226009ge7u82w0dgtszs8x2" format:"string"  json:"native_bridge"         swaggertype:"string"`
	Rollup       string `example:"O0Ia+lPYYMf3iFfxBaWXCSdlhphc6d4ZoBXINov6Tjc="  format:"string"  json:"rollup"                swaggertype:"string"`

	ActionsCount int64     `example:"2"                         format:"integer"   json:"actions_count"      swaggertype:"integer"`
	Size         int64     `example:"1000"                      format:"integer"   json:"size"               swaggertype:"integer"`
	MinSize      int64     `example:"1000"                      format:"integer"   json:"min_size"           swaggertype:"integer"`
	MaxSize      int64     `example:"1000"                      format:"integer"   json:"max_size"           swaggertype:"integer"`
	AvgSize      int64     `example:"1000"                      format:"integer"   json:"avg_size"           swaggertype:"integer"`
	LastAction   time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"last_message_time"  swaggertype:"string"`
	FirstAction  time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"first_message_time" swaggertype:"string"`

	Links []string `json:"links,omitempty"`
}

func NewAppWithStats(r storage.AppWithStats) AppWithStats {
	app := AppWithStats{
		Id:           r.Id,
		Name:         r.Name,
		Description:  r.Description,
		Github:       r.Github,
		Twitter:      r.Twitter,
		Website:      r.Website,
		Logo:         r.Logo,
		L2Beat:       r.L2Beat,
		Explorer:     r.Explorer,
		Links:        r.Links,
		Stack:        r.Stack,
		Slug:         r.Slug,
		ActionsCount: r.ActionsCount,
		Size:         r.Size,
		MinSize:      r.MinSize,
		MaxSize:      r.MaxSize,
		AvgSize:      int64(r.AvgSize),
		LastAction:   r.LastActionTime,
		FirstAction:  r.FirstActionTime,
		Category:     r.Category.String(),
		Type:         r.Type.String(),
		Provider:     r.Provider,
		VM:           r.VM,
	}

	if r.Rollup != nil {
		app.Rollup = base64.StdEncoding.EncodeToString(r.Rollup.AstriaId)
	}
	if r.Bridge != nil {
		app.NativeBridge = r.Bridge.Hash
	}

	return app
}
