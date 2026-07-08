package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"start/internal/database"
	"start/internal/feature_flags"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pgUrl := os.Getenv("DB_POSTGRES_URL")
	pgPool, err := database.NewPostgresPool(ctx, pgUrl)
	if err != nil {
		log.Fatalf("Cannot connect with PostgreSQL: %v", err)
	}
	defer pgPool.Close()

	featureFlagRepo := feature_flags.NewFeatureFlagPostgresRepository(pgPool)

	featureFlagHandler := feature_flags.NewFeatureFlagHandler(featureFlagRepo)

	r.Get("/health", HealthHandler)

	r.Route("/feature-flag", func(r chi.Router) {
		r.Post("/", featureFlagHandler.CreateFlag)
	})

	port := ":8080"
	logger.Info("Hell yeah! Server is starting", "port", port)

	if err := http.ListenAndServe(port, r); err != nil {
		logger.Error("Upsi, ", "error", err)
		os.Exit(1)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}
