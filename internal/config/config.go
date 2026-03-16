package config

import (
	"fmt"
	"os"
	"time"
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

	GitHubCallbackURL string

	JWTAccessTokenTTL  time.Duration
	JWTRefreshTokenTTL time.Duration

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

	cfg.GitHubCallbackURL = getEnvStr("GITHUB_CALLBACK_URL", "http://localhost:8080/api/v1/auth/github/callback")
	cfg.JWTAccessTokenTTL = getEnvDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute)
	cfg.JWTRefreshTokenTTL = getEnvDuration("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour)

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return cfg, nil
}

func getEnvStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
