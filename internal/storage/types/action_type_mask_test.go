// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActionTypeMask(t *testing.T) {

	t.Run("full set", func(t *testing.T) {
		mask := NewActionTypeMask(ActionTypeNames()...)
		require.Equal(t, ActionTypeNames(), mask.Strings())
	})

	t.Run("ibc relay", func(t *testing.T) {
		arr := []string{string(ActionTypeIbcRelay)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("transfer", func(t *testing.T) {
		arr := []string{string(ActionTypeTransfer)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("ics 20 withdrawal", func(t *testing.T) {
		arr := []string{string(ActionTypeIcs20Withdrawal)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("mint", func(t *testing.T) {
		arr := []string{string(ActionTypeMint)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("sequence", func(t *testing.T) {
		arr := []string{string(ActionTypeSequence)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("sudo address change", func(t *testing.T) {
		arr := []string{string(ActionTypeSudoAddressChange)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("validator update", func(t *testing.T) {
		arr := []string{string(ActionTypeValidatorUpdate)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("unknown", func(t *testing.T) {
		arr := []string{"unknown"}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, []string{}, mask.Strings())
	})
}
