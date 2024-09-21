// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/hex"
	"testing"

	primitivev1 "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	astria "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
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
		message := &astria.Action_IbcAction{
			IbcAction: &v1.IbcRelay{
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
		err := parseIbcAction(message, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("ics 20 withdrawal", func(t *testing.T) {
		decodeContext := NewContext()

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
				"timeout_time": uint64(1000),
				"fee_asset":    feeAssetId,
				"memo":         "memo",
				"bridge":       from,
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
					Update:   decimal.RequireFromString("1"),
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

	t.Run("sequence", func(t *testing.T) {
		decodeContext := NewContext()

		from := testsuite.RandomAddress()
		addressModel := decodeContext.Addresses.Set(from, 1000, decimal.Zero, "", 0, 1)

		message := &astria.Action_SequenceAction{
			SequenceAction: &astria.SequenceAction{
				RollupId: &primitivev1.RollupId{Inner: testsuite.RandomHash(10)},
				Data:     testsuite.RandomHash(10),
				FeeAsset: feeAssetId,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeSequence,
			Data: map[string]any{
				"rollup_id": message.SequenceAction.GetRollupId().GetInner(),
				"data":      message.SequenceAction.GetData(),
				"fee_asset": feeAssetId,
			},
			Addresses: make([]*storage.AddressAction, 0),
			RollupAction: &storage.RollupAction{
				Size:   10,
				Height: 1000,
				Rollup: &storage.Rollup{
					AstriaId:     message.SequenceAction.GetRollupId().GetInner(),
					FirstHeight:  1000,
					ActionsCount: 1,
					Size:         10,
				},
				ActionType: types.ActionTypeSequence,
			},
		}
		wantAction.RollupAction.Action = &wantAction
		addressAction := storage.AddressAction{
			Height:     1000,
			Address:    addressModel,
			ActionType: types.ActionTypeSequence,
			Action:     &wantAction,
		}
		wantAction.Addresses = append(wantAction.Addresses, &addressAction)

		action := storage.Action{
			Height: 1000,
		}
		err := parseSequenceAction(message, from, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("sudo address change", func(t *testing.T) {
		decodeContext := NewContext()

		newAddress := testsuite.RandomAddress()
		message := &astria.Action_SudoAddressChangeAction{
			SudoAddressChangeAction: &astria.SudoAddressChangeAction{
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
		err := parseSudoAddressChangeAction(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("transfer", func(t *testing.T) {
		decodeContext := NewContext()

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

		message := &astria.Action_TransferAction{
			TransferAction: &astria.TransferAction{
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
				"asset":     message.TransferAction.GetAsset(),
				"fee_asset": message.TransferAction.GetFeeAsset(),
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
		decodeContext := NewContext()

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

		message := &astria.Action_TransferAction{
			TransferAction: &astria.TransferAction{
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
				"asset":     message.TransferAction.GetAsset(),
				"fee_asset": message.TransferAction.GetFeeAsset(),
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
		decodeContext := NewContext()
		message := &astria.Action_ValidatorUpdateAction{
			ValidatorUpdateAction: &abci.ValidatorUpdate{
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
				"pubkey": message.ValidatorUpdateAction.GetPubKey().GetEd25519(),
			},
			Addresses: make([]*storage.AddressAction, 0),
		}
		address, err := AddressFromPubKey(message.ValidatorUpdateAction.GetPubKey().GetEd25519())
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
		pk := hex.EncodeToString(message.ValidatorUpdateAction.GetPubKey().GetEd25519())
		v, ok := decodeContext.Validators[pk]
		require.True(t, ok)
		require.EqualValues(t, "10", v.Power.String())
	})

	t.Run("fee asset change: addition", func(t *testing.T) {
		message := &astria.Action_FeeAssetChangeAction{
			FeeAssetChangeAction: &astria.FeeAssetChangeAction{
				Value: &astria.FeeAssetChangeAction_Addition{
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
		message := &astria.Action_FeeAssetChangeAction{
			FeeAssetChangeAction: &astria.FeeAssetChangeAction{
				Value: &astria.FeeAssetChangeAction_Removal{
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
		decodeContext := NewContext()

		to := testsuite.RandomAddress()
		dest := testsuite.RandomAddress()

		message := &astria.Action_BridgeLockAction{
			BridgeLockAction: &astria.BridgeLockAction{
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
		decodeContext := NewContext()

		to := testsuite.RandomAddress()
		dest := testsuite.RandomAddress()

		message := &astria.Action_BridgeLockAction{
			BridgeLockAction: &astria.BridgeLockAction{
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
		decodeContext := NewContext()

		to := testsuite.RandomAddress()
		bridge := testsuite.RandomAddress()

		message := &astria.Action_BridgeUnlockAction{
			BridgeUnlockAction: &astria.BridgeUnlockAction{
				FeeAsset:      feeAssetId,
				To:            &primitivev1.Address{Bech32M: to},
				BridgeAddress: &primitivev1.Address{Bech32M: bridge},
				Memo:          []byte("memo"),
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
					Currency: feeAssetId,
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
					Currency: feeAssetId,
					Total:    decimal.RequireFromString("-10"),
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeBridgeUnlock,
			Data: map[string]any{
				"fee_asset": feeAssetId,
				"to":        to,
				"bridge":    bridge,
				"amount":    "10",
				"memo":      "6d656d6f",
			},
			Addresses: make([]*storage.AddressAction, 0),
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Address:  toModel,
					Update:   toModel.Balance[0].Total,
					Currency: feeAssetId,
					Height:   1000,
				},
				{
					Address:  fromModel,
					Update:   fromModel.Balance[0].Total,
					Currency: feeAssetId,
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
		decodeContext := NewContext()

		rollupId := testsuite.RandomHash(10)
		from := testsuite.RandomAddress()
		sudo := testsuite.RandomAddress()
		withdrawer := testsuite.RandomAddress()

		sudoAddr := decodeContext.Addresses.Set(sudo, 1000, decimal.Zero, "", 1, 0)
		wdwAddr := decodeContext.Addresses.Set(withdrawer, 1000, decimal.Zero, "", 1, 0)

		message := &astria.Action_InitBridgeAccountAction{
			InitBridgeAccountAction: &astria.InitBridgeAccountAction{
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
					AstriaId:     message.InitBridgeAccountAction.GetRollupId().GetInner(),
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
	})

	t.Run("init bridge account: the same address", func(t *testing.T) {
		decodeContext := NewContext()

		rollupId := testsuite.RandomHash(10)
		from := testsuite.RandomAddress()
		message := &astria.Action_InitBridgeAccountAction{
			InitBridgeAccountAction: &astria.InitBridgeAccountAction{
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
					AstriaId:     message.InitBridgeAccountAction.GetRollupId().GetInner(),
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
	})

	t.Run("ibc relayer change: addition", func(t *testing.T) {
		decodeContext := NewContext()

		address := testsuite.RandomAddress()
		message := &astria.Action_IbcRelayerChangeAction{
			IbcRelayerChangeAction: &astria.IbcRelayerChangeAction{
				Value: &astria.IbcRelayerChangeAction_Addition{
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
		decodeContext := NewContext()

		address := testsuite.RandomAddress()
		message := &astria.Action_IbcRelayerChangeAction{
			IbcRelayerChangeAction: &astria.IbcRelayerChangeAction{
				Value: &astria.IbcRelayerChangeAction_Removal{
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

	t.Run("fee change: sequence_base_fee", func(t *testing.T) {
		decodeContext := NewContext()

		message := &astria.Action_FeeChangeAction{
			FeeChangeAction: &astria.FeeChangeAction{
				Value: &astria.FeeChangeAction_SequenceBaseFee{
					SequenceBaseFee: &primitivev1.Uint128{
						Hi: 0,
						Lo: 10,
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"sequence_base_fee": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: bridge_lock_byte_cost_multiplier", func(t *testing.T) {
		decodeContext := NewContext()

		message := &astria.Action_FeeChangeAction{
			FeeChangeAction: &astria.FeeChangeAction{
				Value: &astria.FeeChangeAction_BridgeLockByteCostMultiplier{
					BridgeLockByteCostMultiplier: &primitivev1.Uint128{
						Hi: 0,
						Lo: 10,
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"bridge_lock_byte_cost_multiplier": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: bridge_sudo_change_base_fee", func(t *testing.T) {
		decodeContext := NewContext()

		message := &astria.Action_FeeChangeAction{
			FeeChangeAction: &astria.FeeChangeAction{
				Value: &astria.FeeChangeAction_BridgeSudoChangeBaseFee{
					BridgeSudoChangeBaseFee: &primitivev1.Uint128{
						Hi: 0,
						Lo: 10,
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"bridge_sudo_change_base_fee": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: ics20_withdrawal_base_fee", func(t *testing.T) {
		decodeContext := NewContext()

		message := &astria.Action_FeeChangeAction{
			FeeChangeAction: &astria.FeeChangeAction{
				Value: &astria.FeeChangeAction_Ics20WithdrawalBaseFee{
					Ics20WithdrawalBaseFee: &primitivev1.Uint128{
						Hi: 0,
						Lo: 10,
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"ics20_withdrawal_base_fee": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: init_bridge_account_base_fee", func(t *testing.T) {
		decodeContext := NewContext()

		message := &astria.Action_FeeChangeAction{
			FeeChangeAction: &astria.FeeChangeAction{
				Value: &astria.FeeChangeAction_InitBridgeAccountBaseFee{
					InitBridgeAccountBaseFee: &primitivev1.Uint128{
						Hi: 0,
						Lo: 10,
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"init_bridge_account_base_fee": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: sequence_byte_cost_multiplier", func(t *testing.T) {
		decodeContext := NewContext()

		message := &astria.Action_FeeChangeAction{
			FeeChangeAction: &astria.FeeChangeAction{
				Value: &astria.FeeChangeAction_SequenceByteCostMultiplier{
					SequenceByteCostMultiplier: &primitivev1.Uint128{
						Hi: 0,
						Lo: 10,
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"sequence_byte_cost_multiplier": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee change: transfer_base_fee", func(t *testing.T) {
		decodeContext := NewContext()

		message := &astria.Action_FeeChangeAction{
			FeeChangeAction: &astria.FeeChangeAction{
				Value: &astria.FeeChangeAction_TransferBaseFee{
					TransferBaseFee: &primitivev1.Uint128{
						Hi: 0,
						Lo: 10,
					},
				},
			},
		}

		wantAction := storage.Action{
			Type: types.ActionTypeFeeChange,
			Data: map[string]any{
				"transfer_base_fee": "10",
			},
		}

		action := storage.Action{}
		err := parseFeeChange(message, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("bridge sudo change", func(t *testing.T) {
		decodeContext := NewContext()
		bridge := testsuite.RandomAddress()
		sudo := testsuite.RandomAddress()
		withdrawer := testsuite.RandomAddress()

		message := &astria.Action_BridgeSudoChangeAction{
			BridgeSudoChangeAction: &astria.BridgeSudoChangeAction{
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
		decodeContext := NewContext()
		bridge := testsuite.RandomAddress()
		sudo := bridge
		withdrawer := testsuite.RandomAddress()

		message := &astria.Action_BridgeSudoChangeAction{
			BridgeSudoChangeAction: &astria.BridgeSudoChangeAction{
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
		decodeContext := NewContext()
		bridge := testsuite.RandomAddress()
		sudo := testsuite.RandomAddress()
		withdrawer := bridge

		message := &astria.Action_BridgeSudoChangeAction{
			BridgeSudoChangeAction: &astria.BridgeSudoChangeAction{
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
		decodeContext := NewContext()
		bridge := testsuite.RandomAddress()
		sudo := testsuite.RandomAddress()
		withdrawer := sudo

		message := &astria.Action_BridgeSudoChangeAction{
			BridgeSudoChangeAction: &astria.BridgeSudoChangeAction{
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
		decodeContext := NewContext()
		bridge := testsuite.RandomAddress()
		sudo := bridge
		withdrawer := bridge

		message := &astria.Action_BridgeSudoChangeAction{
			BridgeSudoChangeAction: &astria.BridgeSudoChangeAction{
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
}
