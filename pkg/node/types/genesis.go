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
	Balance BigInt  `json:"balance"`
}

type Bech32m struct {
	Value string `json:"bech32m"`
}

type Prefixes struct {
	Base string `json:"base"`
}

type Fees struct {
	TransferBaseFee              BigInt `json:"transfer_base_fee"`
	SequenceBaseFee              BigInt `json:"sequence_base_fee"`
	SequenceByteCostMultiplier   BigInt `json:"sequence_byte_cost_multiplier"`
	InitBridgeAccountBaseFee     BigInt `json:"init_bridge_account_base_fee"`
	BridgeLockByteCostMultiplier BigInt `json:"bridge_lock_byte_cost_multiplier"`
	Ics20WithdrawalBaseFee       BigInt `json:"ics20_withdrawal_base_fee"`
	BridgeSudoChangeFee          BigInt `json:"bridge_sudo_change_fee"`
}

type BigInt struct {
	Low  uint64 `json:"lo"`
	High uint64 `json:"hi"`
}

func (bi BigInt) String() string {
	b := new(big.Int)
	b = b.SetUint64(bi.High)
	b = b.Lsh(b, 64)

	lo := new(big.Int)
	lo = lo.SetUint64(bi.Low)
	b = b.Xor(b, lo)
	return b.String()
}
