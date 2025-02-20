// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"github.com/celenium-io/astria-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/astria-indexer/internal/storage"
)

func blockProcessor(block storage.Block) Notification[*responses.Block] {
	response := responses.NewBlock(block)
	return NewBlockNotification(response)
}

func headProcessor(state storage.State) Notification[*responses.State] {
	response := responses.NewState(state)
	return NewStateNotification(response)
}
