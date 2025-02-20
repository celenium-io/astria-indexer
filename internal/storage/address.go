// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/astria-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type AddressListFilter struct {
	Limit  int
	Offset int
	Sort   storage.SortOrder
	Asset  string
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IAddress interface {
	storage.Table[*Address]

	ByHash(ctx context.Context, hash string) (Address, error)
	ListWithBalance(ctx context.Context, fltrs AddressListFilter) ([]Address, error)
}

// Address -
type Address struct {
	bun.BaseModel `bun:"address" comment:"Table with addresses."`

	Id            uint64      `bun:"id,pk,notnull,autoincrement"  comment:"Unique internal identity"`
	Height        types.Level `bun:"height"                       comment:"Block number of the first address occurrence."`
	Hash          string      `bun:"hash,unique:address_hash"     comment:"Address hash"`
	Nonce         uint32      `bun:"nonce"                        comment:"Nonce"`
	ActionsCount  int64       `bun:"actions_count"                comment:"Count of actions in which the address was involved"`
	SignedTxCount int64       `bun:"signed_tx_count"              comment:"Count of signed transactions"`
	IsBridge      bool        `bun:"is_bridge"                    comment:"Indicate whether the account is a bridge or not"`
	IsIbcRelayer  *bool       `bun:"is_ibc_relayer,default:false" comment:"Indicate whether the account is a IBC realyer or not"`

	Balance []*Balance `bun:"rel:has-many,join:id=id"`
}

// TableName -
func (Address) TableName() string {
	return "address"
}

func (address Address) String() string {
	return address.Hash
}
