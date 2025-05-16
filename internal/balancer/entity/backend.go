package entity

import (
	"net/url"
	"sync"
)

type Backend struct {
	URL       url.URL
	IsDead    bool
	ActiveReq int
	mu        sync.RWMutex
}

func (b *Backend) Inc() {
	b.mu.Lock()
	b.ActiveReq++
	b.mu.Unlock()
}

func (b *Backend) Dec() {
	b.mu.Lock()
	if b.ActiveReq > 0 {
		b.ActiveReq--
	}
	b.mu.Unlock()
}

func (b *Backend) SetDead(dead bool) {
	b.mu.Lock()
	b.IsDead = dead
	b.mu.Unlock()
}

func (b *Backend) RLock()   { b.mu.RLock() }
func (b *Backend) RUnlock() { b.mu.RUnlock() }
