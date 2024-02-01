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
		"3fff1c39b9d163bfb9bcbf9dfea78675f1b4bc2c": {
			Height: 1,
			Hash:   []byte{0x3f, 0xff, 0x1c, 0x39, 0xb9, 0xd1, 0x63, 0xbf, 0xb9, 0xbc, 0xbf, 0x9d, 0xfe, 0xa7, 0x86, 0x75, 0xf1, 0xb4, 0xbc, 0x2c},
			Balance: &storage.Balance{
				Id:       0,
				Total:    decimal.RequireFromString("500000000000000000000"),
				Currency: "nria",
			},
		},
		"2e046327a2ccac7c8f8018ed44e43184b502eb3e": {
			Height: 1,
			Hash:   []byte{0x2e, 0x04, 0x63, 0x27, 0xa2, 0xcc, 0xac, 0x7c, 0x8f, 0x80, 0x18, 0xed, 0x44, 0xe4, 0x31, 0x84, 0xb5, 0x02, 0xeb, 0x3e},
			Balance: &storage.Balance{
				Id:       0,
				Total:    decimal.RequireFromString("500000000000000000000"),
				Currency: "nria",
			},
		},
	}
	require.Equal(t, want, data.addresses)
}
