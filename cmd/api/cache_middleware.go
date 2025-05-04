// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/celenium-io/astria-indexer/cmd/api/cache"
	"github.com/labstack/echo/v4"
	"time"
)

func initCache(url string) cache.ICache {
	if url != "" {
		c, err := cache.NewValKey(url, time.Hour)
		if err != nil {
			panic(err)
		}
		return c
	}

	return nil
}

func newDefaultMiddlewareCache(ttlCache cache.ICache) echo.MiddlewareFunc {
	return cache.Middleware(ttlCache, nil, nil)
}

func newStatMiddlewareCache(ttlCache cache.ICache) echo.MiddlewareFunc {
	return cache.Middleware(ttlCache, nil, func() time.Duration {
		now := time.Now()
		diff := now.Truncate(time.Hour).Add(time.Hour).Sub(now)
		if diff > time.Minute*10 {
			return time.Minute * 10
		}
		return diff
	})
}
