// internal/handlers/clients/clients_handler.go
package clients

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zhavkk/load_balancer_go/internal/ratelimiter/entity"
	"github.com/zhavkk/load_balancer_go/internal/ratelimiter/usecase"
)

type ClientDTO struct {
	IP          string `json:"ip"`
	Capacity    int    `json:"capacity"`
	RefillEvery string `json:"refill_every"`
}

type ClientsHandler struct {
	Repo usecase.RateLimitRepository
}

func NewClientsHandler(repo usecase.RateLimitRepository) *ClientsHandler {
	return &ClientsHandler{Repo: repo}
}

// Create — POST /clients
func (h *ClientsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto ClientDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if dto.IP == "" || dto.Capacity <= 0 || dto.RefillEvery == "" {
		http.Error(w, "fields ip, capacity and refill_every are required", http.StatusBadRequest)
		return
	}
	dur, err := time.ParseDuration(dto.RefillEvery)
	if err != nil {
		http.Error(w, "invalid refill_every format", http.StatusBadRequest)
		return
	}
	rps := int(float64(dto.Capacity) / dur.Seconds())
	if rps < 1 {
		rps = 1
	}

	cfg := entity.LimitConfig{
		ClientID: dto.IP,
		RPS:      rps,
		Burst:    dto.Capacity,
	}
	if err := h.Repo.SaveLimit(r.Context(), &cfg); err != nil {
		http.Error(w, "could not save limit: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete — DELETE /clients
func (h *ClientsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IP string `json:"ip"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if body.IP == "" {
		http.Error(w, "field ip is required", http.StatusBadRequest)
		return
	}
	if err := h.Repo.DeleteLimit(r.Context(), body.IP); err != nil {
		if err == usecase.ErrLimitConfigNotFound {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "could not delete limit: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
