package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        string
	GitHubID  int64
	Username  string
	AvatarURL string
}

func UpsertUser(ctx context.Context, db *pgxpool.Pool, githubID int64, username, avatarURL string) (*User, error) {
	query := `
		INSERT INTO users (github_id, username, avatar_url)
		VALUES ($1, $2, $3)
		ON CONFLICT (github_id) DO UPDATE
			SET username   = EXCLUDED.username,
			    avatar_url = EXCLUDED.avatar_url,
			    updated_at = NOW()
		RETURNING id, github_id, username, avatar_url
	`

	var user User
	err := db.QueryRow(ctx, query, githubID, username, avatarURL).Scan(
		&user.ID,
		&user.GitHubID,
		&user.Username,
		&user.AvatarURL,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert user: %w", err)
	}

	return &user, nil
}
