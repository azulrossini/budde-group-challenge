package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"

	"github.com/azulrossini/budde-group-challenge/backend/internal/config"
	"github.com/azulrossini/budde-group-challenge/backend/internal/handlers"
	"github.com/azulrossini/budde-group-challenge/backend/internal/repository"
	"github.com/azulrossini/budde-group-challenge/backend/internal/service"
	"github.com/azulrossini/budde-group-challenge/backend/internal/web"
	"github.com/azulrossini/budde-group-challenge/backend/migrations"
)

const (
	startupTimeout  = 10 * time.Second
	shutdownTimeout = 10 * time.Second
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open pool: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	svc := service.New(repo)

	spa, err := web.New()
	if err != nil {
		log.Fatalf("web: %v", err)
	}

	handler := handlers.NewRouter(handlers.New(svc), spa)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	stop()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	slog.Info("shutdown complete")
}

func runMigrations(databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), startupTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}
