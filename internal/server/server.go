package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ArcticRay/modern-pokedle/internal/auth"
	"github.com/ArcticRay/modern-pokedle/internal/config"
	"github.com/ArcticRay/modern-pokedle/internal/game"
	"github.com/ArcticRay/modern-pokedle/internal/middleware"
	"github.com/ArcticRay/modern-pokedle/internal/observability"
	"github.com/ArcticRay/modern-pokedle/internal/pokemon"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	cfg            *config.Config
	logger         *observability.Logger
	db             *pgxpool.Pool
	pokemonService *pokemon.Service
	authHandler    *auth.Handler
	authService    *auth.Service
	http           *http.Server
}

func New(cfg *config.Config, db *pgxpool.Pool, pokemonService *pokemon.Service, authHandler *auth.Handler, authService *auth.Service, logger *observability.Logger) *Server {
	s := &Server{
		cfg:            cfg,
		logger:         logger,
		db:             db,
		pokemonService: pokemonService,
		authHandler:    authHandler,
		authService:    authService,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := s.db.Ping(r.Context()); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"status":"unhealthy","db":"unreachable"}`)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok","db":"reachable"}`)
	})

	gameHandler := game.NewHandler(s.db, s.pokemonService)

	r.Route("/api/v1", func(r chi.Router) {
		// public routes
		r.Get("/auth/github", s.authHandler.HandleGitHubLogin)
		r.Get("/auth/github/callback", s.authHandler.HandleGitHubCallback)

		// private routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(s.authService, s.logger))

			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				userID := middleware.UserIDFromContext(r.Context())
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
			})

			r.Route("/games", func(r chi.Router) {
				r.Post("/", gameHandler.HandleStartGame)
				r.Get("/today", gameHandler.HandleGetTodaysGame)
				r.Post("/guess", gameHandler.HandleGuess)
			})
		})
	})

	// Test Endpoints
	r.Get("/test/pokemon/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		p, err := s.pokemonService.GetPokemon(r.Context(), name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	})

	r.Get("/test/guess/{target}/{guess}", func(w http.ResponseWriter, r *http.Request) {
		targetName := chi.URLParam(r, "target")
		guessName := chi.URLParam(r, "guess")
		target, err := s.pokemonService.GetPokemon(r.Context(), targetName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		guess, err := s.pokemonService.GetPokemon(r.Context(), guessName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		result := game.CompareGuess(*guess, *target)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	return s
}

func (s *Server) Start() error {
	s.logger.Info("server listening", map[string]any{
		"port": s.cfg.Port,
	})
	return s.http.ListenAndServe()
}
