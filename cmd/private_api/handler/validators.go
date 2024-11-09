// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/astria-indexer/internal/astria"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ApiValidator struct {
	validator *validator.Validate
}

func NewApiValidator() *ApiValidator {
	v := validator.New()
	if err := v.RegisterValidation("address", addressValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("app_category", categoryValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("app_type", appTypeValidator()); err != nil {
		panic(err)
	}
	return &ApiValidator{validator: v}
}

func (v *ApiValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func isAddress(address string) bool {
	return astria.IsAddress(address)
}

func addressValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return isAddress(fl.Field().String())
	}
}

func categoryValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseAppCategory(fl.Field().String())
		return err == nil
	}
}

func appTypeValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseAppType(fl.Field().String())
		return err == nil
	}
}
