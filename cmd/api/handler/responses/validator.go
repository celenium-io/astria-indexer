// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"fmt"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/pkg/types"
)

type ShortValidator struct {
	Id          uint64 `example:"321"                                      json:"id"      swaggertype:"integer"`
	ConsAddress string `example:"E641C7A2C964833E556AEF934FBF166B712874B6" json:"address" swaggertype:"string"`
	Name        string `example:"Node0"                                    json:"name"    swaggertype:"string"`
}

func NewShortValidator(val *storage.Validator) *ShortValidator {
	if val == nil || val.Id == 0 { // for genesis block
		return nil
	}
	return &ShortValidator{
		Id:          val.Id,
		ConsAddress: val.Address,
		Name:        val.Name,
	}
}

type Validator struct {
	Id          uint64 `example:"321"                                                              json:"id"          swaggertype:"integer"`
	ConsAddress string `example:"E641C7A2C964833E556AEF934FBF166B712874B6"                         json:"address"     swaggertype:"string"`
	Name        string `example:"Node0"                                                            json:"name"        swaggertype:"string"`
	PubkeyType  string `example:"tendermint/PubKeyEd25519"                                         json:"pubkey_type" swaggertype:"string"`
	Pubkey      string `example:"a497aa4a22ca8232876082920b110678988c86194b0c2e12a04dcf6f53688bb2" json:"pubkey"      swaggertype:"string"`
	Power       string `example:"100"                                                              json:"power"       swaggertype:"string"`
}

func NewValidator(val *storage.Validator) *Validator {
	if val == nil || val.Id == 0 { // for genesis block
		return nil
	}
	return &Validator{
		Id:          val.Id,
		ConsAddress: val.Address,
		Name:        val.Name,
		PubkeyType:  val.PubkeyType,
		Pubkey:      hex.EncodeToString(val.PubKey),
		Power:       val.Power.String(),
	}
}

type ValidatorUptime struct {
	Uptime string         `example:"0.97" json:"uptime" swaggertype:"string"`
	Blocks []SignedBlocks `json:"blocks"`
}

type SignedBlocks struct {
	Height types.Level `example:"100"  json:"height" swaggertype:"integer"`
	Signed bool        `example:"true" json:"signed" swaggertype:"boolean"`
}

func NewValidatorUptime(levels []types.Level, currentLevel types.Level, count types.Level) (uptime ValidatorUptime) {
	var (
		levelIndex = 0
		blockIndex = 0
		threshold  = count
	)

	if threshold > currentLevel {
		threshold = currentLevel
	}

	uptime.Blocks = make([]SignedBlocks, threshold)
	for i := currentLevel; i > currentLevel-threshold; i-- {
		if levelIndex < len(levels) && levels[levelIndex] == i {
			levelIndex++
			uptime.Blocks[blockIndex] = SignedBlocks{
				Signed: true,
				Height: i,
			}
		} else {
			uptime.Blocks[blockIndex] = SignedBlocks{
				Signed: false,
				Height: i,
			}
		}
		blockIndex++
	}

	uptime.Uptime = fmt.Sprintf("%.4f", float64(levelIndex)/float64(threshold))
	return uptime
}
