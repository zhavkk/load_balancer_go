package main

import (
	"os"

	"github.com/zhavkk/load_balancer_go/internal/app"
	"github.com/zhavkk/load_balancer_go/internal/config"
	"github.com/zhavkk/load_balancer_go/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	logger.Init(cfg.Env)

	app, err := app.Setup(cfg)
	if err != nil {
		logger.Log.Error("setup failed", "error", err)
		os.Exit(1)
	}
	if err := app.Run(); err != nil {
		logger.Log.Error("run failed", "error", err)
		os.Exit(1)
	}
}
