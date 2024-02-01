// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package currency

import "github.com/shopspring/decimal"

type Denom string

const (
	Nria Denom = "nria"
	Ria  Denom = "ria"
)

const (
	DefaultCurrency = "nria"
)

func StringRia(val decimal.Decimal) string {
	return val.StringFixed(6)
}

func StringNria(val decimal.Decimal) string {
	return val.StringFixed(0)
}
