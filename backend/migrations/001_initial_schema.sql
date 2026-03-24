CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      TEXT UNIQUE NOT NULL,
    rating        INTEGER NOT NULL DEFAULT 1500,
    games_played  INTEGER NOT NULL DEFAULT 0,
    karma         INTEGER NOT NULL DEFAULT 100,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE tournaments (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE games (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status           TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'abandoned')),
    scoring_type     TEXT NOT NULL DEFAULT 'IMP' CHECK (scoring_type IN ('IMP', 'MP', 'BAM', 'rubber')),
    tournament_id    UUID REFERENCES tournaments(id),
    north_user_id    UUID REFERENCES users(id),
    east_user_id     UUID REFERENCES users(id),
    south_user_id    UUID REFERENCES users(id),
    west_user_id     UUID REFERENCES users(id),
    board_number     INTEGER NOT NULL,
    seed             BIGINT NOT NULL,
    deal_pbn         TEXT NOT NULL,
    auction          TEXT NOT NULL DEFAULT '',
    contract         TEXT,
    declarer         SMALLINT CHECK (declarer BETWEEN 0 AND 3),
    declarer_tricks  SMALLINT CHECK (declarer_tricks BETWEEN 0 AND 13),
    score            INTEGER,
    play_sequence    TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at     TIMESTAMPTZ
);

CREATE INDEX idx_games_north_user ON games(north_user_id) WHERE north_user_id IS NOT NULL;
CREATE INDEX idx_games_east_user  ON games(east_user_id)  WHERE east_user_id IS NOT NULL;
CREATE INDEX idx_games_south_user ON games(south_user_id) WHERE south_user_id IS NOT NULL;
CREATE INDEX idx_games_west_user  ON games(west_user_id)  WHERE west_user_id IS NOT NULL;
CREATE INDEX idx_games_tournament ON games(tournament_id)  WHERE tournament_id IS NOT NULL;
CREATE INDEX idx_games_status     ON games(status);
CREATE INDEX idx_games_completed  ON games(completed_at DESC) WHERE completed_at IS NOT NULL;
