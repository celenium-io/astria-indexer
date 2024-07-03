// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.7
// Revision: bf63e108589bbd2327b13ec2c5da532aad234029
// Build Date: 2023-07-25T23:27:55Z
// Built By: goreleaser

package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

const (
	// ActionTypeTransfer is a ActionType of type transfer.
	ActionTypeTransfer ActionType = "transfer"
	// ActionTypeSequence is a ActionType of type sequence.
	ActionTypeSequence ActionType = "sequence"
	// ActionTypeValidatorUpdate is a ActionType of type validator_update.
	ActionTypeValidatorUpdate ActionType = "validator_update"
	// ActionTypeSudoAddressChange is a ActionType of type sudo_address_change.
	ActionTypeSudoAddressChange ActionType = "sudo_address_change"
	// ActionTypeIbcRelay is a ActionType of type ibc_relay.
	ActionTypeIbcRelay ActionType = "ibc_relay"
	// ActionTypeIcs20Withdrawal is a ActionType of type ics20_withdrawal.
	ActionTypeIcs20Withdrawal ActionType = "ics20_withdrawal"
	// ActionTypeIbcRelayerChange is a ActionType of type ibc_relayer_change.
	ActionTypeIbcRelayerChange ActionType = "ibc_relayer_change"
	// ActionTypeFeeAssetChange is a ActionType of type fee_asset_change.
	ActionTypeFeeAssetChange ActionType = "fee_asset_change"
	// ActionTypeInitBridgeAccount is a ActionType of type init_bridge_account.
	ActionTypeInitBridgeAccount ActionType = "init_bridge_account"
	// ActionTypeBridgeLock is a ActionType of type bridge_lock.
	ActionTypeBridgeLock ActionType = "bridge_lock"
	// ActionTypeBridgeUnlock is a ActionType of type bridge_unlock.
	ActionTypeBridgeUnlock ActionType = "bridge_unlock"
	// ActionTypeBridgeSudoChangeAction is a ActionType of type bridge_sudo_change_action.
	ActionTypeBridgeSudoChangeAction ActionType = "bridge_sudo_change_action"
	// ActionTypeFeeChange is a ActionType of type fee_change.
	ActionTypeFeeChange ActionType = "fee_change"
)

var ErrInvalidActionType = fmt.Errorf("not a valid ActionType, try [%s]", strings.Join(_ActionTypeNames, ", "))

var _ActionTypeNames = []string{
	string(ActionTypeTransfer),
	string(ActionTypeSequence),
	string(ActionTypeValidatorUpdate),
	string(ActionTypeSudoAddressChange),
	string(ActionTypeIbcRelay),
	string(ActionTypeIcs20Withdrawal),
	string(ActionTypeIbcRelayerChange),
	string(ActionTypeFeeAssetChange),
	string(ActionTypeInitBridgeAccount),
	string(ActionTypeBridgeLock),
	string(ActionTypeBridgeUnlock),
	string(ActionTypeBridgeSudoChangeAction),
	string(ActionTypeFeeChange),
}

// ActionTypeNames returns a list of possible string values of ActionType.
func ActionTypeNames() []string {
	tmp := make([]string, len(_ActionTypeNames))
	copy(tmp, _ActionTypeNames)
	return tmp
}

// ActionTypeValues returns a list of the values for ActionType
func ActionTypeValues() []ActionType {
	return []ActionType{
		ActionTypeTransfer,
		ActionTypeSequence,
		ActionTypeValidatorUpdate,
		ActionTypeSudoAddressChange,
		ActionTypeIbcRelay,
		ActionTypeIcs20Withdrawal,
		ActionTypeIbcRelayerChange,
		ActionTypeFeeAssetChange,
		ActionTypeInitBridgeAccount,
		ActionTypeBridgeLock,
		ActionTypeBridgeUnlock,
		ActionTypeBridgeSudoChangeAction,
		ActionTypeFeeChange,
	}
}

// String implements the Stringer interface.
func (x ActionType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ActionType) IsValid() bool {
	_, err := ParseActionType(string(x))
	return err == nil
}

var _ActionTypeValue = map[string]ActionType{
	"transfer":                  ActionTypeTransfer,
	"sequence":                  ActionTypeSequence,
	"validator_update":          ActionTypeValidatorUpdate,
	"sudo_address_change":       ActionTypeSudoAddressChange,
	"ibc_relay":                 ActionTypeIbcRelay,
	"ics20_withdrawal":          ActionTypeIcs20Withdrawal,
	"ibc_relayer_change":        ActionTypeIbcRelayerChange,
	"fee_asset_change":          ActionTypeFeeAssetChange,
	"init_bridge_account":       ActionTypeInitBridgeAccount,
	"bridge_lock":               ActionTypeBridgeLock,
	"bridge_unlock":             ActionTypeBridgeUnlock,
	"bridge_sudo_change_action": ActionTypeBridgeSudoChangeAction,
	"fee_change":                ActionTypeFeeChange,
}

// ParseActionType attempts to convert a string to a ActionType.
func ParseActionType(name string) (ActionType, error) {
	if x, ok := _ActionTypeValue[name]; ok {
		return x, nil
	}
	return ActionType(""), fmt.Errorf("%s is %w", name, ErrInvalidActionType)
}

// MarshalText implements the text marshaller method.
func (x ActionType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ActionType) UnmarshalText(text []byte) error {
	tmp, err := ParseActionType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errActionTypeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *ActionType) Scan(value interface{}) (err error) {
	if value == nil {
		*x = ActionType("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseActionType(v)
	case []byte:
		*x, err = ParseActionType(string(v))
	case ActionType:
		*x = v
	case *ActionType:
		if v == nil {
			return errActionTypeNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errActionTypeNilPtr
		}
		*x, err = ParseActionType(*v)
	default:
		return errors.New("invalid type for ActionType")
	}

	return
}

// Value implements the driver Valuer interface.
func (x ActionType) Value() (driver.Value, error) {
	return x.String(), nil
}
