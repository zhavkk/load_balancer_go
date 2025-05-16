package usecase

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/zhavkk/load_balancer_go/internal/logger"
	"github.com/zhavkk/load_balancer_go/internal/ratelimiter/entity"
)

var (
	ErrLimitConfigNotFound = errors.New("limit config not found")
)

type RateLimitRepository interface {
	GetLimit(ctx context.Context, clientID string) (*entity.LimitConfig, error)
	ListLimits(ctx context.Context) ([]*entity.LimitConfig, error)
	SaveLimit(ctx context.Context, cfg *entity.LimitConfig) error
	DeleteLimit(ctx context.Context, clientID string) error
}

type RateLimiter struct {
	repo         RateLimitRepository
	buckets      map[string]*entity.TokenBucket
	mu           sync.RWMutex
	defaultBurst int
	defaultRPS   int
}

func New(repo RateLimitRepository, defaultRPS, defaultBurst int) (*RateLimiter, error) {
	rl := &RateLimiter{
		repo:         repo,
		buckets:      make(map[string]*entity.TokenBucket),
		defaultBurst: defaultBurst,
		defaultRPS:   defaultRPS,
	}
	ctx := context.Background()
	configs, err := repo.ListLimits(ctx)
	if err != nil {
		return nil, err
	}
	for _, cfg := range configs {
		rl.buckets[cfg.ClientID] = entity.NewTokenBucket(cfg.Burst, cfg.RPS)
	}
	return rl, nil
}

func (rl *RateLimiter) getBucket(clientID string) *entity.TokenBucket {
	rl.mu.RLock()
	bkt, ok := rl.buckets[clientID]
	rl.mu.RUnlock()
	if ok {
		return bkt
	}
	ctx := context.Background()
	cfg, err := rl.repo.GetLimit(ctx, clientID)
	if err == ErrLimitConfigNotFound {
		bkt = entity.NewTokenBucket(rl.defaultBurst, rl.defaultRPS)
	} else if err != nil {
		logger.Log.Error("failed to load limit config, using defaults", "client", clientID, "error", err)
		bkt = entity.NewTokenBucket(rl.defaultBurst, rl.defaultRPS)
	} else {
		bkt = entity.NewTokenBucket(cfg.Burst, cfg.RPS)
	}
	rl.mu.Lock()
	rl.buckets[clientID] = bkt
	rl.mu.Unlock()
	return bkt
}

func (rl *RateLimiter) Allow(clientID string) bool {
	bucket := rl.getBucket(clientID)
	allowed := bucket.Allow()
	if !allowed {
		logger.Log.Warn("rate limit exceeded", "client", clientID)
	}
	return allowed
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.RemoteAddr
		if !rl.Allow(clientID) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
