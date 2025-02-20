// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import "github.com/labstack/echo/v4"

type Handler interface {
	InitRoutes(srvr *echo.Group)
}
