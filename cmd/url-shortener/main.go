package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/config"
	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/http-server/handlers/url/save"
	mwLogger "github.com/NeedMoreDoggos/pet-rest-api-go/internal/http-server/middleware/logger"
	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/lib/logger/sl"
	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger, err := setupLogger(cfg.Env)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("starting server...", slog.String("env", cfg.Env))
	logger.Debug("debug message are enebled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("cannot create storage: %w", sl.Err(err))
		os.Exit(1)
	}

	// storage.SaveURL("test", "test")
	fmt.Println(storage.DeleteURL("test"))
	_ = storage

	router := chi.NewRouter()
	//middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(logger, storage))
	router.Get("/{alias}", redirect.New(logger, storage))

	logger.Info("server started", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}

	logger.Error("server stopped")
}

func setupLogger(env string) (*slog.Logger, error) {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return nil, fmt.Errorf("unknown env: %s", env)
	}

	return logger, nil
}
