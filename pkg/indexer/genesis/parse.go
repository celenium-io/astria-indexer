// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package genesis

import (
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type parsedData struct {
	block          storage.Block
	addresses      map[string]*storage.Address
	balanceUpdates []storage.BalanceUpdate
	constants      []storage.Constant
	validators     []*storage.Validator
	supply         decimal.Decimal
}

func newParsedData() parsedData {
	return parsedData{
		addresses:      make(map[string]*storage.Address),
		balanceUpdates: make([]storage.BalanceUpdate, 0),
		constants:      make([]storage.Constant, 0),
		validators:     make([]*storage.Validator, 0),
		supply:         decimal.Zero,
	}
}

func (module *Module) parse(genesis types.Genesis) (parsedData, error) {
	data := newParsedData()
	block := storage.Block{
		Time:    genesis.GenesisTime,
		Height:  pkgTypes.Level(genesis.InitialHeight - 1),
		AppHash: []byte(genesis.AppHash),
		ChainId: genesis.ChainID,
		Txs:     make([]*storage.Tx, 0),
		Stats: &storage.BlockStats{
			Time:         genesis.GenesisTime,
			Height:       pkgTypes.Level(genesis.InitialHeight - 1),
			Fee:          decimal.Zero,
			SupplyChange: decimal.Zero,
		},
	}

	module.parseConstants(genesis.AppState, genesis.ConsensusParams, &data)

	if err := module.parseAccounts(genesis.AppState.Accounts, block.Height, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis accounts")
	}
	if err := module.parseValidators(genesis.Validators, block.Height, &data); err != nil {
		return data, errors.Wrap(err, "parse genesis validators")
	}

	block.Stats.SupplyChange = data.supply
	data.block = block
	return data, nil
}

func (module *Module) parseAccounts(accounts []types.Account, height pkgTypes.Level, data *parsedData) error {
	for i := range accounts {
		address := storage.Address{
			Height: height,
			Balance: &storage.Balance{
				Total:    decimal.RequireFromString(accounts[i].Balance.String()),
				Currency: currency.DefaultCurrency,
			},
		}

		hash, err := pkgTypes.HexFromString(accounts[i].Address)
		if err != nil {
			return err
		}
		address.Hash = hash
		data.addresses[address.String()] = &address

		data.supply = data.supply.Add(address.Balance.Total)

		data.balanceUpdates = append(data.balanceUpdates, storage.BalanceUpdate{
			Address:  &address,
			Update:   address.Balance.Total,
			Currency: address.Balance.Currency,
			Height:   0,
		})
	}
	return nil
}

func (module *Module) parseValidators(validators []types.Validator, height pkgTypes.Level, data *parsedData) error {
	for i := range validators {
		data.validators = append(data.validators, &storage.Validator{
			Address:    validators[i].Address,
			PubkeyType: validators[i].PubKey.Type,
			PubKey:     validators[i].PubKey.Value,
			Name:       validators[i].Name,
			Power:      decimal.RequireFromString(validators[i].Power),
		})

		if _, ok := data.addresses[validators[i].Address]; !ok {
			address := storage.Address{
				Height: height,
				Balance: &storage.Balance{
					Total:    decimal.Zero,
					Currency: currency.DefaultCurrency,
				},
			}

			hash, err := pkgTypes.HexFromString(validators[i].Address)
			if err != nil {
				return err
			}
			address.Hash = hash
			data.addresses[address.String()] = &address
		}
	}
	return nil
}
