// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"
	"regexp"

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
	if err := v.RegisterValidation("status", statusValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("action_type", actionTypeValidator()); err != nil {
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

var regexAddressHash = regexp.MustCompile("[0-9A-Fa-f]{40}")

func isAddress(address string) bool {
	if len(address) != 40 {
		return false
	}
	return regexAddressHash.MatchString(address)
}

func addressValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return isAddress(fl.Field().String())
	}
}

func statusValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseStatus(fl.Field().String())
		return err == nil
	}
}

func actionTypeValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseActionType(fl.Field().String())
		return err == nil
	}
}
