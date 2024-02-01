// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package testsuite

import (
	"crypto/rand"
	"encoding/hex"
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
