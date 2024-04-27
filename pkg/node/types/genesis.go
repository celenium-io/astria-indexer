// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"encoding/json"
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
	AuthoritySudoKey            string    `json:"authority_sudo_key"`
	NativeAssetBaseDenomination string    `json:"native_asset_base_denomination"`
	IbcSudoAddress              string    `json:"ibc_sudo_address"`
}

type Account struct {
	Address string      `json:"address"`
	Balance json.Number `json:"balance"`
}
