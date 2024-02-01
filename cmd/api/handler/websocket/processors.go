// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"github.com/aopoltorzhicky/astria/cmd/api/handler/responses"
	"github.com/aopoltorzhicky/astria/internal/storage"
)

func headProcessor(state storage.State) *responses.State {
	response := responses.NewState(state)
	return &response
}

func blockProcessor(block storage.Block) *responses.Block {
	response := responses.NewBlock(block)
	return &response
}
