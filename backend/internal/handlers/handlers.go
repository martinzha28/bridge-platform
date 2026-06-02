package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martinzha28/bridge-platform/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	DB     *pgxpool.Pool
	Redis  *redis.Client
	Config config.Config
}

func New(db *pgxpool.Pool, rdb *redis.Client, cfg config.Config) *Handler {
	return &Handler{DB: db, Redis: rdb, Config: cfg}
}
