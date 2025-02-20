// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

type Bits uint64

func (b *Bits) Set(flag Bits)     { *b |= flag }
func (b *Bits) Clear(flag Bits)   { *b &^= flag }
func (b *Bits) Toggle(flag Bits)  { *b ^= flag }
func (b Bits) Has(flag Bits) bool { return b&flag != 0 }
