package loadbalancer

import (
	"github.com/zhavkk/load_balancer_go/internal/app"
	"github.com/zhavkk/load_balancer_go/internal/config"
	"github.com/zhavkk/load_balancer_go/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	app, err := app.Setup(cfg)
	if err != nil {
		logger.Log.Error("setup failed", "error", err)
	}
	if err := app.Run(); err != nil {
		logger.Log.Error("run failed", "error", err)
	}
}
