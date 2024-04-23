// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"io"

	"github.com/celenium-io/astria-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/lib/pq"
)

const (
	ChannelBlock = "blocks"
	ChannelHead  = "head"
	ChannelTx    = "tx"
)

var Models = []any{
	&State{},
	&Constant{},
	&Balance{},
	&BalanceUpdate{},
	&Address{},
	&Block{},
	&BlockStats{},
	&Tx{},
	&Action{},
	&Validator{},
	&Rollup{},
	&RollupAction{},
	&RollupAddress{},
	&AddressAction{},
	&BlockSignature{},
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Transaction interface {
	sdk.Transaction

	SaveActions(ctx context.Context, actions ...*Action) error
	SaveAddressActions(ctx context.Context, actions ...*AddressAction) error
	SaveAddresses(ctx context.Context, addresses ...*Address) (int64, error)
	SaveBalances(ctx context.Context, balances ...Balance) error
	SaveBalanceUpdates(ctx context.Context, updates ...BalanceUpdate) error
	SaveBlockSignatures(ctx context.Context, signs ...BlockSignature) error
	SaveConstants(ctx context.Context, constants ...Constant) error
	SaveRollupActions(ctx context.Context, actions ...*RollupAction) error
	SaveRollupAddresses(ctx context.Context, addresses ...*RollupAddress) error
	SaveRollups(ctx context.Context, rollups ...*Rollup) (int64, error)
	SaveTransactions(ctx context.Context, txs ...*Tx) error
	SaveValidators(ctx context.Context, validators ...*Validator) error
	RetentionBlockSignatures(ctx context.Context, height types.Level) error

	RollbackActions(ctx context.Context, height types.Level) (actions []Action, err error)
	RollbackAddressActions(ctx context.Context, height types.Level) (addrActions []AddressAction, err error)
	RollbackAddresses(ctx context.Context, height types.Level) (address []Address, err error)
	RollbackBalances(ctx context.Context, ids []uint64) error
	RollbackBalanceUpdates(ctx context.Context, height types.Level) ([]BalanceUpdate, error)
	RollbackBlockSignatures(ctx context.Context, height types.Level) (err error)
	RollbackBlockStats(ctx context.Context, height types.Level) (stats BlockStats, err error)
	RollbackBlock(ctx context.Context, height types.Level) error
	RollbackRollupActions(ctx context.Context, height types.Level) (rollupActions []RollupAction, err error)
	RollbackRollupAddresses(ctx context.Context, height types.Level) (err error)
	RollbackRollups(ctx context.Context, height types.Level) ([]Rollup, error)
	RollbackTxs(ctx context.Context, height types.Level) (txs []Tx, err error)
	RollbackValidators(ctx context.Context, height types.Level) (err error)
	UpdateAddresses(ctx context.Context, address ...*Address) error
	UpdateRollups(ctx context.Context, rollups ...*Rollup) error

	LastBlock(ctx context.Context) (block Block, err error)
	State(ctx context.Context, name string) (state State, err error)
	LastNonce(ctx context.Context, id uint64) (uint32, error)
	GetProposerId(ctx context.Context, address string) (uint64, error)
	Validators(ctx context.Context) ([]Validator, error)
}

type SearchResult struct {
	Id    uint64 `bun:"id"`
	Value string `bun:"value"`
	Type  string `bun:"type"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ISearch interface {
	Search(ctx context.Context, query string) ([]SearchResult, error)
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Notificator interface {
	Notify(ctx context.Context, channel string, payload string) error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Listener interface {
	io.Closer

	Subscribe(ctx context.Context, channels ...string) error
	Listen() chan *pq.Notification
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ListenerFactory interface {
	CreateListener() Listener
}
