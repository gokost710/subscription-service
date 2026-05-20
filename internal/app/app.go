package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gokost710/subscription-service/internal/config"
	"github.com/gokost710/subscription-service/internal/http/router"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg    *config.Config
	log    *slog.Logger
	db     *pgxpool.Pool
	server *http.Server
}

func New(cfg *config.Config, log *slog.Logger, db *pgxpool.Pool) *App {
	return &App{
		cfg: cfg,
		log: log,
		db:  db,
		server: &http.Server{
			Addr:              cfg.HTTP.Addr(),
			Handler:           router.New(),
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func (a *App) Run(ctx context.Context) error {
	defer a.db.Close()

	errCh := make(chan error, 1)

	go func() {
		a.log.Info("starting http server", "addr", a.cfg.HTTP.Addr(), "config", a.cfg.String())

		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return a.shutdown()
	case err := <-errCh:
		return err
	}
}

func (a *App) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.log.Info("shutting down http server")
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	a.log.Info("http server stopped")

	return nil
}
