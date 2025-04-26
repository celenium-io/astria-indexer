// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

type SearchResult struct {
	Value string `json:"value"`
	Type  string `json:"type"`
	Body  any    `json:"body,omitempty"`
}

func NewSearchResult(value, typ string, body any) SearchResult {
	return SearchResult{
		Value: value,
		Type:  typ,
		Body:  body,
	}
}
