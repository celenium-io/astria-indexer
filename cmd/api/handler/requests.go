// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"strings"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
)

const (
	asc  = "asc"
	desc = "desc"
)

func bindAndValidate[T any](c echo.Context) (*T, error) {
	req := new(T)
	if err := c.Bind(req); err != nil {
		return req, err
	}
	if err := c.Validate(req); err != nil {
		return req, err
	}
	return req, nil
}

func pgSort(sort string) storage.SortOrder {
	switch sort {
	case asc:
		return storage.SortOrderAsc
	case desc:
		return storage.SortOrderDesc
	default:
		return storage.SortOrderAsc
	}
}

type StringArray []string

func (s *StringArray) UnmarshalParam(param string) error {
	*s = StringArray(strings.Split(param, ","))
	return nil
}
