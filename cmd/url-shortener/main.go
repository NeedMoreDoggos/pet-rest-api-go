package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/NeedMoreDoggos/pet-rest-api-go/internal/config"
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

	//TODO: init logger: slog

	//TODO: init storage: sqlite

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
