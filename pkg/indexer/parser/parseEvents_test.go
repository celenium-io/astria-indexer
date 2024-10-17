// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"testing"

	"github.com/celenium-io/astria-indexer/pkg/indexer/decode"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/stretchr/testify/require"
)

func Test_parseTxDeposit(t *testing.T) {
	t.Run("test deposit event", func(t *testing.T) {
		attrs := []types.EventAttribute{
			{
				Key:   "bridgeAddress",
				Value: "astria1yqdjnnmrp7w5ygwj0dkldsgzjhv5vcakp7yeu9",
			}, {
				Key:   "rollupId",
				Value: "xuzPaQ1qQbu1wl5oBdEiP-Xl977xEmJOVjVpqNrktnU=",
			}, {
				Key:   "amount",
				Value: "1000000000",
			}, {
				Key:   "asset",
				Value: "nria",
			}, {
				Key:   "destinationChainAddress",
				Value: "0x9d0CEC7BEB948Ab046e8b64E9aa6Cc9b731A9613",
			}, {
				Key:   "sourceTransactionId",
				Value: "c8c58bfd5d4745c08149da81c24df6738659634507c3eb535fa803c69d6a2b46",
			}, {
				Key:   "sourceActionIndex",
				Value: "0",
			},
		}

		ctx := decode.NewContext(map[string]string{})
		err := parseTxDeposit(attrs, &ctx)
		require.NoError(t, err)
		require.Len(t, ctx.Deposits, 1)

		deposit := ctx.Deposits[0]
		require.EqualValues(t, "astria1yqdjnnmrp7w5ygwj0dkldsgzjhv5vcakp7yeu9", deposit.Bridge.Address.String())
		require.NotNil(t, deposit.Rollup.AstriaId)
		require.EqualValues(t, "1000000000", deposit.Amount.String())
		require.EqualValues(t, "nria", deposit.Asset)
		require.EqualValues(t, "0x9d0CEC7BEB948Ab046e8b64E9aa6Cc9b731A9613", deposit.DestinationChainAddress)
	})
}
