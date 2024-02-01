// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/hex"
	"testing"

	primitivev1 "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/primitive/v1"
	astria "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/sequencer/v1alpha1"
	"buf.build/gen/go/cosmos/cosmos-sdk/protocolbuffers/go/tendermint/abci"
	"buf.build/gen/go/cosmos/cosmos-sdk/protocolbuffers/go/tendermint/crypto"
	ibcv1alpha1 "buf.build/gen/go/penumbra-zone/penumbra/protocolbuffers/go/penumbra/core/component/ibc/v1alpha1"
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
			IbcAction: &ibcv1alpha1.IbcRelay{
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
			ActionsCount:  0,
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
}
