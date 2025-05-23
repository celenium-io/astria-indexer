// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"github.com/dipdup-net/go-lib/database"
	"github.com/pkg/errors"
)

func createScripts(ctx context.Context, conn *database.Bun, dir, subFolder string, split bool) error {
	scriptsDir := filepath.Join(dir, subFolder)

	files, err := os.ReadDir(scriptsDir)
	if err != nil {
		return err
	}

	for i := range files {
		if files[i].IsDir() {
			continue
		}

		path := filepath.Join(scriptsDir, files[i].Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if split {
			queries := bytes.Split(raw, []byte{';'})
			if len(queries) == 0 {
				continue
			}

			for _, query := range queries {
				query = bytes.TrimLeft(query, "\n ")
				if len(query) == 0 {
					continue
				}
				if _, err := conn.DB().NewRaw(string(query)).Exec(ctx); err != nil {
					return errors.Wrapf(err, "creating %s '%s'", subFolder, files[i].Name())
				}
			}
		} else {
			if _, err := conn.DB().NewRaw(string(raw)).Exec(ctx); err != nil {
				return errors.Wrapf(err, "creating %s '%s'", subFolder, files[i].Name())
			}
		}
	}

	return nil
}
