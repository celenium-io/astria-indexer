// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"

	"github.com/aopoltorzhicky/astria/pkg/types"
)

func (api *API) Head(ctx context.Context) (types.ResultBlock, error) {
	return api.Block(ctx, 0)
}
