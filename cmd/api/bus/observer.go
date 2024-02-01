// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package bus

import (
	"github.com/aopoltorzhicky/astria/internal/storage"
	"github.com/dipdup-io/workerpool"
)

type Observer struct {
	blocks chan *storage.Block
	head   chan *storage.State

	listenHead   bool
	listenBlocks bool

	g workerpool.Group
}

func NewObserver(channels ...string) *Observer {
	if len(channels) == 0 {
		return nil
	}

	observer := &Observer{
		blocks: make(chan *storage.Block, 1024),
		head:   make(chan *storage.State, 1024),
		g:      workerpool.NewGroup(),
	}

	for i := range channels {
		switch channels[i] {
		case storage.ChannelBlock:
			observer.listenBlocks = true
		case storage.ChannelHead:
			observer.listenHead = true
		}
	}

	return observer
}

func (observer Observer) Close() error {
	observer.g.Wait()
	close(observer.blocks)
	close(observer.head)
	return nil
}

func (observer Observer) notifyBlocks(block *storage.Block) {
	if observer.listenBlocks {
		observer.blocks <- block
	}
}
func (observer Observer) notifyState(state *storage.State) {
	if observer.listenHead {
		observer.head <- state
	}
}

func (observer Observer) Blocks() <-chan *storage.Block {
	return observer.blocks
}

func (observer Observer) Head() <-chan *storage.State {
	return observer.head
}
