// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/base64"
	"fmt"
	"time"

	primitive "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	astria "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transaction/v1alpha1"
	internalAstria "github.com/celenium-io/astria-indexer/internal/astria"
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/celenium-io/astria-indexer/pkg/types"
	channelTypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"
	"github.com/goccy/go-json"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func parseActions(height types.Level, blockTime time.Time, from string, tx *DecodedTx, ctx *Context) ([]storage.Action, error) {
	var (
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
			err error
		)

		switch val := rawActions[i].GetValue().(type) {
		case *astria.Action_Ibc:
			tx.ActionTypes.Set(storageTypes.ActionTypeIbcRelayBits)
			err = parseIbcAction(val, ctx, &actions[i])

		case *astria.Action_Ics20Withdrawal:
			tx.ActionTypes.Set(storageTypes.ActionTypeIcs20WithdrawalBits)
			err = parseIcs20Withdrawal(val, from, height, ctx, &actions[i])

		case *astria.Action_RollupDataSubmission:
			tx.ActionTypes.Set(storageTypes.ActionTypeRollupDataSubmissionBits)
			err = parseRollupDataSubmission(val, from, height, ctx, &actions[i])

		case *astria.Action_SudoAddressChange:
			tx.ActionTypes.Set(storageTypes.ActionTypeSudoAddressChangeBits)
			err = parseSudoAddressChangeAction(val, ctx, &actions[i])

		case *astria.Action_Transfer:
			tx.ActionTypes.Set(storageTypes.ActionTypeTransferBits)
			err = parseTransferAction(val, from, height, ctx, &actions[i])

		case *astria.Action_ValidatorUpdate:
			tx.ActionTypes.Set(storageTypes.ActionTypeValidatorUpdateBits)
			err = parseValidatorUpdateAction(val, height, ctx, &actions[i])

		case *astria.Action_BridgeLock:
			tx.ActionTypes.Set(storageTypes.ActionTypeBridgeLockBits)
			err = parseBridgeLock(val, from, height, ctx, &actions[i])

		case *astria.Action_FeeAssetChange:
			tx.ActionTypes.Set(storageTypes.ActionTypeFeeAssetChangeBits)
			err = parseFeeAssetChange(val, &actions[i])

		case *astria.Action_IbcRelayerChange:
			tx.ActionTypes.Set(storageTypes.ActionTypeIbcRelayerChangeBits)
			err = parseIbcRelayerChange(val, height, ctx, &actions[i])

		case *astria.Action_InitBridgeAccount:
			tx.ActionTypes.Set(storageTypes.ActionTypeInitBridgeAccountBits)
			err = parseInitBridgeAccount(val, from, height, ctx, &actions[i])

		case *astria.Action_BridgeSudoChange:
			tx.ActionTypes.Set(storageTypes.ActionTypeBridgeSudoChangeBits)
			err = parseBridgeSudoChange(val, height, ctx, &actions[i])

		case *astria.Action_BridgeUnlock:
			tx.ActionTypes.Set(storageTypes.ActionTypeBridgeUnlockBits)
			err = parseBridgeUnlock(val, from, height, ctx, &actions[i])

		case *astria.Action_FeeChange:
			tx.ActionTypes.Set(storageTypes.ActionTypeFeeChangeBits)
			err = parseFeeChange(val, ctx, &actions[i])

		case *astria.Action_IbcSudoChange:
			tx.ActionTypes.Set(storageTypes.ActionTypeIbcSudoChangeBits)
			err = parseIbcSudoChangeAction(val, ctx, &actions[i])

		default:
			return nil, errors.Errorf(
				"unknown action type | position = %d | block = %d: %##v",
				i, height, tx.Actions[i])
		}

		if err != nil {
			return nil, err
		}

		if actionFee, ok := ctx.Fees[int64(i)]; ok {
			actionFee.Height = height
			actionFee.Time = blockTime
			actionFee.Payer = &storage.Address{
				Hash: from,
			}
			actions[i].Fee = actionFee

			fromAmount := actionFee.Amount.Neg()
			addr := ctx.Addresses.Set(from, height, fromAmount, actionFee.Asset, 0, 0)
			actions[i].BalanceUpdates = append(actions[i].BalanceUpdates, storage.BalanceUpdate{
				Address:  addr,
				Height:   actions[i].Height,
				Currency: actionFee.Asset,
				Update:   fromAmount,
			})

			to := ctx.Addresses.Set(ctx.Proposer, height, actionFee.Amount, actionFee.Asset, 0, 0)
			actions[i].BalanceUpdates = append(actions[i].BalanceUpdates, storage.BalanceUpdate{
				Address:  to,
				Height:   actions[i].Height,
				Currency: actionFee.Asset,
				Update:   actionFee.Amount,
			})
		}

		if deposit, ok := ctx.Deposits[int64(i)]; ok {
			deposit.Height = height
			deposit.Time = blockTime
			actions[i].Deposit = deposit
		}
	}

	return actions, nil
}

func parseIbcAction(body *astria.Action_Ibc, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIbcRelay
	action.Data = make(map[string]any)

	if body.Ibc != nil && body.Ibc.GetRawAction() != nil {
		data := body.Ibc.GetRawAction().GetValue()
		typ := body.Ibc.GetRawAction().GetTypeUrl()
		action.Data["raw"] = base64.StdEncoding.EncodeToString(data)
		action.Data["type"] = typ

		switch typ {
		case "/ibc.core.channel.v1.MsgRecvPacket":
			var msg channelTypes.MsgRecvPacket
			if err := proto.Unmarshal(data, &msg); err != nil {
				return err
			}
			var transfer IbcTransfer
			if err := json.Unmarshal(msg.Packet.Data, &transfer); err != nil {
				return nil
			}
			asset := fmt.Sprintf("%s/%s/%s", msg.Packet.GetDestPort(), msg.Packet.GetDestChannel(), transfer.Denom)
			var addr string
			var amount decimal.Decimal

			switch {
			case internalAstria.IsAddress(transfer.Receiver):
				addr = transfer.Receiver
				amount = transfer.Amount.Copy()
			case internalAstria.IsCompatAddress(transfer.Receiver):
				a, err := internalAstria.CompatToAstria(transfer.Receiver)
				if err != nil {
					return err
				}
				addr = a
				amount = transfer.Amount.Copy()
			case internalAstria.IsAddress(transfer.Sender):
				addr = transfer.Sender
				amount = transfer.Amount.Neg()
			case internalAstria.IsCompatAddress(transfer.Sender):
				a, err := internalAstria.CompatToAstria(transfer.Sender)
				if err != nil {
					return err
				}
				addr = a
				amount = transfer.Amount.Neg()
			}

			if addr != "" {
				address := ctx.Addresses.Set(addr, action.Height, amount, asset, 0, 0)
				action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
					Address:  address,
					Height:   action.Height,
					Currency: asset,
					Update:   amount,
				})
			}
		default:
		}
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
		action.Data["use_compat_address"] = body.Ics20Withdrawal.GetUseCompatAddress()

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
					Update:   decAmount.Copy().Neg(),
				})
			} else {
				action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
					Address:  addr,
					Height:   action.Height,
					Currency: body.Ics20Withdrawal.GetDenom(),
					Update:   decAmount.Copy().Neg(),
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
					Update:   decAmount.Copy().Neg(),
				})
			} else {
				action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
					Address:  addr,
					Height:   action.Height,
					Currency: body.Ics20Withdrawal.GetDenom(),
					Update:   decAmount.Copy().Neg(),
				})
			}
		}
	}
	return nil
}

func parseRollupDataSubmission(body *astria.Action_RollupDataSubmission, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeRollupDataSubmission
	action.Data = make(map[string]any)
	if body.RollupDataSubmission != nil {
		rollupId := body.RollupDataSubmission.GetRollupId().GetInner()
		action.Data["rollup_id"] = rollupId
		action.Data["data"] = body.RollupDataSubmission.GetData()
		action.Data["fee_asset"] = body.RollupDataSubmission.GetFeeAsset()
		dataSize := len(body.RollupDataSubmission.GetData())

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

func parseSudoAddressChangeAction(body *astria.Action_SudoAddressChange, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeSudoAddressChange
	action.Data = make(map[string]any)
	if body.SudoAddressChange != nil {
		address := body.SudoAddressChange.GetNewAddress().GetBech32M()
		action.Data["address"] = address

		addr := ctx.Addresses.Set(address, action.Height, decimal.Zero, "", 1, 0)
		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    addr,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})
		ctx.AddGenericConstant("authority_sudo_address", address)
	}
	return nil
}

func parseTransferAction(body *astria.Action_Transfer, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeTransfer
	action.Data = make(map[string]any)
	if body.Transfer != nil {
		amount := uint128ToString(body.Transfer.GetAmount())
		to := body.Transfer.GetTo().GetBech32M()
		asset := body.Transfer.GetAsset()

		action.Data["amount"] = amount
		action.Data["asset"] = asset
		action.Data["fee_asset"] = body.Transfer.GetFeeAsset()
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

func parseValidatorUpdateAction(body *astria.Action_ValidatorUpdate, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeValidatorUpdate
	action.Data = make(map[string]any)
	if body.ValidatorUpdate != nil {
		power := body.ValidatorUpdate.GetPower()
		action.Data["power"] = power
		pubKey := body.ValidatorUpdate.GetPubKey().GetEd25519()
		action.Data["pubkey"] = pubKey

		address, err := AddressFromPubKey(body.ValidatorUpdate.GetPubKey().GetEd25519())
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

func parseIbcRelayerChange(body *astria.Action_IbcRelayerChange, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIbcRelayerChange
	action.Data = make(map[string]any)
	if body.IbcRelayerChange != nil {
		if addition := body.IbcRelayerChange.GetAddition(); len(addition.GetBech32M()) > 0 {
			b32m := addition.GetBech32M()
			action.Data["addition"] = b32m

			addr := ctx.Addresses.Set(b32m, height, decimal.Zero, "", 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
			ctx.Addresses.AddIbcRelayer(b32m)
		}

		if removal := body.IbcRelayerChange.GetRemoval(); len(removal.GetBech32M()) > 0 {
			b32m := removal.GetBech32M()
			action.Data["removal"] = b32m

			addr := ctx.Addresses.Set(b32m, height, decimal.Zero, "", 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
			ctx.Addresses.RemoveIbcRelayer(b32m)
		}
	}
	return nil
}

func parseInitBridgeAccount(body *astria.Action_InitBridgeAccount, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeInitBridgeAccount
	action.Data = make(map[string]any)
	if body.InitBridgeAccount != nil {
		rollupId := body.InitBridgeAccount.GetRollupId().GetInner()
		rollup := ctx.Rollups.Set(rollupId, height, 0)

		bridge := storage.Bridge{
			InitHeight: height,
			Asset:      body.InitBridgeAccount.GetAsset(),
			FeeAsset:   body.InitBridgeAccount.GetFeeAsset(),
			Address:    ctx.Addresses.Set(from, height, decimal.Zero, "", 0, 0),
			Rollup:     rollup,
		}
		bridge.Address.IsBridge = true
		rollup.BridgeCount += 1

		action.Data["rollup_id"] = rollupId
		action.Data["fee_asset"] = bridge.FeeAsset
		action.Data["asset"] = bridge.Asset

		if sudo := body.InitBridgeAccount.GetSudoAddress().GetBech32M(); sudo != "" {
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

		if withdrawer := body.InitBridgeAccount.GetWithdrawerAddress().GetBech32M(); withdrawer != "" {
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

func parseBridgeSudoChange(body *astria.Action_BridgeSudoChange, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeBridgeSudoChangeAction
	action.Data = make(map[string]any)
	if body.BridgeSudoChange != nil {
		bridgeAddress := body.BridgeSudoChange.GetBridgeAddress().GetBech32M()
		sudo := body.BridgeSudoChange.GetNewSudoAddress().GetBech32M()
		withdrawer := body.BridgeSudoChange.GetNewWithdrawerAddress().GetBech32M()
		feeAsset := body.BridgeSudoChange.GetFeeAsset()

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

func parseBridgeLock(body *astria.Action_BridgeLock, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeBridgeLock
	action.Data = make(map[string]any)
	if body.BridgeLock != nil {
		amount := uint128ToString(body.BridgeLock.GetAmount())
		asset := body.BridgeLock.GetAsset()

		action.Data["to"] = body.BridgeLock.GetTo().GetBech32M()
		action.Data["destination_chain_address"] = body.BridgeLock.GetDestinationChainAddress()
		action.Data["asset"] = asset
		action.Data["fee_asset"] = body.BridgeLock.GetFeeAsset()
		action.Data["amount"] = amount

		toAddress := body.BridgeLock.GetTo().GetBech32M()
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

func parseBridgeUnlock(body *astria.Action_BridgeUnlock, from string, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeBridgeUnlock
	action.Data = make(map[string]any)
	if body.BridgeUnlock != nil {
		amount := uint128ToString(body.BridgeUnlock.GetAmount())
		toAddress := body.BridgeUnlock.GetTo().GetBech32M()
		bridge := body.BridgeUnlock.GetBridgeAddress().GetBech32M()
		feeAsset := body.BridgeUnlock.GetFeeAsset()

		action.Data["to"] = toAddress
		action.Data["fee_asset"] = feeAsset
		action.Data["amount"] = amount
		action.Data["rollup_block_number"] = body.BridgeUnlock.GetRollupBlockNumber()
		action.Data["rollup_withdrawal_event_id"] = body.BridgeUnlock.GetRollupWithdrawalEventId()

		if memo := body.BridgeUnlock.GetMemo(); len(memo) > 0 {
			action.Data["memo"] = memo
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

func parseFeeAssetChange(body *astria.Action_FeeAssetChange, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeFeeAssetChange
	action.Data = make(map[string]any)
	if body.FeeAssetChange != nil {
		if addition := body.FeeAssetChange.GetAddition(); len(addition) > 0 {
			action.Data["addition"] = addition
		}

		if removal := body.FeeAssetChange.GetRemoval(); len(removal) > 0 {
			action.Data["removal"] = removal
		}
	}
	return nil
}

func parseFeeChange(body *astria.Action_FeeChange, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeFeeChange
	action.Data = make(map[string]any)
	if body.FeeChange != nil {
		switch t := body.FeeChange.GetFeeComponents().(type) {

		case *astria.FeeChange_BridgeLock:
			processFeeComponent("bridge_lock", t.BridgeLock.GetMultiplier(), t.BridgeLock.GetBase(), action.Data, ctx)

		case *astria.FeeChange_BridgeSudoChange:
			processFeeComponent("bridge_sudo_change", t.BridgeSudoChange.GetMultiplier(), t.BridgeSudoChange.GetBase(), action.Data, ctx)

		case *astria.FeeChange_BridgeUnlock:
			processFeeComponent("bridge_unlock", t.BridgeUnlock.GetMultiplier(), t.BridgeUnlock.GetBase(), action.Data, ctx)

		case *astria.FeeChange_FeeAssetChange:
			processFeeComponent("fee_asset_change", t.FeeAssetChange.GetMultiplier(), t.FeeAssetChange.GetBase(), action.Data, ctx)

		case *astria.FeeChange_FeeChange:
			processFeeComponent("fee_change", t.FeeChange.GetMultiplier(), t.FeeChange.GetBase(), action.Data, ctx)

		case *astria.FeeChange_IbcRelay:
			processFeeComponent("ibc_relay", t.IbcRelay.GetMultiplier(), t.IbcRelay.GetBase(), action.Data, ctx)

		case *astria.FeeChange_IbcRelayerChange:
			processFeeComponent("ibc_relay_change", t.IbcRelayerChange.GetMultiplier(), t.IbcRelayerChange.GetBase(), action.Data, ctx)

		case *astria.FeeChange_IbcSudoChange:
			processFeeComponent("ibc_sudo_change", t.IbcSudoChange.GetMultiplier(), t.IbcSudoChange.GetBase(), action.Data, ctx)

		case *astria.FeeChange_Ics20Withdrawal:
			processFeeComponent("ics20_withdrawal", t.Ics20Withdrawal.GetMultiplier(), t.Ics20Withdrawal.GetBase(), action.Data, ctx)

		case *astria.FeeChange_InitBridgeAccount:
			processFeeComponent("init_bridge_account", t.InitBridgeAccount.GetMultiplier(), t.InitBridgeAccount.GetBase(), action.Data, ctx)

		case *astria.FeeChange_RollupDataSubmission:
			processFeeComponent("rollup_data_submission", t.RollupDataSubmission.GetMultiplier(), t.RollupDataSubmission.GetBase(), action.Data, ctx)

		case *astria.FeeChange_SudoAddressChange:
			processFeeComponent("sudo_address_change", t.SudoAddressChange.GetMultiplier(), t.SudoAddressChange.GetBase(), action.Data, ctx)

		case *astria.FeeChange_Transfer:
			processFeeComponent("transfer", t.Transfer.GetMultiplier(), t.Transfer.GetBase(), action.Data, ctx)

		case *astria.FeeChange_ValidatorUpdate:
			processFeeComponent("validator_update", t.ValidatorUpdate.GetMultiplier(), t.ValidatorUpdate.GetBase(), action.Data, ctx)
		}
	}
	return nil
}

func parseIbcSudoChangeAction(body *astria.Action_IbcSudoChange, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIbcSudoChangeAction
	action.Data = make(map[string]any)
	if body.IbcSudoChange != nil {
		address := body.IbcSudoChange.GetNewAddress().GetBech32M()
		action.Data["address"] = address

		addr := ctx.Addresses.Set(address, action.Height, decimal.Zero, "", 1, 0)
		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    addr,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})

		ctx.AddGenericConstant("ibc_sudo_address", address)
	}

	return nil
}

func processFeeComponent(name string, multiplier, base *primitive.Uint128, data map[string]any, ctx *Context) {
	m := uint128ToString(multiplier)
	mKey := fmt.Sprintf("%s_multiplier", name)
	data[mKey] = m
	ctx.AddGenericConstant(mKey, m)

	b := uint128ToString(base)
	bKey := fmt.Sprintf("%s_base", name)
	data[bKey] = b
	ctx.AddGenericConstant(bKey, b)
}
