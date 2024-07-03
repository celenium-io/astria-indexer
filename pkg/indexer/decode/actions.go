// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/base64"
	"time"

	astria "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func parseActions(height types.Level, blockTime time.Time, from string, tx *DecodedTx, ctx *Context) ([]storage.Action, error) {
	rawActions := tx.UnsignedTx.GetActions()
	actions := make([]storage.Action, len(rawActions))
	for i := range rawActions {
		if tx.UnsignedTx.Actions[i].Value == nil {
			return nil, errors.Errorf("nil action")
		}
		actions[i].Height = height
		actions[i].Time = blockTime
		actions[i].Position = int64(i)
		actions[i].Addresses = make([]*storage.AddressAction, 0)
		actions[i].BalanceUpdates = make([]storage.BalanceUpdate, 0)

		var err error

		switch val := rawActions[i].GetValue().(type) {
		case *astria.Action_IbcAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeIbcRelayBits)
			err = parseIbcAction(val, &actions[i])
		case *astria.Action_Ics20Withdrawal:
			tx.ActionTypes.Set(storageTypes.ActionTypeIcs20WithdrawalBits)
			err = parseIcs20Withdrawal(val, height, ctx, &actions[i])
		case *astria.Action_SequenceAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeSequenceBits)
			err = parseSequenceAction(val, from, height, ctx, &actions[i])
		case *astria.Action_SudoAddressChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeSudoAddressChangeBits)
			err = parseSudoAddressChangeAction(val, height, ctx, &actions[i])
		case *astria.Action_TransferAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeTransferBits)
			err = parseTransferAction(val, from, height, ctx, &actions[i])
		case *astria.Action_ValidatorUpdateAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeValidatorUpdateBits)
			err = parseValidatorUpdateAction(val, height, ctx, &actions[i])
		case *astria.Action_BridgeLockAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeBridgeLockBits)
			err = parseBridgeLock(val, from, height, ctx, &actions[i])
		case *astria.Action_FeeAssetChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeFeeAssetChangeBits)
			err = parseFeeAssetChange(val, &actions[i])
		case *astria.Action_IbcRelayerChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeIbcRelayerChangeBits)
			err = parseIbcRelayerChange(val, height, ctx, &actions[i])
		case *astria.Action_InitBridgeAccountAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeInitBridgeAccountBits)
			err = parseInitBridgeAccount(val, from, height, ctx, &actions[i])
		case *astria.Action_BridgeSudoChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeBridgeSudoChangeBits)
			err = parseBridgeSudoChange(val, height, ctx, &actions[i])
		case *astria.Action_BridgeUnlockAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeBridgeUnlockBits)
			err = parseBridgeUnlock(val, from, height, ctx, &actions[i])
		case *astria.Action_FeeChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeFeeChangeBits)
			err = parseFeeChange(val, &actions[i])

		default:
			return nil, errors.Errorf(
				"unknown action type | position = %d | block = %d: %##v",
				i, height, tx.Actions[i])
		}

		if err != nil {
			return nil, err
		}
	}

	return actions, nil
}

func parseIbcAction(body *astria.Action_IbcAction, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIbcRelay
	action.Data = make(map[string]any)

	if body.IbcAction != nil && body.IbcAction.GetRawAction() != nil {
		action.Data["raw"] = base64.StdEncoding.EncodeToString(body.IbcAction.GetRawAction().GetValue())
	}
	return nil
}

func parseIcs20Withdrawal(body *astria.Action_Ics20Withdrawal, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIcs20Withdrawal
	action.Data = make(map[string]any)
	if body.Ics20Withdrawal != nil {
		amount := uint128ToString(body.Ics20Withdrawal.GetAmount())
		action.Data["amount"] = amount
		action.Data["denom"] = body.Ics20Withdrawal.GetDenom()
		action.Data["fee_asset"] = body.Ics20Withdrawal.GetFeeAsset()
		action.Data["destination_address"] = body.Ics20Withdrawal.GetDestinationChainAddress()
		action.Data["return_address"] = body.Ics20Withdrawal.GetReturnAddress().GetBech32M()
		action.Data["source_channel"] = body.Ics20Withdrawal.GetSourceChannel()
		if memo := body.Ics20Withdrawal.GetMemo(); memo != "" {
			action.Data["memo"] = memo
		}

		if th := body.Ics20Withdrawal.GetTimeoutHeight(); th != nil {
			action.Data["timeout_height"] = map[string]any{
				"revision_number": th.GetRevisionNumber(),
				"revision_height": th.GetRevisionHeight(),
			}
		}
		if body.Ics20Withdrawal.GetTimeoutTime() > 0 {
			action.Data["timeout_time"] = body.Ics20Withdrawal.GetTimeoutTime()
		}
		if bridge := body.Ics20Withdrawal.GetBridgeAddress().GetBech32M(); bridge != "" {
			action.Data["bridge"] = bridge
			addr := ctx.Addresses.Set(bridge, height, decimal.Zero, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
		}

		decAmount := decimal.RequireFromString(amount)
		returnAddress := body.Ics20Withdrawal.GetReturnAddress().GetBech32M()
		addr := ctx.Addresses.Set(returnAddress, height, decAmount, 1, 0)
		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    addr,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})

		action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
			Address:  addr,
			Height:   action.Height,
			Currency: body.Ics20Withdrawal.GetDenom(),
			Update:   decAmount,
		})
	}
	return nil
}

func parseSequenceAction(body *astria.Action_SequenceAction, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeSequence
	action.Data = make(map[string]any)
	if body.SequenceAction != nil {
		rollupId := body.SequenceAction.GetRollupId().GetInner()
		action.Data["rollup_id"] = rollupId
		action.Data["data"] = body.SequenceAction.GetData()
		action.Data["fee_asset"] = body.SequenceAction.GetFeeAsset()
		dataSize := len(body.SequenceAction.GetData())

		rollup := ctx.Rollups.Set(rollupId, height, dataSize)
		fromAddress := ctx.Addresses.Set(from, height, decimal.Zero, 1, 0)

		rollupAddress := &storage.RollupAddress{
			Rollup:  rollup,
			Address: fromAddress,
			Height:  height,
		}
		key := rollupAddress.String()
		if key == "" {
			return errors.Errorf("empty rollup address key")
		}
		if _, ok := ctx.RollupAddress[key]; !ok {
			ctx.RollupAddress[key] = rollupAddress
		}

		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    fromAddress,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})

		action.RollupAction = &storage.RollupAction{
			Time:   action.Time,
			Height: action.Height,
			Size:   int64(dataSize),
			Action: action,
			Rollup: rollup,
		}
		ctx.DataSize += int64(dataSize)
	}
	return nil
}

func parseSudoAddressChangeAction(body *astria.Action_SudoAddressChangeAction, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeSudoAddressChange
	action.Data = make(map[string]any)
	if body.SudoAddressChangeAction != nil {
		address := body.SudoAddressChangeAction.GetNewAddress().GetBech32M()
		action.Data["address"] = address

		addr := ctx.Addresses.Set(address, height, decimal.Zero, 1, 0)
		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    addr,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})
	}
	return nil
}

func parseTransferAction(body *astria.Action_TransferAction, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeTransfer
	action.Data = make(map[string]any)
	if body.TransferAction != nil {
		amount := uint128ToString(body.TransferAction.GetAmount())
		to := body.TransferAction.GetTo().GetBech32M()
		asset := body.TransferAction.GetAsset()

		action.Data["amount"] = amount
		action.Data["asset"] = asset
		action.Data["fee_asset"] = body.TransferAction.GetFeeAsset()
		action.Data["to"] = to

		decAmount := decimal.RequireFromString(amount)

		if from == to {
			addr := ctx.Addresses.Set(from, height, decimal.Zero, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
		} else {
			toAddr := ctx.Addresses.Set(to, height, decAmount, 1, 0)
			fromAddr := ctx.Addresses.Set(from, height, decAmount.Neg(), 1, 0)
			action.Addresses = append(action.Addresses,
				&storage.AddressAction{
					Address:    toAddr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				},
				&storage.AddressAction{
					Address:    fromAddr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				})

			action.BalanceUpdates = append(action.BalanceUpdates,
				storage.BalanceUpdate{
					Address:  toAddr,
					Height:   action.Height,
					Currency: asset,
					Update:   decAmount,
				},
				storage.BalanceUpdate{
					Address:  fromAddr,
					Height:   action.Height,
					Currency: asset,
					Update:   decAmount.Copy().Neg(),
				})
		}
	}
	return nil
}

func parseValidatorUpdateAction(body *astria.Action_ValidatorUpdateAction, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeValidatorUpdate
	action.Data = make(map[string]any)
	if body.ValidatorUpdateAction != nil {
		power := body.ValidatorUpdateAction.GetPower()
		action.Data["power"] = power
		pubKey := body.ValidatorUpdateAction.GetPubKey().GetEd25519()
		action.Data["pubkey"] = pubKey

		address, err := AddressFromPubKey(body.ValidatorUpdateAction.GetPubKey().GetEd25519())
		if err != nil {
			return err
		}
		addr := ctx.Addresses.Set(address, height, decimal.Zero, 1, 0)
		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    addr,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})
		ctx.Validators.Set(pubKey, power)
	}
	return nil
}

func parseIbcRelayerChange(body *astria.Action_IbcRelayerChangeAction, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIbcRelayerChange
	action.Data = make(map[string]any)
	if body.IbcRelayerChangeAction != nil {
		if addition := body.IbcRelayerChangeAction.GetAddition(); len(addition.GetBech32M()) > 0 {
			action.Data["addition"] = addition.GetBech32M()

			addr := ctx.Addresses.Set(addition.GetBech32M(), height, decimal.Zero, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
		}

		if removal := body.IbcRelayerChangeAction.GetRemoval(); len(removal.GetBech32M()) > 0 {
			action.Data["removal"] = removal.GetBech32M()

			addr := ctx.Addresses.Set(removal.GetBech32M(), height, decimal.Zero, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
		}
	}
	return nil
}

func parseInitBridgeAccount(body *astria.Action_InitBridgeAccountAction, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeInitBridgeAccount
	action.Data = make(map[string]any)
	if body.InitBridgeAccountAction != nil {
		action.Data["rollup_id"] = body.InitBridgeAccountAction.GetRollupId().GetInner()
		action.Data["fee_asset"] = body.InitBridgeAccountAction.GetFeeAsset()
		action.Data["asset"] = body.InitBridgeAccountAction.GetAsset()

		if sudo := body.InitBridgeAccountAction.GetSudoAddress().GetBech32M(); sudo != "" {
			action.Data["sudo"] = sudo
			if sudo != from {
				addr := ctx.Addresses.Set(sudo, height, decimal.Zero, 1, 0)
				action.Addresses = append(action.Addresses, &storage.AddressAction{
					Address:    addr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				})
			}
		}

		if withdrawer := body.InitBridgeAccountAction.GetWithdrawerAddress().GetBech32M(); withdrawer != "" {
			action.Data["withdrawer"] = withdrawer
			if withdrawer != from {
				addr := ctx.Addresses.Set(withdrawer, height, decimal.Zero, 1, 0)
				action.Addresses = append(action.Addresses, &storage.AddressAction{
					Address:    addr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				})
			}
		}

		rollup := ctx.Rollups.Set(body.InitBridgeAccountAction.GetRollupId().GetInner(), height, 0)
		action.RollupAction = &storage.RollupAction{
			Time:   action.Time,
			Height: action.Height,
			Action: action,
			Rollup: rollup,
		}

		fromAddress := ctx.Addresses.Set(from, height, decimal.Zero, 1, 0)
		rollup.BridgeAddress = fromAddress
	}
	return nil
}

func parseBridgeSudoChange(body *astria.Action_BridgeSudoChangeAction, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeBridgeSudoChangeAction
	action.Data = make(map[string]any)
	if body.BridgeSudoChangeAction != nil {
		bridge := body.BridgeSudoChangeAction.GetBridgeAddress().GetBech32M()
		sudo := body.BridgeSudoChangeAction.GetNewSudoAddress().GetBech32M()
		withdrawer := body.BridgeSudoChangeAction.GetNewWithdrawerAddress().GetBech32M()

		action.Data["bridge"] = bridge
		bridgeAddr := ctx.Addresses.Set(bridge, height, decimal.Zero, 1, 0)
		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    bridgeAddr,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})

		if sudo != "" {
			action.Data["sudo"] = sudo
			addr := ctx.Addresses.Set(sudo, height, decimal.Zero, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
		}
		if withdrawer != "" {
			action.Data["withdrawer"] = withdrawer
			addr := ctx.Addresses.Set(withdrawer, height, decimal.Zero, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
		}
		action.Data["fee_asset"] = body.BridgeSudoChangeAction.GetFeeAsset()

	}
	return nil
}

func parseBridgeLock(body *astria.Action_BridgeLockAction, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeBridgeLock
	action.Data = make(map[string]any)
	if body.BridgeLockAction != nil {
		amount := uint128ToString(body.BridgeLockAction.GetAmount())
		asset := body.BridgeLockAction.GetAsset()

		action.Data["to"] = body.BridgeLockAction.GetTo().GetBech32M()
		action.Data["destination_chain_address"] = body.BridgeLockAction.GetDestinationChainAddress()
		action.Data["asset"] = asset
		action.Data["fee_asset"] = body.BridgeLockAction.GetFeeAsset()
		action.Data["amount"] = amount

		toAddress := body.BridgeLockAction.GetTo().GetBech32M()
		decAmount := decimal.RequireFromString(amount)
		toAddr := ctx.Addresses.Set(toAddress, height, decAmount, 1, 0)

		if from == toAddress {
			action.Addresses = append(action.Addresses,
				&storage.AddressAction{
					Address:    toAddr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				},
			)
		} else {
			fromAddr := ctx.Addresses.Set(from, height, decAmount.Neg(), 1, 0)

			action.Addresses = append(action.Addresses,
				&storage.AddressAction{
					Address:    toAddr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				},
				&storage.AddressAction{
					Address:    fromAddr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				},
			)

			action.BalanceUpdates = append(action.BalanceUpdates,
				storage.BalanceUpdate{
					Address:  toAddr,
					Height:   action.Height,
					Currency: asset,
					Update:   decAmount,
				},
				storage.BalanceUpdate{
					Address:  fromAddr,
					Height:   action.Height,
					Currency: asset,
					Update:   decAmount.Neg(),
				},
			)

		}
	}
	return nil
}

func parseBridgeUnlock(body *astria.Action_BridgeUnlockAction, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeBridgeUnlock
	action.Data = make(map[string]any)
	if body.BridgeUnlockAction != nil {
		amount := uint128ToString(body.BridgeUnlockAction.GetAmount())
		toAddress := body.BridgeUnlockAction.GetTo().GetBech32M()
		bridge := body.BridgeUnlockAction.GetBridgeAddress().GetBech32M()

		action.Data["to"] = toAddress
		action.Data["fee_asset"] = body.BridgeUnlockAction.GetFeeAsset()
		action.Data["amount"] = amount
		if memo := body.BridgeUnlockAction.GetMemo(); len(memo) > 0 {
			action.Data["memo"] = string(memo)
		}
		if bridge != "" {
			action.Data["bridge"] = bridge
		}

		decAmount := decimal.RequireFromString(amount)
		toAddr := ctx.Addresses.Set(toAddress, height, decAmount, 1, 0)

		var fromAddr *storage.Address
		switch bridge {
		case "", from:
			fromAddr = ctx.Addresses.Set(from, height, decAmount.Neg(), 1, 0)
		default:
			fromAddr = ctx.Addresses.Set(bridge, height, decAmount.Neg(), 1, 0)
		}

		action.Addresses = append(action.Addresses,
			&storage.AddressAction{
				Address:    toAddr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			},
			&storage.AddressAction{
				Address:    fromAddr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			},
		)

		action.BalanceUpdates = append(action.BalanceUpdates,
			storage.BalanceUpdate{
				Address:  toAddr,
				Height:   action.Height,
				Currency: toAddr.Balance.Currency,
				Update:   decAmount,
			},
			storage.BalanceUpdate{
				Address:  fromAddr,
				Height:   action.Height,
				Currency: fromAddr.Balance.Currency,
				Update:   decAmount.Neg(),
			},
		)
	}
	return nil
}

func parseFeeAssetChange(body *astria.Action_FeeAssetChangeAction, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeFeeAssetChange
	action.Data = make(map[string]any)
	if body.FeeAssetChangeAction != nil {
		if addition := body.FeeAssetChangeAction.GetAddition(); len(addition) > 0 {
			action.Data["addition"] = addition
		}

		if removal := body.FeeAssetChangeAction.GetRemoval(); len(removal) > 0 {
			action.Data["removal"] = removal
		}
	}
	return nil
}

func parseFeeChange(body *astria.Action_FeeChangeAction, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeFeeChange
	action.Data = make(map[string]any)
	if body.FeeChangeAction != nil {
		switch t := body.FeeChangeAction.GetValue().(type) {
		case *astria.FeeChangeAction_BridgeLockByteCostMultiplier:
			action.Data["bridge_lock_byte_cost_multiplier"] = uint128ToString(t.BridgeLockByteCostMultiplier)

		case *astria.FeeChangeAction_BridgeSudoChangeBaseFee:
			action.Data["bridge_sudo_change_base_fee"] = uint128ToString(t.BridgeSudoChangeBaseFee)

		case *astria.FeeChangeAction_Ics20WithdrawalBaseFee:
			action.Data["ics20_withdrawal_base_fee"] = uint128ToString(t.Ics20WithdrawalBaseFee)

		case *astria.FeeChangeAction_InitBridgeAccountBaseFee:
			action.Data["init_bridge_account_base_fee"] = uint128ToString(t.InitBridgeAccountBaseFee)

		case *astria.FeeChangeAction_SequenceBaseFee:
			action.Data["sequence_base_fee"] = uint128ToString(t.SequenceBaseFee)

		case *astria.FeeChangeAction_SequenceByteCostMultiplier:
			action.Data["sequence_byte_cost_multiplier"] = uint128ToString(t.SequenceByteCostMultiplier)

		case *astria.FeeChangeAction_TransferBaseFee:
			action.Data["transfer_base_fee"] = uint128ToString(t.TransferBaseFee)
		}
	}
	// TODO: update constants
	return nil
}
