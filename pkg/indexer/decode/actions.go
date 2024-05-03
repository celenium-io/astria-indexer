// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	stdBytes "bytes"
	"encoding/base64"
	"encoding/hex"
	"time"

	astria "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
	"github.com/celenium-io/astria-indexer/internal/storage"
	storageTypes "github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/cometbft/cometbft/libs/bytes"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func parseActions(height types.Level, blockTime time.Time, from bytes.HexBytes, tx *DecodedTx, ctx *Context) ([]storage.Action, error) {
	actions := make([]storage.Action, len(tx.Tx.Transaction.Actions))
	for i := range tx.Tx.Transaction.Actions {
		if tx.Tx.Transaction.Actions[i].Value == nil {
			return nil, errors.Errorf("nil action")
		}
		actions[i].Height = height
		actions[i].Time = blockTime
		actions[i].Position = int64(i)
		actions[i].Addresses = make([]*storage.AddressAction, 0)
		actions[i].BalanceUpdates = make([]storage.BalanceUpdate, 0)

		var err error

		switch val := tx.Tx.Transaction.Actions[i].Value.(type) {
		case *astria.Action_IbcAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeIbcRelayBits)
			err = parseIbcAction(val, &actions[i])
		case *astria.Action_Ics20Withdrawal:
			tx.ActionTypes.Set(storageTypes.ActionTypeIcs20WithdrawalBits)
			err = parseIcs20Withdrawal(val, height, ctx, &actions[i])
		case *astria.Action_MintAction:
			tx.ActionTypes.Set(storageTypes.ActionTypeMintBits)
			err = parseMintAction(val, height, ctx, &actions[i])
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
	if body.IbcAction != nil && body.IbcAction.RawAction != nil {
		action.Data["raw"] = base64.StdEncoding.EncodeToString(body.IbcAction.RawAction.Value)
	}
	return nil
}

func parseIcs20Withdrawal(body *astria.Action_Ics20Withdrawal, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIcs20Withdrawal
	action.Data = make(map[string]any)
	if body.Ics20Withdrawal != nil {
		amount := uint128ToString(body.Ics20Withdrawal.Amount)
		action.Data["amount"] = amount
		action.Data["denom"] = body.Ics20Withdrawal.Denom
		action.Data["destination_address"] = body.Ics20Withdrawal.DestinationChainAddress
		action.Data["return_address"] = hex.EncodeToString(body.Ics20Withdrawal.ReturnAddress)
		action.Data["source_channel"] = body.Ics20Withdrawal.SourceChannel

		if body.Ics20Withdrawal.TimeoutHeight != nil {
			action.Data["timeout_height"] = map[string]any{
				"revision_number": body.Ics20Withdrawal.TimeoutHeight.RevisionNumber,
				"revision_height": body.Ics20Withdrawal.TimeoutHeight.RevisionHeight,
			}
		}
		if body.Ics20Withdrawal.TimeoutTime > 0 {
			action.Data["timeout_time"] = body.Ics20Withdrawal.TimeoutTime
		}

		decAmount := decimal.RequireFromString(amount)
		returnAddress := bytes.HexBytes(body.Ics20Withdrawal.ReturnAddress)
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
			Currency: body.Ics20Withdrawal.Denom,
			Update:   decAmount,
		})
	}
	return nil
}

func parseMintAction(body *astria.Action_MintAction, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeMint
	action.Data = make(map[string]any)
	if body.MintAction != nil {
		amount := uint128ToString(body.MintAction.Amount)
		action.Data["amount"] = amount
		action.Data["to"] = hex.EncodeToString(body.MintAction.To.GetInner())

		toAddress := bytes.HexBytes(body.MintAction.To.GetInner())
		decAmount := decimal.RequireFromString(amount)
		addr := ctx.Addresses.Set(toAddress, height, decAmount, 1, 0)
		action.Addresses = append(action.Addresses, &storage.AddressAction{
			Address:    addr,
			Action:     action,
			Time:       action.Time,
			Height:     action.Height,
			ActionType: action.Type,
		})

		ctx.SupplyChange = ctx.SupplyChange.Add(decAmount)

		action.BalanceUpdates = append(action.BalanceUpdates, storage.BalanceUpdate{
			Address:  addr,
			Height:   action.Height,
			Currency: addr.Balance.Currency,
			Update:   decAmount,
		})
	}
	return nil
}

func parseSequenceAction(body *astria.Action_SequenceAction, from bytes.HexBytes, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeSequence
	action.Data = make(map[string]any)
	if body.SequenceAction != nil {
		action.Data["rollup_id"] = body.SequenceAction.RollupId.GetInner()
		action.Data["data"] = body.SequenceAction.Data
		dataSize := len(body.SequenceAction.Data)

		rollup := ctx.Rollups.Set(body.SequenceAction.RollupId.GetInner(), height, dataSize)
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
		action.Data["address"] = hex.EncodeToString(body.SudoAddressChangeAction.NewAddress.GetInner())

		newAddress := bytes.HexBytes(body.SudoAddressChangeAction.NewAddress.GetInner())
		addr := ctx.Addresses.Set(newAddress, height, decimal.Zero, 1, 0)
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

func parseTransferAction(body *astria.Action_TransferAction, from bytes.HexBytes, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeTransfer
	action.Data = make(map[string]any)
	if body.TransferAction != nil {
		amount := uint128ToString(body.TransferAction.Amount)
		action.Data["amount"] = uint128ToString(body.TransferAction.Amount)
		action.Data["asset_id"] = body.TransferAction.AssetId
		action.Data["to"] = hex.EncodeToString(body.TransferAction.To.GetInner())

		toAddress := bytes.HexBytes(body.TransferAction.To.GetInner())
		decAmount := decimal.RequireFromString(amount)

		if stdBytes.Equal(from, toAddress) {
			addr := ctx.Addresses.Set(from, height, decimal.Zero, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
		} else {
			toAddr := ctx.Addresses.Set(toAddress, height, decAmount, 1, 0)
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
					Currency: toAddr.Balance.Currency,
					Update:   decAmount,
				},
				storage.BalanceUpdate{
					Address:  fromAddr,
					Height:   action.Height,
					Currency: fromAddr.Balance.Currency,
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
		action.Data["power"] = body.ValidatorUpdateAction.Power
		action.Data["pubkey"] = body.ValidatorUpdateAction.PubKey.GetEd25519()

		address := AddressFromPubKey(body.ValidatorUpdateAction.PubKey.GetEd25519())
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

func parseIbcRelayerChange(body *astria.Action_IbcRelayerChangeAction, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeIbcRelayerChange
	action.Data = make(map[string]any)
	if body.IbcRelayerChangeAction != nil {
		if addition := body.IbcRelayerChangeAction.GetAddition(); len(addition.GetInner()) > 0 {
			action.Data["addition"] = hex.EncodeToString(addition.GetInner())

			addrBytes := bytes.HexBytes(addition.GetInner())
			addr := ctx.Addresses.Set(addrBytes, height, decimal.Zero, 1, 0)
			action.Addresses = append(action.Addresses, &storage.AddressAction{
				Address:    addr,
				Action:     action,
				Time:       action.Time,
				Height:     action.Height,
				ActionType: action.Type,
			})
		}

		if removal := body.IbcRelayerChangeAction.GetRemoval(); len(removal.GetInner()) > 0 {
			action.Data["removal"] = hex.EncodeToString(removal.GetInner())

			addrBytes := bytes.HexBytes(removal.GetInner())
			addr := ctx.Addresses.Set(addrBytes, height, decimal.Zero, 1, 0)
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

func parseInitBridgeAccount(body *astria.Action_InitBridgeAccountAction, from bytes.HexBytes, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeInitBridgeAccount
	action.Data = make(map[string]any)
	if body.InitBridgeAccountAction != nil {
		action.Data["rollup_id"] = body.InitBridgeAccountAction.GetRollupId().GetInner()
		action.Data["fee_asset_id"] = body.InitBridgeAccountAction.GetFeeAssetId()
		action.Data["asset_id"] = body.InitBridgeAccountAction.GetAssetId()

		rollup := ctx.Rollups.Set(body.InitBridgeAccountAction.RollupId.GetInner(), height, 0)
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

func parseBridgeLock(body *astria.Action_BridgeLockAction, from bytes.HexBytes, height types.Level, ctx *Context, action *storage.Action) error {
	action.Type = storageTypes.ActionTypeBridgeLock
	action.Data = make(map[string]any)
	if body.BridgeLockAction != nil {
		amount := uint128ToString(body.BridgeLockAction.Amount)

		action.Data["to"] = hex.EncodeToString(body.BridgeLockAction.To.GetInner())
		action.Data["destination_chain_address"] = body.BridgeLockAction.DestinationChainAddress
		action.Data["asset_id"] = body.BridgeLockAction.GetAssetId()
		action.Data["fee_asset_id"] = body.BridgeLockAction.GetFeeAssetId()
		action.Data["amount"] = amount

		toAddress := bytes.HexBytes(body.BridgeLockAction.To.GetInner())
		decAmount := decimal.RequireFromString(amount)
		toAddr := ctx.Addresses.Set(toAddress, height, decAmount, 1, 0)

		if stdBytes.Equal(from, toAddress) {
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
