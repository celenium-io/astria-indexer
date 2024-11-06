// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/base64"

	"github.com/celenium-io/astria-indexer/internal/astria"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type AppHandler struct {
	apps    storage.IApp
	address storage.IAddress
	bridge  storage.IBridge
	rollup  storage.IRollup
	tx      sdk.Transactable
}

func NewAppHandler(
	apps storage.IApp,
	address storage.IAddress,
	bridge storage.IBridge,
	rollup storage.IRollup,
	tx sdk.Transactable,
) AppHandler {
	return AppHandler{
		apps:    apps,
		address: address,
		bridge:  bridge,
		rollup:  rollup,
		tx:      tx,
	}
}

type createAppRequest struct {
	Group       string   `json:"group"       validate:"omitempty,min=1"`
	Name        string   `json:"name"        validate:"required,min=1"`
	Description string   `json:"description" validate:"required,min=1"`
	Website     string   `json:"website"     validate:"omitempty,url"`
	GitHub      string   `json:"github"      validate:"omitempty,url"`
	Twitter     string   `json:"twitter"     validate:"omitempty,url"`
	Logo        string   `json:"logo"        validate:"omitempty,url"`
	L2Beat      string   `json:"l2beat"      validate:"omitempty,url"`
	Explorer    string   `json:"explorer"    validate:"omitempty,url"`
	Stack       string   `json:"stack"       validate:"omitempty"`
	Links       []string `json:"links"       validate:"omitempty,dive,url"`
	Category    string   `json:"category"    validate:"omitempty,app_category"`
	Type        string   `json:"type"        validate:"omitempty,app_type"`
	VM          string   `json:"vm"          validate:"omitempty"`
	Provider    string   `json:"provider"    validate:"omitempty"`

	AppIds  []appId  `json:"app_ids" validate:"required,min=1"`
	Bridges []bridge `json:"bridges" validate:"omitempty"`
}

type appId struct {
	Rollup  string `json:"rollup"  validate:"omitempty,base64"`
	Address string `json:"address" validate:"required,address"`
}

type bridge struct {
	Address string `json:"address" validate:"required,address"`
	Native  bool   `json:"native"  validate:"omitempty"`
}

func (handler AppHandler) Create(c echo.Context) error {
	req, err := bindAndValidate[createAppRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.createApp(c.Request().Context(), req); err != nil {
		return handleError(c, err, handler.apps)
	}

	return success(c)
}

func (handler AppHandler) createApp(ctx context.Context, req *createAppRequest) error {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return err
	}

	rollup := storage.App{
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
	}

	if err := tx.SaveApp(ctx, &rollup); err != nil {
		return tx.HandleError(ctx, err)
	}

	appIds, err := handler.createAppIds(ctx, rollup.Id, req.AppIds...)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.SaveAppId(ctx, appIds...); err != nil {
		return tx.HandleError(ctx, err)
	}

	bridges, err := handler.createBridges(ctx, rollup.Id, req.Bridges...)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.SaveAppBridges(ctx, bridges...); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.RefreshLeaderboard(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}

	return tx.Flush(ctx)
}

func (handler AppHandler) createAppIds(ctx context.Context, id uint64, data ...appId) ([]storage.AppId, error) {
	providers := make([]storage.AppId, len(data))
	for i := range data {
		providers[i].AppId = id

		if !astria.IsAddress(data[i].Address) {
			return nil, errors.Wrap(errInvalidAddress, data[i].Address)
		}

		address, err := handler.address.ByHash(ctx, data[i].Address)
		if err != nil {
			return nil, err
		}
		providers[i].AddressId = address.Id

		if data[i].Rollup != "" {
			hashRollup, err := base64.StdEncoding.DecodeString(data[i].Rollup)
			if err != nil {
				return nil, err
			}
			rollup, err := handler.rollup.ByHash(ctx, hashRollup)
			if err != nil {
				return nil, err
			}
			providers[i].RolllupId = rollup.Id
		}
	}
	return providers, nil
}

func (handler AppHandler) createBridges(ctx context.Context, id uint64, data ...bridge) ([]storage.AppBridge, error) {
	bridges := make([]storage.AppBridge, len(data))
	for i := range data {
		bridges[i].AppId = id

		if !astria.IsAddress(data[i].Address) {
			return nil, errors.Wrap(errInvalidAddress, data[i].Address)
		}

		address, err := handler.address.ByHash(ctx, data[i].Address)
		if err != nil {
			return nil, err
		}

		b, err := handler.bridge.ByAddress(ctx, address.Id)
		if err != nil {
			return nil, err
		}
		bridges[i].BridgeId = b.Id
		bridges[i].Native = data[i].Native
	}
	return bridges, nil
}

type updateAppRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`

	Group       string   `json:"group"       validate:"omitempty,min=1"`
	Name        string   `json:"name"        validate:"omitempty,min=1"`
	Description string   `json:"description" validate:"omitempty,min=1"`
	Website     string   `json:"website"     validate:"omitempty,url"`
	GitHub      string   `json:"github"      validate:"omitempty,url"`
	Twitter     string   `json:"twitter"     validate:"omitempty,url"`
	Logo        string   `json:"logo"        validate:"omitempty,url"`
	L2Beat      string   `json:"l2beat"      validate:"omitempty,url"`
	Explorer    string   `json:"explorer"    validate:"omitempty,url"`
	Stack       string   `json:"stack"       validate:"omitempty"`
	Links       []string `json:"links"       validate:"omitempty,dive,url"`
	Category    string   `json:"category"    validate:"omitempty,app_category"`
	Type        string   `json:"type"        validate:"omitempty,app_type"`
	VM          string   `json:"vm"          validate:"omitempty"`
	Provider    string   `json:"provider"    validate:"omitempty"`

	AppIds  []appId  `json:"app_ids" validate:"omitempty,min=1"`
	Bridges []bridge `json:"bridges" validate:"omitempty"`
}

func (handler AppHandler) Update(c echo.Context) error {
	req, err := bindAndValidate[updateAppRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.updateRollup(c.Request().Context(), req); err != nil {
		return handleError(c, err, handler.apps)
	}

	return success(c)
}

func (handler AppHandler) updateRollup(ctx context.Context, req *updateAppRequest) error {
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

	if err := tx.UpdateApp(ctx, &app); err != nil {
		return tx.HandleError(ctx, err)
	}

	if len(req.AppIds) > 0 {
		if err := tx.DeleteAppId(ctx, req.Id); err != nil {
			return tx.HandleError(ctx, err)
		}

		appIds, err := handler.createAppIds(ctx, app.Id, req.AppIds...)
		if err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := tx.SaveAppId(ctx, appIds...); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if len(req.Bridges) > 0 {
		if err := tx.DeleteAppBridges(ctx, req.Id); err != nil {
			return tx.HandleError(ctx, err)
		}

		bridges, err := handler.createBridges(ctx, app.Id, req.Bridges...)
		if err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := tx.SaveAppBridges(ctx, bridges...); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.RefreshLeaderboard(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}

	return tx.Flush(ctx)
}

type deleteRollupRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`
}

func (handler AppHandler) Delete(c echo.Context) error {
	req, err := bindAndValidate[deleteRollupRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.deleteRollup(c.Request().Context(), req.Id); err != nil {
		return handleError(c, err, handler.apps)
	}

	return success(c)
}

func (handler AppHandler) deleteRollup(ctx context.Context, id uint64) error {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return err
	}

	if err := tx.DeleteAppId(ctx, id); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.DeleteAppBridges(ctx, id); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.DeleteApp(ctx, id); err != nil {
		return tx.HandleError(ctx, err)
	}

	return tx.Flush(ctx)
}
