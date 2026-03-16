package game

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ArcticRay/modern-pokedle/internal/database"
	"github.com/ArcticRay/modern-pokedle/internal/middleware"
	"github.com/ArcticRay/modern-pokedle/internal/pokemon"
	"github.com/jackc/pgx/v5/pgxpool"
)

const MaxGuesses = 6

type Handler struct {
	db             *pgxpool.Pool
	pokemonService *pokemon.Service
}

func NewHandler(db *pgxpool.Pool, pokemonService *pokemon.Service) *Handler {
	return &Handler{
		db:             db,
		pokemonService: pokemonService,
	}
}

func (h *Handler) HandleStartGame(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	pokemonID := DailyPokemonID(721)

	dailyPokemon, err := h.pokemonService.GetPokemon(r.Context(), fmt.Sprintf("%d", pokemonID))
	if err != nil {
		http.Error(w, `{"error":"failed to get daily pokemon"}`, http.StatusInternalServerError)
		return
	}

	game, err := database.GetOrCreateGame(r.Context(), h.db, userID, dailyPokemon.ID, dailyPokemon.Name)
	if err != nil {
		http.Error(w, `{"error":"failed to get or create game"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

type GuessRequest struct {
	PokemonName string `json:"pokemon_name"`
}

func (h *Handler) HandleGuess(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req GuessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.PokemonName == "" {
		http.Error(w, `{"error":"pokemon_name is required"}`, http.StatusBadRequest)
		return
	}

	pokemonID := DailyPokemonID(721)
	targetPokemon, err := h.pokemonService.GetPokemon(r.Context(), fmt.Sprintf("%d", pokemonID))
	if err != nil {
		http.Error(w, `{"error":"failed to get daily pokemon"}`, http.StatusInternalServerError)
		return
	}

	game, err := database.GetOrCreateGame(r.Context(), h.db, userID, targetPokemon.ID, targetPokemon.Name)
	if err != nil {
		http.Error(w, `{"error":"failed to get game"}`, http.StatusInternalServerError)
		return
	}

	if game.Status != "in_progress" {
		http.Error(w, `{"error":"game is already over"}`, http.StatusBadRequest)
		return
	}

	guessPokemon, err := h.pokemonService.GetPokemon(r.Context(), req.PokemonName)
	if err != nil {
		http.Error(w, `{"error":"pokemon not found"}`, http.StatusNotFound)
		return
	}

	result := CompareGuess(*guessPokemon, *targetPokemon)

	newStatus := "in_progress"
	if guessPokemon.ID == targetPokemon.ID {
		newStatus = "won"
	} else if game.GuessesCount+1 >= MaxGuesses {
		newStatus = "lost"
	}

	if err := database.SaveGuess(r.Context(), h.db, game.ID, game.GuessesCount+1, guessPokemon.ID, guessPokemon.Name, result); err != nil {
		http.Error(w, `{"error":"failed to save guess"}`, http.StatusInternalServerError)
		return
	}

	if err := database.UpdateGameStatus(r.Context(), h.db, game.ID, newStatus); err != nil {
		http.Error(w, `{"error":"failed to update game"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"result":       result,
		"status":       newStatus,
		"guesses_left": MaxGuesses - (game.GuessesCount + 1),
	})
}
