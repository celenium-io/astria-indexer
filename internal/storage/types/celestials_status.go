// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum CelestialsStatus
/*
	ENUM(
		NOT_VERIFIED,
		VERIFIED,
		PRIMARY
	)
*/
//go:generate go-enum --marshal --sql --values --names
type CelestialsStatus string
