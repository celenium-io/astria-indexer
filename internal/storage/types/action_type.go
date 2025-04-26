// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum ActionType
/*
	ENUM(
		transfer,
		rollup_data_submission,
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
		fee_change,
		ibc_sudo_change_action,
		bridge_transfer,
		recover_ibc_client,
		currency_pairs_change,
		markets_change
	)
*/
//go:generate go-enum --marshal --sql --values --names
type ActionType string
