package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martinzha28/bridge-platform/backend/internal/handlers"
	"github.com/martinzha28/bridge-platform/backend/internal/ws"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(db *pgxpool.Pool, rdb *redis.Client, hub *ws.Hub) http.Handler {
	mux := http.NewServeMux()

	h := handlers.New(db, rdb)

	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /api/v1/", h.APIRoot)
	mux.HandleFunc("GET /ws", hub.HandleUpgrade)

	return mux
}
