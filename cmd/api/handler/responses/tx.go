// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	pkgTypes "github.com/celenium-io/astria-indexer/pkg/types"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
)

type Tx struct {
	Id           uint64         `example:"321"                                                              format:"int64"     json:"id"                  swaggertype:"integer"`
	Height       pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"              swaggertype:"integer"`
	Position     int64          `example:"11"                                                               format:"int64"     json:"position"            swaggertype:"integer"`
	ActionsCount int64          `example:"1"                                                                format:"int64"     json:"actions_count"       swaggertype:"integer"`
	Nonce        uint32         `example:"1"                                                                format:"int64"     json:"nonce"               swaggertype:"integer"`
	Hash         string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"hash"                swaggertype:"string"`
	Error        string         `example:"some error text"                                                  format:"string"    json:"error,omitempty"     swaggertype:"string"`
	Codespace    string         `example:"sdk"                                                              format:"string"    json:"codespace,omitempty" swaggertype:"string"`
	Signature    string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"string"    json:"signature"           swaggertype:"string"`
	Signer       string         `example:"115F94D8C98FFD73FE65182611140F0EDC7C3C94"                         format:"string"    json:"signer"              swaggertype:"string"`
	Time         time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                swaggertype:"string"`
	Status       types.Status   `example:"success"                                                          format:"string"    json:"status"              swaggertype:"string"`
	ActionTypes  []string       `example:"rollup_data_submission,transfer"                                  format:"string"    json:"action_types"        swaggertype:"string"`

	Actions []Action `json:"actions,omitempty"`
}

func NewTx(tx storage.Tx) Tx {
	result := Tx{
		Id:           tx.Id,
		Height:       tx.Height,
		Time:         tx.Time,
		Position:     tx.Position,
		ActionsCount: tx.ActionsCount,
		Nonce:        tx.Nonce,
		Status:       tx.Status,
		Error:        tx.Error,
		Codespace:    tx.Codespace,
		Hash:         hex.EncodeToString(tx.Hash),
		Signature:    hex.EncodeToString(tx.Signature),
		Actions:      make([]Action, len(tx.Actions)),
		ActionTypes:  types.NewActionTypeMaskBits(tx.ActionTypes).Strings(),
	}

	if tx.Signer != nil {
		result.Signer = tx.Signer.String()
	}

	for i := range tx.Actions {
		result.Actions[i] = NewAction(tx.Actions[i])
	}

	return result
}
