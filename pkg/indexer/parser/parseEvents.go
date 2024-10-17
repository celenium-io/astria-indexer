// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	astria "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/asset/v1alpha1"
	"github.com/celenium-io/astria-indexer/internal/currency"
	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/indexer/decode"
	"github.com/celenium-io/astria-indexer/pkg/node"
	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
)

func parseEvents(ctx context.Context, events []types.Event, decodeCtx *decode.Context, api node.Api) error {
	for i := range events {
		var err error
		switch events[i].Type {
		case "tx.fees":
			err = parseTxFees(ctx, events[i].Attributes, decodeCtx, api)
		case "tx.deposit":
			err = parseTxDeposit(events[i].Attributes, decodeCtx)
		default:
			continue
		}

		if err != nil {
			return errors.Wrap(err, events[i].Type)
		}
	}
	return nil
}

var (
	assets = map[string]string{
		"704031c868fd3d3c84a1cfa8cb45deba4ea746b44697f7f4a6ed1b8f6c239b82": string(currency.Nria),
	}
)

func getAsset(ctx context.Context, api node.Api, val string) (string, error) {
	if !strings.HasPrefix(val, "ibc") {
		return val, nil
	}
	parts := strings.Split(val, "/")
	hash := parts[len(parts)-1]
	if asset, ok := assets[hash]; ok {
		return asset, nil
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	metadata, err := api.GetAssetInfo(timeoutCtx, hash)
	if err != nil {
		return "", errors.Wrap(err, val)
	}

	var response astria.DenomResponse
	if err := proto.Unmarshal(metadata.Response.Value, &response); err != nil {
		return "", errors.Wrap(err, val)
	}
	assets[hash] = response.GetDenom()
	return response.GetDenom(), nil
}

func parseTxFees(ctx context.Context, attrs []types.EventAttribute, decodeCtx *decode.Context, api node.Api) error {
	var (
		fee = new(storage.Fee)
		err error
	)
	for i := range attrs {
		switch attrs[i].Key {
		case "asset":
			asset, err := getAsset(ctx, api, attrs[i].Value)
			if err != nil {
				return err
			}
			fee.Asset = asset
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

	decodeCtx.AddFee(fee)
	return nil
}

func parseTxDeposit(attrs []types.EventAttribute, decodeCtx *decode.Context) error {
	deposit := new(storage.Deposit)
	var idx int64

	for i := range attrs {
		switch attrs[i].Key {
		case "bridgeAddress":
			deposit.Bridge = &storage.Bridge{
				Address: &storage.Address{
					Hash: attrs[i].Value,
				},
			}
		case "rollupId":
			hash, err := base64.URLEncoding.DecodeString(attrs[i].Value)
			if err != nil {
				return err
			}
			deposit.Rollup = &storage.Rollup{
				AstriaId: hash,
			}
		case "amount":
			amount, err := decimal.NewFromString(attrs[i].Value)
			if err != nil {
				return err
			}
			deposit.Amount = amount
		case "asset":
			deposit.Asset = attrs[i].Value
		case "destinationChainAddress":
			deposit.DestinationChainAddress = attrs[i].Value
		case "sourceTransactionId":
		case "sourceActionIndex":
			actionIndex, err := strconv.ParseInt(attrs[i].Value, 10, 64)
			if err != nil {
				return err
			}
			idx = actionIndex
		}
	}

	decodeCtx.AddDeposit(idx, deposit)
	return nil
}
