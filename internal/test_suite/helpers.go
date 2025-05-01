// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package testsuite

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"github.com/celenium-io/astria-indexer/internal/astria"
)

// Ptr - returns pointer of value  for testing purpose
//
//	one := Ptr(1) // one is pointer to int
func Ptr[T any](t T) *T {
	return &t
}

// MustHexDecode - returns decoded hex string, if it can't decode throws panic
//
//	data := MustHexDecode("deadbeaf")
func MustHexDecode(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// RandomHash - returns random hash with fixed size
func RandomHash(length int) []byte {
	hash := make([]byte, length)
	_, _ = rand.Read(hash)
	return hash
}

// RandomAddress - returns random address
func RandomAddress() string {
	hash := make([]byte, 20)
	_, _ = rand.Read(hash)
	val, _ := astria.EncodeAddress(hash)
	return val
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandomText - generates random string with fixed size
//
//	data := RandomText(10)
func RandomText(n int) string {
	b := make([]rune, n)
	for i := range b {
		ids, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		b[i] = letterRunes[ids.Int64()]
	}
	return string(b)
}
