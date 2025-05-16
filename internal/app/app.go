// internal/app/app.go
package app

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	balentity "github.com/zhavkk/load_balancer_go/internal/balancer/entity"
	"github.com/zhavkk/load_balancer_go/internal/balancer/httpadapter"
	balhttp "github.com/zhavkk/load_balancer_go/internal/balancer/httpadapter"
	balusecase "github.com/zhavkk/load_balancer_go/internal/balancer/usecase"
	"github.com/zhavkk/load_balancer_go/internal/config"
	"github.com/zhavkk/load_balancer_go/internal/logger"
	rlus "github.com/zhavkk/load_balancer_go/internal/ratelimiter/usecase"
	"github.com/zhavkk/load_balancer_go/internal/repository/postgresql"
	"github.com/zhavkk/load_balancer_go/internal/server"
	"github.com/zhavkk/load_balancer_go/internal/storage"
)

type App struct {
	srv     *http.Server
	storage *storage.Storage
}

func Setup(cfg *config.Config) (*App, error) {

	st, err := storage.NewStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	repo := postgresql.NewRateLimiterRepository(st)

	rl, err := rlus.New(repo, cfg.RateLimit.DefaultRPS, cfg.RateLimit.DefaultBurst)
	if err != nil {
		st.Close()
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	var backends []balentity.Backend
	for _, bc := range cfg.Backends {
		u, err := url.Parse(bc.URL)
		if err != nil {
			st.Close()
			return nil, fmt.Errorf("invalid backend URL %q: %w", bc.URL, err)
		}
		backends = append(backends, balentity.Backend{URL: *u})
	}

	lb, err := balusecase.NewLoadBalancer(balusecase.Config{
		Backends:  backends,
		Algorithm: cfg.Proxy.Algorithm,
	})
	if err != nil {
		st.Close()
		return nil, fmt.Errorf("balancer: %w", err)
	}

	proxyHandler := balhttp.NewHandler(lb)

	var finalHandler http.Handler = proxyHandler
	if cfg.RateLimit.Enabled {
		finalHandler = rl.Middleware(proxyHandler).(*httpadapter.Handler)
	}

	srv := server.New(cfg.Proxy.Port, finalHandler.(*httpadapter.Handler), repo)

	return &App{srv: srv, storage: st}, nil
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("failed to start server", "error", err)
			panic(err)
		}
	}()
	logger.Log.Info("server started")
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger.Log.Info("server shutting down")
	if err := a.srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	logger.Log.Info("server shut down")
	return a.storage.Close()
}
