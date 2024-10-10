// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

const (
	ActionTypeTransferBits Bits = 1 << iota
	ActionTypeSequenceBits
	ActionTypeValidatorUpdateBits
	ActionTypeSudoAddressChangeBits
	ActionTypeIbcRelayBits
	ActionTypeIcs20WithdrawalBits
	ActionTypeIbcRelayerChangeBits
	ActionTypeFeeAssetChangeBits
	ActionTypeInitBridgeAccountBits
	ActionTypeBridgeLockBits
	ActionTypeBridgeUnlockBits
	ActionTypeBridgeSudoChangeBits
	ActionTypeFeeChangeBits
	ActionTypeIbcSudoChangeBits
)

var (
	actionTypesMap = map[ActionType]Bits{
		ActionTypeIbcRelay:               ActionTypeIbcRelayBits,
		ActionTypeIcs20Withdrawal:        ActionTypeIcs20WithdrawalBits,
		ActionTypeSequence:               ActionTypeSequenceBits,
		ActionTypeSudoAddressChange:      ActionTypeSudoAddressChangeBits,
		ActionTypeTransfer:               ActionTypeTransferBits,
		ActionTypeValidatorUpdate:        ActionTypeValidatorUpdateBits,
		ActionTypeBridgeLock:             ActionTypeBridgeLockBits,
		ActionTypeFeeAssetChange:         ActionTypeFeeAssetChangeBits,
		ActionTypeInitBridgeAccount:      ActionTypeInitBridgeAccountBits,
		ActionTypeIbcRelayerChange:       ActionTypeIbcRelayerChangeBits,
		ActionTypeBridgeUnlock:           ActionTypeBridgeUnlockBits,
		ActionTypeBridgeSudoChangeAction: ActionTypeBridgeSudoChangeBits,
		ActionTypeFeeChange:              ActionTypeFeeChangeBits,
		ActionTypeIbcSudoChangeAction:    ActionTypeIbcSudoChangeBits,
	}
)

type ActionTypeMask struct {
	Bits
}

func NewActionTypeMask(vals ...string) ActionTypeMask {
	mask := ActionTypeMask{Bits: 0}
	for i := range vals {
		switch vals[i] {
		case string(ActionTypeIbcRelay):
			mask.Set(ActionTypeIbcRelayBits)
		case string(ActionTypeIcs20Withdrawal):
			mask.Set(ActionTypeIcs20WithdrawalBits)
		case string(ActionTypeSequence):
			mask.Set(ActionTypeSequenceBits)
		case string(ActionTypeSudoAddressChange):
			mask.Set(ActionTypeSudoAddressChangeBits)
		case string(ActionTypeTransfer):
			mask.Set(ActionTypeTransferBits)
		case string(ActionTypeValidatorUpdate):
			mask.Set(ActionTypeValidatorUpdateBits)
		case string(ActionTypeBridgeLock):
			mask.Set(ActionTypeBridgeLockBits)
		case string(ActionTypeFeeAssetChange):
			mask.Set(ActionTypeFeeAssetChangeBits)
		case string(ActionTypeIbcRelayerChange):
			mask.Set(ActionTypeIbcRelayerChangeBits)
		case string(ActionTypeInitBridgeAccount):
			mask.Set(ActionTypeInitBridgeAccountBits)
		case string(ActionTypeBridgeUnlock):
			mask.Set(ActionTypeBridgeUnlockBits)
		case string(ActionTypeBridgeSudoChangeAction):
			mask.Set(ActionTypeBridgeSudoChangeBits)
		case string(ActionTypeFeeChange):
			mask.Set(ActionTypeFeeChangeBits)
		case string(ActionTypeIbcSudoChangeAction):
			mask.Set(ActionTypeIbcSudoChangeBits)
		}
	}

	return mask
}

func NewActionTypeMaskBits(bits Bits) ActionTypeMask {
	return ActionTypeMask{Bits: bits}
}

func (mask ActionTypeMask) Strings() []string {
	if mask.Bits == 0 {
		return []string{}
	}

	vals := make([]string, 0)
	for val := ActionTypeTransferBits; val <= ActionTypeIbcSudoChangeBits; val <<= 1 {
		if !mask.Has(val) {
			continue
		}
		switch val {
		case ActionTypeIbcRelayBits:
			vals = append(vals, string(ActionTypeIbcRelay))
		case ActionTypeIcs20WithdrawalBits:
			vals = append(vals, string(ActionTypeIcs20Withdrawal))
		case ActionTypeSequenceBits:
			vals = append(vals, string(ActionTypeSequence))
		case ActionTypeSudoAddressChangeBits:
			vals = append(vals, string(ActionTypeSudoAddressChange))
		case ActionTypeTransferBits:
			vals = append(vals, string(ActionTypeTransfer))
		case ActionTypeValidatorUpdateBits:
			vals = append(vals, string(ActionTypeValidatorUpdate))
		case ActionTypeBridgeLockBits:
			vals = append(vals, string(ActionTypeBridgeLock))
		case ActionTypeFeeAssetChangeBits:
			vals = append(vals, string(ActionTypeFeeAssetChange))
		case ActionTypeIbcRelayerChangeBits:
			vals = append(vals, string(ActionTypeIbcRelayerChange))
		case ActionTypeInitBridgeAccountBits:
			vals = append(vals, string(ActionTypeInitBridgeAccount))
		case ActionTypeBridgeSudoChangeBits:
			vals = append(vals, string(ActionTypeBridgeSudoChangeAction))
		case ActionTypeBridgeUnlockBits:
			vals = append(vals, string(ActionTypeBridgeUnlock))
		case ActionTypeFeeChangeBits:
			vals = append(vals, string(ActionTypeFeeChange))
		case ActionTypeIbcSudoChangeBits:
			vals = append(vals, string(ActionTypeIbcSudoChangeAction))
		}
	}

	return vals
}

func (mask ActionTypeMask) Empty() bool {
	return mask.Bits == 0
}

func (mask *ActionTypeMask) SetType(typ ActionType) {
	value, ok := actionTypesMap[typ]
	if !ok {
		return
	}
	mask.Set(value)
}
