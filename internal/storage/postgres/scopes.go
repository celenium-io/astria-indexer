// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"fmt"

	"github.com/celenium-io/astria-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

const (
	columnSize         = "size"
	columnActionsCount = "actions_count"
	columnTime         = "time"
)

func limitScope(q *bun.SelectQuery, limit int) *bun.SelectQuery {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return q.Limit(limit)
}

func offsetScope(q *bun.SelectQuery, offset int) *bun.SelectQuery {
	if offset > 0 {
		return q.Offset(offset)
	}
	return q
}

func sortScope(q *bun.SelectQuery, field string, sort sdk.SortOrder) *bun.SelectQuery {
	if sort != sdk.SortOrderAsc && sort != sdk.SortOrderDesc {
		sort = sdk.SortOrderAsc
	}
	return q.OrderExpr("? ?", bun.Ident(field), bun.Safe(sort))
}

func addressListFilter(query *bun.SelectQuery, fltrs storage.AddressListFilter) *bun.SelectQuery {
	query = limitScope(query, fltrs.Limit)
	query = sortScope(query, "id", fltrs.Sort)
	return query
}

func txFilter(query *bun.SelectQuery, fltrs storage.TxFilter) *bun.SelectQuery {
	query = limitScope(query, fltrs.Limit)
	query = sortScope(query, "tx.id", fltrs.Sort)
	query = offsetScope(query, fltrs.Offset)

	if !fltrs.ActionTypes.Empty() {
		query = query.Where("action_types & ? > 0", fltrs.ActionTypes.Bits)
	}

	if len(fltrs.Status) > 0 {
		query = query.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			for i := range fltrs.Status {
				sq = sq.WhereOr("status = ?", fltrs.Status[i])
			}
			return sq
		})
	}
	if fltrs.Height > 0 {
		query = query.Where("tx.height = ?", fltrs.Height)
	}

	if !fltrs.TimeFrom.IsZero() {
		query = query.Where("tx.time >= ?", fltrs.TimeFrom)
	}
	if !fltrs.TimeTo.IsZero() {
		query = query.Where("tx.time < ?", fltrs.TimeTo)
	}
	if fltrs.WithActions {
		query = query.Relation("Actions")
	}
	return query
}

func joinCelestials(query *bun.SelectQuery, prefix string, where string) *bun.SelectQuery {
	tableName := bun.Safe(fmt.Sprintf("%scelestials", prefix))
	return query.
		ColumnExpr("?4.id as ?, ?4.image_url as ?, ?4.change_id as ?, ?4.status as ?",
			makeCelestialsColumnName(prefix, "id"),
			makeCelestialsColumnName(prefix, "image_url"),
			makeCelestialsColumnName(prefix, "change_id"),
			makeCelestialsColumnName(prefix, "status"),
			tableName,
		).
		Join(`left join lateral (
		select * from celestial where celestial.address_id = ?
		order by celestial.change_id desc limit 1
	) ? on true`, bun.Safe(where), tableName)
}

func makeCelestialsColumnName(prefix, name string) bun.Safe {
	return bun.Safe(fmt.Sprintf("%scelestials__%s", prefix, name))
}
