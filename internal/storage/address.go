// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/hex"

	"github.com/celenium-io/astria-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

type AddressListFilter struct {
	Limit  int
	Offset int
	Sort   storage.SortOrder
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IAddress interface {
	storage.Table[*Address]

	ByHash(ctx context.Context, hash []byte) (Address, error)
	ListWithBalance(ctx context.Context, fltrs AddressListFilter) ([]Address, error)
}

// Address -
type Address struct {
	bun.BaseModel `bun:"address" comment:"Table with addresses."`

	Id            uint64      `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	Height        types.Level `bun:"height"                      comment:"Block number of the first address occurrence."`
	Hash          []byte      `bun:"hash,unique:address_hash"    comment:"Address hash."`
	Nonce         uint32      `bun:"nonce"                       comment:"Nonce"`
	ActionsCount  int64       `bun:"actions_count"               comment:"Count of actions in which the address was involved"`
	SignedTxCount int64       `bun:"signed_tx_count"             comment:"Count of signed transactions"`

	Balance *Balance `bun:"rel:has-one,join:id=id"`
}

// TableName -
func (Address) TableName() string {
	return "address"
}

func (address Address) String() string {
	return hex.EncodeToString(address.Hash)
}
