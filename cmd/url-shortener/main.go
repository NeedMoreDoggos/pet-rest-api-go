package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/config"
	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/lib/logger/sl"
	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/storage/sqlite"
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

	//TODO: init router: chi, chi rander

	//TODO: run server
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
