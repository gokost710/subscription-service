package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gokost710/subscription-service/internal/app"
	"github.com/gokost710/subscription-service/internal/config"
	"github.com/gokost710/subscription-service/internal/logger"
	"github.com/gokost710/subscription-service/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Log.Level)

	db, err := postgres.New(context.Background(), cfg.DB.DSN())
	if err != nil {
		log.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}

	log.Info("connected to postgres", "host", cfg.DB.Host, "port", cfg.DB.Port, "db", cfg.DB.Name)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	application := app.New(cfg, log, db)
	if err := application.Run(ctx); err != nil {
		log.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}
