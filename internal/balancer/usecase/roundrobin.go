package usecase

import (
	"sync/atomic"

	"github.com/zhavkk/load_balancer_go/internal/balancer/entity"
)

type RoundRobin struct {
	backends []*entity.Backend
	idx      uint64
}

func NewRoundRobin(b []*entity.Backend) *RoundRobin {
	return &RoundRobin{backends: b}
}

func (rr *RoundRobin) Next() *entity.Backend {
	n := uint64(len(rr.backends))
	for i := uint64(0); i < n; i++ {
		pos := atomic.AddUint64(&rr.idx, 1)
		be := rr.backends[pos%n]
		be.RLock()
		dead := be.IsDead
		be.RUnlock()
		if !dead {
			return be
		}
	}
	return nil
}

func (rr *RoundRobin) MarkDead(b *entity.Backend) {
	b.SetDead(true)
}
