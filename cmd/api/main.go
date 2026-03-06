package main

import (
	"fmt"
	"os"

	"github.com/ArcticRay/modern-pokedle/internal/config"
	"github.com/ArcticRay/modern-pokedle/internal/observability"
	"github.com/ArcticRay/modern-pokedle/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := observability.NewLogger(string(cfg.Env))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting pokedle", map[string]any{
		"env":  cfg.Env,
		"port": cfg.Port,
	})

	srv := server.New(cfg, logger)
	if err := srv.Start(); err != nil {
		logger.Fatal("server error", map[string]any{"error": err})
	}
}
