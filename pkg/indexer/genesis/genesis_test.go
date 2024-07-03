// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package genesis

import (
	"os"
	"testing"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/postgres"
	"github.com/celenium-io/astria-indexer/pkg/indexer/config"
	"github.com/celenium-io/astria-indexer/pkg/node/types"
	"github.com/goccy/go-json"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestParseAccounts(t *testing.T) {
	f, err := os.Open("../../../test/json/genesis.json")
	require.NoError(t, err)
	defer f.Close()

	var g types.Genesis
	err = json.NewDecoder(f).Decode(&g)
	require.NoError(t, err)

	data := newParsedData()

	module := NewModule(postgres.Storage{}, config.Indexer{})

	err = module.parseAccounts(g.AppState.Accounts, 1, &data)
	require.NoError(t, err)

	want := map[string]*storage.Address{
		"astria1lhpxecq5ffhq68dgu9s8y2g5h53jqw5cvudrkk": {
			Height: 1,
			Hash:   "astria1lhpxecq5ffhq68dgu9s8y2g5h53jqw5cvudrkk",
			Balance: &storage.Balance{
				Id:       0,
				Total:    decimal.RequireFromString("333333333333333333"),
				Currency: "nria",
			},
		},
		"astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p": {
			Height: 1,
			Hash:   "astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p",
			Balance: &storage.Balance{
				Id:       0,
				Total:    decimal.RequireFromString("60"),
				Currency: "nria",
			},
		},
		"astria1c220qfmjrwqlk939ca5a5z2rjxryyr9m3ah8gl": {
			Height: 1,
			Hash:   "astria1c220qfmjrwqlk939ca5a5z2rjxryyr9m3ah8gl",
			Balance: &storage.Balance{
				Id:       0,
				Total:    decimal.RequireFromString("333333333333333333"),
				Currency: "nria",
			},
		},
		"astria1475jkpuvznd44szgfz8wwdf9w6xh5dx9jwqgvz": {
			Height: 1,
			Hash:   "astria1475jkpuvznd44szgfz8wwdf9w6xh5dx9jwqgvz",
			Balance: &storage.Balance{
				Id:       0,
				Total:    decimal.RequireFromString("333333333333333333"),
				Currency: "nria",
			},
		},
		"astria16rgmx2s86kk2r69rhjnvs9y44ujfhadc7yav9a": {
			Height: 1,
			Hash:   "astria16rgmx2s86kk2r69rhjnvs9y44ujfhadc7yav9a",
			Balance: &storage.Balance{
				Id:       0,
				Total:    decimal.RequireFromString("340282366920938463463374607431768211455"),
				Currency: "nria",
			},
		},
	}
	require.Equal(t, want, data.addresses)
}
