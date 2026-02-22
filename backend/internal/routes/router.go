package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martinzha28/bridge-platform/backend/internal/handlers"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(db *pgxpool.Pool, rdb *redis.Client) http.Handler {
	mux := http.NewServeMux()

	h := handlers.New(db, rdb)

	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /api/v1/", h.APIRoot)

	return mux
}
