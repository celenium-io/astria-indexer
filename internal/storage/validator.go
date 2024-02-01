// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	pkgTypes "github.com/aopoltorzhicky/astria/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IValidator interface {
	storage.Table[*Validator]
}

type Validator struct {
	bun.BaseModel `bun:"validator" comment:"Table with validators"`

	Id         uint64          `bun:"id,pk,notnull,autoincrement"                comment:"Unique internal identity"`
	Address    string          `bun:"address,unique:validator_address,type:text" comment:"Validator address"`
	PubkeyType string          `bun:"pubkey_type,type:text"                      comment:"Validator public key type"`
	PubKey     []byte          `bun:"pubkey,unique:validator_pubkey"             comment:"Validator public key"`
	Name       string          `bun:"name,type:text"                             comment:"Human-readable name for the validator"`
	Power      decimal.Decimal `bun:"power,type:numeric"                         comment:"Validator power"`
	Height     pkgTypes.Level  `bun:"height"                                     comment:"Height when validator was created"`
}

func (Validator) TableName() string {
	return "validator"
}
