package routes

import (
	"net/http"

	"github.com/martinzha28/bridge-platform/backend/internal/handlers"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", handlers.HealthCheck)

	// API routes
	mux.HandleFunc("GET /api/v1/", handlers.APIRoot)

	return mux
}
