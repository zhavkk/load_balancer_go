package loadbalancer

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/zhavkk/load_balancer_go/internal/config"
	"github.com/zhavkk/load_balancer_go/internal/logger"
)

func main() {

	//TODO: config
	if err := godotenv.Load(); err != nil {
		logger.Log.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}
	cfg := config.MustLoad()
	//TODO: logger
	logger.Init(cfg.Env)

	//TODO init load balancer

	//TODO init database

	//TODO: init chi router

	//TODO: start server

	//TODO: graceful shutdown
}
