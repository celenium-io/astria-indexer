// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum AppCategory
/*
	ENUM(
		uncategorized,
		finance,
		gaming,
		nft,
		social
	)
*/
//go:generate go-enum --marshal --sql --values --names
type AppCategory string

// swagger:enum AppType
/*
	ENUM(
		sovereign,
		settled
	)
*/
//go:generate go-enum --marshal --sql --values --names
type AppType string
