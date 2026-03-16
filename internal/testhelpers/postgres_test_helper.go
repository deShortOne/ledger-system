package testhelpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// PostgresTestContainer wraps an integration test Postgres container and connection pool.
type PostgresTestContainer struct {
	Ctx       context.Context
	Container *postgres.PostgresContainer
	Pool      *pgxpool.Pool
	DBURL     string
}

// NewPostgresTestContainer starts a testcontainers Postgres instance and returns a pool and cleanup.
func NewPostgresTestContainer(ctx context.Context, initScriptPath string) (*PostgresTestContainer, func(), error) {
	container, err := postgres.Run(
		ctx,
		"postgres:18-alpine",
		postgres.WithInitScripts(initScriptPath),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("starting postgres container: %w", err)
	}

	dbURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("getting postgres connection string: %w", err)
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("creating pgx pool: %w", err)
	}

	tc := &PostgresTestContainer{
		Ctx:       ctx,
		Container: container,
		Pool:      pool,
		DBURL:     dbURL,
	}

	cleanup := func() {
		pool.Close()
		_ = container.Terminate(ctx)
	}

	return tc, cleanup, nil
}

const defaultMigrationScriptRelPath = "migration/postgres/001_initial.up.sql"

func getRepoRootPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working dir: %w", err)
	}

	cur := wd
	for {
		if _, err := os.Stat(filepath.Join(cur, "go.mod")); err == nil {
			return cur, nil
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}

	return "", fmt.Errorf("repo root not found from cwd %s", wd)
}

func GetDefaultMigrationScriptPath() (string, error) {
	repoRoot, err := getRepoRootPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(repoRoot, defaultMigrationScriptRelPath), nil
}

func NewPostgresTestContainerWithDefaultMigration(ctx context.Context) (*PostgresTestContainer, func(), error) {
	migrationPath, err := GetDefaultMigrationScriptPath()
	if err != nil {
		return nil, nil, err
	}
	return NewPostgresTestContainer(ctx, migrationPath)
}
