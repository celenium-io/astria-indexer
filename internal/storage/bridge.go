// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBridge interface {
	storage.Table[*Bridge]

	ByAddress(ctx context.Context, addressId uint64) (Bridge, error)
	ByRollup(ctx context.Context, rollupId uint64, limit, offset int) ([]Bridge, error)
	ByRoles(ctx context.Context, addressId uint64, limit, offset int) ([]Bridge, error)
}

type Bridge struct {
	bun.BaseModel `bun:"table:bridge" comment:"Table with bridges"`

	Id           uint64         `bun:"id,pk,notnull,autoincrement"         comment:"Unique internal identity"`
	RollupId     uint64         `bun:"rollup_id,unique:bridge_rollup_id"   comment:"Rollup id"`
	AddressId    uint64         `bun:"address_id,unique:bridge_address_id" comment:"Address id"`
	Asset        string         `bun:"asset"                               comment:"Asset"`
	FeeAsset     string         `bun:"fee_asset"                           comment:"Fee asset"`
	SudoId       uint64         `bun:"sudo_id"                             comment:"Address which is authorized to change the bridge"`
	WithdrawerId uint64         `bun:"withdrawer_id"                       comment:"Address which is used to make withdrawals from the bridge account"`
	InitHeight   pkgTypes.Level `bun:"init_height"                         comment:"Height when bridge was initialized"`

	Rollup     *Rollup  `bun:"rel:has-one,join:rollup_id=id"`
	Address    *Address `bun:"rel:has-one,join:address_id=id"`
	Sudo       *Address `bun:"rel:has-one,join:sudo_id=id"`
	Withdrawer *Address `bun:"rel:has-one,join:withdrawer_id=id"`
}

func (Bridge) TableName() string {
	return "bridge"
}
