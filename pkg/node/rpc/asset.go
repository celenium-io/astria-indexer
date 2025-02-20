// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"fmt"

	"github.com/celenium-io/astria-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

const pathAbciQuery = "abci_query"

func (api *API) GetAssetInfo(ctx context.Context, asset string) (types.DenomMetadataResponse, error) {
	args := make(map[string]string)
	args["path"] = fmt.Sprintf(`"asset/denom/%s"`, asset)

	var gbr types.Response[types.DenomMetadataResponse]
	if err := api.get(ctx, pathAbciQuery, args, &gbr); err != nil {
		return gbr.Result, errors.Wrap(err, "api.get")
	}

	if gbr.Error != nil {
		return gbr.Result, errors.Wrapf(types.ErrRequest, "request %d error: %s", gbr.Id, gbr.Error.Error())
	}

	return gbr.Result, nil
}
