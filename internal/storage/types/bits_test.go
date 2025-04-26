// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBits(t *testing.T) {

	t.Run("set", func(t *testing.T) {
		b := Bits(0)
		b.Set(1)
		require.EqualValues(t, 1, b)

		b.Set(1)
		require.EqualValues(t, 1, b)

		b.Set(2)
		require.EqualValues(t, 3, b)
	})

	t.Run("clear", func(t *testing.T) {
		b := Bits(255)
		b.Clear(2)
		require.EqualValues(t, 253, b)

		b.Clear(2)
		require.EqualValues(t, 253, b)

		b.Clear(255)
		require.EqualValues(t, 0, b)
	})

	t.Run("toggle", func(t *testing.T) {
		b := Bits(255)
		b.Toggle(2)
		require.EqualValues(t, 253, b)

		b.Toggle(2)
		require.EqualValues(t, 255, b)
	})

	t.Run("has", func(t *testing.T) {
		b := Bits(255)

		require.True(t, b.Has(1))
		require.True(t, b.Has(2))
		require.True(t, b.Has(4))
		require.True(t, b.Has(8))
		require.True(t, b.Has(16))
		require.True(t, b.Has(32))
		require.True(t, b.Has(64))
		require.True(t, b.Has(128))
		require.False(t, b.Has(256))
	})
}
