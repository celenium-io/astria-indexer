// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// swagger:enum ModuleName
/*
	ENUM(
		block,
		evidence,
		validator,
		version,
		generic
	)
*/
//go:generate go-enum --marshal --sql --values
type ModuleName string
