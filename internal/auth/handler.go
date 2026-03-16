package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type Handler struct {
	oauthConfig *oauth2.Config
	authService *Service
}

func NewHandler(oauthConfig *oauth2.Config, authService *Service) *Handler {
	return &Handler{
		oauthConfig: oauthConfig,
		authService: authService,
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

	accessToken, err := h.authService.GenerateAccessToken(fmt.Sprintf("%d", githubUser.ID))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
		"user":         githubUser.Login,
	})
}
