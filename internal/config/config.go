package config

import (
	"fmt"
	"os"
)

type Env string

const (
	EnvDevelopment Env = "development"
	EnvProduction  Env = "production"
	EnvTest        Env = "test"
)

type Config struct {
	Port int
	Env  Env

	DatabaseURL string

	RedisURL string

	KafkaBrokers []string

	GitHubClientID     string
	GitHubClientSecret string

	JWTSecret string
}

func Load() (*Config, error) {
	cfg := &Config{}
	var missing []string

	cfg.Port = 8080
	cfg.Env = Env(os.Getenv("APP_ENV"))

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}

	cfg.RedisURL = os.Getenv("REDIS_URL")
	if cfg.RedisURL == "" {
		missing = append(missing, "REDIS_URL")
	}

	cfg.GitHubClientID = os.Getenv("GITHUB_CLIENT_ID")
	if cfg.GitHubClientID == "" {
		missing = append(missing, "GITHUB_CLIENT_ID")
	}

	cfg.GitHubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	if cfg.GitHubClientSecret == "" {
		missing = append(missing, "GITHUB_CLIENT_SECRET")
	}

	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return cfg, nil
}
