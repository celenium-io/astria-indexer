// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/hex"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/shopspring/decimal"
)

type Validators map[string]*storage.Validator

func NewValidators() Validators {
	return make(map[string]*storage.Validator)
}

func (v Validators) Set(pubKey []byte, power int64) *storage.Validator {
	sPubKey := hex.EncodeToString(pubKey)

	pow := decimal.NewFromInt(power)
	if validator, ok := v[sPubKey]; ok {
		validator.Power = pow
		return validator
	}

	validator := &storage.Validator{
		PubKey: pubKey,
		Power:  pow,
	}
	v[sPubKey] = validator
	return validator
}
