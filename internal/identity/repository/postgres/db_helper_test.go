package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/deshortone/ledger-system/internal/testhelpers"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbURL string
	pool  *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	tc, cleanup, err := testhelpers.NewPostgresTestContainerWithDefaultMigration(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	dbURL = tc.DBURL
	pool = tc.Pool

	code := m.Run()
	os.Exit(code)
}
