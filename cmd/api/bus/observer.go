// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package bus

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-io/workerpool"
)

type Observer struct {
	blocks    chan *storage.Block
	head      chan *storage.State
	constants chan *storage.Constant

	listenHead      bool
	listenBlocks    bool
	listenConstants bool

	g workerpool.Group
}

func NewObserver(channels ...string) *Observer {
	if len(channels) == 0 {
		return nil
	}

	observer := &Observer{
		blocks:    make(chan *storage.Block, 1024),
		head:      make(chan *storage.State, 1024),
		constants: make(chan *storage.Constant, 1024),
		g:         workerpool.NewGroup(),
	}

	for i := range channels {
		switch channels[i] {
		case storage.ChannelBlock:
			observer.listenBlocks = true
		case storage.ChannelHead:
			observer.listenHead = true
		case storage.ChannelConstant:
			observer.listenConstants = true
		}
	}

	return observer
}

func (observer Observer) Close() error {
	observer.g.Wait()
	close(observer.blocks)
	close(observer.head)
	close(observer.constants)
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
func (observer Observer) notifyConstants(constant *storage.Constant) {
	if observer.listenConstants {
		observer.constants <- constant
	}
}

func (observer Observer) Blocks() <-chan *storage.Block {
	return observer.blocks
}

func (observer Observer) Head() <-chan *storage.State {
	return observer.head
}

func (observer Observer) Constants() <-chan *storage.Constant {
	return observer.constants
}
