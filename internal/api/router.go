package api

import (
	"context"
	"fmt"
	"time"

	"github.com/deshortone/ledger-system/internal/identity"
	"github.com/deshortone/ledger-system/internal/ledger"
	"github.com/deshortone/ledger-system/internal/middleware"
	"github.com/deshortone/ledger-system/internal/platform"
	"github.com/deshortone/ledger-system/internal/transfer"
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

	readinessGate := middleware.NewReadinessGate()

	app, err := setupPostgresDatabase(ctx, cfg, readinessGate)
	if err != nil {
		return App{}, err
	}

	// setup modules
	ledgerModule := ledger.SetupLedgerModule(app.Pool)
	identityModule := identity.SetupIdentityModule(app.Pool, ledgerModule.LedgerService)
	transferModule := transfer.SetupTransferModule(app.Pool, ledgerModule.LedgerService)
	platformModule := platform.SetupPlatformModule(app.Pool)

	// setup routes
	version1 := r.Group("/api")
	version1.Use(middleware.ReadinessGateMiddleware(readinessGate))
	identityModule.Handler.RegisterRoutes(version1)
	transferModule.Handler.RegisterRoutes(version1)

	platformModule.PlatformHandler.RegisterRoutes(r.Group("/health"))

	return app, nil
}

func (a *App) Close() {
	if a.Pool != nil {
		a.Pool.Close()
	}
}

func setupPostgresDatabase(ctx context.Context, cfg *Config, readinessGate *middleware.ReadinessGate) (App, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return App{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		readinessGate.SetReady(false)

		for range ticker.C {
			if err := runMigrations(pool); err == nil {
				break
			}
			fmt.Println("Migration failed to run. App is live but not ready.")
		}

		for range ticker.C {
			err = pool.Ping(ctx)
			readinessGate.SetReady(err == nil)
		}
	}()

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
