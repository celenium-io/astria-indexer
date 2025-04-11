package responses

import (
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
)

type Price struct {
	CurrencyPair string    `example:"BTC/USDT"                  format:"string"    json:"currency_pair" swaggertype:"string"`
	Price        string    `example:"50.00"                     format:"string"    json:"price"         swaggertype:"string"`
	Time         time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"          swaggertype:"string"`
}

func NewPrice(price storage.Price) Price {
	return Price{
		CurrencyPair: price.CurrencyPair,
		Price:        price.Price.String(),
		Time:         price.Time,
	}
}
