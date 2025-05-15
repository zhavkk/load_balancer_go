package config

import (
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Proxy     ProxyConfig     `yaml:"proxy"`
	Backends  []BackendConfig `yaml:"backends"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	DB        DBConfig        `yaml:"db"`
	Env       string          `yaml:"env"`
}

type ProxyConfig struct {
	Port      string `yaml:"port"`
	Algorithm string `yaml:"algorithm"`
}

type BackendConfig struct {
	URL string `yaml:"url"`
}

type RateLimitConfig struct {
	Enabled bool `yaml:"enabled"`
}

var logger *slog.Logger // to logger package

type DBConfig struct {
	DSN            string `yaml:"dsn"`
	UpdateInterval string `yaml:"update_interval"`
}

func MustLoad() *Config {
	var cfg Config
	configpath := os.Getenv("CONFIG_PATH")
	if configpath == "" {
		panic("config path is not set")
	}
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		panic("config file does not exist")
	}
	if err := cleanenv.ReadConfig(configpath, &cfg); err != nil {
		logger.Error("failed to read config", "error", err)
		os.Exit(1)
	}
	return &cfg
}
