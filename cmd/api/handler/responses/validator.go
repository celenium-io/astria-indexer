// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"fmt"

	"github.com/aopoltorzhicky/astria/internal/storage"
	"github.com/aopoltorzhicky/astria/pkg/types"
)

type Validator struct {
	Id          uint64 `example:"321"                                      json:"id"      swaggertype:"integer"`
	ConsAddress string `example:"E641C7A2C964833E556AEF934FBF166B712874B6" json:"address" swaggertype:"string"`
	Name        string `example:"Node0"                                    json:"name"    swaggertype:"string"`
}

func NewValidator(val *storage.Validator) *Validator {
	if val == nil || val.Id == 0 { // for genesis block
		return nil
	}
	return &Validator{
		Id:          val.Id,
		ConsAddress: val.Address,
		Name:        val.Name,
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
