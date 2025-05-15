package config

import (
	"errors"
)

type LimitConfig struct {
	ClientID string
	RPS      int
	Burst    int
}

var (
	ErrLimitConfigNotFound = errors.New("limit config not found")
)
