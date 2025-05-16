package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zhavkk/load_balancer_go/internal/config"
)

type Storage struct {
	db *pgxpool.Pool
}

var (
	ErrFailedToConnect = errors.New("failed to connect to database")
	ErrDBNotConnected  = errors.New("database is not connected")
)

func NewStorage(cfg *config.Config) (*Storage, error) {
	db, err := pgxpool.New(context.Background(), cfg.DB.DSN)

	if err != nil {
		return nil, ErrFailedToConnect
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, ErrDBNotConnected
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		s.db.Close()
	}
	return nil
}

func (s *Storage) GetPool() *pgxpool.Pool {
	return s.db
}
