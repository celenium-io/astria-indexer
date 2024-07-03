// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package astria

import (
	"encoding/hex"

	"github.com/pactus-project/pactus/util/bech32m"
)

const (
	prefix = "astria"
)

func IsAddress(s string) bool {
	if len(s) != 45 {
		return false
	}
	p, _, err := bech32m.Decode(s)
	if err != nil {
		return false
	}
	return p == prefix
}

func EncodeAddress(b []byte) (string, error) {
	return bech32m.EncodeFromBase256(prefix, b)
}

func DecodeAddress(s string) ([]byte, error) {
	_, b, err := bech32m.DecodeToBase256(s)
	return b, err
}

func EncodeFromHex(s string) (string, error) {
	hash, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	return EncodeAddress(hash)
}
