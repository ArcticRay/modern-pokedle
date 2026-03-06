package server

import (
	"fmt"
	"net/http"

	"github.com/ArcticRay/modern-pokedle/internal/config"
	"github.com/ArcticRay/modern-pokedle/internal/observability"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	cfg    *config.Config
	logger *observability.Logger
	http   *http.Server
}

func New(cfg *config.Config, logger *observability.Logger) *Server {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	return &Server{
		cfg:    cfg,
		logger: logger,
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: r,
		},
	}
}

func (s *Server) Start() error {
	s.logger.Info("server listening", map[string]any{
		"port": s.cfg.Port,
	})
	return s.http.ListenAndServe()
}
