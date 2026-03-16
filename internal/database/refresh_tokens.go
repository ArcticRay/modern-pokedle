package database

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshToken struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

func SaveRefreshToken(ctx context.Context, db *pgxpool.Pool, userID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := db.Exec(ctx, query, userID, hashToken(token), expiresAt)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}

	return nil
}

func GetRefreshToken(ctx context.Context, db *pgxpool.Pool, token string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var rt RefreshToken
	err := db.QueryRow(ctx, query, hashToken(token)).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.ExpiresAt,
		&rt.RevokedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	return &rt, nil
}

func RevokeRefreshToken(ctx context.Context, db *pgxpool.Pool, token string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1
	`

	_, err := db.Exec(ctx, query, hashToken(token))
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}
