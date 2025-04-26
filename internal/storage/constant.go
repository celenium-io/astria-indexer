// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IConstant interface {
	Get(ctx context.Context, module types.ModuleName, name string) (Constant, error)
	ByModule(ctx context.Context, module types.ModuleName) ([]Constant, error)
	All(ctx context.Context) ([]Constant, error)
	IsNoRows(err error) bool
}

type Constant struct {
	bun.BaseModel `bun:"table:constant" comment:"Table with constants"`

	Module types.ModuleName `bun:"module,pk,type:module_name" comment:"Module name which declares constant" json:"module"`
	Name   string           `bun:"name,pk,type:text"          comment:"Constant name"                       json:"name"`
	Value  string           `bun:"value,type:text"            comment:"Constant value"                      json:"value"`
}

func (Constant) TableName() string {
	return "constant"
}
