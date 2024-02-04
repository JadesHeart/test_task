package main

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"test_task/internal/config"
	"test_task/internal/pkg/handlers/auth"
	"test_task/internal/pkg/handlers/ping"
	"test_task/internal/pkg/lib/sl"
	"test_task/internal/pkg/repository"
	"test_task/internal/pkg/service"
	"test_task/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		logger.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	logger.Info("Successful init database")

	repos := repository.NewRepository(storage, logger)
	services := service.NewService(repos, logger)

	r := gin.Default()

	r.POST("/auth", auth.New(logger, services))
	r.POST("/ping", ping.New(logger, services))

	err = r.Run(cfg.Address)
	if err != nil {
		logger.Error("error start server", sl.Err(err))
	}

	logger.Info("starting server", slog.String("address", cfg.Address))

}

func setupLogger(env string) *slog.Logger {

	var logger *slog.Logger

	switch env {
	case "local":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return logger
}
