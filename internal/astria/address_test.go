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
