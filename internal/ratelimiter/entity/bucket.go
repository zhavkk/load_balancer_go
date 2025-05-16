package entity

import (
	"sync"
	"time"
)

type LimitConfig struct {
	ClientID string
	RPS      int
	Burst    int
}

type TokenBucket struct {
	capacity   int
	tokens     int
	fillRate   float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity int, rps int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		fillRate:   float64(rps),
		lastRefill: time.Now(),
	}
}

func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	add := int(elapsed * b.fillRate)
	if add > 0 {
		b.tokens += add
		if b.tokens > b.capacity {
			b.tokens = b.capacity
		}
		b.lastRefill = now
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}
