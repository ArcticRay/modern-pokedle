package server

import (
	"fmt"
	"net/http"

	"github.com/ArcticRay/modern-pokedle/internal/config"
	"github.com/ArcticRay/modern-pokedle/internal/middleware"
	"github.com/ArcticRay/modern-pokedle/internal/observability"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	cfg    *config.Config
	logger *observability.Logger
	db     *pgxpool.Pool
	http   *http.Server
}

func New(cfg *config.Config, db *pgxpool.Pool, logger *observability.Logger) *Server {
	s := &Server{
		cfg:    cfg,
		logger: logger,
		db:     db,
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
