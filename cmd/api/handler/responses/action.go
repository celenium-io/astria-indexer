// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"
)

type Action struct {
	Id       uint64           `example:"1"                                                                format:"int64"     json:"id"                swaggertype:"integer"`
	Height   pkgTypes.Level   `example:"1000"                                                             format:"int64"     json:"height"            swaggertype:"integer"`
	Time     time.Time        `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"              swaggertype:"string"`
	Position int64            `example:"1"                                                                format:"int64"     json:"position"          swaggertype:"integer"`
	Type     types.ActionType `example:"rollup_data_submission"                                           format:"string"    json:"type"              swaggertype:"string"`
	TxHash   string           `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty" swaggertype:"string"`

	Fee  *Fee           `json:"fee,omitempty"`
	Data map[string]any `json:"data"`
}

type Fee struct {
	Amount string `example:"1000" format:"string" json:"amount" swaggertype:"string"`
	Asset  string `example:"nria" format:"string" json:"asset"  swaggertype:"string"`
}

func NewFee(fee *storage.Fee) *Fee {
	if fee == nil {
		return nil
	}
	return &Fee{
		Amount: fee.Amount.String(),
		Asset:  fee.Asset,
	}
}

func NewAction(action storage.Action) Action {
	result := Action{
		Id:       action.Id,
		Height:   action.Height,
		Time:     action.Time,
		Position: action.Position,
		Type:     action.Type,
		Data:     action.Data,
		Fee:      NewFee(action.Fee),
	}

	return result
}

func NewActionWithTx(action storage.ActionWithTx) Action {
	result := Action{
		Id:       action.Id,
		Height:   action.Height,
		Time:     action.Time,
		Position: action.Position,
		Type:     action.Type,
		Data:     action.Data,
		Fee:      NewFee(action.Fee),
	}

	if action.Tx != nil {
		result.TxHash = hex.EncodeToString(action.Tx.Hash)
	}

	return result
}

func NewAddressAction(action storage.AddressAction) Action {
	result := Action{
		Id:     action.ActionId,
		Height: action.Height,
		Time:   action.Time,
		Type:   action.ActionType,
	}

	if action.Tx != nil {
		result.TxHash = hex.EncodeToString(action.Tx.Hash)
	}
	if action.Action != nil {
		result.Data = action.Action.Data
		result.Position = action.Action.Position
		result.Fee = NewFee(action.Action.Fee)
	}

	return result
}
