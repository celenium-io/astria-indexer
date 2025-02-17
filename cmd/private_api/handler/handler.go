package handler

import "github.com/labstack/echo/v4"

type Handler interface {
	InitRoutes(srvr *echo.Group)
}
