// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"math/big"
	"time"

	"github.com/celenium-io/astria-indexer/pkg/types"
)

type Genesis struct {
	GenesisTime     time.Time             `json:"genesis_time"`
	ChainID         string                `json:"chain_id"`
	InitialHeight   int64                 `json:"initial_height,string"`
	ConsensusParams types.ConsensusParams `json:"consensus_params"`
	AppHash         types.Hex             `json:"app_hash"`
	AppState        AppState              `json:"app_state"`
	Validators      []Validator           `json:"validators"`
}

type Validator struct {
	Address string `json:"address"`
	PubKey  PubKey `json:"pub_key"`
	Power   string `json:"power"`
	Name    string `json:"name"`
}

type AppState struct {
	Accounts                    []Account `json:"accounts"`
	AddressesPrefixes           Prefixes  `json:"address_prefixes"`
	AuthoritySudoAddress        Bech32m   `json:"authority_sudo_address"`
	IbcSudoAddress              Bech32m   `json:"ibc_sudo_address"`
	IbcRelayerAddresses         []Bech32m `json:"ibc_relayer_addresses"`
	NativeAssetBaseDenomination string    `json:"native_asset_base_denomination"`
	Fees                        Fees      `json:"fees"`
	AllowedFeeAssets            []string  `json:"allowed_fee_assets"`
}

type Account struct {
	Address Bech32m `json:"address"`
	Balance UInt128 `json:"balance"`
}

type Bech32m struct {
	Value string `json:"bech32m"`
}

type Prefixes struct {
	Base string `json:"base"`
}

type Fees struct {
	Transfer             Fee `json:"transfer"`
	RollupDataSubmission Fee `json:"rollup_data_submission"`
	Ics20Withdrawal      Fee `json:"ics20_withdrawal"`
	InitBridgeAccount    Fee `json:"init_bridge_account"`
	BridgeLock           Fee `json:"bridge_lock"`
	BridgeUnlock         Fee `json:"bridge_unlock"`
	BridgeSudoChange     Fee `json:"bridge_sudo_change"`
	BridgeTransfer       Fee `json:"bridge_transfer"`

	IbcRelay          Fee `json:"ibc_relay"`
	ValidatorUpdate   Fee `json:"validator_update"`
	FeeAssetChange    Fee `json:"fee_asset_change"`
	FeeChange         Fee `json:"fee_change"`
	IbcRelayerChange  Fee `json:"ibc_relayer_change"`
	SudoAddressChange Fee `json:"sudo_address_change"`
	IbcSudoChange     Fee `json:"ibc_sudo_change"`
	PriceFeed         Fee `json:"price_feed"`
}

type Fee struct {
	Base       UInt128 `json:"base"`
	Multiplier UInt128 `json:"multiplier"`
}

type UInt128 struct {
	Low  uint64 `json:"lo"`
	High uint64 `json:"hi"`
}

func (bi UInt128) String() string {
	b := new(big.Int)
	b = b.SetUint64(bi.High)
	b = b.Lsh(b, 64)

	lo := new(big.Int)
	lo = lo.SetUint64(bi.Low)
	b = b.Xor(b, lo)
	return b.String()
}
