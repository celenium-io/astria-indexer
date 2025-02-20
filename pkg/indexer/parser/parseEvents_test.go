// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"
	"testing"

	"github.com/celenium-io/astria-indexer/pkg/indexer/decode"
	"github.com/celenium-io/astria-indexer/pkg/node/mock"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
		err := parseTxDeposit(attrs, 100, &ctx)
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

func Test_parseTxFee(t *testing.T) {
	t.Run("test fee event", func(t *testing.T) {
		attrs := []types.EventAttribute{
			{
				Key:   "actionName",
				Value: "astria.protocol.transaction.v1.RollupDataSubmission",
			}, {
				Key:   "asset",
				Value: "ibc/704031c868fd3d3c84a1cfa8cb45deba4ea746b44697f7f4a6ed1b8f6c239b82",
			}, {
				Key:   "feeAmount",
				Value: "321",
			}, {
				Key:   "positionInTransaction",
				Value: "0",
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		api := mock.NewMockApi(ctrl)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		decodeCtx := decode.NewContext(map[string]string{})
		err := parseTxFees(ctx, attrs, &decodeCtx, api)
		require.NoError(t, err)
		require.Len(t, decodeCtx.Fees, 1)

		fee, ok := decodeCtx.Fees[0]
		require.True(t, ok)
		require.EqualValues(t, "321", fee.Amount.String())
		require.EqualValues(t, "nria", fee.Asset)
		require.EqualValues(t, "astria.protocol.transaction.v1.RollupDataSubmission", fee.ActionType)
	})
}

func Test_parseWriteAck(t *testing.T) {
	t.Run("test write_acknowledgement event", func(t *testing.T) {
		attrs := []types.EventAttribute{
			{
				Key:   "packet_ack",
				Value: `{\"error\":\"{\"0\": \"failed to execute ics20 transfer to bridge account\", \"1\": \"failed to parse memo as Ics20TransferDepositMemo\", \"2\": \"EOF while parsing a value at line 1 column 0\"}\"}`,
			}, {
				Key:   "packet_data",
				Value: `{"amount":"100000","denom":"uusdc","receiver":"astriacompat13vptdafyttpmlwppt0s844efey2cpc0mw6dhm2","sender":"noble1rmhdkl3aaw95zdecnj5paaqcjavm8sylftznrs"}`,
			}, {
				Key:   "packet_sequence",
				Value: "1",
			}, {
				Key:   "packet_dst_port",
				Value: "transfer",
			},
		}

		decodeCtx := decode.NewContext(map[string]string{})
		err := parseWriteAck(attrs, &decodeCtx)
		require.NoError(t, err)
		require.True(t, decodeCtx.HasWriteAckError)
	})

	t.Run("test write_acknowledgement success event", func(t *testing.T) {
		attrs := []types.EventAttribute{
			{
				Key:   "packet_ack",
				Value: `{\"result\":\"AQ==\"}`,
			}, {
				Key:   "packet_data",
				Value: `{"amount":"100000","denom":"uusdc","receiver":"astriacompat13vptdafyttpmlwppt0s844efey2cpc0mw6dhm2","sender":"noble1rmhdkl3aaw95zdecnj5paaqcjavm8sylftznrs"}`,
			}, {
				Key:   "packet_sequence",
				Value: "1",
			}, {
				Key:   "packet_dst_port",
				Value: "transfer",
			},
		}

		decodeCtx := decode.NewContext(map[string]string{})
		err := parseWriteAck(attrs, &decodeCtx)
		require.NoError(t, err)
		require.False(t, decodeCtx.HasWriteAckError)
	})
}
