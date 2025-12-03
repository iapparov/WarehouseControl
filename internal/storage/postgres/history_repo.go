package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/wb-go/wbf/retry"
	wbzlog "github.com/wb-go/wbf/zlog"

	"warehousecontrol/internal/domain/history"
	"warehousecontrol/internal/domain/item"
)

func (p *Postgres) GetItemsHistory(id string, from, to time.Time, action string, login string) ([]*history.History, error) {
	ctx := context.Background()

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
        SELECT id, item_id, action, changed_by, changed_by_login, changed_at,
               old_data, new_data
        FROM history
        WHERE 1 = 1
    `)

	args := []interface{}{}
	argIndex := 1

	if !from.IsZero() {
		queryBuilder.WriteString(fmt.Sprintf(" AND changed_at >= $%d", argIndex))
		args = append(args, from)
		argIndex++
	}
	if !to.IsZero() {
		queryBuilder.WriteString(fmt.Sprintf(" AND changed_at <= $%d", argIndex))
		args = append(args, to)
		argIndex++
	}
	if id != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND item_id = $%d", argIndex))
		args = append(args, id)
		argIndex++
	}
	if action != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND action = $%d", argIndex))
		args = append(args, action)
		argIndex++
	}
	if login != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND changed_by_login ILIKE $%d", argIndex))
		args = append(args, "%"+login+"%")
	}

	query := queryBuilder.String()

	rows, err := p.db.QueryWithRetry(ctx, retry.Strategy{Attempts: p.cfg.Attempts, Delay: p.cfg.Delay, Backoff: p.cfg.Backoffs}, query, args...)
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to execute get items history query")
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			wbzlog.Logger.Error().Err(err).Msg("Failed to close history rows")
		}
	}()

	histories := []*history.History{}

	for rows.Next() {
		h := &history.History{}

		var oldJSON, newJSON []byte

		err := rows.Scan(
			&h.ID,
			&h.ItemID,
			&h.Action,
			&h.ChangedBy,
			&h.ChangedByLogin,
			&h.ChangedAt,
			&oldJSON,
			&newJSON,
		)
		if err != nil {
			wbzlog.Logger.Error().Err(err).Msg("Failed to scan history row")
			return nil, err
		}

		if len(oldJSON) > 0 {
			var old item.Item
			if err := json.Unmarshal(oldJSON, &old); err != nil {
				return nil, err
			}
			h.OldItemSnapshot = old
		}

		if len(newJSON) > 0 {
			var nw item.Item
			if err := json.Unmarshal(newJSON, &nw); err != nil {
				return nil, err
			}
			h.NewItemSnapshot = nw
		}
		h.ItemDiff = h.OldItemSnapshot.Diff(h.NewItemSnapshot)
		histories = append(histories, h)
	}

	return histories, nil
}
