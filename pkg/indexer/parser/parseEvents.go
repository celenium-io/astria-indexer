// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/indexer/decode"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func parseEvents(events []types.Event, ctx *decode.Context) error {
	for i := range events {
		var err error
		switch events[i].Type {
		case "tx.fees":
			err = parseTxFees(events[i].Attributes, ctx)
		default:
			continue
		}

		if err != nil {
			return errors.Wrap(err, events[i].Type)
		}
	}
	return nil
}

func parseTxFees(attrs []types.EventAttribute, ctx *decode.Context) error {
	var (
		fee = new(storage.Fee)
		err error
	)
	for i := range attrs {
		switch attrs[i].Key {
		case "asset":
			fee.Asset = attrs[i].Value
			// TODO: think about general logic with IBC channels
			if fee.Asset == "ibc/704031c868fd3d3c84a1cfa8cb45deba4ea746b44697f7f4a6ed1b8f6c239b82" {
				fee.Asset = string(currency.Nria)
			}
		case "feeAmount":
			fee.Amount, err = decimal.NewFromString(attrs[i].Value)
			if err != nil {
				return err
			}
		case "actionType":
			fee.ActionType = attrs[i].Value
		default:
		}
	}

	ctx.AddFee(fee)
	return nil
}
