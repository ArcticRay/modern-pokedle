package main

import (
	"fmt"
	"os"

	"github.com/ArcticRay/modern-pokedle/internal/auth"
	"github.com/ArcticRay/modern-pokedle/internal/config"
	"github.com/ArcticRay/modern-pokedle/internal/database"
	"github.com/ArcticRay/modern-pokedle/internal/observability"
	"github.com/ArcticRay/modern-pokedle/internal/pokemon"
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

	db, err := database.NewPool(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to database", map[string]any{"error": err})
	}
	defer db.Close()
	logger.Info("database connected", map[string]any{})

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		logger.Fatal("failed to run migrations", map[string]any{"error": err})
	}
	logger.Info("migrations complete", map[string]any{})

	pokemonCache, err := pokemon.NewCache(cfg.RedisURL)
	if err != nil {
		logger.Fatal("failed to connect to redis", map[string]any{"error": err})
	}

	pokemonClient := pokemon.NewClient("https://pokeapi.co/api/v2")
	pokemonService := pokemon.NewService(pokemonClient, pokemonCache)

	githubOAuthConfig := auth.NewGitHubOAuthConfig(auth.GitHubConfig{
		ClientID:     cfg.GitHubClientID,
		ClientSecret: cfg.GitHubClientSecret,
		CallbackURL:  cfg.GitHubCallbackURL,
	})

	authService := auth.NewService(cfg.JWTSecret, cfg.JWTAccessTokenTTL, cfg.JWTRefreshTokenTTL)
	authHandler := auth.NewHandler(githubOAuthConfig, authService)

	srv := server.New(cfg, db, pokemonService, authHandler, logger)
	if err := srv.Start(); err != nil {
		logger.Fatal("server error", map[string]any{"error": err})
	}
}
