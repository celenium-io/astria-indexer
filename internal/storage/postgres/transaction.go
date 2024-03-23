// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/uptrace/bun"

	models "github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

type Transaction struct {
	storage.Transaction
}

func BeginTransaction(ctx context.Context, tx storage.Transactable) (models.Transaction, error) {
	t, err := tx.BeginTransaction(ctx)
	return Transaction{t}, err
}

func (tx Transaction) SaveConstants(ctx context.Context, constants ...models.Constant) error {
	if len(constants) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&constants).Exec(ctx)
	return err
}

func (tx Transaction) SaveTransactions(ctx context.Context, txs ...*models.Tx) error {
	if len(txs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&txs).Returning("id").Exec(ctx)
	return err
}

type addedAddress struct {
	bun.BaseModel `bun:"address"`
	*models.Address

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveAddresses(ctx context.Context, addresses ...*models.Address) (int64, error) {
	if len(addresses) == 0 {
		return 0, nil
	}

	addr := make([]addedAddress, len(addresses))
	for i := range addresses {
		addr[i].Address = addresses[i]
	}

	_, err := tx.Tx().NewInsert().Model(&addr).
		Column("height", "hash", "nonce", "actions_count", "signed_tx_count").
		On("CONFLICT ON CONSTRAINT address_hash DO UPDATE").
		Set("actions_count = added_address.actions_count + EXCLUDED.actions_count").
		Set("signed_tx_count = added_address.signed_tx_count + EXCLUDED.signed_tx_count").
		Set("nonce = GREATEST(EXCLUDED.nonce, added_address.nonce)").
		Returning("xmax, id").
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	for i := range addr {
		if addr[i].Xmax == 0 {
			count++
		}
	}

	return count, err
}

func (tx Transaction) SaveBalances(ctx context.Context, balances ...models.Balance) error {
	if len(balances) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&balances).
		Column("id", "currency", "total").
		On("CONFLICT (id, currency) DO UPDATE").
		Set("total = EXCLUDED.total + balance.total").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveActions(ctx context.Context, actions ...*models.Action) error {
	if len(actions) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&actions).Returning("id").Exec(ctx)
	return err
}

func (tx Transaction) SaveValidators(ctx context.Context, validators ...*models.Validator) error {
	if len(validators) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&validators).
		On("CONFLICT ON CONSTRAINT validator_pubkey DO UPDATE").
		Set("power = EXCLUDED.power").
		Returning("id").
		Exec(ctx)
	return err
}

type addedRollup struct {
	bun.BaseModel `bun:"rollup"`
	*models.Rollup

	Xmax uint64 `bun:"xmax"`
}

func (tx Transaction) SaveRollups(ctx context.Context, rollups ...*models.Rollup) (int64, error) {
	if len(rollups) == 0 {
		return 0, nil
	}

	rs := make([]addedRollup, len(rollups))
	for i := range rollups {
		rs[i].Rollup = rollups[i]
	}

	query := tx.Tx().NewInsert().Model(&rs).
		Column("first_height", "astria_id", "actions_count", "size", "bridge_address_id").
		On("CONFLICT ON CONSTRAINT rollup_id DO UPDATE").
		Set("actions_count = added_rollup.actions_count + EXCLUDED.actions_count").
		Set("size = added_rollup.size + EXCLUDED.size").
		Set("bridge_address_id = case when EXCLUDED.bridge_address_id is not null then EXCLUDED.bridge_address_id else added_rollup.bridge_address_id end")

	if _, err := query.Returning("xmax, id").Exec(ctx); err != nil {
		return 0, err
	}

	var count int64
	for i := range rs {
		if rs[i].Xmax == 0 {
			count++
		}
	}

	return count, nil
}

func (tx Transaction) SaveRollupActions(ctx context.Context, actions ...*models.RollupAction) error {
	if len(actions) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&actions).Exec(ctx)
	return err
}

func (tx Transaction) SaveRollupAddresses(ctx context.Context, addresses ...*models.RollupAddress) error {
	if len(addresses) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&addresses).
		On("CONFLICT (rollup_id, address_id) DO NOTHING").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveAddressActions(ctx context.Context, actions ...*models.AddressAction) error {
	if len(actions) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&actions).Exec(ctx)
	return err
}

func (tx Transaction) SaveBlockSignatures(ctx context.Context, signs ...models.BlockSignature) error {
	if len(signs) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&signs).Exec(ctx)
	return err
}

func (tx Transaction) SaveBalanceUpdates(ctx context.Context, updates ...models.BalanceUpdate) error {
	if len(updates) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&updates).Exec(ctx)
	return err
}

func (tx Transaction) LastBlock(ctx context.Context) (block models.Block, err error) {
	err = tx.Tx().NewSelect().Model(&block).Order("id desc").Limit(1).Scan(ctx)
	return
}

func (tx Transaction) State(ctx context.Context, name string) (state models.State, err error) {
	err = tx.Tx().NewSelect().Model(&state).Where("name = ?", name).Scan(ctx)
	return
}

func (tx Transaction) RollbackBlock(ctx context.Context, height types.Level) error {
	_, err := tx.Tx().NewDelete().
		Model((*models.Block)(nil)).
		Where("height = ?", height).
		Exec(ctx)
	return err
}

func (tx Transaction) RollbackBlockStats(ctx context.Context, height types.Level) (stats models.BlockStats, err error) {
	_, err = tx.Tx().NewDelete().Model(&stats).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackAddresses(ctx context.Context, height types.Level) (address []models.Address, err error) {
	_, err = tx.Tx().NewDelete().Model(&address).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackTxs(ctx context.Context, height types.Level) (txs []models.Tx, err error) {
	_, err = tx.Tx().NewDelete().Model(&txs).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackActions(ctx context.Context, height types.Level) (actions []models.Action, err error) {
	_, err = tx.Tx().NewDelete().Model(&actions).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackValidators(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.Validator)(nil)).Where("height = ?", height).Exec(ctx)
	return
}
func (tx Transaction) RollbackBlockSignatures(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.BlockSignature)(nil)).
		Where("height = ?", height).Exec(ctx)
	return
}
func (tx Transaction) RollbackBalanceUpdates(ctx context.Context, height types.Level) (updates []models.BalanceUpdate, err error) {
	_, err = tx.Tx().NewDelete().Model(&updates).Where("height = ?", height).Returning("*").Exec(ctx)
	return
}

func (tx Transaction) RollbackAddressActions(ctx context.Context, height types.Level) (addrActions []models.AddressAction, err error) {
	_, err = tx.Tx().NewDelete().Model(&addrActions).
		Where("height = ?", height).
		Returning("*").
		Exec(ctx)
	return
}

func (tx Transaction) RollbackRollupActions(ctx context.Context, height types.Level) (rActions []models.RollupAction, err error) {
	_, err = tx.Tx().NewDelete().Model(&rActions).
		Where("height = ?", height).
		Returning("*").
		Exec(ctx)
	return
}

func (tx Transaction) RollbackRollupAddresses(ctx context.Context, height types.Level) (err error) {
	_, err = tx.Tx().NewDelete().Model((*models.RollupAddress)(nil)).
		Where("height = ?", height).Exec(ctx)
	return
}

func (tx Transaction) RollbackRollups(ctx context.Context, height types.Level) (rollups []models.Rollup, err error) {
	_, err = tx.Tx().NewDelete().
		Model(&rollups).
		Where("first_height = ?", height).
		Returning("*").
		Exec(ctx)
	return
}

func (tx Transaction) RollbackBalances(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := tx.Tx().NewDelete().
		Model((*models.Balance)(nil)).
		Where("id IN (?)", bun.In(ids)).
		Exec(ctx)
	return err
}

func (tx Transaction) UpdateAddresses(ctx context.Context, addresses ...*models.Address) error {
	if len(addresses) == 0 {
		return nil
	}
	values := tx.Tx().NewValues(&addresses)

	_, err := tx.Tx().NewUpdate().
		With("_data", values).
		Model((*models.Address)(nil)).
		TableExpr("_data").
		Set("actions_count = address.actions_count + _data.actions_count").
		Set("signed_tx_count = address.signed_tx_count + _data.signed_tx_count").
		Set("nonce = GREATEST(_data.nonce, address.nonce)").
		Where("address.id = _data.id").
		Exec(ctx)
	return err
}

func (tx Transaction) UpdateRollups(ctx context.Context, rollups ...*models.Rollup) error {
	if len(rollups) == 0 {
		return nil
	}
	values := tx.Tx().NewValues(&rollups)

	_, err := tx.Tx().NewUpdate().
		With("_data", values).
		Model((*models.Rollup)(nil)).
		TableExpr("_data").
		Set("actions_count = rollup.actions_count + _data.actions_count").
		Set("size = rollup.size + _data.size").
		Where("rollup.id = _data.id").
		Exec(ctx)
	return err
}

func (tx Transaction) LastNonce(ctx context.Context, id uint64) (uint32, error) {
	var nonce uint32
	_, err := tx.Tx().NewSelect().
		Model((*models.Tx)(nil)).
		Column("nonce").
		Where("signer_id = ?", id).
		Order("id desc").
		Limit(1).
		Exec(ctx, &nonce)
	return nonce, err
}

func (tx Transaction) GetProposerId(ctx context.Context, address string) (id uint64, err error) {
	err = tx.Tx().NewSelect().
		Model((*models.Validator)(nil)).
		Column("id").
		Where("address = ?", address).
		Limit(1).
		Scan(ctx, &id)
	return
}

func (tx Transaction) Validators(ctx context.Context) (validators []models.Validator, err error) {
	err = tx.Tx().NewSelect().
		Model(&validators).
		Column("id", "address", "pubkey").
		Scan(ctx)
	return
}

func (tx Transaction) RetentionBlockSignatures(ctx context.Context, height types.Level) error {
	_, err := tx.Tx().NewDelete().Model((*models.BlockSignature)(nil)).
		Where("height <= ?", height).
		Exec(ctx)
	return err
}
