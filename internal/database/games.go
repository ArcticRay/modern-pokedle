package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Game struct {
	ID           string
	UserID       string
	GameDate     time.Time
	Status       string
	PokemonID    int
	PokemonName  string
	GuessesCount int
	StartedAt    time.Time
	CompletedAt  *time.Time
}

type Guess struct {
	ID          string
	GameID      string
	GuessNumber int
	PokemonID   int
	PokemonName string
	Result      json.RawMessage
	CreatedAt   time.Time
}

func GetOrCreateGame(ctx context.Context, db *pgxpool.Pool, userID string, pokemonID int, pokemonName string) (*Game, error) {
	query := `
		INSERT INTO games (user_id, pokemon_id, pokemon_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, game_date) DO UPDATE
			SET user_id = EXCLUDED.user_id
		RETURNING id, user_id, game_date, status, pokemon_id, pokemon_name, guesses_count, started_at, completed_at
	`

	var game Game
	err := db.QueryRow(ctx, query, userID, pokemonID, pokemonName).Scan(
		&game.ID,
		&game.UserID,
		&game.GameDate,
		&game.Status,
		&game.PokemonID,
		&game.PokemonName,
		&game.GuessesCount,
		&game.StartedAt,
		&game.CompletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get or create game: %w", err)
	}

	return &game, nil
}

func SaveGuess(ctx context.Context, db *pgxpool.Pool, gameID string, guessNumber, pokemonID int, pokemonName string, result any) error {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	query := `
		INSERT INTO guesses (game_id, guess_number, pokemon_id, pokemon_name, result)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = db.Exec(ctx, query, gameID, guessNumber, pokemonID, pokemonName, resultJSON)
	if err != nil {
		return fmt.Errorf("save guess: %w", err)
	}

	return nil
}

func UpdateGameStatus(ctx context.Context, db *pgxpool.Pool, gameID, status string) error {
	query := `
		UPDATE games
		SET status       = $1,
		    guesses_count = guesses_count + 1,
		    completed_at  = CASE WHEN $1 != 'in_progress' THEN NOW() ELSE NULL END
		WHERE id = $2
	`

	_, err := db.Exec(ctx, query, status, gameID)
	if err != nil {
		return fmt.Errorf("update game status: %w", err)
	}

	return nil
}

func GetGameWithGuesses(ctx context.Context, db *pgxpool.Pool, gameID string) (*Game, []Guess, error) {
	game := &Game{}
	err := db.QueryRow(ctx, `
		SELECT id, user_id, game_date, status, pokemon_id, pokemon_name, guesses_count, started_at, completed_at
		FROM games WHERE id = $1
	`, gameID).Scan(
		&game.ID, &game.UserID, &game.GameDate, &game.Status,
		&game.PokemonID, &game.PokemonName, &game.GuessesCount,
		&game.StartedAt, &game.CompletedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("get game: %w", err)
	}

	rows, err := db.Query(ctx, `
		SELECT id, game_id, guess_number, pokemon_id, pokemon_name, result, created_at
		FROM guesses WHERE game_id = $1 ORDER BY guess_number
	`, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("get guesses: %w", err)
	}
	defer rows.Close()

	var guesses []Guess
	for rows.Next() {
		var g Guess
		if err := rows.Scan(&g.ID, &g.GameID, &g.GuessNumber, &g.PokemonID, &g.PokemonName, &g.Result, &g.CreatedAt); err != nil {
			return nil, nil, fmt.Errorf("scan guess: %w", err)
		}
		guesses = append(guesses, g)
	}

	return game, guesses, nil
}
