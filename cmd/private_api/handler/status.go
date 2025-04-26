// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type StatusHandler struct {
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

func (handler *StatusHandler) Status(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

var _ Handler = (*StatusHandler)(nil)

func (handler *StatusHandler) InitRoutes(srvr *echo.Group) {
	srvr.GET("/status", handler.Status)
}
