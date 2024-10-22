package decode

import "github.com/shopspring/decimal"

type IbcTransfer struct {
	Amount   decimal.Decimal `json:"amount"`
	Denom    string          `json:"denom"`
	Receiver string          `json:"receiver"`
	Sender   string          `json:"sender"`
}
