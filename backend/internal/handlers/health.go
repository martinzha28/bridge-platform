package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pgStatus := "up"
	if err := h.DB.Ping(ctx); err != nil {
		pgStatus = "down"
	}

	redisStatus := "up"
	if err := h.Redis.Ping(ctx).Err(); err != nil {
		redisStatus = "down"
	}

	status := "ok"
	code := http.StatusOK
	if pgStatus == "down" || redisStatus == "down" {
		status = "degraded"
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   status,
		"postgres": pgStatus,
		"redis":    redisStatus,
	})
}

func (h *Handler) APIRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Bridge Platform API v1",
	})
}
