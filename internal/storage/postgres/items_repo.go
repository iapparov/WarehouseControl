package postgres

import (
	"context"
	"database/sql"

	"warehousecontrol/internal/domain/item"

	"github.com/wb-go/wbf/retry"
	wbzlog "github.com/wb-go/wbf/zlog"
)

func (p *Postgres) CreateItem(item *item.Item, userID string, login string) error {
	ctx := context.Background()

	tx, err := p.setHistoryConfig(ctx, userID, login)
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to set history config")
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		INSERT INTO items (id, name, count, price)
		VALUES ($1, $2, $3, $4)
	`

	err = retry.DoContext(ctx, retry.Strategy{Attempts: p.cfg.Attempts, Delay: p.cfg.Delay, Backoff: p.cfg.Backoffs}, func() error {
		_, err := tx.ExecContext(ctx, query,
			item.ID,
			item.Name,
			item.Count,
			item.Price,
		)
		return err
	})

	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to execute create item query")
		return err
	}

	err = tx.Commit()
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to commit transaction")
		return err
	}

	return nil
}

func (p *Postgres) GetItems() ([]*item.Item, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, count, price
		FROM items
	`

	rows, err := p.db.QueryWithRetry(ctx, retry.Strategy{Attempts: p.cfg.Attempts, Delay: p.cfg.Delay, Backoff: p.cfg.Backoffs}, query)
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to execute get items query")
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			wbzlog.Logger.Error().Err(err).Msg("Failed to close history rows")
		}
	}()

	var items []*item.Item
	for rows.Next() {
		var it item.Item
		err := rows.Scan(
			&it.ID,
			&it.Name,
			&it.Count,
			&it.Price,
		)
		if err != nil {
			wbzlog.Logger.Error().Err(err).Msg("Failed to scan item row")
			return nil, err
		}
		items = append(items, &it)
	}

	return items, nil
}
func (p *Postgres) GetItem(uuid string) (*item.Item, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, count, price
		FROM items
		WHERE id = $1
	`

	var it item.Item
	row, err := p.db.QueryRowWithRetry(ctx, retry.Strategy{Attempts: p.cfg.Attempts, Delay: p.cfg.Delay, Backoff: p.cfg.Backoffs}, query, uuid)
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to execute get item query")
		return nil, err
	}

	err = row.Scan(
		&it.ID,
		&it.Name,
		&it.Count,
		&it.Price,
	)
	if err != nil && err != sql.ErrNoRows {
		wbzlog.Logger.Error().Err(err).Msg("Failed to scan item row")
		return nil, err
	}
	return &it, nil
}
func (p *Postgres) PutItem(item *item.Item, userID string, login string) error {
	ctx := context.Background()

	tx, err := p.setHistoryConfig(ctx, userID, login)
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to set history config")
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		UPDATE items
		SET name = $2, count = $3, price = $4
		WHERE id = $1
	`

	err = retry.DoContext(ctx, retry.Strategy{Attempts: p.cfg.Attempts, Delay: p.cfg.Delay, Backoff: p.cfg.Backoffs}, func() error {
		_, err := tx.ExecContext(ctx, query,
			item.ID,
			item.Name,
			item.Count,
			item.Price,
		)
		return err
	})

	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to execute update item query")
		return err
	}

	err = tx.Commit()
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to commit transaction")
		return err
	}

	return nil
}
func (p *Postgres) DeleteItem(uuid string, userID string, login string) error {
	ctx := context.Background()

	tx, err := p.setHistoryConfig(ctx, userID, login)
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to set history config")
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		DELETE FROM items
		WHERE id = $1
	`

	err = retry.DoContext(ctx, retry.Strategy{Attempts: p.cfg.Attempts, Delay: p.cfg.Delay, Backoff: p.cfg.Backoffs}, func() error {
		_, err := tx.ExecContext(ctx, query,
			uuid,
		)
		return err
	})
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to execute delete item query")
		return err
	}

	err = tx.Commit()
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to commit transaction")
		return err
	}

	return nil
}

func (p *Postgres) setHistoryConfig(ctx context.Context, userID string, login string) (*sql.Tx, error) {
	tx, err := p.db.Master.BeginTx(ctx, nil)
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("cant start transaction in create_booking")
		return nil, err
	}

	queryUser := `SELECT set_config('app.current_user', $1, true)`

	err = retry.DoContext(ctx, retry.Strategy{Attempts: p.cfg.Attempts, Delay: p.cfg.Delay, Backoff: p.cfg.Backoffs}, func() error {
		_, err := tx.ExecContext(ctx, queryUser, userID)
		return err
	})

	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to execute set current user query")
		return nil, err
	}

	queryLogin := `SELECT set_config('app.current_user_login', $1, true)`
	err = retry.DoContext(ctx, retry.Strategy{Attempts: p.cfg.Attempts, Delay: p.cfg.Delay, Backoff: p.cfg.Backoffs}, func() error {
		_, err := tx.ExecContext(ctx, queryLogin, login)
		return err
	})

	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Failed to execute set current user login query")
		return nil, err
	}
	return tx, nil
}
