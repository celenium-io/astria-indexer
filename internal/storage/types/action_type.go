// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum ActionType
/*
	ENUM(
		transfer,
		sequence,
		validator_update,
		sudo_address_change,
		ibc_relay,
		ics20_withdrawal,
		ibc_relayer_change,
		fee_asset_change,
		init_bridge_account,
		bridge_lock,
		bridge_unlock,
		bridge_sudo_change_action,
		fee_change
	)
*/
//go:generate go-enum --marshal --sql --values --names
type ActionType string
