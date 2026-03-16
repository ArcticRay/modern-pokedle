package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ArcticRay/modern-pokedle/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
)

type Handler struct {
	oauthConfig *oauth2.Config
	authService *Service
	db          *pgxpool.Pool
}

func NewHandler(oauthConfig *oauth2.Config, authService *Service, db *pgxpool.Pool) *Handler {
	return &Handler{
		oauthConfig: oauthConfig,
		authService: authService,
		db:          db,
	}
}

func (h *Handler) HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "failed to exchange code", http.StatusInternalServerError)
		return
	}

	githubUser, err := GetGitHubUser(r.Context(), token.AccessToken)
	if err != nil {
		http.Error(w, "failed to get github user", http.StatusInternalServerError)
		return
	}

	user, err := database.UpsertUser(r.Context(), h.db, githubUser.ID, githubUser.Login, githubUser.AvatarURL)
	if err != nil {
		http.Error(w, "failed to upsert user", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.authService.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	if err := database.SaveRefreshToken(r.Context(), h.db, user.ID, refreshToken, time.Now().Add(h.authService.refreshTokenTTL)); err != nil {
		http.Error(w, "failed to save refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          githubUser.Login,
	})
}
