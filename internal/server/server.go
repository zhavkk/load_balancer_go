// internal/server/server.go
package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/zhavkk/load_balancer_go/internal/handlers/clients"
	"github.com/zhavkk/load_balancer_go/internal/handlers/proxy"
	rl "github.com/zhavkk/load_balancer_go/internal/ratelimiter/usecase"
)

func New(
	port string,
	lb http.Handler,
	repo rl.RateLimitRepository,
) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	ch := clients.NewClientsHandler(repo)
	r.Post("/clients", ch.Create)
	r.Delete("/clients", ch.Delete)

	r.Handle("/*", proxy.NewProxyHandler(lb))

	return &http.Server{Addr: ":" + port, Handler: r}
}
