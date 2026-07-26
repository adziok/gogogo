package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"start/internal/auth"
	"start/internal/database"
	externalapi "start/internal/external_api"
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

	// Load Auth0 configuration
	cfg, err := auth.LoadAuthConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create JWT validator
	jwtValidator, err := auth.NewValidator(cfg.Domain, cfg.Audience)
	if err != nil {
		log.Fatalf("Failed to create validator: %v", err)
	}

	// Create HTTP middleware
	middleware, err := auth.NewMiddleware(jwtValidator)
	if err != nil {
		log.Fatalf("Failed to create middleware: %v", err)
	}

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
		r.Use(func(next http.Handler) http.Handler {
			return middleware.CheckJWT(next)
		})
		r.Use(auth.UserDetailsMiddleware)

		r.Post("/", featureFlagHandler.CreateFlag)
		r.Get("/", featureFlagHandler.DisplayFlags)
		r.Delete("/{id}", featureFlagHandler.DeleteFlag)
		r.Put("/{id}", featureFlagHandler.UpdateFlag)
	})

	externalApiRepo := externalapi.NewFeatureFlagExternalPostgresReposiotory(pgPool)
	externalApiHandler := externalapi.CreateExternalApiHandler(externalApiRepo)

	r.Route("/api", func(r chi.Router) {
		// Create JWT validator
		jwtValidator, err := auth.NewValidator(cfg.Domain, cfg.AudienceApi)
		if err != nil {
			log.Fatalf("Failed to create validator: %v", err)
		}

		// Create HTTP middleware
		middleware, err := auth.NewMiddleware(jwtValidator)
		if err != nil {
			log.Fatalf("Failed to create middleware: %v", err)
		}

		r.Use(func(next http.Handler) http.Handler {
			return middleware.CheckJWT(next)
		})
		r.Use(auth.UserDetailsMiddleware)

		r.Get("/{id}", externalApiHandler.GetByTenantAndName)
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
