package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/zhavkk/load_balancer_go/internal/ratelimiter/entity"
	"github.com/zhavkk/load_balancer_go/internal/ratelimiter/usecase"
	"github.com/zhavkk/load_balancer_go/internal/storage"
)

type RateLimiterRepository interface {
	GetLimit(ctx context.Context, clientID string) (*entity.LimitConfig, error)
	ListLimits(ctx context.Context) ([]*entity.LimitConfig, error)
	SaveLimit(ctx context.Context, limit *entity.LimitConfig) error
	DeleteLimit(ctx context.Context, clientID string) error
}

type RateLimiterPostgres struct {
	storage *storage.Storage
}

func NewRateLimiterRepository(storage *storage.Storage) *RateLimiterPostgres {
	return &RateLimiterPostgres{storage: storage}
}

func (r *RateLimiterPostgres) GetLimit(ctx context.Context, clientID string) (*entity.LimitConfig, error) {
	const sqlGet = `
        SELECT client_id, req_per_sec, burst_capacity
          FROM client_limits
         WHERE client_id = $1;
    `
	row := r.storage.GetPool().QueryRow(ctx, sqlGet, clientID)

	var cfg entity.LimitConfig
	err := row.Scan(&cfg.ClientID, &cfg.RPS, &cfg.Burst)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, usecase.ErrLimitConfigNotFound
		}
		return nil, fmt.Errorf("GetLimit scan: %w", err)
	}
	return &cfg, nil
}

func (r *RateLimiterPostgres) ListLimits(ctx context.Context) ([]*entity.LimitConfig, error) {
	const sqlList = `
        SELECT client_id, req_per_sec, burst_capacity
          FROM client_limits;
    `
	rows, err := r.storage.GetPool().Query(ctx, sqlList)
	if err != nil {
		return nil, fmt.Errorf("ListLimits query: %w", err)
	}
	defer rows.Close()

	var result []*entity.LimitConfig
	for rows.Next() {
		var cfg entity.LimitConfig
		if err := rows.Scan(&cfg.ClientID, &cfg.RPS, &cfg.Burst); err != nil {
			return nil, fmt.Errorf("ListLimits scan: %w", err)
		}
		result = append(result, &cfg)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("ListLimits rows: %w", rows.Err())
	}
	return result, nil
}

func (r *RateLimiterPostgres) SaveLimit(ctx context.Context, limit *entity.LimitConfig) error {
	const sqlUpsert = `
        INSERT INTO client_limits(client_id, req_per_sec, burst_capacity)
        VALUES($1, $2, $3)
        ON CONFLICT (client_id) DO UPDATE
          SET req_per_sec    = EXCLUDED.req_per_sec,
              burst_capacity = EXCLUDED.burst_capacity;
    `
	if _, err := r.storage.GetPool().Exec(ctx, sqlUpsert, limit.ClientID, limit.RPS, limit.Burst); err != nil {
		return fmt.Errorf("SaveLimit exec: %w", err)
	}
	return nil
}

func (r *RateLimiterPostgres) DeleteLimit(ctx context.Context, clientID string) error {
	const sqlDelete = `
        DELETE FROM client_limits
         WHERE client_id = $1;
    `
	cmd, err := r.storage.GetPool().Exec(ctx, sqlDelete, clientID)
	if err != nil {
		return fmt.Errorf("DeleteLimit exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return usecase.ErrLimitConfigNotFound
	}
	return nil
}
