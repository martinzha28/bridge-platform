package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martinzha28/bridge-platform/backend/internal/auth"
	"github.com/martinzha28/bridge-platform/backend/internal/config"
	"github.com/martinzha28/bridge-platform/backend/internal/handlers"
	"github.com/martinzha28/bridge-platform/backend/internal/ws"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(db *pgxpool.Pool, rdb *redis.Client, hub *ws.Hub, cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	h := handlers.New(db, rdb, cfg)
	protected := auth.Middleware(cfg.JWTSecret)

	// Public
	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /api/v1/", h.APIRoot)
	mux.HandleFunc("POST /api/v1/auth/register", h.Register)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("POST /api/v1/auth/logout", h.Logout)

	// Protected
	mux.Handle("GET /api/v1/auth/me", protected(http.HandlerFunc(h.GetMe)))
	mux.Handle("PUT /api/v1/auth/me", protected(http.HandlerFunc(h.UpdateMe)))

	// The WebSocket is open to guests. A valid token still attaches the
	// player's identity; without one, HandleUpgrade assigns a throwaway ID.
	optionalAuth := auth.OptionalMiddleware(cfg.JWTSecret)
	mux.Handle("GET /ws", optionalAuth(http.HandlerFunc(hub.HandleUpgrade)))

	return mux
}
