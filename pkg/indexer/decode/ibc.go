// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import "github.com/shopspring/decimal"

type IbcTransfer struct {
	Amount   decimal.Decimal `json:"amount"`
	Denom    string          `json:"denom"`
	Receiver string          `json:"receiver"`
	Sender   string          `json:"sender"`
}
