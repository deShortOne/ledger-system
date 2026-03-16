package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/deshortone/ledger-system/internal/testhelpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var (
	dbURL string
	pool  *pgxpool.Pool

	account1Id uuid.UUID
	account2Id uuid.UUID
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

	transferDbGlobalSeed(ctx)

	code := m.Run()
	os.Exit(code)
}

func transferDbGlobalSeed(ctx context.Context) {
	var err error
	// add user
	_, err = pool.Exec(ctx, `
        INSERT INTO identity.users (id, identifier, first_name, last_name)
		OVERRIDING SYSTEM VALUE
        VALUES ($1, $2, $3, $4);
    `, 1, uuid.New(), "first", "last")
	if err != nil {
		panic(err)
	}

	// add accounts, only id and identifier are different here
	account1Id = uuid.New()
	account2Id = uuid.New()
	_, err = pool.Exec(ctx, `
        INSERT INTO identity.accounts (id, identifier, user_id, created_at, account_type, currency, status)
		OVERRIDING SYSTEM VALUE
        VALUES ($1, $2, $3,  NOW(), $4, $5, $6), ($7, $8, $3, NOW(), $4, $5, $6);
    `, 1, account1Id, 1, "account type", "GBP", "available", 2, account2Id)
	if err != nil {
		panic(err)
	}
}
