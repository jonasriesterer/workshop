// Package main ist der Einstiegspunkt der Serveranwendung.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/myorg/myservice/config"
	"github.com/myorg/myservice/internal/seed"
	"github.com/myorg/myservice/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()
	slog.Info("configuration loaded", "addr", cfg.Addr(), "testMode", cfg.TestMode)

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		slog.Error("failed to create db pool", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("db ping failed", "err", err)
		os.Exit(1)
	}
	slog.Info("database connected", "dsn_host", cfg.DBHost, "db", cfg.DBName)

	if cfg.TestMode {
		slog.Info("TEST_MODE enabled – running migrations and seeding")
		if err := runMigrations(ctx, pool); err != nil {
			slog.Error("migrations failed", "err", err)
			os.Exit(1)
		}
		csvDir := csvSeedsDir()
		if err := seed.Load(ctx, pool, csvDir); err != nil {
			slog.Error("seeding failed", "err", err)
			os.Exit(1)
		}
	}

	srv := server.New(cfg, pool)

	go func() {
		slog.Info("server listening", "addr", cfg.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutCtx); err != nil {
		slog.Error("forced shutdown", "err", err)
	}
	slog.Info("server stopped")
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrationsDir := migrationDir()
	for _, file := range []string{"001_drop.sql", "002_create.sql"} {
		path := filepath.Join(migrationsDir, file)
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return err
		}
		slog.Info("migration executed", "file", file)
	}
	return nil
}

// migrationDir returns the absolute path to db/migrations relative to the binary's source.
func migrationDir() string {
	_, filename, _, _ := runtime.Caller(0)
	// filename is .../cmd/server/main.go → go up 3 levels to project root
	root := filepath.Join(filepath.Dir(filename), "..", "..")
	return filepath.Join(root, "db", "migrations")
}

// csvSeedsDir returns the absolute path to db/seeds.
func csvSeedsDir() string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "..", "..")
	return filepath.Join(root, "db", "seeds")
}
