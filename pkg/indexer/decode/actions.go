// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/base64"
	"encoding/hex"
	"time"

	astria "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func parseActions(height types.Level, blockTime time.Time, from string, tx *DecodedTx, ctx *Context) ([]storage.Action, error) {
	var (
		feeCounter = 0
		rawActions = tx.UnsignedTx.GetActions()
		actions    = make([]storage.Action, len(rawActions))
	)

	for i := range rawActions {
		if tx.UnsignedTx.Actions[i].Value == nil {
			return nil, errors.Errorf("nil action")
		}
		actions[i].Height = height
		actions[i].Time = blockTime
		actions[i].Position = int64(i)
		actions[i].Addresses = make([]*storage.AddressAction, 0)
		actions[i].BalanceUpdates = make([]storage.BalanceUpdate, 0)

		var (
			err     error
			feeType string
		)

		switch val := rawActions[i].GetValue().(type) {
		case *astria.Action_IbcAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeIbcRelayBits)
			err = parseIbcAction(val, &actions[i])
			feeType = "penumbra.core.component.ibc.v1.IbcAction"

		case *astria.Action_Ics20Withdrawal:
			tx.ActionTypes.Set(storageTypes.ActionTypeIcs20WithdrawalBits)
			err = parseIcs20Withdrawal(val, from, height, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.Ics20Withdrawal"

		case *astria.Action_SequenceAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeSequenceBits)
			err = parseSequenceAction(val, from, height, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.SequenceAction"

		case *astria.Action_SudoAddressChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeSudoAddressChangeBits)
			err = parseSudoAddressChangeAction(val, height, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.SudoAddressChangeAction"

		case *astria.Action_TransferAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeTransferBits)
			err = parseTransferAction(val, from, height, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.TransferAction"

		case *astria.Action_ValidatorUpdateAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeValidatorUpdateBits)
			err = parseValidatorUpdateAction(val, height, ctx, &actions[i])
			feeType = "tendermint.abci.ValidatorUpdateAction"

		case *astria.Action_BridgeLockAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeBridgeLockBits)
			err = parseBridgeLock(val, from, height, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.BridgeLockAction"

		case *astria.Action_FeeAssetChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeFeeAssetChangeBits)
			err = parseFeeAssetChange(val, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.FeeAssetChangeAction"

		case *astria.Action_IbcRelayerChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeIbcRelayerChangeBits)
			err = parseIbcRelayerChange(val, height, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.IbcRelayerChangeAction"

		case *astria.Action_InitBridgeAccountAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeInitBridgeAccountBits)
			err = parseInitBridgeAccount(val, from, height, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.InitBridgeAccountAction"

		case *astria.Action_BridgeSudoChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeBridgeSudoChangeBits)
			err = parseBridgeSudoChange(val, height, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.BridgeSudoChangeAction"

		case *astria.Action_BridgeUnlockAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeBridgeUnlockBits)
			err = parseBridgeUnlock(val, from, height, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.BridgeUnlockAction"

		case *astria.Action_FeeChangeAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeFeeChangeBits)
			err = parseFeeChange(val, ctx, &actions[i])
			feeType = "astria.protocol.transactions.v1alpha1.FeeChangeAction"

		default:
			return nil, errors.Errorf(
				"unknown action type | position = %d | block = %d: %##v",
				i, height, tx.Actions[i])
		}

		if err != nil {
			return nil, err
		}

		// merge fees
		if len(ctx.Fees) > feeCounter {
			if ctx.Fees[feeCounter].ActionType == feeType {
				ctx.Fees[feeCounter].Height = height
				ctx.Fees[feeCounter].Time = blockTime
				ctx.Fees[feeCounter].Payer = &storage.Address{
					Hash: from,
				}
				actions[i].Fee = ctx.Fees[feeCounter]
				fromAmount := ctx.Fees[feeCounter].Amount.Neg()
				addr := ctx.Addresses.Set(from, height, fromAmount, ctx.Fees[feeCounter].Asset, 0, 0)
				actions[i].BalanceUpdates = append(actions[i].BalanceUpdates, storage.BalanceUpdate{
					Address:  addr,
					Height:   actions[i].Height,
					Currency: ctx.Fees[feeCounter].Asset,
					Update:   fromAmount,
				})

				to := ctx.Addresses.Set(ctx.Proposer, height, ctx.Fees[feeCounter].Amount, ctx.Fees[feeCounter].Asset, 0, 0)
				actions[i].BalanceUpdates = append(actions[i].BalanceUpdates, storage.BalanceUpdate{
					Address:  to,
					Height:   actions[i].Height,
					Currency: ctx.Fees[feeCounter].Asset,
					Update:   ctx.Fees[feeCounter].Amount,
				})
				feeCounter++
			}
		}
	}

	return actions, nil
}

func parseIbcAction(body *astria.Action_IbcAction, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIbcRelay
	action.Data = make(map[string]any)

	if body.IbcAction != nil && body.IbcAction.GetRawAction() != nil {
		data := body.IbcAction.GetRawAction().GetValue()
		action.Data["raw"] = base64.StdEncoding.EncodeToString(data)
		action.Data["type"] = body.IbcAction.GetRawAction().GetTypeUrl()
	}
	return nil
}

func parseIcs20Withdrawal(body *astria.Action_Ics20Withdrawal, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIcs20Withdrawal
	action.Data = make(map[string]any)
	if body.Ics20Withdrawal != nil {
		amount := uint128ToString(body.Ics20Withdrawal.GetAmount())
		asset := body.Ics20Withdrawal.GetDenom()
		action.Data["amount"] = amount
		action.Data["denom"] = asset
		action.Data["fee_asset"] = body.Ics20Withdrawal.GetFeeAsset()
		action.Data["destination_address"] = body.Ics20Withdrawal.GetDestinationChainAddress()
		action.Data["return_address"] = body.Ics20Withdrawal.GetReturnAddress().GetBech32M()
		action.Data["source_channel"] = body.Ics20Withdrawal.GetSourceChannel()

		decAmount := decimal.RequireFromString(amount)

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

		returnAddress := body.Ics20Withdrawal.GetReturnAddress().GetBech32M()
		if bridge := body.Ics20Withdrawal.GetBridgeAddress().GetBech32M(); bridge != "" {
			action.Data["bridge"] = bridge
			addr := ctx.Addresses.Set(bridge, height, decAmount.Copy().Neg(), asset, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
			if returnAddress != bridge {
				returnAddr := ctx.Addresses.Set(returnAddress, height, decAmount, asset, 1, 0)
				action.Addresses = append(action.Addresses, &storage.AddressAction{
					Address:    returnAddr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				})

				action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
					Address:  returnAddr,
					Height:   action.Height,
					Currency: body.Ics20Withdrawal.GetDenom(),
					Update:   decAmount,
				})
			} else {
				action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
					Address:  addr,
					Height:   action.Height,
					Currency: body.Ics20Withdrawal.GetDenom(),
					Update:   decAmount,
				})
			}
		} else {
			addr := ctx.Addresses.Set(from, height, decAmount.Copy().Neg(), asset, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
			if returnAddress != from {
				returnAddr := ctx.Addresses.Set(returnAddress, height, decAmount, asset, 1, 0)
				action.Addresses = append(action.Addresses, &storage.AddressAction{
					Address:    returnAddr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				})
				action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
					Address:  returnAddr,
					Height:   action.Height,
					Currency: body.Ics20Withdrawal.GetDenom(),
					Update:   decAmount,
				})
			} else {
				action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
					Address:  addr,
					Height:   action.Height,
					Currency: body.Ics20Withdrawal.GetDenom(),
					Update:   decAmount,
				})
			}
		}
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
		fromAddress := ctx.Addresses.Set(from, height, decimal.Zero, "", 1, 0)

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
			Time:       action.Time,
			Height:     action.Height,
			Size:       int64(dataSize),
			Action:     action,
			Rollup:     rollup,
			ActionType: action.Type,
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

		addr := ctx.Addresses.Set(address, height, decimal.Zero, "", 1, 0)
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

		transfer := storage.Transfer{
			Height: height,
			Time:   action.Time,
			Asset:  asset,
			Amount: decAmount,
		}

		if from == to {
			addr := ctx.Addresses.Set(from, height, decimal.Zero, "", 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})

			transfer.Source = addr
			transfer.Destination = addr
		} else {
			toAddr := ctx.Addresses.Set(to, height, decAmount, asset, 1, 0)
			fromAddr := ctx.Addresses.Set(from, height, decAmount.Neg(), asset, 1, 0)
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

			transfer.Source = fromAddr
			transfer.Destination = toAddr
		}

		ctx.AddTransfer(&transfer)
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
		addr := ctx.Addresses.Set(address, height, decimal.Zero, "", 1, 0)
		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    addr,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})
		ctx.Validators.Set(pubKey, power, address, height)
	}
	return nil
}

func parseIbcRelayerChange(body *astria.Action_IbcRelayerChangeAction, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIbcRelayerChange
	action.Data = make(map[string]any)
	if body.IbcRelayerChangeAction != nil {
		if addition := body.IbcRelayerChangeAction.GetAddition(); len(addition.GetBech32M()) > 0 {
			action.Data["addition"] = addition.GetBech32M()

			addr := ctx.Addresses.Set(addition.GetBech32M(), height, decimal.Zero, "", 1, 0)
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

			addr := ctx.Addresses.Set(removal.GetBech32M(), height, decimal.Zero, "", 1, 0)
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
		rollupId := body.InitBridgeAccountAction.GetRollupId().GetInner()
		rollup := ctx.Rollups.Set(rollupId, height, 0)

		bridge := storage.Bridge{
			InitHeight: height,
			Asset:      body.InitBridgeAccountAction.GetAsset(),
			FeeAsset:   body.InitBridgeAccountAction.GetFeeAsset(),
			Address:    ctx.Addresses.Set(from, height, decimal.Zero, "", 0, 0),
			Rollup:     rollup,
		}
		bridge.Address.IsBridge = true
		rollup.BridgeCount += 1

		action.Data["rollup_id"] = rollupId
		action.Data["fee_asset"] = bridge.FeeAsset
		action.Data["asset"] = bridge.Asset

		if sudo := body.InitBridgeAccountAction.GetSudoAddress().GetBech32M(); sudo != "" {
			action.Data["sudo"] = sudo
			if sudo != from {
				addr := ctx.Addresses.Set(sudo, height, decimal.Zero, "", 1, 0)
				action.Addresses = append(action.Addresses, &storage.AddressAction{
					Address:    addr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				})
				bridge.Sudo = addr
			}
		}

		if bridge.Sudo == nil {
			bridge.Sudo = bridge.Address
		}

		if withdrawer := body.InitBridgeAccountAction.GetWithdrawerAddress().GetBech32M(); withdrawer != "" {
			action.Data["withdrawer"] = withdrawer
			if withdrawer != from {
				addr := ctx.Addresses.Set(withdrawer, height, decimal.Zero, "", 1, 0)
				action.Addresses = append(action.Addresses, &storage.AddressAction{
					Address:    addr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				})
				bridge.Withdrawer = addr
			}
		}

		if bridge.Withdrawer == nil {
			bridge.Withdrawer = bridge.Address
		}

		action.RollupAction = &storage.RollupAction{
			Time:       action.Time,
			Height:     action.Height,
			Action:     action,
			Rollup:     rollup,
			ActionType: action.Type,
		}
		ctx.AddBridge(&bridge)
		ctx.AddBridgeAsset(from, bridge.Asset)
	}
	return nil
}

func parseBridgeSudoChange(body *astria.Action_BridgeSudoChangeAction, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeBridgeSudoChangeAction
	action.Data = make(map[string]any)
	if body.BridgeSudoChangeAction != nil {
		bridgeAddress := body.BridgeSudoChangeAction.GetBridgeAddress().GetBech32M()
		sudo := body.BridgeSudoChangeAction.GetNewSudoAddress().GetBech32M()
		withdrawer := body.BridgeSudoChangeAction.GetNewWithdrawerAddress().GetBech32M()
		feeAsset := body.BridgeSudoChangeAction.GetFeeAsset()

		action.Data["bridge"] = bridgeAddress
		bridgeAddr := ctx.Addresses.Set(bridgeAddress, height, decimal.Zero, "", 1, 0)
		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    bridgeAddr,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})

		bridge := storage.Bridge{
			Address:    bridgeAddr,
			Sudo:       bridgeAddr,
			Withdrawer: bridgeAddr,
		}

		if sudo != "" {
			action.Data["sudo"] = sudo

			if bridgeAddress != sudo {
				addr := ctx.Addresses.Set(sudo, height, decimal.Zero, "", 1, 0)
				action.Addresses = append(action.Addresses, &storage.AddressAction{
					Address:    addr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				})
				bridge.Sudo = addr
			}
		}

		if withdrawer != "" {
			action.Data["withdrawer"] = withdrawer

			actions := 1
			if sudo == withdrawer || bridgeAddress == withdrawer {
				actions = 0
			}
			addr := ctx.Addresses.Set(withdrawer, height, decimal.Zero, "", actions, 0)
			bridge.Withdrawer = addr

			if bridgeAddress != withdrawer && sudo != withdrawer {
				action.Addresses = append(action.Addresses, &storage.AddressAction{
					Address:    addr,
					Action:     action,
					Time:       action.Time,
					Height:     action.Height,
					ActionType: action.Type,
				})
			}

		}

		if feeAsset != "" {
			action.Data["fee_asset"] = feeAsset
			bridge.FeeAsset = feeAsset
		}

		ctx.AddBridge(&bridge)
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
		toAddr := ctx.Addresses.Set(toAddress, height, decAmount, asset, 1, 0)

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
			fromAddr := ctx.Addresses.Set(from, height, decAmount.Neg(), asset, 1, 0)

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
		feeAsset := body.BridgeUnlockAction.GetFeeAsset()

		action.Data["to"] = toAddress
		action.Data["fee_asset"] = feeAsset
		action.Data["amount"] = amount
		if memo := body.BridgeUnlockAction.GetMemo(); len(memo) > 0 {
			action.Data["memo"] = hex.EncodeToString(memo)
		}
		if bridge != "" {
			action.Data["bridge"] = bridge
		}

		var (
			decAmount   = decimal.RequireFromString(amount)
			fromAddr    *storage.Address
			unlockAsset string
		)

		if bridge == "" {
			fromAddr = ctx.Addresses.Set(from, height, decAmount.Neg(), "", 1, 0)
			unlockAsset = currency.DefaultCurrency
		} else {
			asset, ok := ctx.bridgeAssets[bridge]
			if !ok {
				return errors.Errorf("unknown bridge asset: %s", bridge)
			}
			fromAddr = ctx.Addresses.Set(bridge, height, decAmount.Neg(), asset, 1, 0)
			unlockAsset = asset
		}

		toAddr := ctx.Addresses.Set(toAddress, height, decAmount, unlockAsset, 1, 0)

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
				Currency: unlockAsset,
				Update:   decAmount,
			},
			storage.BalanceUpdate{
				Address:  fromAddr,
				Height:   action.Height,
				Currency: unlockAsset,
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

func parseFeeChange(body *astria.Action_FeeChangeAction, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeFeeChange
	action.Data = make(map[string]any)
	if body.FeeChangeAction != nil {
		switch t := body.FeeChangeAction.GetValue().(type) {
		case *astria.FeeChangeAction_BridgeLockByteCostMultiplier:
			val := uint128ToString(t.BridgeLockByteCostMultiplier)
			action.Data["bridge_lock_byte_cost_multiplier"] = val
			ctx.AddGenericConstant("bridge_lock_byte_cost_multiplier", val)

		case *astria.FeeChangeAction_BridgeSudoChangeBaseFee:
			val := uint128ToString(t.BridgeSudoChangeBaseFee)
			action.Data["bridge_sudo_change_base_fee"] = val
			ctx.AddGenericConstant("bridge_sudo_change_fee", val)

		case *astria.FeeChangeAction_Ics20WithdrawalBaseFee:
			val := uint128ToString(t.Ics20WithdrawalBaseFee)
			action.Data["ics20_withdrawal_base_fee"] = val
			ctx.AddGenericConstant("ics20_withdrawal_base_fee", val)

		case *astria.FeeChangeAction_InitBridgeAccountBaseFee:
			val := uint128ToString(t.InitBridgeAccountBaseFee)
			action.Data["init_bridge_account_base_fee"] = val
			ctx.AddGenericConstant("init_bridge_account_base_fee", val)

		case *astria.FeeChangeAction_SequenceBaseFee:
			val := uint128ToString(t.SequenceBaseFee)
			action.Data["sequence_base_fee"] = val
			ctx.AddGenericConstant("sequence_base_fee", val)

		case *astria.FeeChangeAction_SequenceByteCostMultiplier:
			val := uint128ToString(t.SequenceByteCostMultiplier)
			action.Data["sequence_byte_cost_multiplier"] = val
			ctx.AddGenericConstant("sequence_byte_cost_multiplier", val)

		case *astria.FeeChangeAction_TransferBaseFee:
			val := uint128ToString(t.TransferBaseFee)
			action.Data["transfer_base_fee"] = val
			ctx.AddGenericConstant("transfer_base_fee", val)
		}
	}
	return nil
}
