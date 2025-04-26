// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/shopspring/decimal"
)

type Constants struct {
	Module map[string]Params `json:"module"`
}

type Params map[string]string

func roundCounstant(val string) string {
	d, err := decimal.NewFromString(val)
	if err != nil {
		return val
	}
	return d.String()
}

func NewConstants(consts []storage.Constant) Constants {
	response := Constants{
		Module: make(map[string]Params),
	}

	for i := range consts {
		if params, ok := response.Module[string(consts[i].Module)]; ok {
			params[consts[i].Name] = roundCounstant(consts[i].Value)
		} else {
			response.Module[string(consts[i].Module)] = Params{
				consts[i].Name: roundCounstant(consts[i].Value),
			}
		}
	}

	return response
}

type Enums struct {
	Status     []string `json:"status"`
	ActionType []string `json:"action_type"`
}

func NewEnums() Enums {
	return Enums{
		Status:     types.StatusNames(),
		ActionType: types.ActionTypeNames(),
	}
}
