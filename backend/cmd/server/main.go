package main

import (
	"context"
	"log"
	"net/http"

	"github.com/martinzha28/bridge-platform/backend/internal/config"
	"github.com/martinzha28/bridge-platform/backend/internal/database"
	"github.com/martinzha28/bridge-platform/backend/internal/repository"
	"github.com/martinzha28/bridge-platform/backend/internal/routes"
	"github.com/martinzha28/bridge-platform/backend/internal/ws"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	db, err := database.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	rdb, err := database.NewRedis(ctx, cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer rdb.Close()
	log.Println("Connected to Redis")

	gameRepo := repository.NewGameRepository(db)
	hub := ws.NewHub(gameRepo)
	router := routes.SetupRouter(db, rdb, hub, cfg)

	log.Printf("Server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
