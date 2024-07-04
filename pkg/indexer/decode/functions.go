// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"math/big"

	v1 "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	"github.com/celenium-io/astria-indexer/internal/astria"
	"github.com/cometbft/cometbft/crypto/ed25519"
)

func uint128ToString(u *v1.Uint128) string {
	val := new(big.Int)
	val = val.SetUint64(u.GetHi())
	val = val.Lsh(val, 64)

	low := new(big.Int).SetUint64(u.GetLo())
	val = val.Add(val, low)
	return val.Text(10)
}

func AddressFromPubKey(pk []byte) (string, error) {
	return astria.EncodeAddress(ed25519.PubKey(pk).Address())
}
