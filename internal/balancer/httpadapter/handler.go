package httpadapter

import (
	"net/http"

	"github.com/zhavkk/load_balancer_go/internal/balancer/entity"
	"github.com/zhavkk/load_balancer_go/internal/logger"
	"github.com/zhavkk/load_balancer_go/internal/proxy"
)

type LoadBalancer interface {
	Next() *entity.Backend
	MarkDead(b *entity.Backend)
}

type Handler struct {
	lb LoadBalancer
}

func NewHandler(lb LoadBalancer) *Handler {
	return &Handler{lb: lb}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := h.lb.Next()
	if backend == nil {
		http.Error(w, "no healthy backends", http.StatusBadGateway)
		return
	}

	rp := proxy.NewReverseProxy(&backend.URL)

	rp.SetErrorHandler(func(w http.ResponseWriter, req *http.Request, err error) {
		logger.Log.Error("backend error, marking dead", "error", err, "backend", backend.URL.String())
		h.lb.MarkDead(backend)
		http.Error(w, "bad gateway", http.StatusBadGateway)
	})

	backend.Inc()
	defer backend.Dec()

	rp.ServeHTTP(w, r)
}
