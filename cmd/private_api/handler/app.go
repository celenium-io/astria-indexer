// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type AppHandler struct {
	apps    storage.IApp
	address storage.IAddress
	rollup  storage.IRollup
	tx      sdk.Transactable
}

func NewAppHandler(
	apps storage.IApp,
	address storage.IAddress,
	rollup storage.IRollup,
	tx sdk.Transactable,
) *AppHandler {
	return &AppHandler{
		apps:    apps,
		address: address,
		rollup:  rollup,
		tx:      tx,
	}
}

var _ Handler = (*AppHandler)(nil)

func (handler *AppHandler) InitRoutes(srvr *echo.Group) {
	keyMiddleware := middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:Authorization",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("PRIVATE_API_AUTH_KEY"), nil
		},
	})

	app := srvr.Group("/app")
	{
		app.POST("", handler.Create, keyMiddleware)
		app.PATCH("/:id", handler.Update, keyMiddleware)
		app.DELETE("/:id", handler.Delete, keyMiddleware)
	}
}

type createAppRequest struct {
	Group        string   `json:"group"         validate:"omitempty,min=1"`
	Name         string   `json:"name"          validate:"required,min=1"`
	Description  string   `json:"description"   validate:"required,min=1"`
	Website      string   `json:"website"       validate:"omitempty,url"`
	GitHub       string   `json:"github"        validate:"omitempty,url"`
	Twitter      string   `json:"twitter"       validate:"omitempty,url"`
	Logo         string   `json:"logo"          validate:"omitempty,url"`
	L2Beat       string   `json:"l2beat"        validate:"omitempty,url"`
	Explorer     string   `json:"explorer"      validate:"omitempty,url"`
	Stack        string   `json:"stack"         validate:"omitempty"`
	Links        []string `json:"links"         validate:"omitempty,dive,url"`
	Category     string   `json:"category"      validate:"omitempty,app_category"`
	Type         string   `json:"type"          validate:"omitempty,app_type"`
	VM           string   `json:"vm"            validate:"omitempty"`
	Provider     string   `json:"provider"      validate:"omitempty"`
	Rollup       string   `json:"rollup"        validate:"required,base64"`
	NativeBridge string   `json:"native_bridge" validate:"omitempty,address"`
}

func (handler *AppHandler) Create(c echo.Context) error {
	req, err := bindAndValidate[createAppRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.createApp(c.Request().Context(), req); err != nil {
		return handleError(c, err, handler.apps)
	}

	return success(c)
}

func (handler *AppHandler) createApp(ctx context.Context, req *createAppRequest) error {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return err
	}

	hash, err := base64.StdEncoding.DecodeString(req.Rollup)
	if err != nil {
		return err
	}
	rollup, err := handler.rollup.ByHash(ctx, hash)
	if err != nil {
		return err
	}

	app := storage.App{
		Group:       req.Group,
		Name:        req.Name,
		Description: req.Description,
		Website:     req.Website,
		Github:      req.GitHub,
		Twitter:     req.Twitter,
		Logo:        req.Logo,
		L2Beat:      req.L2Beat,
		Explorer:    req.Explorer,
		Stack:       req.Stack,
		Links:       req.Links,
		Provider:    req.Provider,
		VM:          req.VM,
		Type:        types.AppType(req.Type),
		Category:    types.AppCategory(req.Category),
		Slug:        slug.Make(req.Name),
		RollupId:    rollup.Id,
	}

	if req.NativeBridge != "" {
		addr, err := handler.address.ByHash(ctx, req.NativeBridge)
		if err != nil {
			return err
		}
		if !addr.IsBridge {
			return tx.HandleError(ctx, errors.Errorf("address %s is not a bridge", req.NativeBridge))
		}

		app.NativeBridgeId = addr.Id
	}

	if err := tx.SaveApp(ctx, &app); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.RefreshLeaderboard(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}

	return tx.Flush(ctx)
}

type updateAppRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`

	Group        string   `json:"group"         validate:"omitempty,min=1"`
	Name         string   `json:"name"          validate:"omitempty,min=1"`
	Description  string   `json:"description"   validate:"omitempty,min=1"`
	Website      string   `json:"website"       validate:"omitempty,url"`
	GitHub       string   `json:"github"        validate:"omitempty,url"`
	Twitter      string   `json:"twitter"       validate:"omitempty,url"`
	Logo         string   `json:"logo"          validate:"omitempty,url"`
	L2Beat       string   `json:"l2beat"        validate:"omitempty,url"`
	Explorer     string   `json:"explorer"      validate:"omitempty,url"`
	Stack        string   `json:"stack"         validate:"omitempty"`
	Links        []string `json:"links"         validate:"omitempty,dive,url"`
	Category     string   `json:"category"      validate:"omitempty,app_category"`
	Type         string   `json:"type"          validate:"omitempty,app_type"`
	VM           string   `json:"vm"            validate:"omitempty"`
	Provider     string   `json:"provider"      validate:"omitempty"`
	Rollup       string   `json:"rollup"        validate:"omitempty,base64"`
	NativeBridge string   `json:"native_bridge" validate:"omitempty,address"`
}

func (handler *AppHandler) Update(c echo.Context) error {
	req, err := bindAndValidate[updateAppRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.updateRollup(c.Request().Context(), req); err != nil {
		return handleError(c, err, handler.apps)
	}

	return success(c)
}

func (handler *AppHandler) updateRollup(ctx context.Context, req *updateAppRequest) error {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return err
	}

	if _, err := handler.apps.GetByID(ctx, req.Id); err != nil {
		return err
	}

	app := storage.App{
		Id:          req.Id,
		Name:        req.Name,
		Slug:        slug.Make(req.Name),
		Description: req.Description,
		Website:     req.Website,
		Github:      req.GitHub,
		Twitter:     req.Twitter,
		Logo:        req.Logo,
		L2Beat:      req.L2Beat,
		Explorer:    req.Explorer,
		Stack:       req.Stack,
		Provider:    req.Provider,
		VM:          req.VM,
		Type:        types.AppType(req.Type),
		Category:    types.AppCategory(req.Category),
		Links:       req.Links,
	}

	if req.Rollup != "" {
		hash, err := base64.StdEncoding.DecodeString(req.Rollup)
		if err != nil {
			return err
		}
		rollup, err := handler.rollup.ByHash(ctx, hash)
		if err != nil {
			return err
		}
		app.RollupId = rollup.Id
	}

	if req.NativeBridge != "" {
		addr, err := handler.address.ByHash(ctx, req.NativeBridge)
		if err != nil {
			return err
		}
		if !addr.IsBridge {
			return tx.HandleError(ctx, errors.Errorf("address %s is not a bridge", req.NativeBridge))
		}

		app.NativeBridgeId = addr.Id
	}

	if err := tx.UpdateApp(ctx, &app); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.RefreshLeaderboard(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}

	return tx.Flush(ctx)
}

type deleteRollupRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`
}

func (handler *AppHandler) Delete(c echo.Context) error {
	req, err := bindAndValidate[deleteRollupRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.deleteRollup(c.Request().Context(), req.Id); err != nil {
		return handleError(c, err, handler.apps)
	}

	return success(c)
}

func (handler *AppHandler) deleteRollup(ctx context.Context, id uint64) error {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return err
	}

	if err := tx.DeleteApp(ctx, id); err != nil {
		return tx.HandleError(ctx, err)
	}

	return tx.Flush(ctx)
}
