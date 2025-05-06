// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func websocketSkipper(c echo.Context) bool {
	return strings.Contains(c.Request().URL.Path, "ws")
}

func postSkipper(c echo.Context) bool {
	if strings.Contains(c.Request().URL.Path, "blob") {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "auth/rollup") {
		return true
	}
	return false
}

func gzipSkipper(c echo.Context) bool {
	if strings.Contains(c.Request().URL.Path, "swagger") {
		return true
	}
	if strings.Contains(c.Request().URL.Path, "metrics") {
		return true
	}
	return websocketSkipper(c)
}
