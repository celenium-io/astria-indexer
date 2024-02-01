// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package node

import (
	"context"

	pkgTypes "github.com/aopoltorzhicky/astria/pkg/types"

	"github.com/aopoltorzhicky/astria/pkg/node/types"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Api interface {
	Status(ctx context.Context) (types.Status, error)
	Head(ctx context.Context) (pkgTypes.ResultBlock, error)
	Block(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlock, error)
	BlockResults(ctx context.Context, level pkgTypes.Level) (pkgTypes.ResultBlockResults, error)
	Genesis(ctx context.Context) (types.Genesis, error)
	BlockData(ctx context.Context, level pkgTypes.Level) (pkgTypes.BlockData, error)
	BlockDataGet(ctx context.Context, level pkgTypes.Level) (pkgTypes.BlockData, error)
}
