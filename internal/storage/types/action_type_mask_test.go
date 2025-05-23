// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
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

	t.Run("rollup data submission", func(t *testing.T) {
		arr := []string{string(ActionTypeRollupDataSubmission)}

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

	t.Run("bridge unlock", func(t *testing.T) {
		arr := []string{string(ActionTypeBridgeUnlock)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("bridge lock", func(t *testing.T) {
		arr := []string{string(ActionTypeBridgeLock)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("bridge sudo change action", func(t *testing.T) {
		arr := []string{string(ActionTypeBridgeSudoChangeAction)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("fee change", func(t *testing.T) {
		arr := []string{string(ActionTypeFeeChange)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("init bridge account", func(t *testing.T) {
		arr := []string{string(ActionTypeInitBridgeAccount)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("ibc relayer change", func(t *testing.T) {
		arr := []string{string(ActionTypeIbcRelayerChange)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("ibc sudo change", func(t *testing.T) {
		arr := []string{string(ActionTypeIbcSudoChangeAction)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("bridge transfer", func(t *testing.T) {
		arr := []string{string(ActionTypeBridgeTransfer)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("recover ibc client", func(t *testing.T) {
		arr := []string{string(ActionTypeRecoverIbcClient)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("currency pairs change", func(t *testing.T) {
		arr := []string{string(ActionTypeCurrencyPairsChange)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("markets change", func(t *testing.T) {
		arr := []string{string(ActionTypeMarketsChange)}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, arr, mask.Strings())
	})

	t.Run("unknown", func(t *testing.T) {
		arr := []string{"unknown"}

		mask := NewActionTypeMask(arr...)
		require.Equal(t, []string{}, mask.Strings())
	})
}
