// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"time"

	"github.com/celenium-io/astria-indexer/pkg/types"
	"github.com/pkg/errors"
)

func (r *Module) sync(ctx context.Context) {
	var blocksCtx context.Context
	blocksCtx, r.cancelReadBlocks = context.WithCancel(ctx)
	if err := r.readBlocks(blocksCtx); err != nil {
		r.Log.Err(err).Msg("while reading blocks")
		r.stopAll()
		return
	}

	if ctx.Err() != nil {
		return
	}

	ticker := time.NewTicker(time.Second * time.Duration(r.cfg.BlockPeriod))
	defer ticker.Stop()

	for {
		r.rollbackSync.Wait()

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			blocksCtx, r.cancelReadBlocks = context.WithCancel(ctx)
			if err := r.readBlocks(blocksCtx); err != nil && !errors.Is(err, context.Canceled) {
				r.Log.Err(err).Msg("while reading blocks by timer")
				r.stopAll()
				return
			}
		}
	}
}

func (r *Module) readBlocks(ctx context.Context) error {
	for {
		headLevel, err := r.headLevel(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		if level, _ := r.Level(); level == headLevel {
			time.Sleep(time.Second)
			continue
		}

		r.passBlocks(ctx, headLevel)
		return nil
	}
}

func (r *Module) passBlocks(ctx context.Context, head types.Level) {
	level, _ := r.Level()
	level += 1

	for ; level <= head; level++ {
		select {
		case <-ctx.Done():
			return
		default:
			if _, ok := r.taskQueue.Get(level); !ok {
				r.taskQueue.Set(level, struct{}{})
				r.pool.AddTask(level)
			}
		}
	}
}

func (r *Module) headLevel(ctx context.Context) (types.Level, error) {
	status, err := r.api.Status(ctx)
	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}
