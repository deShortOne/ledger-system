package api

import (
	"context"
	"fmt"

	"github.com/deshortone/ledger-system/internal/identity"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Pool *pgxpool.Pool
}

func RegisterRoutes(ctx context.Context, r *gin.Engine) (App, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return App{}, fmt.Errorf("failed to load config: %w", err)
	}

	app, err := setupPostgresDatabase(ctx, cfg)
	if err != nil {
		return App{}, err
	}

	// setup modules
	identityModule := identity.SetupIdentityModule(app.Pool)

	// setup routes
	version1 := r.Group("/api")
	identityModule.Handler.RegisterRoutes(version1)

	return app, nil
}

func (a *App) Close() {
	if a.Pool != nil {
		a.Pool.Close()
	}
}

func setupPostgresDatabase(ctx context.Context, cfg *Config) (App, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return App{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := runMigrations(pool); err != nil {
		pool.Close()
		return App{}, fmt.Errorf("failed to run migrations: %w", err)
	}
	return App{
		Pool: pool,
	}, nil
}

func runMigrations(pool *pgxpool.Pool) error {
	m, err := migrate.New(
		"file://./migration/postgres",
		pool.Config().ConnString())
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	//nolint
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migration: %w", err)
	}
	return nil
}
