// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	primitivev1 "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	feesv1alpha1 "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/fees/v1alpha1"
	astria "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transaction/v1alpha1"
	v1 "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria_vendored/penumbra/core/component/ibc/v1"
	abci "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria_vendored/tendermint/abci"
	crypto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria_vendored/tendermint/crypto"
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/astria-indexer/internal/test_suite"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	assetId    = "nria"
	feeAssetId = "fee_asset"
)

func TestDecodeActions(t *testing.T) {

	t.Run("ibc action", func(t *testing.T) {
		message := &astria.Action_Ibc{
			Ibc: &v1.IbcRelay{
				RawAction: &anypb.Any{
					Value:   []byte{0, 0, 0, 0},
					TypeUrl: "msg_type",
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcRelay,
			Data: map[string]any{
				"raw":  "AAAAAA==",
				"type": "msg_type",
			},
		}

		action := storage.Action{
			Height: 1000,
		}

		ctx := NewContext(nil)
		err := parseIbcAction(message, &ctx, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("ibc action: MsgRecvPacket", func(t *testing.T) {
		raw := "CtwBCA4SCHRyYW5zZmVyGgtjaGFubmVsLTE2MCIIdHJhbnNmZXIqCWNoYW5uZWwtMDKZAXsiYW1vdW50IjoiMTAwMDAwMCIsImRlbm9tIjoidXRpYSIsInJlY2VpdmVyIjoiYXN0cmlhMTNuazlnZnA3dDVxZDc3dXZmZ3RwZG04czZ0a2xzMmE5c2ZhN2NnIiwic2VuZGVyIjoiY2VsZXN0aWExNzNreDM1ZmxsbHZ1N2NuOWpnc2hxaHQ0NXpjeG55dXBnOWZ6M2QifToGCAEQhYcKQMi7iqqh96WAGBKhCQqdBwqaBwo8Y29tbWl0bWVudHMvcG9ydHMvdHJhbnNmZXIvY2hhbm5lbHMvY2hhbm5lbC0xNjAvc2VxdWVuY2VzLzE0EiDF/8Kt2IYZocp1frOVz7thvZ9lSc04VMXjjLd6TK3h2RoOCAEYASABKgYAAry/6gIiLAgBEigCBLy/6gIgLuSKhXE31+NU/81ooA5xyV/6RDhuiyoe6ISY557QdaUgIiwIARIoBAi8v+oCILmZik8suYT5PZBngmPEpVktKyMnXMQzQE8YiER7SQOZICIuCAESBwYMvL/qAiAaISDgSzqseZ2lK1iNp5HRNfuNi1E9C+afiTBIs9ElhhtWwCIsCAESKAgYvL/qAiD4SeyplPW835TS+ME0s7y0cPoXAiZWt9JDmJDVva/kfCAiLAgBEigKMLy/6gIgnReiffi6CTiFzdWHcVLV/bhDaWAVHmVUXkEQMtF4tYcgIi4IARIHDEi8v+oCIBohIA/CHKbG7a3mZSMnQnluPclvDtFlK0MQOEOkEE11kMMyIiwIARIoDmy8v+oCIHGSSADA9gkQgvKY7BmqfmKhpI94B8otqcgsKK8W9a0JICItCAESKRD2Aby/6gIgVrDJq95UwgiLSMHQ6bozBYHacAjcybOGddwm4ud/1nogIi8IARIIEvoCvL/qAiAaISBF3QfuSqRYINaQrZr/Ryka+Ee8MfoX3wgs1Jdw2In1JiIvCAESCBS8B7y/6gIgGiEgaY4SsHsxSMaQXR9+PNkPiML7WAIVH6ZGO4aLcyeM60YiLQgBEikW6gu8v+oCIID4Tiv4C3/tNFDmn4z5X6L9yVKsBfWyuaem+X52G6pPICIvCAESCBjMEry/6gIgGiEgxR4SceAkaiHu5rJg+IRmypwDwhDXEVQlGh4q1gKBmdgiLQgBEikcukC8v+oCIEMbaafJ5iGvjw98CnJvovD9+tBTxW06GkmCDWIa9BFpICIuCAESKh6CkAG8v+oCIBhX60I7dKkHLhCsblc4i25nsRH8aoEvnlosHfLiLgjMICIwCAESCSCOsQK8v+oCIBohIGjAlyFlC1aoLcdXCj8BZ6fbvAH8RsZHqiNTNzTf6UXrIi4IARIqJOS3Bby/6gIgD1M0oseDXxLVDSRkj7fEJCd7U4dg4uWMNwpGaxNKqeMgIjAIARIJJpq/Cby/6gIgGiEgYCM0Ij2+xQpYwRKJ52NzrFcVppMbJ1/VvNDDvNHZHuAK/gEK+wEKA2liYxIgDKSZICeXaQVzlcNE0Cv/QOlU2nDreTxbpslXZibxQTYaCQgBGAEgASoBACIlCAESIQFc8LCXIH/ReynuGPU2hd+7XpwozF25Cg2lAmKb/BmnHCInCAESAQEaIGCg+C171gjmvrp9OkKH5uaMMYv2ZUNybiq88hDdeRCWIicIARIBARogJ4+4wexaAIVpe77Od5f2bnkfbkxl8rV3F3GeFMa/QlMiJQgBEiEBbwCFwIFOC8GpODyVT+os/fDZA2h1eLy2LHVjXhdhVGciJwgBEgEBGiAO17o0UZMScdg7UmBLYNIruHT5zwW1KsRdxrWjek9IoRoHCAQQ35+1ASItYXN0cmlhMWp4ZjY3bTlhOHcyMG55anI0ZTZtY3N5NTZjbWNrdnpqMGQ4ZzA3"
		data, err := base64.StdEncoding.DecodeString(raw)
		require.NoError(t, err)

		message := &astria.Action_Ibc{
			Ibc: &v1.IbcRelay{
				RawAction: &anypb.Any{
					Value:   data,
					TypeUrl: "/ibc.core.channel.v1.MsgRecvPacket",
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcRelay,
			Data: map[string]any{
				"raw":  raw,
				"type": "/ibc.core.channel.v1.MsgRecvPacket",
			},
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Height:   1000,
					Currency: "transfer/channel-0/utia",
					Address: &storage.Address{
						Hash:   "astria13nk9gfp7t5qd77uvfgtpdm8s6tkls2a9sfa7cg",
						Height: 1000,
						Balance: []*storage.Balance{
							{
								Currency: "transfer/channel-0/utia",
								Total:    decimal.RequireFromString("1000000"),
							},
						},
					},
					Update: decimal.RequireFromString("1000000"),
				},
			},
		}

		action := storage.Action{
			Height: 1000,
		}
		ctx := NewContext(nil)
		err = parseIbcAction(message, &ctx, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
		require.Len(t, ctx.Addresses, 1)
	})

	t.Run("ibc action: MsgRecvPacket without astria address", func(t *testing.T) {
		raw := "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"
		data, err := base64.StdEncoding.DecodeString(raw)
		require.NoError(t, err)

		message := &astria.Action_Ibc{
			Ibc: &v1.IbcRelay{
				RawAction: &anypb.Any{
					Value:   data,
					TypeUrl: "/ibc.core.channel.v1.MsgRecvPacket",
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcRelay,
			Data: map[string]any{
				"raw":  raw,
				"type": "/ibc.core.channel.v1.MsgRecvPacket",
			},
		}

		action := storage.Action{
			Height: 1000,
		}
		ctx := NewContext(nil)
		err = parseIbcAction(message, &ctx, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
		require.Len(t, ctx.Addresses, 0)
	})

	t.Run("ibc action: MsgRecvPacket with compat astria address", func(t *testing.T) {
		raw := "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"
		data, err := base64.StdEncoding.DecodeString(raw)
		require.NoError(t, err)

		message := &astria.Action_Ibc{
			Ibc: &v1.IbcRelay{
				RawAction: &anypb.Any{
					Value:   data,
					TypeUrl: "/ibc.core.channel.v1.MsgRecvPacket",
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcRelay,
			Data: map[string]any{
				"raw":  raw,
				"type": "/ibc.core.channel.v1.MsgRecvPacket",
			},
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Height:   1000,
					Update:   decimal.NewFromInt(300000000),
					Currency: "transfer/channel-1/uusdc",
					Address: &storage.Address{
						Hash:   "astria1eg8hhey0n4untdvqqdvlyl0e7zx8wfca48kglh",
						Height: 1000,
						Balance: []*storage.Balance{
							{
								Currency: "transfer/channel-1/uusdc",
								Total:    decimal.NewFromInt(300000000),
							},
						},
					},
				},
			},
		}

		action := storage.Action{
			Height: 1000,
		}
		ctx := NewContext(nil)
		err = parseIbcAction(message, &ctx, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
		require.Len(t, ctx.Addresses, 1)
	})

	t.Run("ics 20 withdrawal", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		address := testsuite.RandomAddress()
		from := testsuite.RandomAddress()

		message := &astria.Action_Ics20Withdrawal{
			Ics20Withdrawal: &astria.Ics20Withdrawal{
				Amount: &primitivev1.Uint128{
					Lo: 1,
					Hi: 0,
				},
				Denom:                   currency.DefaultCurrency,
				DestinationChainAddress: "celestia1lx7dfjp20shd6y5f4tauvy8cv4pjhvszfrh9ah",
				ReturnAddress: &primitivev1.Address{
					Bech32M: address,
				},
				SourceChannel: "channel-12",
				TimeoutHeight: &astria.IbcHeight{
					RevisionNumber: 1000,
					RevisionHeight: 1001,
				},
				TimeoutTime: 1000,
				FeeAsset:    feeAssetId,
				Memo:        "memo",
				BridgeAddress: &primitivev1.Address{
					Bech32M: from,
				},
				UseCompatAddress: true,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIcs20Withdrawal,
			Data: map[string]any{
				"amount":              "1",
				"denom":               currency.DefaultCurrency,
				"destination_address": "celestia1lx7dfjp20shd6y5f4tauvy8cv4pjhvszfrh9ah",
				"return_address":      address,
				"source_channel":      "channel-12",
				"timeout_height": map[string]any{
					"revision_number": uint64(1000),
					"revision_height": uint64(1001),
				},
				"timeout_time":       uint64(1000),
				"fee_asset":          feeAssetId,
				"memo":               "memo",
				"bridge":             from,
				"use_compat_address": true,
			},
			Addresses: []*storage.AddressAction{},
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Height: 1000,
					Address: &storage.Address{
						Height:       1000,
						Hash:         address,
						ActionsCount: 1,
						Balance: []*storage.Balance{
							{
								Currency: currency.DefaultCurrency,
								Total:    decimal.RequireFromString("1"),
							},
						},
					},
					Currency: currency.DefaultCurrency,
					Update:   decimal.RequireFromString("-1"),
				},
			},
		}
		addressAction := storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         address,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.RequireFromString("1"),
					},
				},
			},
			ActionType: types.ActionTypeIcs20Withdrawal,
			Action:     &wantAction,
		}
		addressActionFrom := storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         from,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.RequireFromString("-1"),
					},
				},
			},
			ActionType: types.ActionTypeIcs20Withdrawal,
			Action:     &wantAction,
		}
		wantAction.Addresses = append(wantAction.Addresses, &addressActionFrom, &addressAction)

		action := storage.Action{
			Height: 1000,
		}
		err := parseIcs20Withdrawal(message, from, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("rollup data submission", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		from := testsuite.RandomAddress()
		addressModel := decodeContext.Addresses.Set(from, 1000, decimal.Zero, "", 0, 1)

		message := &astria.Action_RollupDataSubmission{
			RollupDataSubmission: &astria.RollupDataSubmission{
				RollupId: &primitivev1.RollupId{Inner: testsuite.RandomHash(10)},
				Data:     testsuite.RandomHash(10),
				FeeAsset: feeAssetId,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeRollupDataSubmission,
			Data: map[string]any{
				"rollup_id": message.RollupDataSubmission.GetRollupId().GetInner(),
				"data":      message.RollupDataSubmission.GetData(),
				"fee_asset": feeAssetId,
			},
			Addresses: make([]*storage.AddressAction, 0),
			RollupAction: &storage.RollupAction{
				Size:   10,
				Height: 1000,
				Rollup: &storage.Rollup{
					AstriaId:     message.RollupDataSubmission.GetRollupId().GetInner(),
					FirstHeight:  1000,
					ActionsCount: 1,
					Size:         10,
				},
				ActionType: types.ActionTypeRollupDataSubmission,
			},
		}
		wantAction.RollupAction.Action = &wantAction
		addressAction := storage.AddressAction{
			Height:     1000,
			Address:    addressModel,
			ActionType: types.ActionTypeRollupDataSubmission,
			Action:     &wantAction,
		}
		wantAction.Addresses = append(wantAction.Addresses, &addressAction)

		action := storage.Action{
			Height: 1000,
		}
		err := parseRollupDataSubmission(message, from, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("sudo address change", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		newAddress := testsuite.RandomAddress()
		message := &astria.Action_SudoAddressChange{
			SudoAddressChange: &astria.SudoAddressChange{
				NewAddress: &primitivev1.Address{Bech32M: newAddress},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeSudoAddressChange,
			Data: map[string]any{
				"address": newAddress,
			},
			Addresses: make([]*storage.AddressAction, 0),
		}
		balance := storage.EmptyBalance()
		addressAction := storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				ActionsCount: 1,
				Hash:         newAddress,
				Balance:      []*storage.Balance{&balance},
			},
			ActionType: types.ActionTypeSudoAddressChange,
			Action:     &wantAction,
		}
		wantAction.Addresses = append(wantAction.Addresses, &addressAction)

		action := storage.Action{
			Height: 1000,
		}
		err := parseSudoAddressChangeAction(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("transfer", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		from := testsuite.RandomAddress()
		fromModel := &storage.Address{
			Height:        1000,
			Hash:          from,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: currency.DefaultCurrency,
					Total:    decimal.RequireFromString("-10"),
				},
			},
		}

		to := testsuite.RandomAddress()
		toModel := &storage.Address{
			Height:        1000,
			Hash:          to,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: currency.DefaultCurrency,
					Total:    decimal.RequireFromString("10"),
				},
			},
		}

		message := &astria.Action_Transfer{
			Transfer: &astria.Transfer{
				To: &primitivev1.Address{Bech32M: to},
				Amount: &primitivev1.Uint128{
					Lo: 10,
				},
				Asset:    assetId,
				FeeAsset: feeAssetId,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeTransfer,
			Data: map[string]any{
				"to":        to,
				"amount":    "10",
				"asset":     message.Transfer.GetAsset(),
				"fee_asset": message.Transfer.GetFeeAsset(),
			},
			Addresses: make([]*storage.AddressAction, 0),
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Address:  toModel,
					Update:   toModel.Balance[0].Total,
					Currency: toModel.Balance[0].Currency,
					Height:   1000,
				}, {
					Address:  fromModel,
					Update:   fromModel.Balance[0].Total,
					Currency: fromModel.Balance[0].Currency,
					Height:   1000,
				},
			},
		}

		wantAction.Addresses = append(wantAction.Addresses, &storage.AddressAction{
			Height:     1000,
			Address:    toModel,
			ActionType: types.ActionTypeTransfer,
			Action:     &wantAction,
		}, &storage.AddressAction{
			Height:     1000,
			Address:    fromModel,
			ActionType: types.ActionTypeTransfer,
			Action:     &wantAction,
		})

		action := storage.Action{
			Height: 1000,
		}
		err := parseTransferAction(message, from, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("transfer to myself", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		from := testsuite.RandomAddress()
		fromModel := &storage.Address{
			Height:        1000,
			Hash:          from,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: currency.DefaultCurrency,
					Total:    decimal.Zero,
				},
			},
		}

		message := &astria.Action_Transfer{
			Transfer: &astria.Transfer{
				To: &primitivev1.Address{Bech32M: from},
				Amount: &primitivev1.Uint128{
					Lo: 10,
				},
				Asset:    assetId,
				FeeAsset: feeAssetId,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeTransfer,
			Data: map[string]any{
				"to":        from,
				"amount":    "10",
				"asset":     message.Transfer.GetAsset(),
				"fee_asset": message.Transfer.GetFeeAsset(),
			},
			Addresses: make([]*storage.AddressAction, 0),
		}

		wantAction.Addresses = append(wantAction.Addresses, &storage.AddressAction{
			Height:     1000,
			Address:    fromModel,
			ActionType: types.ActionTypeTransfer,
			Action:     &wantAction,
		})

		action := storage.Action{
			Height: 1000,
		}
		err := parseTransferAction(message, from, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("validator update", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})
		message := &astria.Action_ValidatorUpdate{
			ValidatorUpdate: &abci.ValidatorUpdate{
				PubKey: &crypto.PublicKey{
					Sum: &crypto.PublicKey_Ed25519{
						Ed25519: testsuite.RandomHash(32),
					},
				},
				Power: 10,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeValidatorUpdate,
			Data: map[string]any{
				"power":  int64(10),
				"pubkey": message.ValidatorUpdate.GetPubKey().GetEd25519(),
			},
			Addresses: make([]*storage.AddressAction, 0),
		}
		address, err := AddressFromPubKey(message.ValidatorUpdate.GetPubKey().GetEd25519())
		require.NoError(t, err)

		balance := storage.EmptyBalance()
		addressAction := storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				ActionsCount: 1,
				Hash:         address,
				Balance:      []*storage.Balance{&balance},
			},
			ActionType: types.ActionTypeValidatorUpdate,
			Action:     &wantAction,
		}
		wantAction.Addresses = append(wantAction.Addresses, &addressAction)

		action := storage.Action{
			Height: 1000,
		}
		err = parseValidatorUpdateAction(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)

		require.Len(t, decodeContext.Validators, 1)
		pk := hex.EncodeToString(message.ValidatorUpdate.GetPubKey().GetEd25519())
		v, ok := decodeContext.Validators[pk]
		require.True(t, ok)
		require.EqualValues(t, "10", v.Power.String())
	})

	t.Run("fee asset change: addition", func(t *testing.T) {
		message := &astria.Action_FeeAssetChange{
			FeeAssetChange: &astria.FeeAssetChange{
				Value: &astria.FeeAssetChange_Addition{
					Addition: assetId,
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeFeeAssetChange,
			Data: map[string]any{
				"addition": assetId,
			},
		}

		action := storage.Action{
			Height: 1000,
		}
		err := parseFeeAssetChange(message, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee asset change: removal", func(t *testing.T) {
		message := &astria.Action_FeeAssetChange{
			FeeAssetChange: &astria.FeeAssetChange{
				Value: &astria.FeeAssetChange_Removal{
					Removal: assetId,
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeFeeAssetChange,
			Data: map[string]any{
				"removal": assetId,
			},
		}

		action := storage.Action{
			Height: 1000,
		}
		err := parseFeeAssetChange(message, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("bridge lock", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		to := testsuite.RandomAddress()
		dest := testsuite.RandomAddress()

		message := &astria.Action_BridgeLock{
			BridgeLock: &astria.BridgeLock{
				FeeAsset:                feeAssetId,
				Asset:                   assetId,
				To:                      &primitivev1.Address{Bech32M: to},
				DestinationChainAddress: dest,
				Amount: &primitivev1.Uint128{
					Lo: 10,
					Hi: 0,
				},
			},
		}

		toModel := &storage.Address{
			Height:        1000,
			Hash:          to,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: currency.DefaultCurrency,
					Total:    decimal.RequireFromString("10"),
				},
			},
		}

		from := testsuite.RandomAddress()
		fromModel := &storage.Address{
			Height:        1000,
			Hash:          from,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: currency.DefaultCurrency,
					Total:    decimal.RequireFromString("-10"),
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeBridgeLock,
			Data: map[string]any{
				"asset":                     assetId,
				"fee_asset":                 feeAssetId,
				"to":                        to,
				"destination_chain_address": dest,
				"amount":                    "10",
			},
			Addresses: make([]*storage.AddressAction, 0),
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Address:  toModel,
					Update:   toModel.Balance[0].Total,
					Currency: toModel.Balance[0].Currency,
					Height:   1000,
				},
				{
					Address:  fromModel,
					Update:   fromModel.Balance[0].Total,
					Currency: fromModel.Balance[0].Currency,
					Height:   1000,
				},
			},
		}
		wantAction.Addresses = append(wantAction.Addresses,
			&storage.AddressAction{
				Height:     1000,
				Address:    toModel,
				ActionType: types.ActionTypeBridgeLock,
				Action:     &wantAction,
			},
			&storage.AddressAction{
				Height:     1000,
				Address:    fromModel,
				ActionType: types.ActionTypeBridgeLock,
				Action:     &wantAction,
			},
		)

		action := storage.Action{
			Height: 1000,
		}
		err := parseBridgeLock(message, from, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("bridge lock the same address", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		to := testsuite.RandomAddress()
		dest := testsuite.RandomAddress()

		message := &astria.Action_BridgeLock{
			BridgeLock: &astria.BridgeLock{
				FeeAsset:                feeAssetId,
				Asset:                   assetId,
				To:                      &primitivev1.Address{Bech32M: to},
				DestinationChainAddress: dest,
				Amount: &primitivev1.Uint128{
					Lo: 10,
					Hi: 0,
				},
			},
		}

		toModel := &storage.Address{
			Height:        1000,
			Hash:          to,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: currency.DefaultCurrency,
					Total:    decimal.RequireFromString("10"),
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeBridgeLock,
			Data: map[string]any{
				"asset":                     assetId,
				"fee_asset":                 feeAssetId,
				"to":                        to,
				"destination_chain_address": dest,
				"amount":                    "10",
			},
			Addresses: make([]*storage.AddressAction, 0),
		}
		wantAction.Addresses = append(wantAction.Addresses,
			&storage.AddressAction{
				Height:     1000,
				Address:    toModel,
				ActionType: types.ActionTypeBridgeLock,
				Action:     &wantAction,
			},
		)

		action := storage.Action{
			Height: 1000,
		}
		err := parseBridgeLock(message, to, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("bridge unlock", func(t *testing.T) {
		bridge := testsuite.RandomAddress()
		decodeContext := NewContext(map[string]string{
			bridge: assetId,
		})

		to := testsuite.RandomAddress()

		message := &astria.Action_BridgeUnlock{
			BridgeUnlock: &astria.BridgeUnlock{
				FeeAsset:      feeAssetId,
				To:            &primitivev1.Address{Bech32M: to},
				BridgeAddress: &primitivev1.Address{Bech32M: bridge},
				Memo:          "memo",
				Amount: &primitivev1.Uint128{
					Lo: 10,
					Hi: 0,
				},
				RollupBlockNumber:       101,
				RollupWithdrawalEventId: "event_id",
			},
		}

		toModel := &storage.Address{
			Height:        1000,
			Hash:          to,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: assetId,
					Total:    decimal.RequireFromString("10"),
				},
			},
		}
		fromModel := &storage.Address{
			Height:        1000,
			Hash:          bridge,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: assetId,
					Total:    decimal.RequireFromString("-10"),
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeBridgeUnlock,
			Data: map[string]any{
				"fee_asset":                  feeAssetId,
				"to":                         to,
				"bridge":                     bridge,
				"amount":                     "10",
				"memo":                       "memo",
				"rollup_block_number":        uint64(101),
				"rollup_withdrawal_event_id": "event_id",
			},
			Addresses: make([]*storage.AddressAction, 0),
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Address:  toModel,
					Update:   toModel.Balance[0].Total,
					Currency: assetId,
					Height:   1000,
				},
				{
					Address:  fromModel,
					Update:   fromModel.Balance[0].Total,
					Currency: assetId,
					Height:   1000,
				},
			},
		}
		wantAction.Addresses = append(wantAction.Addresses,
			&storage.AddressAction{
				Height:     1000,
				Address:    toModel,
				ActionType: types.ActionTypeBridgeUnlock,
				Action:     &wantAction,
			},
			&storage.AddressAction{
				Height:     1000,
				Address:    fromModel,
				ActionType: types.ActionTypeBridgeUnlock,
				Action:     &wantAction,
			},
		)

		action := storage.Action{
			Height: 1000,
		}
		err := parseBridgeUnlock(message, to, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("init bridge account", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		rollupId := testsuite.RandomHash(10)
		from := testsuite.RandomAddress()
		sudo := testsuite.RandomAddress()
		withdrawer := testsuite.RandomAddress()

		sudoAddr := decodeContext.Addresses.Set(sudo, 1000, decimal.Zero, "", 1, 0)
		wdwAddr := decodeContext.Addresses.Set(withdrawer, 1000, decimal.Zero, "", 1, 0)

		message := &astria.Action_InitBridgeAccount{
			InitBridgeAccount: &astria.InitBridgeAccount{
				RollupId:          &primitivev1.RollupId{Inner: rollupId},
				FeeAsset:          feeAssetId,
				Asset:             assetId,
				SudoAddress:       &primitivev1.Address{Bech32M: sudo},
				WithdrawerAddress: &primitivev1.Address{Bech32M: withdrawer},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeInitBridgeAccount,
			Data: map[string]any{
				"rollup_id":  rollupId,
				"asset":      assetId,
				"fee_asset":  feeAssetId,
				"sudo":       sudo,
				"withdrawer": withdrawer,
			},
			RollupAction: &storage.RollupAction{
				Height: 1000,
				Rollup: &storage.Rollup{
					AstriaId:     message.InitBridgeAccount.GetRollupId().GetInner(),
					FirstHeight:  1000,
					ActionsCount: 1,
					BridgeCount:  1,
				},
				ActionType: types.ActionTypeInitBridgeAccount,
			},
			Addresses: make([]*storage.AddressAction, 0),
		}
		wantAction.RollupAction.Action = &wantAction
		wantAction.Addresses = append(wantAction.Addresses,
			&storage.AddressAction{
				Address:    sudoAddr,
				Height:     1000,
				Time:       wantAction.Time,
				Action:     &wantAction,
				ActionType: types.ActionTypeInitBridgeAccount,
			}, &storage.AddressAction{
				Address:    wdwAddr,
				Height:     1000,
				Time:       wantAction.Time,
				Action:     &wantAction,
				ActionType: types.ActionTypeInitBridgeAccount,
			})

		action := storage.Action{
			Height: 1000,
		}
		err := parseInitBridgeAccount(message, from, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
		require.Len(t, decodeContext.bridgeAssets, 1)
		require.Contains(t, decodeContext.bridgeAssets, from)
	})

	t.Run("init bridge account: the same address", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		rollupId := testsuite.RandomHash(10)
		from := testsuite.RandomAddress()
		message := &astria.Action_InitBridgeAccount{
			InitBridgeAccount: &astria.InitBridgeAccount{
				RollupId:          &primitivev1.RollupId{Inner: rollupId},
				FeeAsset:          feeAssetId,
				Asset:             assetId,
				SudoAddress:       &primitivev1.Address{Bech32M: from},
				WithdrawerAddress: &primitivev1.Address{Bech32M: from},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeInitBridgeAccount,
			Data: map[string]any{
				"rollup_id":  rollupId,
				"asset":      assetId,
				"fee_asset":  feeAssetId,
				"sudo":       from,
				"withdrawer": from,
			},
			RollupAction: &storage.RollupAction{
				Height: 1000,
				Rollup: &storage.Rollup{
					AstriaId:     message.InitBridgeAccount.GetRollupId().GetInner(),
					FirstHeight:  1000,
					ActionsCount: 1,
					BridgeCount:  1,
				},
				ActionType: types.ActionTypeInitBridgeAccount,
			},
		}
		wantAction.RollupAction.Action = &wantAction

		action := storage.Action{
			Height: 1000,
		}
		err := parseInitBridgeAccount(message, from, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
		require.Len(t, decodeContext.bridgeAssets, 1)
		require.Contains(t, decodeContext.bridgeAssets, from)
	})

	t.Run("ibc relayer change: addition", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		address := testsuite.RandomAddress()
		message := &astria.Action_IbcRelayerChange{
			IbcRelayerChange: &astria.IbcRelayerChange{
				Value: &astria.IbcRelayerChange_Addition{
					Addition: &primitivev1.Address{Bech32M: address},
				},
			},
		}

		model := storage.Address{
			Height:        1000,
			Hash:          address,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: currency.DefaultCurrency,
					Total:    decimal.Zero,
				},
			},
			IsIbcRelayer: testsuite.Ptr(true),
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcRelayerChange,
			Data: map[string]any{
				"addition": address,
			},
			Addresses: make([]*storage.AddressAction, 0),
		}
		wantAction.Addresses = append(wantAction.Addresses, &storage.AddressAction{
			Height:     1000,
			Address:    &model,
			ActionType: types.ActionTypeIbcRelayerChange,
			Action:     &wantAction,
		})

		action := storage.Action{
			Height: 1000,
		}
		err := parseIbcRelayerChange(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("ibc relayer change: removal", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		address := testsuite.RandomAddress()
		message := &astria.Action_IbcRelayerChange{
			IbcRelayerChange: &astria.IbcRelayerChange{
				Value: &astria.IbcRelayerChange_Removal{
					Removal: &primitivev1.Address{Bech32M: address},
				},
			},
		}

		model := storage.Address{
			Height:        1000,
			Hash:          address,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: []*storage.Balance{
				{
					Currency: currency.DefaultCurrency,
					Total:    decimal.Zero,
				},
			},
			IsIbcRelayer: testsuite.Ptr(false),
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcRelayerChange,
			Data: map[string]any{
				"removal": address,
			},
			Addresses: make([]*storage.AddressAction, 0),
		}
		wantAction.Addresses = append(wantAction.Addresses, &storage.AddressAction{
			Height:     1000,
			Address:    &model,
			ActionType: types.ActionTypeIbcRelayerChange,
			Action:     &wantAction,
		})

		action := storage.Action{
			Height: 1000,
		}
		err := parseIbcRelayerChange(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: rollup_data_submission", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		message := &astria.Action_FeeChange{
			FeeChange: &astria.FeeChange{
				FeeComponents: &astria.FeeChange_RollupDataSubmission{
					RollupDataSubmission: &feesv1alpha1.RollupDataSubmissionFeeComponents{
						Base: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
						Multiplier: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"rollup_data_submission_base":       "10",
				"rollup_data_submission_multiplier": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: bridge_lock", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		message := &astria.Action_FeeChange{
			FeeChange: &astria.FeeChange{
				FeeComponents: &astria.FeeChange_BridgeLock{
					BridgeLock: &feesv1alpha1.BridgeLockFeeComponents{
						Base: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
						Multiplier: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"bridge_lock_base":       "10",
				"bridge_lock_multiplier": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: bridge_sudo_change", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		message := &astria.Action_FeeChange{
			FeeChange: &astria.FeeChange{
				FeeComponents: &astria.FeeChange_BridgeSudoChange{
					BridgeSudoChange: &feesv1alpha1.BridgeSudoChangeFeeComponents{
						Base: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
						Multiplier: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"bridge_sudo_change_base":       "10",
				"bridge_sudo_change_multiplier": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: ics20_withdrawal", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		message := &astria.Action_FeeChange{
			FeeChange: &astria.FeeChange{
				FeeComponents: &astria.FeeChange_Ics20Withdrawal{
					Ics20Withdrawal: &feesv1alpha1.Ics20WithdrawalFeeComponents{
						Base: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
						Multiplier: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"ics20_withdrawal_base":       "10",
				"ics20_withdrawal_multiplier": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: init_bridge_account", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		message := &astria.Action_FeeChange{
			FeeChange: &astria.FeeChange{
				FeeComponents: &astria.FeeChange_InitBridgeAccount{
					InitBridgeAccount: &feesv1alpha1.InitBridgeAccountFeeComponents{
						Base: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
						Multiplier: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"init_bridge_account_base":       "10",
				"init_bridge_account_multiplier": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: transfer", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		message := &astria.Action_FeeChange{
			FeeChange: &astria.FeeChange{
				FeeComponents: &astria.FeeChange_Transfer{
					Transfer: &feesv1alpha1.TransferFeeComponents{
						Base: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
						Multiplier: &primitivev1.Uint128{
							Hi: 0,
							Lo: 10,
						},
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"transfer_base":       "10",
				"transfer_multiplier": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("bridge sudo change", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})
		bridge := testsuite.RandomAddress()
		sudo := testsuite.RandomAddress()
		withdrawer := testsuite.RandomAddress()

		message := &astria.Action_BridgeSudoChange{
			BridgeSudoChange: &astria.BridgeSudoChange{
				FeeAsset:             feeAssetId,
				BridgeAddress:        &primitivev1.Address{Bech32M: bridge},
				NewWithdrawerAddress: &primitivev1.Address{Bech32M: withdrawer},
				NewSudoAddress:       &primitivev1.Address{Bech32M: sudo},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeBridgeSudoChangeAction,
			Data: map[string]any{
				"fee_asset":  feeAssetId,
				"withdrawer": withdrawer,
				"sudo":       sudo,
				"bridge":     bridge,
			},
			Height:    1000,
			Addresses: make([]*storage.AddressAction, 0),
		}

		wantAction.Addresses = append(wantAction.Addresses, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         bridge,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		}, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         sudo,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		}, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         withdrawer,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		})

		action := storage.Action{
			Height: 1000,
		}
		err := parseBridgeSudoChange(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("bridge sudo change: bridge is suor", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})
		bridge := testsuite.RandomAddress()
		sudo := bridge
		withdrawer := testsuite.RandomAddress()

		message := &astria.Action_BridgeSudoChange{
			BridgeSudoChange: &astria.BridgeSudoChange{
				FeeAsset:             feeAssetId,
				BridgeAddress:        &primitivev1.Address{Bech32M: bridge},
				NewWithdrawerAddress: &primitivev1.Address{Bech32M: withdrawer},
				NewSudoAddress:       &primitivev1.Address{Bech32M: sudo},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeBridgeSudoChangeAction,
			Data: map[string]any{
				"fee_asset":  feeAssetId,
				"withdrawer": withdrawer,
				"sudo":       sudo,
				"bridge":     bridge,
			},
			Height:    1000,
			Addresses: make([]*storage.AddressAction, 0),
		}

		wantAction.Addresses = append(wantAction.Addresses, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         bridge,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		}, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         withdrawer,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		})

		action := storage.Action{
			Height: 1000,
		}
		err := parseBridgeSudoChange(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("bridge sudo change: bridge is withdrawer", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})
		bridge := testsuite.RandomAddress()
		sudo := testsuite.RandomAddress()
		withdrawer := bridge

		message := &astria.Action_BridgeSudoChange{
			BridgeSudoChange: &astria.BridgeSudoChange{
				FeeAsset:             feeAssetId,
				BridgeAddress:        &primitivev1.Address{Bech32M: bridge},
				NewWithdrawerAddress: &primitivev1.Address{Bech32M: withdrawer},
				NewSudoAddress:       &primitivev1.Address{Bech32M: sudo},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeBridgeSudoChangeAction,
			Data: map[string]any{
				"fee_asset":  feeAssetId,
				"withdrawer": withdrawer,
				"sudo":       sudo,
				"bridge":     bridge,
			},
			Height:    1000,
			Addresses: make([]*storage.AddressAction, 0),
		}

		wantAction.Addresses = append(wantAction.Addresses, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         bridge,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		}, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         sudo,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		})

		action := storage.Action{
			Height: 1000,
		}
		err := parseBridgeSudoChange(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("bridge sudo change: sudo is withdrawer", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})
		bridge := testsuite.RandomAddress()
		sudo := testsuite.RandomAddress()
		withdrawer := sudo

		message := &astria.Action_BridgeSudoChange{
			BridgeSudoChange: &astria.BridgeSudoChange{
				FeeAsset:             feeAssetId,
				BridgeAddress:        &primitivev1.Address{Bech32M: bridge},
				NewWithdrawerAddress: &primitivev1.Address{Bech32M: withdrawer},
				NewSudoAddress:       &primitivev1.Address{Bech32M: sudo},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeBridgeSudoChangeAction,
			Data: map[string]any{
				"fee_asset":  feeAssetId,
				"withdrawer": withdrawer,
				"sudo":       sudo,
				"bridge":     bridge,
			},
			Height:    1000,
			Addresses: make([]*storage.AddressAction, 0),
		}

		wantAction.Addresses = append(wantAction.Addresses, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         bridge,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		}, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         sudo,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		})

		action := storage.Action{
			Height: 1000,
		}
		err := parseBridgeSudoChange(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("bridge sudo change: all equals", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})
		bridge := testsuite.RandomAddress()
		sudo := bridge
		withdrawer := bridge

		message := &astria.Action_BridgeSudoChange{
			BridgeSudoChange: &astria.BridgeSudoChange{
				FeeAsset:             feeAssetId,
				BridgeAddress:        &primitivev1.Address{Bech32M: bridge},
				NewWithdrawerAddress: &primitivev1.Address{Bech32M: withdrawer},
				NewSudoAddress:       &primitivev1.Address{Bech32M: sudo},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeBridgeSudoChangeAction,
			Data: map[string]any{
				"fee_asset":  feeAssetId,
				"withdrawer": withdrawer,
				"sudo":       sudo,
				"bridge":     bridge,
			},
			Height:    1000,
			Addresses: make([]*storage.AddressAction, 0),
		}

		wantAction.Addresses = append(wantAction.Addresses, &storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				Hash:         bridge,
				ActionsCount: 1,
				Balance: []*storage.Balance{
					{
						Currency: currency.DefaultCurrency,
						Total:    decimal.Zero,
					},
				},
			},
			ActionType: types.ActionTypeBridgeSudoChangeAction,
			Action:     &wantAction,
		})

		action := storage.Action{
			Height: 1000,
		}
		err := parseBridgeSudoChange(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("ibc sudo change", func(t *testing.T) {
		decodeContext := NewContext(map[string]string{})

		newAddress := testsuite.RandomAddress()
		message := &astria.Action_IbcSudoChange{
			IbcSudoChange: &astria.IbcSudoChange{
				NewAddress: &primitivev1.Address{Bech32M: newAddress},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcSudoChangeAction,
			Data: map[string]any{
				"address": newAddress,
			},
			Addresses: make([]*storage.AddressAction, 0),
		}
		balance := storage.EmptyBalance()
		addressAction := storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				ActionsCount: 1,
				Hash:         newAddress,
				Balance:      []*storage.Balance{&balance},
			},
			ActionType: types.ActionTypeIbcSudoChangeAction,
			Action:     &wantAction,
		}
		wantAction.Addresses = append(wantAction.Addresses, &addressAction)

		action := storage.Action{
			Height: 1000,
		}
		err := parseIbcSudoChangeAction(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})
}
