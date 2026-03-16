package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/deshortone/ledger-system/internal/testhelpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
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

	ledgerDbGlobalSeed(ctx)

	code := m.Run()
	os.Exit(code)
}

func ledgerDbGlobalSeed(ctx context.Context) {
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
	_, err = pool.Exec(ctx, `
        INSERT INTO identity.accounts (id, identifier, user_id, created_at, account_type, currency, status)
		OVERRIDING SYSTEM VALUE
        VALUES ($1, $2, $3,  NOW(), $4, $5, $6), ($7, $8, $3, NOW(), $4, $5, $6);
    `, 1, uuid.New(), 1, "account type", "GBP", "available", 2, uuid.New())
	if err != nil {
		panic(err)
	}

	// add transfers
	_, err = pool.Exec(ctx, `
        INSERT INTO transfer.transfers (id, identifier, from_account_id, to_account_id, amount, status, created_at)
		OVERRIDING SYSTEM VALUE
        VALUES ($1, $2, $3, $4, $5, $6, NOW());
    `, 1, uuid.New(), 1, 2, 100, "posted")
	if err != nil {
		panic(err)
	}

	// add account balances
	timee, err := time.Parse("2006-01-02 15:04:05 -0700", "2026-03-15 12:00:00 +0000")
	if err != nil {
		panic(err)
	}
	_, err = pool.Exec(ctx, `
        INSERT INTO ledger.account_balances (account_id, available_balance, updated_at)
        VALUES ($1, $2, $3), ($4, $2, $3);
    `, 1, 100, timee, 2)
	if err != nil {
		panic(err)
	}
}
