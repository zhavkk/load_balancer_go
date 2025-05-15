package postgresql

import (
	"context"

	"github.com/zhavkk/load_balancer_go/internal/ratelimiter/config"
)

type RateLimiterRepository interface {
	GetLimit(ctx context.Context, clientID string) (*config.LimitConfig, error)
	ListLimits(ctx context.Context) ([]*config.LimitConfig, error)
	SaveLimit(ctx context.Context, limit *config.LimitConfig) error
	DeleteLimit(ctx context.Context, clientID string) error
}
