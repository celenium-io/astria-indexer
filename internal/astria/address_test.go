// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package astria

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsAddress(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "astria1lhpxecq5ffhq68dgu9s8y2g5h53jqw5cvudrkk",
			s:    "astria1lhpxecq5ffhq68dgu9s8y2g5h53jqw5cvudrkk",
			want: true,
		}, {
			name: "astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p",
			s:    "astria1lm45urgugesyhaymn68xww0m6g49zreqa32w7p",
			want: true,
		}, {
			name: "astria1475jkpuvznd44szgfz8wwdf9w6xh5dx9jwqgvz",
			s:    "astria1475jkpuvznd44szgfz8wwdf9w6xh5dx9jwqgvz",
			want: true,
		}, {
			name: "astria16rgmx2s86kk2r69rhjnvs9y44ujfhadc7yav9a",
			s:    "astria16rgmx2s86kk2r69rhjnvs9y44ujfhadc7yav9a",
			want: true,
		}, {
			name: "prefix16rgmx2s86kk2r69rhjnvs9y44ujfhadc7yav9a",
			s:    "prefix16rgmx2s86kk2r69rhjnvs9y44ujfhadc7yav9a",
			want: false,
		}, {
			name: "invalid",
			s:    "invalid",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAddress(tt.s)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIsCompatAddress(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "astriacompat1eg8hhey0n4untdvqqdvlyl0e7zx8wfcaz3l6wu",
			s:    "astriacompat1eg8hhey0n4untdvqqdvlyl0e7zx8wfcaz3l6wu",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCompatAddress(tt.s)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestEncodeFromHex(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "test 1",
			s:    "0DC9BAF2CB94F4897F2A569EF2A33EE1D4E7B50B",
			want: "astria1phym4uktjn6gjle226009ge7u82w0dgtszs8x2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeFromHex(tt.s)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCompatToAstria(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "astriacompat1eg8hhey0n4untdvqqdvlyl0e7zx8wfcaz3l6wu",
			s:    "astriacompat1eg8hhey0n4untdvqqdvlyl0e7zx8wfcaz3l6wu",
			want: "astria1eg8hhey0n4untdvqqdvlyl0e7zx8wfca48kglh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompatToAstria(tt.s)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
