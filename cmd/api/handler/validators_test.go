// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "test 1",
			address: "115F94D8C98FFD73FE65182611140F0EDC7C3C94",
			want:    true,
		}, {
			name:    "test 2",
			address: "B385E68E3A3A2D250C7C4024972576D698B9E748",
			want:    true,
		}, {
			name:    "test 3",
			address: "B385E68E3A3A2D250C7C4024972576D698B9E74811",
			want:    false,
		}, {
			name:    "test 4",
			address: "some_strange_address",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAddress(tt.address)
			require.Equal(t, tt.want, got, tt.name)
		})
	}
}

func TestRollupIdValidate(t *testing.T) {
	t.Run("validate rollup id", func(t *testing.T) {
		require.True(t, isHash("19BA8ABB3E4B56A309DF6756C47B97E298E3A72D88449D36A0FADB1CA7366539"))
		require.False(t, isHash("LLBA8ABB3E4B56A309DF6756C47B97E298E3A72D88449D36A0FADB1CA7366539"))
		require.False(t, isHash("invalid"))
		require.False(t, isHash(""))
	})
}
