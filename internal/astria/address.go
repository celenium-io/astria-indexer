// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package astria

import (
	"encoding/hex"

	"github.com/cosmos/btcutil/bech32"
	"github.com/pactus-project/pactus/util/bech32m"
)

const (
	Prefix       = "astria"
	PrefixCompat = "astriacompat"
)

func IsAddress(s string) bool {
	if len(s) != 45 {
		return false
	}
	p, _, err := bech32m.Decode(s)
	if err != nil {
		return false
	}
	return p == Prefix
}

func IsCompatAddress(s string) bool {
	if len(s) != 51 {
		return false
	}
	p, _, err := bech32.DecodeToBase256(s)
	if err != nil {
		return false
	}
	return p == PrefixCompat
}

func EncodeAddress(b []byte) (string, error) {
	return bech32m.EncodeFromBase256(Prefix, b)
}

func DecodeAddress(s string) (string, []byte, error) {
	p, b, err := bech32m.DecodeToBase256(s)
	return p, b, err
}

func CompatToAstria(s string) (string, error) {
	_, b, err := bech32.DecodeToBase256(s)
	if err != nil {
		return "", err
	}
	return EncodeAddress(b)
}

func EncodeFromHex(s string) (string, error) {
	hash, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	return EncodeAddress(hash)
}
