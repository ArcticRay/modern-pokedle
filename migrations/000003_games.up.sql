CREATE TABLE games (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    game_date     DATE        NOT NULL DEFAULT CURRENT_DATE,
    status        TEXT        NOT NULL DEFAULT 'in_progress'
                              CHECK (status IN ('in_progress', 'won', 'lost')),
    pokemon_id    INTEGER     NOT NULL,
    pokemon_name  TEXT        NOT NULL,
    guesses_count INTEGER     NOT NULL DEFAULT 0,
    started_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at  TIMESTAMPTZ,

    UNIQUE (user_id, game_date)
);

CREATE TABLE guesses (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id      UUID        NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    guess_number INTEGER     NOT NULL,
    pokemon_id   INTEGER     NOT NULL,
    pokemon_name TEXT        NOT NULL,
    result       JSONB       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (game_id, guess_number)
);

CREATE INDEX idx_games_user_id ON games(user_id);
CREATE INDEX idx_games_game_date ON games(game_date);
CREATE INDEX idx_guesses_game_id ON guesses(game_id);