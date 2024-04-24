// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/hex"
	"testing"

	primitivev1 "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/primitive/v1"
	astria "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/sequencer/v1alpha1"
	v1 "buf.build/gen/go/astria/astria/protocolbuffers/go/astria_vendored/penumbra/core/component/ibc/v1"
	abci "buf.build/gen/go/astria/astria/protocolbuffers/go/astria_vendored/tendermint/abci"
	crypto "buf.build/gen/go/astria/astria/protocolbuffers/go/astria_vendored/tendermint/crypto"
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/astria-indexer/internal/test_suite"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestDecodeActions(t *testing.T) {

	t.Run("ibc action", func(t *testing.T) {
		message := &astria.Action_IbcAction{
			IbcAction: &v1.IbcRelay{
				RawAction: &anypb.Any{
					Value: []byte{0, 0, 0, 0},
				},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcRelay,
			Data: map[string]any{
				"raw": "AAAAAA==",
			},
		}

		action := storage.Action{
			Height: 1000,
		}
		err := parseIbcAction(message, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("ibc 20 withdrawal", func(t *testing.T) {
		decodeContext := NewContext()

		address := testsuite.RandomHash(20)

		message := &astria.Action_Ics20Withdrawal{
			Ics20Withdrawal: &astria.Ics20Withdrawal{
				Amount: &primitivev1.Uint128{
					Lo: 1,
					Hi: 0,
				},
				Denom:                   currency.DefaultCurrency,
				DestinationChainAddress: "celestia1lx7dfjp20shd6y5f4tauvy8cv4pjhvszfrh9ah",
				ReturnAddress:           address,
				SourceChannel:           "channel-12",
				TimeoutHeight: &astria.IbcHeight{
					RevisionNumber: 1000,
					RevisionHeight: 1001,
				},
				TimeoutTime: 1000,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIcs20Withdrawal,
			Data: map[string]any{
				"amount":              "1",
				"denom":               currency.DefaultCurrency,
				"destination_address": "celestia1lx7dfjp20shd6y5f4tauvy8cv4pjhvszfrh9ah",
				"return_address":      hex.EncodeToString(address),
				"source_channel":      "channel-12",
				"timeout_height": map[string]any{
					"revision_number": uint64(1000),
					"revision_height": uint64(1001),
				},
				"timeout_time": uint64(1000),
			},
			Addresses: []*storage.AddressAction{},
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Height: 1000,
					Address: &storage.Address{
						Height:       1000,
						Hash:         address,
						ActionsCount: 1,
						Balance: &storage.Balance{
							Currency: currency.DefaultCurrency,
							Total:    decimal.RequireFromString("1"),
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
				Balance: &storage.Balance{
					Currency: currency.DefaultCurrency,
					Total:    decimal.RequireFromString("1"),
				},
			},
			ActionType: types.ActionTypeIcs20Withdrawal,
			Action:     &wantAction,
		}
		wantAction.Addresses = append(wantAction.Addresses, &addressAction)

		action := storage.Action{
			Height: 1000,
		}
		err := parseIcs20Withdrawal(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("mint", func(t *testing.T) {
		decodeContext := NewContext()

		address := testsuite.RandomHash(20)
		addressModel := &storage.Address{
			Height:       1000,
			Hash:         address,
			ActionsCount: 1,
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.RequireFromString("10"),
			},
		}

		message := &astria.Action_MintAction{
			MintAction: &astria.MintAction{
				Amount: &primitivev1.Uint128{
					Lo: 10,
				},
				To: address,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeMint,
			Data: map[string]any{
				"amount": "10",
				"to":     hex.EncodeToString(address),
			},
			Addresses: make([]*storage.AddressAction, 0),
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Address:  addressModel,
					Currency: currency.DefaultCurrency,
					Update:   addressModel.Balance.Total,
					Height:   1000,
				},
			},
		}
		addressAction := storage.AddressAction{
			Height:     1000,
			Address:    addressModel,
			ActionType: types.ActionTypeMint,
			Action:     &wantAction,
		}
		wantAction.Addresses = append(wantAction.Addresses, &addressAction)

		action := storage.Action{
			Height: 1000,
		}
		err := parseMintAction(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("sequence", func(t *testing.T) {
		decodeContext := NewContext()

		from := testsuite.RandomHash(20)
		decodeContext.Addresses.Set(from, 1000, decimal.Zero, 0, 1)

		addressModel := &storage.Address{
			Height:        1000,
			Hash:          from,
			ActionsCount:  1,
			SignedTxCount: 1,
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.Zero,
			},
		}

		message := &astria.Action_SequenceAction{
			SequenceAction: &astria.SequenceAction{
				RollupId: testsuite.RandomHash(10),
				Data:     testsuite.RandomHash(10),
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeSequence,
			Data: map[string]any{
				"rollup_id": message.SequenceAction.RollupId,
				"data":      message.SequenceAction.Data,
			},
			Addresses: make([]*storage.AddressAction, 0),
			RollupAction: &storage.RollupAction{
				Size:   10,
				Height: 1000,
				Rollup: &storage.Rollup{
					AstriaId:     message.SequenceAction.RollupId,
					FirstHeight:  1000,
					ActionsCount: 1,
					Size:         10,
				},
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

		newAddress := testsuite.RandomHash(20)
		message := &astria.Action_SudoAddressChangeAction{
			SudoAddressChangeAction: &astria.SudoAddressChangeAction{
				NewAddress: newAddress,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeSudoAddressChange,
			Data: map[string]any{
				"address": hex.EncodeToString(newAddress),
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
				Balance:      &balance,
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

		from := testsuite.RandomHash(20)
		fromModel := &storage.Address{
			Height:        1000,
			Hash:          from,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.RequireFromString("-10"),
			},
		}

		to := testsuite.RandomHash(20)
		toModel := &storage.Address{
			Height:        1000,
			Hash:          to,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.RequireFromString("10"),
			},
		}

		message := &astria.Action_TransferAction{
			TransferAction: &astria.TransferAction{
				To: to,
				Amount: &primitivev1.Uint128{
					Lo: 10,
				},
				AssetId: testsuite.RandomHash(20),
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeTransfer,
			Data: map[string]any{
				"to":       hex.EncodeToString(to),
				"amount":   "10",
				"asset_id": message.TransferAction.AssetId,
			},
			Addresses: make([]*storage.AddressAction, 0),
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Address:  toModel,
					Update:   toModel.Balance.Total,
					Currency: toModel.Balance.Currency,
					Height:   1000,
				}, {
					Address:  fromModel,
					Update:   fromModel.Balance.Total,
					Currency: fromModel.Balance.Currency,
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

		from := testsuite.RandomHash(20)
		fromModel := &storage.Address{
			Height:        1000,
			Hash:          from,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.Zero,
			},
		}

		message := &astria.Action_TransferAction{
			TransferAction: &astria.TransferAction{
				To: from,
				Amount: &primitivev1.Uint128{
					Lo: 10,
				},
				AssetId: testsuite.RandomHash(20),
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeTransfer,
			Data: map[string]any{
				"to":       hex.EncodeToString(from),
				"amount":   "10",
				"asset_id": message.TransferAction.AssetId,
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
				"pubkey": message.ValidatorUpdateAction.PubKey.GetEd25519(),
			},
			Addresses: make([]*storage.AddressAction, 0),
		}
		address := AddressFromPubKey(message.ValidatorUpdateAction.PubKey.GetEd25519())
		balance := storage.EmptyBalance()
		addressAction := storage.AddressAction{
			Height: 1000,
			Address: &storage.Address{
				Height:       1000,
				ActionsCount: 1,
				Hash:         address,
				Balance:      &balance,
			},
			ActionType: types.ActionTypeValidatorUpdate,
			Action:     &wantAction,
		}
		wantAction.Addresses = append(wantAction.Addresses, &addressAction)

		action := storage.Action{
			Height: 1000,
		}
		err := parseValidatorUpdateAction(message, 1000, &decodeContext, &action)
		require.NoError(t, err)
		require.Equal(t, wantAction, action)
	})

	t.Run("fee asset change: addition", func(t *testing.T) {
		assetId := testsuite.RandomHash(32)
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
		assetId := testsuite.RandomHash(32)
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

		assetId := testsuite.RandomHash(32)
		feeAssetId := testsuite.RandomHash(32)
		to := testsuite.RandomHash(20)
		dest := "random_address"

		message := &astria.Action_BridgeLockAction{
			BridgeLockAction: &astria.BridgeLockAction{
				FeeAssetId:              feeAssetId,
				AssetId:                 assetId,
				To:                      to,
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
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.RequireFromString("10"),
			},
		}

		from := testsuite.RandomHash(20)
		fromModel := &storage.Address{
			Height:        1000,
			Hash:          from,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.RequireFromString("-10"),
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeBridgeLock,
			Data: map[string]any{
				"asset_id":                  assetId,
				"fee_asset_id":              feeAssetId,
				"to":                        hex.EncodeToString(to),
				"destination_chain_address": dest,
				"amount":                    "10",
			},
			Addresses: make([]*storage.AddressAction, 0),
			BalanceUpdates: []storage.BalanceUpdate{
				{
					Address:  toModel,
					Update:   toModel.Balance.Total,
					Currency: toModel.Balance.Currency,
					Height:   1000,
				},
				{
					Address:  fromModel,
					Update:   fromModel.Balance.Total,
					Currency: fromModel.Balance.Currency,
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

		assetId := testsuite.RandomHash(32)
		feeAssetId := testsuite.RandomHash(32)
		to := testsuite.RandomHash(20)
		dest := "random_address"

		message := &astria.Action_BridgeLockAction{
			BridgeLockAction: &astria.BridgeLockAction{
				FeeAssetId:              feeAssetId,
				AssetId:                 assetId,
				To:                      to,
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
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.RequireFromString("10"),
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeBridgeLock,
			Data: map[string]any{
				"asset_id":                  assetId,
				"fee_asset_id":              feeAssetId,
				"to":                        hex.EncodeToString(to),
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

	t.Run("init bridge account", func(t *testing.T) {
		decodeContext := NewContext()

		assetId := testsuite.RandomHash(32)
		rollupId := testsuite.RandomHash(10)
		feAssetId := testsuite.RandomHash(32)
		from := testsuite.RandomHash(20)
		fromAddr := decodeContext.Addresses.Set(from, 1000, decimal.Zero, 0, 1)

		message := &astria.Action_InitBridgeAccountAction{
			InitBridgeAccountAction: &astria.InitBridgeAccountAction{
				RollupId:   rollupId,
				FeeAssetId: feAssetId,
				AssetIds:   [][]byte{assetId},
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeInitBridgeAccount,
			Data: map[string]any{
				"rollup_id":    rollupId,
				"asset_ids":    [][]byte{assetId},
				"fee_asset_id": feAssetId,
			},
			RollupAction: &storage.RollupAction{
				Height: 1000,
				Rollup: &storage.Rollup{
					AstriaId:      message.InitBridgeAccountAction.RollupId,
					FirstHeight:   1000,
					ActionsCount:  1,
					BridgeAddress: fromAddr,
				},
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

		address := testsuite.RandomHash(32)
		message := &astria.Action_IbcRelayerChangeAction{
			IbcRelayerChangeAction: &astria.IbcRelayerChangeAction{
				Value: &astria.IbcRelayerChangeAction_Addition{
					Addition: address,
				},
			},
		}

		model := storage.Address{
			Height:        1000,
			Hash:          address,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.Zero,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcRelayerChange,
			Data: map[string]any{
				"addition": hex.EncodeToString(address),
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

		address := testsuite.RandomHash(20)
		message := &astria.Action_IbcRelayerChangeAction{
			IbcRelayerChangeAction: &astria.IbcRelayerChangeAction{
				Value: &astria.IbcRelayerChangeAction_Removal{
					Removal: address,
				},
			},
		}

		model := storage.Address{
			Height:        1000,
			Hash:          address,
			ActionsCount:  1,
			SignedTxCount: 0,
			Balance: &storage.Balance{
				Currency: currency.DefaultCurrency,
				Total:    decimal.Zero,
			},
		}

		wantAction := storage.Action{
			Height: 1000,
			Type:   types.ActionTypeIbcRelayerChange,
			Data: map[string]any{
				"removal": hex.EncodeToString(address),
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
}
