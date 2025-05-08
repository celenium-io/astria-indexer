// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"time"

	"github.com/labstack/echo/v4"
)

func InitCache(url string) ICache {
	if url != "" {
		c, err := NewValKey(url, time.Hour)
		if err != nil {
			panic(err)
		}
		return c
	}

	return nil
}

func NewDefaultMiddlewareCache(ttlCache ICache) echo.MiddlewareFunc {
	if ttlCache == nil {
		return nil
	}

	return Middleware(ttlCache, nil, nil)
}

func NewStatMiddlewareCache(ttlCache ICache) echo.MiddlewareFunc {
	return Middleware(ttlCache, nil, func() time.Duration {
		now := time.Now()
		diff := now.Truncate(time.Hour).Add(time.Hour).Sub(now)
		if diff > time.Minute*10 {
			return time.Minute * 10
		}
		return diff
	})
}
