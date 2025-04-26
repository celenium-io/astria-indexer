// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IState interface {
	storage.Table[*State]

	ByName(ctx context.Context, name string) (State, error)
}

// State -
type State struct {
	bun.BaseModel `bun:"state" comment:"Current indexer state"`

	Id              uint64          `bun:",pk,autoincrement"         comment:"Unique internal identity" json:"id"`
	Name            string          `bun:",unique:state_name"        comment:"Indexer name"             json:"name"`
	LastHeight      types.Level     `bun:"last_height"               comment:"Last block height"        json:"height"`
	LastHash        []byte          `bun:"last_hash"                 comment:"Last block hash"          json:"hash"`
	LastTime        time.Time       `bun:"last_time"                 comment:"Time of last block"       json:"time"`
	ChainId         string          `bun:"chain_id"                  comment:"Astria chain id"          json:"chain_id"`
	TotalTx         int64           `bun:"total_tx"                  comment:"Transactions count"       json:"tx"`
	TotalAccounts   int64           `bun:"total_accounts"            comment:"Accounts count"           json:"accounts"`
	TotalRollups    int64           `bun:"total_rollups"             comment:"Rollups count"            json:"rollups"`
	TotalValidators int             `bun:"total_validators"          comment:"Validators count"         json:"validators"`
	TotalSupply     decimal.Decimal `bun:"total_supply,type:numeric" comment:"Total supply"             json:"supply"`
	TotalBridges    int64           `bun:"total_bridges"             comment:"Count of bridges"         json:"bridges"`
	TotalBytes      int64           `bun:"total_bytes"               comment:"Total rollup bytes"       json:"bytes"`
}

// TableName -
func (State) TableName() string {
	return "state"
}
