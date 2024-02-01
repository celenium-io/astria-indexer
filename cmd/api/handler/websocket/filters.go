// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"github.com/aopoltorzhicky/astria/cmd/api/handler/responses"
)

type Filterable[M any] interface {
	Filter(c client, msg M) bool
}

type HeadFilter struct{}

func (hf HeadFilter) Filter(c client, msg *responses.State) bool {
	if msg == nil {
		return false
	}
	fltrs := c.Filters()
	if fltrs == nil {
		return false
	}
	return fltrs.head
}

type BlockFilter struct{}

func (hf BlockFilter) Filter(c client, msg *responses.Block) bool {
	if msg == nil {
		return false
	}
	fltrs := c.Filters()
	if fltrs == nil {
		return false
	}
	return fltrs.blocks
}

type Filters struct {
	head   bool
	blocks bool
}
