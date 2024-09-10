// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import "encoding/json"

type DenomMetadata struct {
	Code      int         `json:"code"`
	Log       string      `json:"log"`
	Info      string      `json:"info"`
	Index     json.Number `json:"index"`
	Key       []byte      `json:"key"`
	Value     []byte      `json:"value"`
	Height    json.Number `json:"height"`
	Codespace string      `json:"codespace"`
}

type DenomMetadataResponse struct {
	Response DenomMetadata `json:"response"`
}
