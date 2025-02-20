// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rollback

import (
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/celenium-io/astria-indexer/internal/storage/types"
	"github.com/stretchr/testify/require"
)

func Test_getActionSize(t *testing.T) {
	tests := []struct {
		name    string
		action  storage.Action
		want    int64
		wantErr bool
	}{
		{
			name: "test 1",
			action: storage.Action{
				Id:       3,
				Height:   1000,
				Time:     time.Now(),
				Position: 2,
				Type:     types.ActionTypeRollupDataSubmission,
				TxId:     1,
				Data: map[string]any{
					"data":      "+G6AhDuaygeCUgiUaN0ig7sPHLWZae8gW9rtKb4FEKSIiscjBInoAACAgxvZgqDlaFLJ2rb9OUtQRsM/meiHSoW2nSkIGJiW6fhUti+v16Ani2wgQDfXhYkgZylMwLhCXtawIhnoA8eVSnnsg/7jGQ==",
					"rollup_id": "GbqKuz5LVqMJ32dWxHuX4pjjpy2IRJ02oPrbHKc2ZTk=",
				},
			},
			want: 112,
		}, {
			name: "test 2",
			action: storage.Action{
				Id:       3,
				Height:   1000,
				Time:     time.Now(),
				Position: 2,
				Type:     types.ActionTypeRollupDataSubmission,
				TxId:     1,
				Data: map[string]any{
					"rollup_id": "GbqKuz5LVqMJ32dWxHuX4pjjpy2IRJ02oPrbHKc2ZTk=",
				},
			},
			wantErr: true,
		}, {
			name: "test 3",
			action: storage.Action{
				Id:       3,
				Height:   1000,
				Time:     time.Now(),
				Position: 2,
				Type:     types.ActionTypeRollupDataSubmission,
				TxId:     1,
				Data: map[string]any{
					"rollup_id": "GbqKuz5LVqMJ32dWxHuX4pjjpy2IRJ02oPrbHKc2ZTk=",
					"data":      123,
				},
			},
			wantErr: true,
		}, {
			name: "test 4",
			action: storage.Action{
				Id:       3,
				Height:   1000,
				Time:     time.Now(),
				Position: 2,
				Type:     types.ActionTypeRollupDataSubmission,
				TxId:     1,
				Data: map[string]any{
					"rollup_id": "GbqKuz5LVqMJ32dWxHuX4pjjpy2IRJ02oPrbHKc2ZTk=",
					"data":      "wrong string",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getActionSize(tt.action)
			require.True(t, (tt.wantErr == (err != nil)))
			require.Equal(t, tt.want, got)
		})
	}
}
