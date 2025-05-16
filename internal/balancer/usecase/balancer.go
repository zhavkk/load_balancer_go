package usecase

import (
	"fmt"

	"github.com/zhavkk/load_balancer_go/internal/balancer/entity"
)

type LoadBalancer interface {
	Next() *entity.Backend
	MarkDead(b *entity.Backend)
}

type Config struct {
	Backends  []entity.Backend
	Algorithm string
}

func NewLoadBalancer(cfg Config) (LoadBalancer, error) {
	bs := make([]*entity.Backend, len(cfg.Backends))
	for i, be := range cfg.Backends {
		copyB := be
		bs[i] = &copyB
	}
	switch cfg.Algorithm {
	case "round-robin":
		return NewRoundRobin(bs), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", cfg.Algorithm)
	}
}
