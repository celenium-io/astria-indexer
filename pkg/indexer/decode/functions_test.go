// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/hex"
	"testing"

	v1 "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	"github.com/stretchr/testify/require"
)

func Test_uint128ToString(t *testing.T) {
	tests := []struct {
		name string
		u    *v1.Uint128
		want string
	}{
		{
			name: "test 1",
			u: &v1.Uint128{
				Hi: 0,
				Lo: 1,
			},
			want: "1",
		}, {
			name: "test 2",
			u: &v1.Uint128{
				Hi: 0,
				Lo: 10,
			},
			want: "10",
		}, {
			name: "test 3",
			u: &v1.Uint128{
				Hi: 1,
				Lo: 0,
			},
			want: "18446744073709551616",
		}, {
			name: "test 4",
			u: &v1.Uint128{
				Hi: 1,
				Lo: 1,
			},
			want: "18446744073709551617",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uint128ToString(tt.u)
			require.EqualValues(t, tt.want, got)
		})
	}
}

func TestAddressFromPubKey(t *testing.T) {

	tests := []struct {
		name string
		pk   string
		want string
	}{
		{
			name: "test 1",
			pk:   "32415F09DBEE4297CC9A841C2C2312BF903FC53C48860D788AE66097355A585F",
			want: "astria1yvzeyceqqmdjwv6yfwmduywm8a9jlxhyj5tlx2",
		}, {
			name: "test 2",
			pk:   "09E29331B2FAD4CBD367986803484A2F544441485E8E736112D2AD49B83656CA",
			want: "astria1du65j67v3ncwl832czg04m6hs9f9f9gchacdam",
		}, {
			name: "test 3",
			pk:   "96F43A8448928F1E580864D69FE44E093C5A82A1D4A80C59086D7E67976CDA45",
			want: "astria1z90efkxf3l7h8ln9rqnpz9q0pmw8c0y5dvfdhe",
		}, {
			name: "test 4",
			pk:   "352b09264c7ca6e2b40845f589973eeeb1c1068fc336ba571714ec018760be06",
			want: "astria13cyhkel6fkkzskxtjntulcqtpwpvjvslfclvtg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk, err := hex.DecodeString(tt.pk)
			require.NoError(t, err)

			got, err := AddressFromPubKey(pk)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
