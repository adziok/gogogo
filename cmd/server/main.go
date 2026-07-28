package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"start/internal/auth"
	"start/internal/database"
	externalapi "start/internal/external_api"
	"start/internal/feature_flags"
	usagelog "start/internal/usage_log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	cfg, err := auth.LoadAuthConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	initCtx, cancelInit := context.WithTimeout(appCtx, 5*time.Second)

	pgUrl := os.Getenv("DB_POSTGRES_URL")
	pgPool, err := database.NewPostgresPool(initCtx, pgUrl)
	if err != nil {
		log.Fatalf("Cannot connect with PostgreSQL: %v", err)
	}
	defer pgPool.Close()

	clickHouseURL := os.Getenv("DB_CLICKHOUSE_URL")
	clickHouse, err := database.NewClickHouseClient(initCtx, clickHouseURL)
	if err != nil {
		log.Fatalf("Cannot connect to ClickHouse: %v", err)
	}
	cancelInit()
	defer clickHouse.Close()

	featureFlagRepo := feature_flags.NewFeatureFlagPostgresRepository(pgPool)
	featureFlagHandler := feature_flags.NewFeatureFlagHandler(featureFlagRepo)

	r.Get("/health", HealthHandler)

	r.Route("/feature-flag", func(r chi.Router) {
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
	externalApiHandler, err := externalapi.CreateExternalApiHandler(externalApiRepo)

	if err != nil {
		log.Fatalf("Cannot create Kafka producer: %v", err)
	}

	defer externalApiHandler.Close()

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

	usageLogRepository := usagelog.NewClickHouseUsageLogRepository(clickHouse)
	usageLogProcessor := usagelog.NewClickHouseProcessor(usageLogRepository)
	usagelog.RegisterUsageLogHandler(appCtx, usageLogProcessor)

	server := &http.Server{Addr: ":8080", Handler: r}
	logger.Info("Hell yeah! Server is starting", "port", server.Addr)

	go func() {
		<-appCtx.Done()
		// appCtx is already cancelled, so this deadline needs an independent parent.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP server shutdown failed", "error", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("HTTP server failed", "error", err)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}
