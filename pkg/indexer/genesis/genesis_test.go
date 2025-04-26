// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package genesis

import (
	"os"
	"testing"

	"github.com/celenium-io/astria-indexer/internal/storage"
	testsuite "github.com/celenium-io/astria-indexer/internal/test_suite"
	"github.com/celenium-io/astria-indexer/pkg/indexer/config"
	"github.com/celenium-io/astria-indexer/pkg/node/types"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
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

	module := NewModule(&postgres.Transactable{}, config.Indexer{})

	err = module.parseAccounts(g.AppState.Accounts, 1, &data)
	require.NoError(t, err)

	err = module.parseValidators(g.Validators, 1, &data)
	require.NoError(t, err)

	err = module.parseIbcRelayerAddresses(g.AppState.IbcRelayerAddresses, 1, &data)
	require.NoError(t, err)

	want := map[string]*storage.Address{
		"astria12znlg972crgwspaw2ljsnglsvs0yk2waf7wxqw": {
			Height: 1,
			Hash:   "astria12znlg972crgwspaw2ljsnglsvs0yk2waf7wxqw",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.RequireFromString("333333333333333333"),
					Currency: "nria",
				},
			},
		},
		"astria14arwxsq84vjj0eeywah8e80lt6ehgunkn7pxf7": {
			Height: 1,
			Hash:   "astria14arwxsq84vjj0eeywah8e80lt6ehgunkn7pxf7",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.RequireFromString("333333333333333333"),
					Currency: "nria",
				},
			},
		},
		"astria1umqxfss8wua2ptcclwf7z3ly9phvsjt0hxpcjm": {
			Height: 1,
			Hash:   "astria1umqxfss8wua2ptcclwf7z3ly9phvsjt0hxpcjm",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.RequireFromString("333333333333333333"),
					Currency: "nria",
				},
			},
		},
		"astria1hrukdtz6xxk4x0hjrvc2wlx9wqc86s6u6tjmhy": {
			Height: 1,
			Hash:   "astria1hrukdtz6xxk4x0hjrvc2wlx9wqc86s6u6tjmhy",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.RequireFromString("333333333333333333"),
					Currency: "nria",
				},
			},
		},
		"astria1je48an7fhxfl70nv5x4s4mphe7nszm6t0fryy7": {
			Height: 1,
			Hash:   "astria1je48an7fhxfl70nv5x4s4mphe7nszm6t0fryy7",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.RequireFromString("60"),
					Currency: "nria",
				},
			},
			IsIbcRelayer: testsuite.Ptr(true),
		},
		"astria10cgc54dxh3sdetsr03rkhhkt3vsn3r7j46yvqh": {
			Height: 1,
			Hash:   "astria10cgc54dxh3sdetsr03rkhhkt3vsn3r7j46yvqh",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.RequireFromString("1000000000000"),
					Currency: "nria",
				},
			},
		},
		"astria1lhd67thmzakayk37mh53qmnc6v568ptzhptunx": {
			Height: 1,
			Hash:   "astria1lhd67thmzakayk37mh53qmnc6v568ptzhptunx",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.Zero,
					Currency: "nria",
				},
			},
		},
		"astria1eradtvny7l2m23rf0r5k4mp4f8jp7y6njpc6k7": {
			Height: 1,
			Hash:   "astria1eradtvny7l2m23rf0r5k4mp4f8jp7y6njpc6k7",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.Zero,
					Currency: "nria",
				},
			},
		},
		"astria19adwu5yhtat0t8zv2fyzwl93uj5rypvnrfsyw0": {
			Height: 1,
			Hash:   "astria19adwu5yhtat0t8zv2fyzwl93uj5rypvnrfsyw0",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.Zero,
					Currency: "nria",
				},
			},
		},
		"astria1phym4uktjn6gjle226009ge7u82w0dgtszs8x2": {
			Height: 1,
			Hash:   "astria1phym4uktjn6gjle226009ge7u82w0dgtszs8x2",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.Zero,
					Currency: "nria",
				},
			},
		},
		"astria1x62tjjddjspjquk503ww6l2nck46vxjaz6nq4f": {
			Height: 1,
			Hash:   "astria1x62tjjddjspjquk503ww6l2nck46vxjaz6nq4f",
			Balance: []*storage.Balance{
				{
					Id:       0,
					Total:    decimal.Zero,
					Currency: "nria",
				},
			},
			IsIbcRelayer: testsuite.Ptr(true),
		},
	}
	require.Equal(t, want, data.addresses)
}
