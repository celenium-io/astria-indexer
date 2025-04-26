// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"context"
	"sync"

	"github.com/celenium-io/astria-indexer/cmd/api/bus"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
)

type ConstantsCache struct {
	data     map[string]map[string]string
	observer *bus.Observer

	wg *sync.WaitGroup
	mx *sync.RWMutex
}

func NewConstantsCache(observer *bus.Observer) *ConstantsCache {
	return &ConstantsCache{
		data:     make(map[string]map[string]string),
		observer: observer,
		wg:       new(sync.WaitGroup),
		mx:       new(sync.RWMutex),
	}
}

func (c *ConstantsCache) Start(ctx context.Context, repo storage.IConstant) error {
	constants, err := repo.All(ctx)
	if err != nil {
		return err
	}

	for i := range constants {
		c.addConstant(&constants[i])
	}

	c.wg.Add(1)
	go c.listen(ctx)

	return nil
}

func (c *ConstantsCache) listen(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case constant := <-c.observer.Constants():
			c.addConstant(constant)
		}
	}
}

func (c *ConstantsCache) addConstant(constant *storage.Constant) {
	c.mx.Lock()
	{
		module := string(constant.Module)
		if _, ok := c.data[module]; !ok {
			c.data[module] = make(map[string]string)
		}
		c.data[module][constant.Name] = constant.Value
	}
	c.mx.Unlock()
}

func (c *ConstantsCache) Get(module types.ModuleName, name string) (string, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if m, ok := c.data[string(module)]; ok {
		val, ok := m[name]
		return val, ok
	}

	return "", false
}

func (c *ConstantsCache) Close() error {
	c.wg.Wait()
	return nil
}
